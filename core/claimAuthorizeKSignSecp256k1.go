package core

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3-core/merkletree"
)

// ClaimAuthorizeKSignSecp256k1 is a claim to autorize a public key for signing.
type ClaimAuthorizeKSignSecp256k1 struct {
	// Version is the claim version.
	Version uint32
	// PubKey is the ECDSA public key.
	PubKey *ecdsa.PublicKey
}

// NewClaimAuthorizeKSignSecp256k1 returns a ClaimAuthorizeKSignSecp256k1 with the given elliptic
// public key parameters.
func NewClaimAuthorizeKSignSecp256k1(pk *ecdsa.PublicKey) *ClaimAuthorizeKSignSecp256k1 {
	return &ClaimAuthorizeKSignSecp256k1{
		Version: 0,
		PubKey:  pk,
	}
}

// NewClaimAuthorizeKSignSecp256k1FromEntry deserializes a ClaimAuthorizeKSignSecp256k1 from an Entry.
func NewClaimAuthorizeKSignSecp256k1FromEntry(e *merkletree.Entry) (*ClaimAuthorizeKSignSecp256k1, error) {
	c := &ClaimAuthorizeKSignSecp256k1{}
	_, c.Version = GetClaimTypeVersion(e)
	var cpk [33]byte
	copyFromElemBytes(cpk[len(cpk)-2:], ClaimTypeVersionLen, &e.Data[3])
	copyFromElemBytes(cpk[:len(cpk)-2], 0, &e.Data[2])
	var err error
	c.PubKey, err = crypto.DecompressPubkey(cpk[:])
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSignSecp256k1) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	SetClaimTypeVersion(e, c.Type(), c.Version)
	cpk := crypto.CompressPubkey(c.PubKey)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, cpk[len(cpk)-2:])
	copyToElemBytes(&e.Data[2], 0, cpk[:len(cpk)-2])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeKSignSecp256k1) Type() ClaimType {
	return *ClaimTypeAuthorizeKSignSecp256k1
}
