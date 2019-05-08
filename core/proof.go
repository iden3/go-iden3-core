package core

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

var (
	ErrRevokedClaim = errors.New("the claim is revoked: the next version exists")
)

// ProofClaimPartial is a proof of existence and non-existence of a claim in
// a single tree (only one level).
type ProofClaimPartial struct {
	Mtp0 *merkletree.Proof `json:"mtp0" binding:"required"`
	Mtp1 *merkletree.Proof `json:"mtp1" binding:"required"`
	Root *merkletree.Hash  `json:"root" binding:"required"`
	Aux  *SetRootAux       `json:"aux" binding:"required"`
}

func (pcp *ProofClaimPartial) String() string {
	buf := bytes.NewBufferString("ProofClaimPartial:\n")
	fmt.Fprintf(buf, "mtp0: %v\n", pcp.Mtp0)
	fmt.Fprintf(buf, "mtp0: %v\n", pcp.Mtp1)
	fmt.Fprintf(buf, "root: %v\n", pcp.Root)
	if pcp.Aux != nil {
		fmt.Fprintf(buf, "aux: Version:%v Era:%v IdAddr:%v\n", pcp.Aux.Version, pcp.Aux.Era,
			common3.HexEncode(pcp.Aux.IdAddr[:]))
	}
	return buf.String()
}

// SetRootAux is the auxiliary data to build the set root claim from a root in
// a partial proof of claim.
type SetRootAux struct {
	Version uint32 `json:"version" binding:"required"`
	Era     uint32 `json:"era" binding:"required"`
	IdAddr  ID     `json:"idAddr" binding:"required"`
}

// ProofClaim is a complete proof of a claim that includes all the proofs of
// existence and non-existence for mutliple levels from the leaf of a tree to
// the signed root of possibly another tree whose root binding:"required".
type ProofClaim struct {
	Proofs    []ProofClaimPartial    `json:"proofs" binding:"required"`
	Leaf      *merkletree.Data       `json:"leaf" binding:"required"`
	Date      int64                  `json:"date" binding:"required"`
	Signature *utils.SignatureEthMsg `json:"signature" binding:"required"` // signature of the Root of the Relay
	Signer    common.Address         `json:"signer" binding:"required"`
}

func (pc *ProofClaim) String() string {
	buf := bytes.NewBufferString("ProofClaim:\n")
	fmt.Fprintf(buf, "signature: %v\n", common3.HexEncode(pc.Signature[:]))
	fmt.Fprintf(buf, "date: %v\n", time.Unix(pc.Date, 0))
	fmt.Fprintf(buf, "leaf: %v\n", pc.Leaf)
	fmt.Fprintf(buf, "proofs:\n")
	for i, proof := range pc.Proofs {
		fmt.Fprintf(buf, "%v: { %v}\n", i, &proof)
	}
	return buf.String()
}

// CheckProofClaim checks the claim proofs from the bottom to the top are valid and not revoked, and that the top root is signed by relayAddr.
func VerifyProofClaim(relayAddr common.Address, pc *ProofClaim) (bool, error) {
	// For now we only allow proof verification of Nameserver (one level) and
	// Relay (two levels: relay + user)
	if len(pc.Proofs) > 2 || len(pc.Proofs) < 1 {
		return false, fmt.Errorf("Invalid number of partial proofs")
	}
	// Top root signature (by Relay) verification
	if !utils.VerifySigEthMsgDate(relayAddr, pc.Signature, pc.Proofs[len(pc.Proofs)-1].Root[:], pc.Date) {
		return false, fmt.Errorf("Invalid signature")
	}

	leaf := &merkletree.Entry{Data: *pc.Leaf}
	leafNext := &merkletree.Entry{}
	rootKey := &merkletree.Hash{}
	for i, proof := range pc.Proofs {
		mtpEx := proof.Mtp0
		mtpNoEx := proof.Mtp1
		rootKey = proof.Root

		*leafNext = *leaf

		// Proof of existence verification
		if !mtpEx.Existence {
			return false, fmt.Errorf("Mtp0 at lvl %v is a non-existence proof", i)
		}
		if !merkletree.VerifyProof(rootKey, mtpEx, leaf.HIndex(), leaf.HValue()) {
			return false, fmt.Errorf("Mtp0 at lvl %v doesn't match with the root", i)
		}

		// Proof of non-existence of next version (revocation) verification
		if mtpNoEx.Existence {
			return false, fmt.Errorf("Mtp1 at lvl %v is an existence proof", i)
		}
		claimType, claimVer := GetClaimTypeVersionFromData(&leafNext.Data)
		SetClaimTypeVersionInData(&leafNext.Data, claimType, claimVer+1)
		if !merkletree.VerifyProof(rootKey, mtpNoEx, leafNext.HIndex(), leafNext.HValue()) {
			return false, fmt.Errorf("Mtp1 at lvl %v doesn't match with the root", i)
		}

		if i == len(pc.Proofs)-1 {
			break
		} else if proof.Aux == nil {
			return false, fmt.Errorf("partial proof at lvl %v doesn't contain auxiliary data", i)
		}

		// Create the set root key claim for the next level
		claim := NewClaimSetRootKey(proof.Aux.IdAddr, *rootKey)
		claim.Version = proof.Aux.Version
		claim.Era = proof.Aux.Era
		leaf = claim.Entry()
	}
	return true, nil
}

// GetNonRevocationMTProof is a helper function to return a proof of non
// existence of the following version of a given claim (leafData).  If the
// following version exists, an error is returned.
func GetNonRevocationMTProof(mt *merkletree.MerkleTree, leafData *merkletree.Data, hi *merkletree.Hash) (*merkletree.Proof, error) {
	claimType, claimVersion := GetClaimTypeVersionFromData(leafData)

	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], leafData[:])
	SetClaimTypeVersionInData(leafDataCpy, claimType, claimVersion+1)
	entry := merkletree.Entry{
		Data: *leafDataCpy,
	}
	proof, err := mt.GenerateProof(entry.HIndex(), nil)
	if err != nil {
		return nil, err
	}
	if proof.Existence {
		return nil, ErrRevokedClaim
	}
	return proof, nil
}

// GetClaimProofByHi given a Hash(index) (Hi) and an idAddr, returns the Claim
// in that Hi position inside the User merkletree, it's proof of existence and
// of non-revocation, and the proof of existence and of non-revocation for the
// set root claim in the relay tree, all in the form of a ProofClaim.  The
// result is not yet signed and has no timestamp.
func GetClaimProofByHi(mt *merkletree.MerkleTree, hi *merkletree.Hash) (*ProofClaim, error) {
	// get the value in the hi position
	leafData, err := mt.GetDataByIndex(hi)
	if err != nil {
		return nil, err
	}

	// get the MT proof of existence of the claim and the non-existence of
	// the claim's next version in the Relay Tree
	mtpExist, err := mt.GenerateProof(hi, nil)
	if err != nil {
		return nil, err
	}
	mtpNonExist, err := GetNonRevocationMTProof(mt, leafData, hi)
	if err != nil {
		return nil, err
	}

	rootKey := mt.RootKey()

	proofClaimPartial := ProofClaimPartial{
		Mtp0: mtpExist,
		Mtp1: mtpNonExist,
		Root: rootKey,
		Aux:  nil,
	}
	proofClaim := ProofClaim{
		Proofs:    []ProofClaimPartial{proofClaimPartial},
		Leaf:      leafData,
		Date:      0,
		Signature: nil,
	}

	return &proofClaim, nil
}
