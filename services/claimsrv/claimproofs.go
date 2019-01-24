package claimsrv

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// CheckProofOfClaim checks the claim proofs from the bottom to the top are valid and not revoked, and that the top root is signed by relayAddr.
func VerifyProofOfClaim(relayAddr common.Address, pc *ProofOfClaim) (bool, error) {
	// For now we only allow proof verification of Nameserver (one level) and
	// Relay (two levels: relay + user)
	if len(pc.Proofs) > 2 || len(pc.Proofs) < 1 {
		return false, fmt.Errorf("Invalid number of partial proofs")
	}
	// Top root signature (by Relay) verification
	if !utils.VerifySigBytesDate(relayAddr, pc.Signature, pc.Proofs[len(pc.Proofs)-1].Root[:], pc.Date) {
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
		claimType, claimVer := core.GetClaimTypeVersionFromData(&leafNext.Data)
		core.SetClaimTypeVersionInData(&leafNext.Data, claimType, claimVer+1)
		if !merkletree.VerifyProof(rootKey, mtpNoEx, leafNext.HIndex(), leafNext.HValue()) {
			return false, fmt.Errorf("Mtp1 at lvl %v doesn't match with the root", i)
		}

		if i == len(pc.Proofs)-1 {
			break
		} else if proof.Aux == nil {
			return false, fmt.Errorf("partial proof at lvl %v doesn't contain auxiliary data", i)
		}

		// Create the set root key claim for the next level
		claim := core.NewClaimSetRootKey(proof.Aux.EthAddr, *rootKey)
		claim.Version = proof.Aux.Version
		claim.Era = proof.Aux.Era
		leaf = claim.Entry()
	}
	return true, nil
}

// CheckProofOfClaimUser checks the Merkle Proof of the Claim, the SetRootClaim,
// and the non revocation proof of both claims
func CheckProofOfClaimUser(relayAddr common.Address, pc ProofOfClaimUser, numLevels int) bool {
	node, err := merkletree.NewNodeFromBytes(pc.ClaimProof.Leaf)
	if err != nil {
		return false
	}
	node.Entry.HIndex()
	pf, err := merkletree.NewProofFromBytes(pc.ClaimProof.Proof)
	if err != nil {
		return false
	}
	vClaimProof := merkletree.VerifyProof(&pc.ClaimProof.Root, pf,
		node.Entry.HIndex(), node.Entry.HValue())

	node, err = merkletree.NewNodeFromBytes(pc.SetRootClaimProof.Leaf)
	if err != nil {
		return false
	}
	node.Entry.HIndex()
	pf, err = merkletree.NewProofFromBytes(pc.SetRootClaimProof.Proof)
	if err != nil {
		return false
	}
	vSetRootClaimProof := merkletree.VerifyProof(&pc.SetRootClaimProof.Root, pf,
		node.Entry.HIndex(), node.Entry.HValue())

	node, err = merkletree.NewNodeFromBytes(pc.ClaimNonRevocationProof.Leaf)
	if err != nil {
		return false
	}
	node.Entry.HIndex()
	pf, err = merkletree.NewProofFromBytes(pc.ClaimNonRevocationProof.Proof)
	if err != nil {
		return false
	}
	vClaimNonRevocationProof := merkletree.VerifyProof(&pc.ClaimNonRevocationProof.Root, pf,
		node.Entry.HIndex(), node.Entry.HValue())

	node, err = merkletree.NewNodeFromBytes(pc.SetRootClaimNonRevocationProof.Leaf)
	if err != nil {
		return false
	}
	node.Entry.HIndex()
	pf, err = merkletree.NewProofFromBytes(pc.SetRootClaimNonRevocationProof.Proof)
	if err != nil {
		return false
	}
	vSetRootClaimNonRevocationProof := merkletree.VerifyProof(&pc.SetRootClaimNonRevocationProof.Root, pf,
		node.Entry.HIndex(), node.Entry.HValue())

	// additional, check caducity of the pc.Date

	// check signature of the ProofOfClaim.SetRootClaimProof.Root with the identity of the Relay
	// checking this Root and the four Merkle Proofs, we check the full ProofOfClaim
	if !utils.VerifySigBytesDate(relayAddr, pc.Signature, pc.SetRootClaimProof.Root[:], pc.Date) {
		return false
	}

	if vClaimProof && vSetRootClaimProof && vClaimNonRevocationProof && vSetRootClaimNonRevocationProof {
		return true
	}
	return false
}

// CheckKSignInIDdb checks that a given KSign is in an AuthorizeKSignClaim in the Identity Merkle Tree (in this version, as the Merkle Tree don't allows to delete data, the verification only needs to check if the AuthorizeKSignClaim is in the key-value)
func CheckKSignInIDdb(mt *merkletree.MerkleTree, kSignPk *ecdsa.PublicKey) bool {
	claimAuthorizeKSign := core.NewClaimAuthorizeKSignSecp256k1(kSignPk)
	entry := claimAuthorizeKSign.Entry()
	node := merkletree.NewNodeLeaf(entry)
	nodeGot, err := mt.GetNode(node.Key())
	if err != nil {
		return false
	}
	if !bytes.Equal(node.Value(), nodeGot.Value()) {
		return false
	}

	// non revocation
	claimAuthorizeKSign.Version++
	entry = claimAuthorizeKSign.Entry()
	node = merkletree.NewNodeLeaf(entry)
	_, err = mt.GetNode(node.Key())
	if err != db.ErrNotFound {
		return false
	}

	return true
}
