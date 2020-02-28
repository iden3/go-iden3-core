package claims

import (
	"math/big"

	"github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// ClaimAuthorizeKSignBabyJub is a claim to authorize a baby jub public key for
// signing.
type ClaimAuthorizeKSignBabyJub struct {
	metadata Metadata
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
		metadata: NewMetadata(ClaimHeaderAuthorizeKSignBabyJub),
		Sign:     babyjub.PointCoordSign(pk.X),
		Ay:       pk.Y,
	}
}

// NewClaimAuthorizeKSignBabyJubFromEntry deserializes a
// ClaimAuthorizeKSignBabyJubFrom from an Entry.
func NewClaimAuthorizeKSignBabyJubFromEntry(e *merkletree.Entry) *ClaimAuthorizeKSignBabyJub {
	c := &ClaimAuthorizeKSignBabyJub{}
	c.metadata.Unmarshal(e)
	sign := []byte{0}
	copy(sign, e.Data[1][:])
	if sign[0] == 1 {
		c.Sign = true
	}
	c.Ay = new(big.Int).SetBytes(common.SwapEndianness(e.Data[2][:]))
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimAuthorizeKSignBabyJub) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	index := e.Index()
	sign := []byte{0}
	if c.Sign {
		sign = []byte{1}
	}
	copy(index[1][:], sign)
	ayBytes := c.Ay.Bytes()
	copy(index[2][:], common.SwapEndianness(ayBytes))
	c.metadata.Marshal(e)
	return e
}

func (c *ClaimAuthorizeKSignBabyJub) Metadata() *Metadata {
	return &c.metadata
}

// TODO: Keep the PublicKey in the Claim and only compress it when calling
// Entry() so that the key is available any time without this extra function.
// PublicKeyComp returns the compressed form of the public key in this claim
func (c *ClaimAuthorizeKSignBabyJub) PublicKeyComp() *babyjub.PublicKeyComp {
	pkc := babyjub.PublicKeyComp(
		babyjub.PackPoint(c.Ay, c.Sign))
	return &pkc
}
