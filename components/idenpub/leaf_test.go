package idenpub

import (
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/testgen"
	"github.com/stretchr/testify/assert"
)

func TestLeafRoT(t *testing.T) {
	root := hexStringToHash(testgen.GetTestValue("root0").(string))

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
