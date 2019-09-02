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
		"0x1da8cb501998d7cdeb95ab0c65afc0277632ac87d2d74e7a68a0cdca4287649d",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"2d9e82263b94a343ee95d56c810a5a0adb63a439cd5b4944dfb56f09e28c6f04"+
		"0000000000000000000000000000000000000001000000010000000000000001",
		e.Data.String())
	c1 := NewClaimAuthorizeKSignBabyJubFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}
