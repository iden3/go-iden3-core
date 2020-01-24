package core

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestClaimAuthorizeKSignSecp256k1(t *testing.T) {
	// ClaimAuthorizeKSignSecp256k1
	secKeyHex := "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
	secKey, err := crypto.HexToECDSA(secKeyHex)
	if err != nil {
		panic(err)
	}
	pubKey := secKey.Public().(*ecdsa.PublicKey)
	assert.Equal(t,
		"036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd4d59",
		hex.EncodeToString(crypto.CompressPubkey(pubKey)))
	c0 := NewClaimAuthorizeKSignSecp256k1(pubKey)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x011006c8eb9349422ee77c7721563788ab5fb05e83fb7f070f4355e2e5fbe6e9",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"00036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd"+
		"0000000000000000000000000000000000004d59000000010000000000000004",
		e.Data.String())
	c1, err := NewClaimAuthorizeKSignSecp256k1FromEntry(e)
	if err != nil {
		panic(err)
	}
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}
