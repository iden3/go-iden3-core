package core

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/iden3/go-iden3-crypto/babyjub"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/merkletree"
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
		fmt.Fprintf(buf, "aux: Version:%v Era:%v Id:%v\n", pcp.Aux.Version, pcp.Aux.Era,
			common3.HexEncode(pcp.Aux.Id[:]))
	}
	return buf.String()
}

// SetRootAux is the auxiliary data to build the set root claim from a root in
// a partial proof of claim.
type SetRootAux struct {
	Version uint32 `json:"version" binding:"required"`
	Era     uint32 `json:"era" binding:"required"`
	Id      ID     `json:"id" binding:"required"`
}

// ProofClaim is a complete proof of a claim that includes all the proofs of
// existence and non-existence for mutliple levels from the leaf of a tree to
// the signed root of possibly another tree whose root binding:"required".
type ProofClaim struct {
	Proofs    []ProofClaimPartial    `json:"proofs" binding:"required"`
	Leaf      *merkletree.Data       `json:"leaf" binding:"required"`
	Date      int64                  `json:"date" binding:"required"`
	Signature *babyjub.SignatureComp `json:"signature" binding:"required"` // signature of the Root of the Relay
	Signer    ID                     `json:"signer" binding:"required"`
}

func (pc *ProofClaim) String() string {
	buf := bytes.NewBufferString("ProofClaim:\n")
	if pc.Signature != nil {
		fmt.Fprintf(buf, "signature: %v\n", common3.HexEncode(pc.Signature[:]))
	}
	fmt.Fprintf(buf, "date: %v\n", time.Unix(pc.Date, 0))
	fmt.Fprintf(buf, "leaf: %v\n", pc.Leaf)
	fmt.Fprintf(buf, "proofs:\n")
	for i, proof := range pc.Proofs {
		fmt.Fprintf(buf, "%v: { %v}\n", i, &proof)
	}
	return buf.String()
}

// CheckProofClaim checks the claim proofs from the bottom to the top are valid and not revoked, and that the top root is signed by relayAddr.
// WARNING TODO currently the Root signature verification is disabled, see comment in line 82
func VerifyProofClaim(operationalPk *babyjub.PublicKey, pc *ProofClaim) (bool, error) {
	// For now we only allow proof verification of Nameserver (one level) and
	// Relay (two levels: relay + user)
	if len(pc.Proofs) > 2 || len(pc.Proofs) < 1 {
		return false, fmt.Errorf("Invalid number of partial proofs")
	}

	// TODO currently this is verifying with relayAddr, that comes directly from privateKey
	// in next iteration needs to verify that the signature is performed
	// by a private key authorized in a claim under the ID merkle tree
	// or even not verify the signature in this function, and check that the Root is in the blockchain for the relayID, if is there will mean that is made by that relay (as the relay needs to sign it to perform the transaction)
	// if is this last option, somewhere need to check that 'relayAddr' is equal to the ProofClaim emiter address, or remove that input as is checked outside this function
	/*
		if pc.Signature == nil {
			return false, fmt.Errorf("No signature in the ProofClaim")
		}
		// Top root signature (by Relay) verification
			if !utils.VerifySigEthMsgDate(relayAddr, pc.Signature, pc.Proofs[len(pc.Proofs)-1].Root[:], pc.Date) {
				return false, fmt.Errorf("Invalid signature")
			}
	*/

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
		claim, err := NewClaimSetRootKey(proof.Aux.Id, *rootKey)
		if err != nil {
			return false, err
		}
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

// GetClaimProofByHi given a Hash(index) (Hi) and an id, returns the Claim
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

type PredicateProof struct {
	LeafEntry               *merkletree.Entry
	MtpNonExistInOldRoot    *merkletree.Proof
	MtpExist                *merkletree.Proof
	MtpNonExistNextVersion  *merkletree.Proof
	MtpExistPreviousVersion *merkletree.Proof
	OldRoot                 *merkletree.Hash
	Root                    *merkletree.Hash
}

func GetPreviousVersionEntry(entry *merkletree.Entry) (*merkletree.Entry, error) {
	claimType, claimVer := GetClaimTypeVersionFromData(&entry.Data)
	if claimVer == 0 {
		return nil, errors.New("claim is in version 0, can not exist a previous version")
	}
	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], entry.Data[:])
	SetClaimTypeVersionInData(leafDataCpy, claimType, claimVer-1)
	entry1 := merkletree.Entry{
		Data: *leafDataCpy,
	}
	return &entry1, nil
}
func GetNextVersionEntry(entry *merkletree.Entry) *merkletree.Entry {
	claimType, claimVer := GetClaimTypeVersionFromData(&entry.Data)
	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], entry.Data[:])
	SetClaimTypeVersionInData(leafDataCpy, claimType, claimVer+1)
	entry1 := merkletree.Entry{
		Data: *leafDataCpy,
	}
	return &entry1
}

