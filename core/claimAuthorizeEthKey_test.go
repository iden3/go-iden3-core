package core

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestClaimAuthEthKey(t *testing.T) {
	ethKey := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	ethKeyType := EthKeyTypeUpgrade

	c0 := NewClaimAuthEthKey(ethKey, ethKeyType)

	c1 := NewClaimAuthEthKeyFromEntry(c0.Entry())
	c2, err := NewClaimFromEntry(c0.Entry())
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.Equal(t, c0.EthKey, ethKey)
	assert.Equal(t, c0.EthKeyType, binary.BigEndian.Uint32(ethKeyType[:]))
	assert.Equal(t, c0.EthKey, c1.EthKey)
	assert.Equal(t, c0.EthKeyType, c1.EthKeyType)
	assert.Equal(t, c0.Type(), c1.Type())
	assert.Equal(t, c0.Type(), *ClaimTypeAuthEthKey)

	assert.Equal(t, c0.Entry().Bytes(), c1.Entry().Bytes())
	assert.Equal(t, c0.Entry().Bytes(), c2.Entry().Bytes())

	e := c0.Entry()
	assert.Equal(t,
		"0x0ce27cf2190dfa6ee36276e79335942c28a08dbc5ef8c564ed2f337d5c85b666",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"000000000000000000000002e0fbce58cfaa72812103f003adce3f284fe5fc7c"+
		"0000000000000000000000000000000000000000000000000000000000000009",
		e.Data.String())
}
