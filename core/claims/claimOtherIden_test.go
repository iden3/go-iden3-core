package claims

import (
	"encoding/hex"
	"testing"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-merkletree"
	"github.com/stretchr/testify/assert"
)

func TestClaimOtherIden(t *testing.T) {
	// ClaimOtherIden
	var indexSlot [IndexSubjectSlotLen]byte
	var valueSlot [ValueSlotLen]byte
	indexSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_indexSubjectSlot").(string))
	assert.Nil(t, err)
	valueSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_valueSlot").(string))
	assert.Nil(t, err)
	copy(indexSlot[:], indexSlotHex[:IndexSubjectSlotLen])
	copy(valueSlot[:], valueSlotHex[:ValueSlotLen])
	id, err := core.IDFromString(testgen.GetTestValue("0_subject").(string))
	assert.Nil(t, err)
	c0 := NewClaimOtherIden(&id, indexSlot, valueSlot)
	c0.Metadata().RevNonce = 5678
	e := c0.Entry()
	// Check claim against test vector
	hi, hv, err := e.HiHv()
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "ClaimOtherIden0_HIndex", hi.Hex())
	testgen.CheckTestValue(t, "ClaimOtherIden0_HValue", hv.Hex())
	testgen.CheckTestValue(t, "ClaimOtherIden0_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimOtherIdenFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0.Metadata(), c1.Metadata())
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))
}

func TestClaimOtherIden1(t *testing.T) {
	indexData := []byte(testgen.GetTestValue("1_indexData").(string))
	data := []byte(testgen.GetTestValue("1_valueData").(string))
	var indexSlot [IndexSubjectSlotLen]byte
	var valueSlot [ValueSlotLen]byte
	copy(indexSlot[:], indexData[:])
	copy(valueSlot[:], data[:])
	id, err := core.IDFromString(testgen.GetTestValue("0_subject").(string))
	assert.Nil(t, err)

	// ClaimOtherIden
	c0 := NewClaimOtherIden(&id, indexSlot, valueSlot)
	e := c0.Entry()
	hi, hv, err := e.HiHv()
	assert.Nil(t, err)
	// Check claim against test vector
	testgen.CheckTestValue(t, "ClaimOtherIden1_HIndex", hi.Hex())
	testgen.CheckTestValue(t, "ClaimOtherIden1_HValue", hv.Hex())
	testgen.CheckTestValue(t, "ClaimOtherIden1_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimOtherIdenFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0.Metadata(), c1.Metadata())
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))
}
