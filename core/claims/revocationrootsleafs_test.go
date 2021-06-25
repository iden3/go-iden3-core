package claims

import (
	"encoding/hex"
	"github.com/iden3/go-merkletree/db/memory"
	"testing"

	"github.com/iden3/go-iden3-core/testgen"
	"github.com/iden3/go-merkletree"
	"github.com/stretchr/testify/assert"
)

func TestLeafRootsTree(t *testing.T) {
	root := HexStringToHash(testgen.GetTestValue("root0").(string))

	l0 := NewLeafRootsTree(root)
	e := l0.Entry()

	hi, hv, err := e.HiHv()
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "Leaf0_HIndex", hi.Hex())
	testgen.CheckTestValue(t, "Leaf0_HValue", hv.Hex())
	testgen.CheckTestValue(t, "Leaf0_dataString", e.Data.String())
	l1 := NewLeafRootsTreeFromEntry(e)
	assert.Equal(t, l0, l1)
	assert.True(t, merkletree.CheckEntryInField(*e))
	assert.Equal(t, l0.Root[:31], root[:31])
	assert.Equal(t, l1.Root[:31], root[:31])
}

func TestLeafRevocationsTree(t *testing.T) {
	nonce := uint32(testgen.GetTestValue("nonce0").(float64))
	version := uint32(testgen.GetTestValue("version0").(float64))

	l0 := NewLeafRevocationsTree(nonce, version)
	e := l0.Entry()

	hi, hv, err := e.HiHv()
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "Leaf1_HIndex", hi.Hex())
	testgen.CheckTestValue(t, "Leaf1_HValue", hv.Hex())
	testgen.CheckTestValue(t, "Leaf1_dataString", e.Data.String())
	l1 := NewLeafRevocationsTreeFromEntry(e)
	assert.Equal(t, l0, l1)
	assert.True(t, merkletree.CheckEntryInField(*e))
	assert.Equal(t, l0.Nonce, nonce)
	assert.Equal(t, l1.Nonce, nonce)
	assert.Equal(t, l0.Version, version)
	assert.Equal(t, l1.Version, version)
}

func TestAddLeafRootsTree(t *testing.T) {
	root := HexStringToHash(testgen.GetTestValue("root0").(string))

	mt, err := merkletree.NewMerkleTree(memory.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	err = AddLeafRootsTree(mt, &root)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "rootRootsTree0", mt.Root().Hex())

	hi, err := NewLeafRootsTree(root).Entry().HIndex()
	assert.Nil(t, err)
	proof, _, err := mt.GenerateProof(hi.BigInt(), nil)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "proofLeafRootsTree", hex.EncodeToString(proof.Bytes()))
}

func TestAddLeafRevocationsTree(t *testing.T) {
	nonce := uint32(testgen.GetTestValue("nonce0").(float64))
	version := uint32(testgen.GetTestValue("version0").(float64))

	mt, err := merkletree.NewMerkleTree(memory.NewMemoryStorage(), 140)
	assert.Nil(t, err)

	err = AddLeafRevocationsTree(mt, nonce, version)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "rootRevocationsTree0", mt.Root().Hex())

	hi, err := NewLeafRevocationsTree(nonce, version).Entry().HIndex()
	assert.Nil(t, err)
	proof, _, err := mt.GenerateProof(hi.BigInt(), nil)
	assert.Nil(t, err)
	testgen.CheckTestValue(t, "proofRevocationsTree", hex.EncodeToString(proof.Bytes()))
}
