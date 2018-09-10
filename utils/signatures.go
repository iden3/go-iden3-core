package utils

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3/merkletree"
)

// Sign performs the signature over a Hash
func Sign(msgHash merkletree.Hash, privK *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(msgHash[:], privK)
	if err != nil {
		return []byte{}, err
	}
	return sig, nil
}

// VerifySig verifies a given signature and the msgHash with the expected address
func VerifySig(addr common.Address, sig, msgHash []byte) bool {
	recoveredPub, err := crypto.Ecrecover(msgHash, sig)
	if err != nil {
		fmt.Printf("ECRecover error: %s", err)
		return false
	}
	pubKey, _ := crypto.UnmarshalPubkey(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return bytes.Equal(addr.Bytes(), recoveredAddr.Bytes())
}
