package idenpub

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

// If generateTest is true, the checked values will be used to generate a test vector
var generateTest = false

func TestLeafRoT(t *testing.T) {
	root := hexStringToKey(testgen.GetTestValue("root0").(string))

	l0 := NewLeafRoT(root)
	e := l0.Entry()

	testgen.CheckTestValue(t, "Leaf0_HIndex", e.HIndex().Hex())
	testgen.CheckTestValue(t, "Leaf0_HValue", e.HValue().Hex())
	testgen.CheckTestValue(t, "Leaf0_dataString", e.Data.String())
	l1 := NewLeafRoTFromEntry(e)
	assert.Equal(t, l0, l1)
	assert.True(t, merkletree.CheckEntryInField(*e))
	assert.Equal(t, l0.Root[:31], root[:31])
	assert.Equal(t, l1.Root[:31], root[:31])
}

func TestLeafReT(t *testing.T) {
	nonce := uint32(testgen.GetTestValue("nonce0").(float64))
	version := uint32(testgen.GetTestValue("version0").(float64))

	l0 := NewLeafReT(nonce, version)
	e := l0.Entry()

	testgen.CheckTestValue(t, "Leaf1_HIndex", e.HIndex().Hex())
	testgen.CheckTestValue(t, "Leaf1_HValue", e.HValue().Hex())
	testgen.CheckTestValue(t, "Leaf1_dataString", e.Data.String())
	l1 := NewLeafReTFromEntry(e)
	assert.Equal(t, l0, l1)
	assert.True(t, merkletree.CheckEntryInField(*e))
	assert.Equal(t, l0.Nonce, nonce)
	assert.Equal(t, l1.Nonce, nonce)
	assert.Equal(t, l0.Version, version)
	assert.Equal(t, l1.Version, version)
}

func initTest() {
	// Init test
	err := testgen.InitTest("idenpub", generateTest)
	if err != nil {
		fmt.Println("error initializing test data:", err)
		return
	}
	// Add input data to the test vector
	if generateTest {
		testgen.SetTestValue("root0", "0x2718b18e6743a777501accab821bf348e7c8a44ef23eb4da9f15c546092d302b")
		testgen.SetTestValue("nonce0", float64(5))
		testgen.SetTestValue("version0", float64(5))
	}
}

func TestMain(m *testing.M) {
	initTest()
	result := m.Run()
	if err := testgen.StopTest(); err != nil {
		panic(fmt.Errorf("Error stopping test: %w", err))
	}
	os.Exit(result)
}

func hexStringToKey(s string) merkletree.Hash {
	var keyBytes [merkletree.ElemBytesLen]byte
	keyBytesHex, _ := hex.DecodeString(s)
	copy(keyBytes[:], keyBytesHex[:merkletree.ElemBytesLen])
	return merkletree.Hash(merkletree.ElemBytes(keyBytes))
}
