package claims

import (
	"encoding/binary"
	"math/big"

	"github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree"
)

const (
	BabyJubKeyTypeGeneric        BabyJubKeyType = 0
	BabyJubKeyTypeAuthorizeKSign BabyJubKeyType = 1
)

type BabyJubKeyType uint64

// ClaimKeyBabyJub is a claim to authorize a baby jub public key for
// signing.
type ClaimKeyBabyJub struct {
	metadata Metadata

	// KeyType is used to specify if the type of usage given to the Key
	// if 0 is the generic one
	// if 1 is to AuthorizeKSign Key
	KeyType BabyJubKeyType

	// Ax is the x coordinate of the BabyJubJub curve point which
	// corresponds to the public key.
	Ax *big.Int

	// Ay is the y coordinate of the BabyJubJub curve point which
	// corresponds to the public key.
	Ay *big.Int
}

// NewClaimKeyBabyJub returns a ClaimKeyBabyJub with the
// given elliptic public key parameters.
func NewClaimKeyBabyJub(pk *babyjub.PublicKey, subType BabyJubKeyType) *ClaimKeyBabyJub {
	return &ClaimKeyBabyJub{
		metadata: NewMetadata(ClaimHeaderKeyBabyJub),
		KeyType:  subType,
		Ax:       pk.X,
		Ay:       pk.Y,
	}
}

// NewClaimKeyBabyJubFromEntry deserializes a
// ClaimKeyBabyJubFrom from an Entry.
func NewClaimKeyBabyJubFromEntry(e *merkletree.Entry) *ClaimKeyBabyJub {
	c := &ClaimKeyBabyJub{}
	c.metadata.Unmarshal(e)

	c.KeyType = BabyJubKeyType(binary.BigEndian.Uint64(e.Data[1][:]))

	c.Ax = new(big.Int).SetBytes(common.SwapEndianness(e.Data[2][:]))
	c.Ay = new(big.Int).SetBytes(common.SwapEndianness(e.Data[3][:]))
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimKeyBabyJub) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	index := e.Index()

	var keyType [8]byte
	binary.BigEndian.PutUint64(keyType[:], uint64(c.KeyType))
	copy(index[1][:], keyType[:])

	axBytes := c.Ax.Bytes()
	copy(index[2][:], common.SwapEndianness(axBytes))

	ayBytes := c.Ay.Bytes()
	copy(index[3][:], common.SwapEndianness(ayBytes))

	c.metadata.Marshal(e)
	return e
}

func (c *ClaimKeyBabyJub) Metadata() *Metadata {
	return &c.metadata
}
