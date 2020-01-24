package core

import (
	"testing"

	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/assert"
)

func TestClaimSetRootKey(t *testing.T) {
	// ClaimSetRootKey
	id, err := IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	assert.Nil(t, err)

	rootKey := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c})
	c0, err := NewClaimSetRootKey(&id, &rootKey)
	assert.Nil(t, err)
	c0.Version = 1
	c0.Era = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x15447e94331c3d175a77e51c008f6e7ec684f4734df06d91a44a4c9b4f98893a",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"0000003cc1c968fa000000000000000000000000000000000000000000000328"+
		"0000000000000000000000000000000000000001000000010000000000000002",
		e.Data.String())
	c1 := NewClaimSetRootKeyFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}
