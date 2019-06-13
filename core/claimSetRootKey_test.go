package core

import (
	"testing"

	"github.com/iden3/go-iden3/merkletree"
	"github.com/stretchr/testify/assert"
)

func TestClaimSetRootKey(t *testing.T) {
	// ClaimSetRootKey
	id, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)

	rootKey := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c})
	c0 := NewClaimSetRootKey(id, rootKey)
	c0.Version = 1
	c0.Era = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x12bf59ff4171debe81321c04a52298e62650ca8514e9a7a8a64c23cb55eeaa2e",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x01705b25f2cf7cda34d836f09e9b0dd1777bdc16752657cd9d1ae5f6286525ba",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"0000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f4"+
		"0000000000000000000000000000000000000001000000010000000000000002",
		e.Data.String())
	c1 := NewClaimSetRootKeyFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}
