package core

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestClaimEthId(t *testing.T) {
	ethId := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	identityFactoryAddr := common.HexToAddress("0x66D0c2F85F1B717168cbB508AfD1c46e07227130")

	c0 := NewClaimEthId(ethId, identityFactoryAddr)

	c1 := NewClaimEthIdFromEntry(c0.Entry())
	c2, err := NewClaimFromEntry(c0.Entry())
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.Equal(t, c0.Address, ethId)
	assert.Equal(t, c0.IdentityFactory, identityFactoryAddr)
	assert.Equal(t, c0.Address, c1.Address)
	assert.Equal(t, c0.IdentityFactory, c1.IdentityFactory)

	assert.Equal(t, c0.Entry().Bytes(), c1.Entry().Bytes())
	assert.Equal(t, c0.Entry().Bytes(), c2.Entry().Bytes())

	e := c0.Entry()
	assert.Equal(t,
		"0x21c4885f4574ea1713b74656cfedef76402b8c8d83a8c8959dadff5b00384130",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"00000000000000000000000066d0c2f85f1b717168cbb508afd1c46e07227130"+
		"000000000000000000000000e0fbce58cfaa72812103f003adce3f284fe5fc7c"+
		"0000000000000000000000000000000000000000000000000000000000000008",
		e.Data.String())
}
