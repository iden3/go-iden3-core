package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

func TestClaimSetRootKey(t *testing.T) {
	// If generateTest is true, the checked values will be used to generate a test vector
	generateTest := true
	// Init test
	err := testgen.InitTest("claimSetRootKey", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("idString", "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
		testgen.SetTestValue("rootKey", hex.EncodeToString([]byte{
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
			0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c}))
	}
	// ClaimSetRootKey
	id, err := IDFromString(testgen.GetTestValue("idString").(string))
	assert.Nil(t, err)
	rootKey := hexStringToKey(testgen.GetTestValue("rootKey").(string))
	c0, err := NewClaimSetRootKey(&id, &rootKey)
	assert.Nil(t, err)
	c0.Version = 1
	c0.Era = 1
	e := c0.Entry()
	// Check claim against test vector
	checkClaim(e, t)
	dataTestOutput(&e.Data)
	c1 := NewClaimSetRootKeyFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
	// Stop test (write new test vector if needed)
	err = testgen.StopTest()
	if err != nil {
		panic(fmt.Sprint("Error stopping test:", err))
	}
}

func hexStringToKey(s string) merkletree.Hash {
	var keyBytes [merkletree.ElemBytesLen]byte
	keyBytesHex, _ := hex.DecodeString(s)
	copy(keyBytes[:], keyBytesHex[:merkletree.ElemBytesLen])
	return merkletree.Hash(merkletree.ElemBytes(keyBytes))
}
