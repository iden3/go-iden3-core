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
		"0x16b32b9bc822ab5c1136eb099c7b05864914d7ee2cc531f932e6264c2d4b65e2",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x155f5c3e252fe5d439c197260fd9eddd627a4972f5e974482e62d2a81ba94264",
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
