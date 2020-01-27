package claims

import (
	"github.com/iden3/go-iden3-core/merkletree"
)

// ClaimBasic is a simple claim that can be used for anything.
type ClaimBasic struct {
	// Version is the claim version.
	Version uint32
	// IndexSlot is data that goes into the remaining space used for the index.
	IndexSlot [400 / 8]byte
	// DataSlot is the data that goes into the remaining space not used for the index.
	DataSlot [496 / 8]byte
}

// NewClaimBasic returns a ClaimBasic with the provided data.
func NewClaimBasic(indexSlot [400 / 8]byte, dataSlot [496 / 8]byte) *ClaimBasic {
	return &ClaimBasic{
		Version:   0,
		IndexSlot: indexSlot,
		DataSlot:  dataSlot,
	}
}

// NewClaimBasicFromEntry deserializes a ClaimBasic from an Entry.
func NewClaimBasicFromEntry(e *merkletree.Entry) *ClaimBasic {
	c := &ClaimBasic{}
	_, c.Version = GetClaimTypeVersion(e)
	copyFromElemBytes(c.IndexSlot[len(c.IndexSlot)-152/8:], ClaimTypeVersionLen, &e.Data[3])
	copyFromElemBytes(c.IndexSlot[:248/8], 0, &e.Data[2])
	copyFromElemBytes(c.DataSlot[248/8:], 0, &e.Data[1])
	copyFromElemBytes(c.DataSlot[:248/8], 0, &e.Data[0])
	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimBasic) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	SetClaimTypeVersion(e, c.Type(), c.Version)
	copyToElemBytes(&e.Data[3], ClaimTypeVersionLen, c.IndexSlot[len(c.IndexSlot)-152/8:])
	copyToElemBytes(&e.Data[2], 0, c.IndexSlot[:248/8])
	copyToElemBytes(&e.Data[1], 0, c.DataSlot[248/8:])
	copyToElemBytes(&e.Data[0], 0, c.DataSlot[:248/8])
	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimBasic) Type() ClaimType {
	return *ClaimTypeBasic
}
