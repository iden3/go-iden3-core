package claims

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-merkletree-sql"
)

const (
	IndexSubjectSlotLen = EntryFullBytesLen - ClaimTypeLen - ClaimFlagsLen + EntryFullBytesLen*2
)

// ClaimOtherIden is a simple claim that can be used for anything.
type ClaimOtherIden struct {
	metadata Metadata
	// IndexSlot is data that goes into the remaining space used for the index.
	IndexSlot [IndexSubjectSlotLen]byte
	// ValueSlot is the data that goes into the remaining space not used for the index.
	ValueSlot [ValueSlotLen]byte
}

// NewClaimOtherIden returns a ClaimOtherIden with the provided data.
func NewClaimOtherIden(subject *core.ID, indexSlot [IndexSubjectSlotLen]byte,
	valueSlot [ValueSlotLen]byte) *ClaimOtherIden {
	metadata := NewMetadata(ClaimHeaderOtherIden)
	metadata.Subject = subject
	return &ClaimOtherIden{
		metadata:  metadata,
		IndexSlot: indexSlot,
		ValueSlot: valueSlot,
	}
}

// NewClaimOtherIdenFromEntry deserializes a ClaimOtherIden from an Entry.
func NewClaimOtherIdenFromEntry(e *merkletree.Entry) *ClaimOtherIden {
	c := &ClaimOtherIden{}
	c.metadata.Unmarshal(e)

	n := 0
	for i, start := range []int{ClaimHeaderLen, 0, 0, 0} {
		if i == 1 {
			continue
		}
		n += copy(c.IndexSlot[n:], e.Index()[i][start:EntryFullBytesLen])
	}
	n = 0
	for i, start := range []int{ClaimRevNonceLen, 0, 0, 0} {
		n += copy(c.ValueSlot[n:], e.Value()[i][start:EntryFullBytesLen])
	}

	return c
}

// Entry serializes the claim into an Entry.
func (c *ClaimOtherIden) Entry() *merkletree.Entry {
	e := &merkletree.Entry{}

	n := 0
	for i, start := range []int{ClaimHeaderLen, 0, 0, 0} {
		if i == 1 {
			continue
		}
		n += copy(e.Index()[i][start:], c.IndexSlot[n:n+EntryFullBytesLen-start])
	}
	n = 0
	for i, start := range []int{ClaimRevNonceLen, 0, 0, 0} {
		n += copy(e.Value()[i][start:], c.ValueSlot[n:n+EntryFullBytesLen-start])
	}

	c.metadata.Marshal(e)
	return e
}

func (c *ClaimOtherIden) Metadata() *Metadata {
	return &c.metadata
}
