package claims

import (
	"encoding/binary"
	"math/big"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// ClaimAuthorizeKSignBabyJub is a claim to authorize a baby jub public key for
// signing.
type ClaimAuthorizeKSignBabyJub struct {
	// Version is the claim version.
	Version uint32
	// RevocationNonce is used to revocate the claim
	RevocationNonce uint32
	// Sign means positive if false, negative if true.
	Sign bool
	// Ay is the y coordinate of the baby jub curve point which corresponds
	// to the public key.
	Ay *big.Int
}

// NewClaimAuthorizeKSignBabyJub returns a ClaimAuthorizeKSignBabyJub with the
// given elliptic public key parameters.
func NewClaimAuthorizeKSignBabyJub(pk *babyjub.PublicKey, revocationNonce uint32) *ClaimAuthorizeKSignBabyJub {
	return &ClaimAuthorizeKSignBabyJub{
		Version:         0,
		RevocationNonce: revocationNonce,
		Sign:            babyjub.PointCoordSign(pk.X),
		Ay:              pk.Y,
	}
}

// NewClaimAuthorizeKSignBabyJubFromEntry deserializes a
// ClaimAuthorizeKSignBabyJubFrom from an Entry.
func NewClaimAuthorizeKSignBabyJubFromEntry(e *merkletree.Entry) *ClaimAuthorizeKSignBabyJub {
	c := &ClaimAuthorizeKSignBabyJub{}
	_, c.Version = GetClaimTypeVersion(e)
	sign := []byte{0}
	copy(sign, e.Data[1][:])
	if sign[0] == 1 {
		c.Sign = true
	}
	c.Ay = new(big.Int).SetBytes(merkletree.SwapEndianness(e.Data[2][:]))
	c.RevocationNonce = binary.BigEndian.Uint32(e.Data[4][:4])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSignBabyJub) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	index := e.Index()
	SetClaimTypeVersion(e, c.Type(), c.Version)
	sign := []byte{0}
	if c.Sign {
		sign = []byte{1}
	}
	copy(index[1][:], sign)
	ayBytes := c.Ay.Bytes()
	copy(index[2][:], merkletree.SwapEndianness(ayBytes))

	binary.BigEndian.PutUint32(e.Data[4][:4], c.RevocationNonce)

	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimAuthorizeKSignBabyJub) Type() ClaimType {
	return *ClaimTypeAuthorizeKSignBabyJub
}

// TODO: Keep the PublicKey in the Claim and only compress it when calling
// Entry() so that the key is available any time without this extra function.
// PublicKeyComp returns the compressed form of the public key in this claim
func (c *ClaimAuthorizeKSignBabyJub) PublicKeyComp() *babyjub.PublicKeyComp {
	pkc := babyjub.PublicKeyComp(
		babyjub.PackPoint(c.Ay, c.Sign))
	return &pkc
}
