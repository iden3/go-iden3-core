package claimsrv

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// CheckProofOfClaim checks the Merkle Proof of the Claim, the SetRootClaim, and the non revocation proof of both claims
func CheckProofOfClaim(relayAddr common.Address, pc ProofOfClaimUser, numLevels int) bool {
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
	dateBytes, err := utils.Uint64ToEthBytes(pc.Date)
	if err != nil {
		return false
	}
	rootdate := pc.SetRootClaimProof.Root[:]
	rootdate = append(rootdate, dateBytes...)
	rootdateHash := utils.HashBytes(rootdate)
	if !utils.VerifySig(relayAddr, pc.Signature, rootdateHash[:]) {
		return false
	}

	if vClaimProof && vSetRootClaimProof && vClaimNonRevocationProof && vSetRootClaimNonRevocationProof {
		return true
	}
	return false
}

// CheckKSignInIDdb checks that a given KSign is in an AuthorizeKSignClaim in the Identity Merkle Tree (in this version, as the Merkle Tree don't allows to delete data, the verification only needs to check if the AuthorizeKSignClaim is in the key-value)
func CheckKSignInIDdb(mt *merkletree.MerkleTree, ksign common.Address) bool {
	// generate the AuthorizeKSignClaim
	var tmpFakeAy merkletree.ElemBytes
	claimAuthorizeKSign := core.NewClaimAuthorizeKSign(false, tmpFakeAy) // TODO ethAddress to pubK
	entry := claimAuthorizeKSign.Entry()
	node := merkletree.NewNodeLeaf(entry)
	nodeGetted, err := mt.GetNode(node.Key())
	if err != nil {
		return false
	}
	if !bytes.Equal(node.Value(), nodeGetted.Value()) {
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
