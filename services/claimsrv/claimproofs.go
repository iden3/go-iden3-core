package claimsrv

import (
	"bytes"
	"crypto/ecdsa"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
)

// CheckKSignInIddb checks that a given KSign is in an AuthorizeKSignClaim in the Identity Merkle Tree (in this version, as the Merkle Tree don't allows to delete data, the verification only needs to check if the AuthorizeKSignClaim is in the key-value)
func CheckKSignInIddb(mt *merkletree.MerkleTree, kSignPk *ecdsa.PublicKey) bool {
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

// CheckKSignBabyJubInIddb checks that a given KSign is in an AuthorizeKSignClaim in the Identity Merkle Tree (in this version, as the Merkle Tree don't allows to delete data, the verification only needs to check if the AuthorizeKSignClaim is in the key-value)
func CheckKSignBabyJubInIddb(mt *merkletree.MerkleTree, kSignPk *babyjub.PublicKey) bool {
	claimAuthorizeKSign := core.NewClaimAuthorizeKSignBabyJub(kSignPk)
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
