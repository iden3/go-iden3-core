package core

import (
	"encoding/hex"
	"testing"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

func TestClaimAuthorizeKSignBabyJub(t *testing.T) {
	// ClaimAuthorizeKSignBabyJub
	var k babyjub.PrivateKey
	hex.Decode(k[:], []byte("28156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"))
	pk := k.Public()

	c0 := NewClaimAuthorizeKSignBabyJub(pk)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x04f41fdac3240e7b68905df19a2394e4a4f1fb7eaeb310e39e1bb0b225b7763f",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x06d4571fb9634e4bed32e265f91a373a852c476656c5c13b09bc133ac61bc5a6",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"2b05184c7195b259c95169348434f3a7228fbcfb187d3b07649f3791330cf05c"+
		"0000000000000000000000000000000000000001000000010000000000000001",
		e.Data.String())
	c1 := NewClaimAuthorizeKSignBabyJubFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}
