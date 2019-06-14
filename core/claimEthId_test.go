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
		"0x07f139661fb6c709668e5cbbaad48f63358c096fbb60f17c84b26004ef2624ba",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x1063d112ef1c8b044c29c351aa0e84a234e44979624673aa76982b87c293ecf3",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"00000000000000000000000066d0c2f85f1b717168cbb508afd1c46e07227130"+
		"000000000000000000000000e0fbce58cfaa72812103f003adce3f284fe5fc7c"+
		"0000000000000000000000000000000000000000000000000000000000000008",
		e.Data.String())
}
