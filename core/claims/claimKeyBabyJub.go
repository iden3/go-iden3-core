package claims

import (
	"math/big"

	"github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// ClaimKeyBabyJub is a claim to authorize a baby jub public key for
// signing.
type ClaimKeyBabyJub struct {
	metadata Metadata
	// Ax is the x coordinate of the BabyJubJub curve point which corresponds
	// to the public key.
	Ax *big.Int
	// Ay is the y coordinate of the BabyJubJub curve point which corresponds
	// to the public key.
	Ay *big.Int
}

// NewClaimKeyBabyJub returns a ClaimKeyBabyJub with the
// given elliptic public key parameters.
func NewClaimKeyBabyJub(pk *babyjub.PublicKey) *ClaimKeyBabyJub {
	return &ClaimKeyBabyJub{
		metadata: NewMetadata(ClaimHeaderKeyBabyJub),
		Ax:       pk.X,
		Ay:       pk.Y,
	}
}

// NewClaimKeyBabyJubFromEntry deserializes a
// ClaimKeyBabyJubFrom from an Entry.
func NewClaimKeyBabyJubFromEntry(e *merkletree.Entry) *ClaimKeyBabyJub {
	c := &ClaimKeyBabyJub{}
	c.metadata.Unmarshal(e)
	c.Ax = new(big.Int).SetBytes(common.SwapEndianness(e.Data[1][:]))
	c.Ay = new(big.Int).SetBytes(common.SwapEndianness(e.Data[2][:]))
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimKeyBabyJub) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	index := e.Index()
	axBytes := c.Ax.Bytes()
	copy(index[1][:], common.SwapEndianness(axBytes))
	ayBytes := c.Ay.Bytes()
	copy(index[2][:], common.SwapEndianness(ayBytes))
	c.metadata.Marshal(e)
	return e
}

func (c *ClaimKeyBabyJub) Metadata() *Metadata {
	return &c.metadata
}
