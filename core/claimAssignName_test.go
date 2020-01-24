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
	id, err := IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	assert.Nil(t, err)
	c0 := NewClaimAssignName(name, id)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x1ba59bde7b10fb76b401495acfe3063268d6f7b4d3cfea62f2d2c39bd016c02d",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000003cc1c968fa000000000000000000000000000000000000000000000328"+
		"00d67b05d8e2d1ace8f3e84b8451dd2e9da151578c3c6be23e7af11add5a807a"+
		"0000000000000000000000000000000000000000000000010000000000000003",
		e.Data.String())
	c1 := NewClaimAssignNameFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}