// GetPredicateProof, ϕ_min
// checks that:
// - 0: tree is updated incrementally
// 	- claim position was empty in oldRoot
// - 1: claim is added correctly
// 	- claim position contains the claim in currentRoot
// - 2: claim is not revocated
// 	- claim (version+1) is empty in currentRoot
// in case that the claim version != 0:
// - 3: claim is at the expected version
//	- claim (version-1) exist in oldRoot
// - 4: current version is incremental from the last one
//	- siblings of check_0 are inside siblings of check_1
//
// *TODO The output format will depend on the zkSnark inputs format (not specified yet)
func GetPredicateProof(mt *merkletree.MerkleTree, oldRoot, hi *merkletree.Hash) (*PredicateProof, error) {
	// proof_0: that claim position was empty in oldRoot
	mtpNonExistInOldRoot, err := mt.GenerateProof(hi, oldRoot)
	if err != nil {
		return nil, err
	}

	// proof_1: that claim position contains the claim in newRoot
	mtpExist, err := mt.GenerateProof(hi, nil)
	if err != nil {
		return nil, err
	}

	// proof_2: that claim (version+1) is empty in newRoot
	// get the value in the hi position
	leafData, err := mt.GetDataByIndex(hi)
	if err != nil {
		return nil, err
	}
	mtpNonExistNextVersion, err := GetNonRevocationMTProof(mt, leafData, hi)
	if err != nil {
		return nil, err
	}

	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], leafData[:])
	entry := merkletree.Entry{
		Data: *leafDataCpy,
	}

	// checks 3 and 4 are necessary if the claim.Version != 0
	var mtpExistPrevVersion *merkletree.Proof
	_, v := getClaimTypeVersion(&entry)
	if v != 0 {
		// version is not 0 (v!=0), so we need to provide proof_3 and proof_4

		// proof_3: that claim (version-1) is not empty in oldRoot
		entryPrevVersion, err := GetPreviousVersionEntry(&entry)
		if err != nil {
			return nil, err
		}
		mtpExistPrevVersion, err = mt.GenerateProof(entryPrevVersion.HIndex(), oldRoot)
		if err != nil {
			return nil, err
		}

		// proof_4: that siblings of check_0 are inside of siblings of check_1
		// this proof don't needs more additional data
		// is something that the verifier needs to verify with the data from the other proofs
	}

	predicateProof := &PredicateProof{
		LeafEntry:               &entry,
		MtpNonExistInOldRoot:    mtpNonExistInOldRoot,   // proof_0
		MtpExist:                mtpExist,               // proof_1
		MtpNonExistNextVersion:  mtpNonExistNextVersion, // proof_2
		MtpExistPreviousVersion: mtpExistPrevVersion,    // proof_3
		OldRoot:                 oldRoot,
		Root:                    mt.RootKey(),
	}
	return predicateProof, nil
}

// VerifyPredicateProof, ϕ_min
// checks that:
// - 0: tree is updated incrementally
// 	- claim position was empty in oldRoot
// - 1: claim is added correctly
// 	- claim position contains the claim in currentRoot
// - 2: claim is not revocated
// 	- claim (version+1) is empty in currentRoot
// in case that the claim version != 0:
// - 3: claim is at the expected version
//	- claim (version-1) exist in oldRoot
// - 4: current version is incremental from the last one
//	- siblings of check_0 are inside of check_1
//
// *TODO The input format will depend on the zkSnark inputs format (not specified yet)
func VerifyPredicateProof(p *PredicateProof) bool {
	if bytes.Equal(p.Root.Bytes(), p.OldRoot.Bytes()) {
		return false
	}

	// check_0: that claim position was empty in oldRoot
	if p.MtpNonExistInOldRoot.Existence {
		// should be a proof of non existence, if not, verification fails
		return false
	}
	if !merkletree.VerifyProof(p.OldRoot, p.MtpNonExistInOldRoot, p.LeafEntry.HIndex(), p.LeafEntry.HValue()) {
		return false
	}

	// check_1: that claim position contains the claim in currentRoot
	if !p.MtpExist.Existence {
		// should be a proof of existence, if not, verification fails
		return false
	}
	if !merkletree.VerifyProof(p.Root, p.MtpExist, p.LeafEntry.HIndex(), p.LeafEntry.HValue()) {
		return false
	}

	// check_2: that claim (version+1) is empty in currentRoot
	if p.MtpNonExistNextVersion.Existence {
		// should be a proof of non existence, if not, verification fails
		return false
	}
	entry1 := GetNextVersionEntry(p.LeafEntry)
	if !merkletree.VerifyProof(p.Root, p.MtpNonExistNextVersion, entry1.HIndex(), entry1.HValue()) {
		return false
	}

	// checks 3 and 4 are necessary if the claim.Version != 0
	_, v := getClaimTypeVersion(p.LeafEntry)
	if v == 0 {
		// if version == 0, return true expected checks have passed
		return true
	}

	// check_3: that claim (version-1) is not empty in oldRoot
	if !p.MtpExistPreviousVersion.Existence {
		// should be a proof of existence, if not, verification fails
		return false
	}
	entryPrevVersion, err := GetPreviousVersionEntry(p.LeafEntry)
	if err != nil {
		// if err!=nil means that there is no previous version possible, as the current version is 0
		return false
	}
	if !merkletree.VerifyProof(p.OldRoot, p.MtpExistPreviousVersion, entryPrevVersion.HIndex(), entryPrevVersion.HValue()) {
		return false
	}

	// check_4: check that siblings of check_0 are inside the siblings of check_1
	// p.MtpNonExistInOldRoot.Siblings == p.MtpExist.Siblings[:len(p.MtpNonExistInOldRoot.Siblings)]
	for i := 0; i < len(p.MtpNonExistInOldRoot.Siblings); i++ {
		if !bytes.Equal(p.MtpNonExistInOldRoot.Siblings[i].Bytes(), p.MtpExist.Siblings[i].Bytes()) {
			return false
		}
	}

	return true
}
