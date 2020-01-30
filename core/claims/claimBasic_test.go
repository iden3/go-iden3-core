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
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false
	// Init test
	if err := testgen.InitTest("claimBasic", generateTest); err != nil {
		panic(fmt.Errorf("error initializing test data: %w", err))
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("0_indexSlot", hex.EncodeToString([]byte{
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
			0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2b}))
		testgen.SetTestValue("0_dataSlot", hex.EncodeToString([]byte{
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40,
			0x56, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x59}))
	}

	// ClaimBasic
	var indexSlot [IndexSlotBytes]byte
	var dataSlot [DataSlotBytes]byte
	indexSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_indexSlot").(string))
	assert.Nil(t, err)
	dataSlotHex, err := hex.DecodeString(testgen.GetTestValue("0_dataSlot").(string))
	assert.Nil(t, err)
	copy(indexSlot[:], indexSlotHex[:IndexSlotBytes])
	copy(dataSlot[:], dataSlotHex[:DataSlotBytes])
	c0 := NewClaimBasic(indexSlot, dataSlot)
	c0.Version = 1
	e := c0.Entry()
	// Check claim against test vector
	testgen.CheckTestValue(t, "0_HIndex", e.HIndex().Hex())
	testgen.CheckTestValue(t, "0_HValue", e.HValue().Hex())
	testgen.CheckTestValue(t, "0_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimBasicFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))

	// Stop test (write new test vector if needed)
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
}

func TestClaimBasic1(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := false
	// Init test
	if err := testgen.InitTest("claimBasic", generateTest); err != nil {
		panic(fmt.Errorf("error initializing test data: %w", err))
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("1_indexData", "c1")
		testgen.SetTestValue("1_valueData", "")
	}
	indexData := []byte(testgen.GetTestValue("1_indexData").(string))
	data := []byte(testgen.GetTestValue("1_valueData").(string))
	var indexSlot [IndexSlotBytes]byte
	var dataSlot [DataSlotBytes]byte
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])

	// ClaimBasic
	c0 := NewClaimBasic(indexSlot, dataSlot)
	e := c0.Entry()
	// Check claim against test vector
	testgen.CheckTestValue(t, "1_HIndex", e.HIndex().Hex())
	testgen.CheckTestValue(t, "1_HValue", e.HValue().Hex())
	testgen.CheckTestValue(t, "1_dataString", e.Data.String())
	dataTestOutput(&e.Data)
	c1 := NewClaimBasicFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.True(t, merkletree.CheckEntryInField(*e))

	// Stop test (write new test vector if needed)
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
}
