package core

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestClaimAuthorizeService(t *testing.T) {
	// ClaimAuthorizeService
	ethAddr := common.BytesToAddress([]byte{
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x3a})
	pubKstr := "af048ddcc131d526699d928e8b8548c5c85fb7d407fc408bb543e4e58f305347f67942a7e56d7dc90bbcecca865f2fbde3118c91516594262f62857136f71dbc"
	c0 := NewClaimAuthorizeService(ServiceTypeRelay, ethAddr.Hex(), pubKstr, "relay.iden3.io")
	e := c0.Entry()
	assert.Equal(t,
		"0x0ee7fb1c970abca8667607eca3974704783f8812bc7f745c1c7ee49a2faf7927",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x2ac15c5d5a255d7d92d84580c9a19b2e6beed42cfd26978b448d6b4abfa6d017",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"00f3b1c89978c483ef94f9ecff889cbef9db68036f3b2dc251e72b7960b8529d"+
		"00f28abb0b5b73fdcc8eed8e707f33d8dd9b50b3e2c6e1957a585903ae3b729a"+
		"00f54d900e54dfb5d19c0e19e5e3abca0d744fee18b72cb8b9cc05f655495983"+
		"0000000000000000000000000000000000000000000000000000000000000006",
		e.Data.String())
	c1 := NewClaimAuthorizeServiceFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
	assert.Equal(t, c0.ServiceType, ServiceTypeRelay)
}
