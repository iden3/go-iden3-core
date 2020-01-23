package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaimLinkObjectIdentity(t *testing.T) {
	// ClaimLinkObjectIdentity
	const objectType = ObjectTypeAddress
	var indexType uint16
	id, err := IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	assert.Nil(t, err)

	objectHash := [256 / 8]byte{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c}

	auxData := [256 / 8]byte{
		0x0f, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x09,
		0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x01, 0x02}

	claim, err := NewClaimLinkObjectIdentity(objectType, indexType, id, objectHash, auxData)
	assert.Nil(t, err)
	claim.Version = 1
	entry := claim.Entry()
	assert.Equal(t,
		"0x1b43ee36452e725d8e9990e888cb6bff71ad9071d773f2be52feabe78c3da143",
		entry.HIndex().Hex())
	assert.Equal(t,
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		entry.HValue().Hex())
	dataTestOutput(&entry.Data)
	assert.Equal(t, ""+
		"0f0102030405060708090a0b0c0d0e0f01020304050607090a0b0c0d0e0f0102"+
		"0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"0000003cc1c968fa000000000000000000000000000000000000000000000328"+
		"0000000000000000000000000000000000000001000000010000000000000005",
		entry.Data.String())
	c1 := NewClaimLinkObjectIdentityFromEntry(entry)
	c2, err := NewClaimFromEntry(entry)
	assert.Nil(t, err)
	assert.Equal(t, claim, c1)
	assert.Equal(t, claim, c2)
}
