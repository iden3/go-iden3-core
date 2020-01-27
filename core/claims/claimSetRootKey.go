package claims

import (
	"encoding/binary"
	"errors"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	cryptoConstants "github.com/iden3/go-iden3-crypto/constants"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
)

// ClaimSetRootKey is a claim of the root key of a merkle tree that goes into the relay.
type ClaimSetRootKey struct {
	// Version is the claim version.
	Version uint32
	// Era is used for labeling epochs.
	Era uint32
	// Id is the ID related to the root key.
	Id core.ID
	// RootKey is the root of the mekrlee tree.
	RootKey merkletree.Hash
}

// NewClaimSetRootKey returns a ClaimSetRootKey with the given Eth ID and
// merklee tree root key.
func NewClaimSetRootKey(id *core.ID, rootKey *merkletree.Hash) (*ClaimSetRootKey, error) {
	if ok := cryptoUtils.CheckBigIntArrayInField(merkletree.ElemBytesToBigInts(merkletree.ElemBytes(*rootKey)), cryptoConstants.Q); !ok {
		return nil, errors.New("Elements not in the Finite Field over R")
	}
	return &ClaimSetRootKey{
		Version: 0,
		Era:     0,
		Id:      *id,
		RootKey: *rootKey,
	}, nil
}

// NewClaimSetRootKeyFromEntry deserializes a ClaimSetRootKey from an Entry.
func NewClaimSetRootKeyFromEntry(e *merkletree.Entry) *ClaimSetRootKey {
	c := &ClaimSetRootKey{}
	_, c.Version = GetClaimTypeVersion(e)
	var era [32 / 8]byte
	copyFromElemBytes(era[:], ClaimTypeVersionLen, &e.Data[3])
	c.Era = binary.BigEndian.Uint32(era[:])
	copyFromElemBytes(c.Id[:], 0, &e.Data[2])
	c.RootKey = merkletree.Hash(e.Data[1])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimSetRootKey) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	SetClaimTypeVersion(e, c.Type(), c.Version)
	var era [32 / 8]byte
	binary.BigEndian.PutUint32(era[:], c.Era)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, era[:])
	copyToElemBytes(&e.Data[2], 0, c.Id[:])
	e.Data[1] = merkletree.ElemBytes(c.RootKey)
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimSetRootKey) Type() ClaimType {
	return *ClaimTypeSetRootKey
}
