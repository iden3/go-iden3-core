package claims

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

func TestClaimBasic(t *testing.T) {
	// ClaimBasic
	var indexSlot [IndexSlotBytes]byte
	var dataSlot [DataSlotBytes]byte
	indexSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_indexSlot").(string))
	assert.Nil(t, err)
	dataSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_dataSlot").(string))
	assert.Nil(t, err)
	copy(indexSlot[:], indexSlotHex[:IndexSlotBytes])
	copy(dataSlot[:], dataSlotHex[:DataSlotBytes])
	c0 := NewClaimBasic(indexSlot, dataSlot, 5678)
	c0.Version = 1
	e := c0.Entry()
	// Check claim against test vector
	testgen.CheckTestValue(t, "ClaimBasic0_HIndex", e.HIndex().Hex())
	testgen.CheckTestValue(t, "ClaimBasic0_HValue", e.HValue().Hex())
	testgen.CheckTestValue(t, "ClaimBasic0_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimBasicFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))

	// revocation nonce
	c3 := NewClaimBasic(indexSlot, dataSlot, 3)
	assert.Equal(t, c3.RevocationNonce, uint32(3))
	c3.Version = 1
	c1.RevocationNonce = 3
	assert.Equal(t, c3, c1)
}

func TestClaimBasic1(t *testing.T) {
	indexData := []byte(testgen.GetTestValue("1_indexData").(string))
	data := []byte(testgen.GetTestValue("1_valueData").(string))
	var indexSlot [IndexSlotBytes]byte
	var dataSlot [DataSlotBytes]byte
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])

	// ClaimBasic
	c0 := NewClaimBasic(indexSlot, dataSlot, 0)
	e := c0.Entry()
	// Check claim against test vector
	testgen.CheckTestValue(t, "ClaimBasic1_HIndex", e.HIndex().Hex())
	testgen.CheckTestValue(t, "ClaimBasic1_HValue", e.HValue().Hex())
	testgen.CheckTestValue(t, "ClaimBasic1_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimBasicFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))

	// revocation nonce
	c3 := NewClaimBasic(indexSlot, dataSlot, 3)
	assert.Equal(t, c3.RevocationNonce, uint32(3))
	c1.RevocationNonce = 3
	assert.Equal(t, c3, c1)

	// Stop test (write new test vector if needed)
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
}
