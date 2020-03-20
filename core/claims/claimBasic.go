package claims

import (
	"github.com/iden3/go-iden3-core/merkletree"
)

const (
	IndexSlotLen = EntryFullBytesLen - ClaimTypeLen - ClaimFlagsLen + EntryFullBytesLen*3
	ValueSlotLen = EntryFullBytesLen - ClaimRevNonceLen + EntryFullBytesLen*3
)

// ClaimBasic is a simple claim that can be used for anything.
type ClaimBasic struct {
	metadata Metadata
	// IndexSlot is data that goes into the remaining space used for the index.
	IndexSlot [IndexSlotLen]byte
	// ValueSlot is the data that goes into the remaining space not used for the index.
	ValueSlot [ValueSlotLen]byte
}

// NewClaimBasic returns a ClaimBasic with the provided data.
func NewClaimBasic(indexSlot [IndexSlotLen]byte, valueSlot [ValueSlotLen]byte) *ClaimBasic {
	return &ClaimBasic{
		metadata:  NewMetadata(ClaimHeaderBasic),
		IndexSlot: indexSlot,
		ValueSlot: valueSlot,
	}
}

// NewClaimBasicFromEntry deserializes a ClaimBasic from an Entry.
func NewClaimBasicFromEntry(e *merkletree.Entry) *ClaimBasic {
	c := &ClaimBasic{}
	c.metadata.Unmarshal(e)

	n := 0
	for i, start := range []int{ClaimHeaderLen, 0, 0, 0} {
		n += copy(c.IndexSlot[n:], e.Index()[i].Bytes()[start:EntryFullBytesLen])
	}
	n = 0
	for i, start := range []int{ClaimRevNonceLen, 0, 0, 0} {
		n += copy(c.ValueSlot[n:], e.Value()[i].Bytes()[start:EntryFullBytesLen])
	}

	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimBasic) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}

	n := 0
	for i, start := range []int{ClaimHeaderLen, 0, 0, 0} {
		n += copy(e.Index()[i].Bytes()[start:], c.IndexSlot[n:n+EntryFullBytesLen-start])
	}
	n = 0
	for i, start := range []int{ClaimRevNonceLen, 0, 0, 0} {
		n += copy(e.Value()[i].Bytes()[start:], c.ValueSlot[n:n+EntryFullBytesLen-start])
	}

	c.metadata.Marshal(e)
	return e
}

func (c *ClaimBasic) Metadata() *Metadata {
	return &c.metadata
}
