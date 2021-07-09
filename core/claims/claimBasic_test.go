package claims

import (
	"encoding/hex"
	"testing"

	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-merkletree-sql"
	"github.com/stretchr/testify/assert"
)

func TestClaimBasic(t *testing.T) {
	// ClaimBasic
	var indexSlot [IndexSlotLen]byte
	var valueSlot [ValueSlotLen]byte
	indexSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_indexSlot").(string))
	assert.Nil(t, err)
	valueSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_valueSlot").(string))
	assert.Nil(t, err)
	copy(indexSlot[:], indexSlotHex[:IndexSlotLen])
	copy(valueSlot[:], valueSlotHex[:ValueSlotLen])
	c0 := NewClaimBasic(indexSlot, valueSlot)
	c0.Metadata().RevNonce = 5678
	e := c0.Entry()
	// Check claim against test vector
	hi, hv, err := e.HiHv()
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "ClaimBasic0_HIndex", hi.Hex())
	testgen.CheckTestValue(t, "ClaimBasic0_HValue", hv.Hex())
	testgen.CheckTestValue(t, "ClaimBasic0_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimBasicFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0.Metadata(), c1.Metadata())
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))
}

func TestClaimBasic1(t *testing.T) {
	indexData := []byte(testgen.GetTestValue("1_indexData").(string))
	data := []byte(testgen.GetTestValue("1_valueData").(string))
	var indexSlot [IndexSlotLen]byte
	var valueSlot [ValueSlotLen]byte
	copy(indexSlot[:], indexData[:])
	copy(valueSlot[:], data[:])

	// ClaimBasic
	c0 := NewClaimBasic(indexSlot, valueSlot)
	e := c0.Entry()
	hi, hv, err := e.HiHv()
	assert.Nil(t, err)
	// Check claim against test vector
	testgen.CheckTestValue(t, "ClaimBasic1_HIndex", hi.Hex())
	testgen.CheckTestValue(t, "ClaimBasic1_HValue", hv.Hex())
	testgen.CheckTestValue(t, "ClaimBasic1_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimBasicFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0.Metadata(), c1.Metadata())
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))
}
