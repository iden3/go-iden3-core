package claims

import (
	"github.com/iden3/go-iden3-core/merkletree"
)

const (
	// 56+248+248+248=800 bits
	indexSlotBits  = 800
	IndexSlotBytes = indexSlotBits / 8
	// 216+248+248+248=960 bits
	dataSlotBits  = 960
	DataSlotBytes = dataSlotBits / 8
)

// ClaimBasic is a simple claim that can be used for anything.
type ClaimBasic struct {
	// Version is the claim version.
	Version uint32
	// IndexSlot is data that goes into the remaining space used for the index.
	IndexSlot [IndexSlotBytes]byte
	// DataSlot is the data that goes into the remaining space not used for the index.
	DataSlot [DataSlotBytes]byte
}

// NewClaimBasic returns a ClaimBasic with the provided data.
func NewClaimBasic(indexSlot [IndexSlotBytes]byte, dataSlot [DataSlotBytes]byte) *ClaimBasic {
	// TODO: at this moment, revocation nonce is not defined, neither other
	// claim options.  So, for now, the ClaimBasic just holds two static
	// blocks of data (IndexSlot and DataSlot).  Once the rest of the claim
	// parameters are defined, this claim will be updated with all the
	// options on the construction
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
	copy(c.IndexSlot[:56/8], e.Data[0][merkletree.ElemBytesLen-(64/8):]) // last 56 bits of the index_slot[0]
	copy(c.IndexSlot[56/8:304/8], e.Data[1][:])                          // first 248 bits of index_slot[2]
	copy(c.IndexSlot[304/8:552/8], e.Data[2][:])                         // first 248 bits of index_slot[2]
	copy(c.IndexSlot[552/8:800/8], e.Data[3][:])                         // first 248 bits of index_slot[3]

	copy(c.DataSlot[:216/8], e.Data[4][:])      // last 216 bits of data_slot[0]
	copy(c.DataSlot[216/8:464/8], e.Data[5][:]) // first 248 bits of data_slot[1]
	copy(c.DataSlot[464/8:712/8], e.Data[6][:]) // first 248 bits of data_slot[2]
	copy(c.DataSlot[712/8:960/8], e.Data[7][:]) // first 248 bits of data_slot[3]

	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimBasic) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}
	SetClaimTypeVersion(e, c.Type(), c.Version)

	copy(e.Data[0][merkletree.ElemBytesLen-(64/8):], c.IndexSlot[0:56/8])
	copy(e.Data[1][0:], c.IndexSlot[56/8:304/8])
	copy(e.Data[2][0:], c.IndexSlot[304/8:552/8])
	copy(e.Data[3][0:], c.IndexSlot[552/8:800/8])

	copy(e.Data[4][0:], c.DataSlot[:216/8])
	copy(e.Data[5][0:], c.DataSlot[216/8:464/8])
	copy(e.Data[6][0:], c.DataSlot[464/8:712/8])
	copy(e.Data[7][0:], c.DataSlot[712/8:960/8])

	return e
}

// Type returns the ClaimType of the claim.
func (c *ClaimBasic) Type() ClaimType {
	return *ClaimTypeBasic
}
