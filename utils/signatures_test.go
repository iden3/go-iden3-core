package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/stretchr/testify/assert"
)

const (
	testPrivKHex = "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"
)

func TestVerifySig(t *testing.T) {
	signatureHex := "0xd45d0a89d5bbe9770ce3241cf8672aefdcdd2f204b5d63c8500e9770335314c532d0b16e3b41caabd1dde37c62a7cb6273d97c09b7394080ae7d1d8d211e05fb00"
	signature, err := common3.HexToBytes(signatureHex)
	assert.Nil(t, err)
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	assert.Nil(t, err)
	msgHash := HashBytes([]byte("to sign"))
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)
	assert.True(t, VerifySig(testAddr, signature, msgHash[:]))
}

func TestVerifySigFromJS(t *testing.T) {
	// verify signature performed in iden3js
	signatureHex := "0x5413b44384531e9e92bdd80ff21cea7449441dcfff6f4ed0f90864583e3fcade3d5c8857672b473f71d09355e034dba11bb2ca4aa73c55c534293fdca68941041c"
	signature, err := common3.HexToBytes(signatureHex)
	signature[64] -= 27
	assert.Nil(t, err)
	testAddr := common.HexToAddress("0xBc8C480E68d0895f1E410f4e4eA6E2d6b160Ca9F")
	msgHash := EthHash([]byte("test"))
	assert.True(t, VerifySig(testAddr, signature, msgHash[:]))
}

func TestSignAndVerify(t *testing.T) {
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	assert.Nil(t, err)
	msgHash := EthHash([]byte("test"))
	signature, err := crypto.Sign(msgHash[:], testPrivK)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)
	assert.True(t, VerifySig(testAddr, signature, msgHash[:]))
}
