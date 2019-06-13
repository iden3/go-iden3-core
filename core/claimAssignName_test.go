package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaimAssignName(t *testing.T) {
	// ClaimAssignName
	name := "example.iden3.eth"
	// genesis := common.BytesToAddress([]byte{
	//         0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
	//         0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
	//         0x39, 0x39, 0x39, 0x3a})
	id, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)
	c0 := NewClaimAssignName(name, id)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x106d1a898d4503f4cb20be6ce9aeb2ac1e65d522579805e3633408a4b9ffcb53",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x25867e06233f276f39e298775245bad077eb0852b4eaac8dbf646a95bd3f8625",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f4"+
		"00d67b05d8e2d1ace8f3e84b8451dd2e9da151578c3c6be23e7af11add5a807a"+
		"0000000000000000000000000000000000000000000000010000000000000003",
		e.Data.String())
	c1 := NewClaimAssignNameFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}
