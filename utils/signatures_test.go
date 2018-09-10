package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/stretchr/testify/assert"
)

const (
	testPrivKHex = "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"
)

func TestSign(t *testing.T) {
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	assert.Nil(t, err)
	msgHash := merkletree.HashBytes([]byte("to sign"))
	signature, err := Sign(msgHash, testPrivK)
	assert.Equal(t, "0xd45d0a89d5bbe9770ce3241cf8672aefdcdd2f204b5d63c8500e9770335314c532d0b16e3b41caabd1dde37c62a7cb6273d97c09b7394080ae7d1d8d211e05fb00", common3.BytesToHex(signature))
}

func TestVerifySig(t *testing.T) {
	signatureHex := "0xd45d0a89d5bbe9770ce3241cf8672aefdcdd2f204b5d63c8500e9770335314c532d0b16e3b41caabd1dde37c62a7cb6273d97c09b7394080ae7d1d8d211e05fb00"
	signature, err := common3.HexToBytes(signatureHex)
	assert.Nil(t, err)
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	assert.Nil(t, err)
	msgHash := merkletree.HashBytes([]byte("to sign"))
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)
	assert.True(t, VerifySig(testAddr, signature, msgHash[:]))
}

func TestSignAndVerify(t *testing.T) {
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	assert.Nil(t, err)
	msgHash := merkletree.HashBytes([]byte("to sign"))
	signature, err := Sign(msgHash, testPrivK)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)
	assert.True(t, VerifySig(testAddr, signature, msgHash[:]))
}
