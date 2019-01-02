package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	Base  string `json:"base"`
	Nonce int    `json:"nonce"`
}

func (td testData) IncrementNonce() PoWData {
	td.Nonce++
	return td
}

func TestPoW(t *testing.T) {
	testData := testData{
		"test",
		0,
	}
	data, err := PoW(testData, 2)
	assert.Nil(t, err)
	b, err := json.Marshal(data)
	assert.Nil(t, err)
	hash := HashBytes(b)
	assert.True(t, CheckPoW(hash, 2))
	assert.True(t, !CheckPoW(hash, 3))

	assert.Equal(t, hash.Hex(), "0x0000dcc15e80ddbf1aad5d2c207084f3058f7353eacad6bcabe795a318c50698")
}

func TestCheckPoW(t *testing.T) {
	testData := testData{
		"test",
		129451,
	}
	b, err := json.Marshal(testData)
	assert.Nil(t, err)
	hash := HashBytes(b)
	assert.True(t, CheckPoW(hash, 2))
	assert.True(t, !CheckPoW(hash, 3))

	assert.Equal(t, hash.Hex(), "0x0000dcc15e80ddbf1aad5d2c207084f3058f7353eacad6bcabe795a318c50698")
}
