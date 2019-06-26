package core

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3/merkletree"
)

// ClaimAuthorizeKSignBabyJub is a claim to authorize a baby jub public key for
// signing.
type ClaimAuthorizeKSignBabyJub struct {
	// Version is the claim version.
	Version uint32
	// Sign means positive if false, negative if true.
	Sign bool
	// Ay is the y coordinate of the baby jub curve point which corresponds
	// to the public key.
	Ay *big.Int
}

// NewClaimAuthorizeKSignBabyJub returns a ClaimAuthorizeKSignBabyJub with the
// given elliptic public key parameters.
func NewClaimAuthorizeKSignBabyJub(pk *babyjub.PublicKey) *ClaimAuthorizeKSignBabyJub {
	return &ClaimAuthorizeKSignBabyJub{
		Version: 0,
		Sign:    babyjub.PointCoordSign(pk.X),
		Ay:      pk.Y,
	}
}

// NewClaimAuthorizeKSignBabyJubFromEntry deserializes a
// ClaimAuthorizeKSignBabyJubFrom from an Entry.
func NewClaimAuthorizeKSignBabyJubFromEntry(e *merkletree.Entry) *ClaimAuthorizeKSignBabyJub {
	c := &ClaimAuthorizeKSignBabyJub{}
	_, c.Version = getClaimTypeVersion(e)
	sign := []byte{0}
	copyFromElemBytes(sign, ClaimTypeVersionLen, &e.Data[3])
	if sign[0] == 1 {
		c.Sign = true
	}
	c.Ay = new(big.Int).SetBytes(e.Data[2][:])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSignBabyJub) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	setClaimTypeVersion(e, c.Type(), c.Version)
	sign := []byte{0}
	if c.Sign {
		sign = []byte{1}
	}
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, sign)
	copy(e.Data[2][:], c.Ay.Bytes())
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeKSignBabyJub) Type() ClaimType {
	return *ClaimTypeAuthorizeKSignBabyJub
}
