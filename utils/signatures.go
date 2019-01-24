package utils

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Sign performs the signature over a Hash
func Sign(h Hash, ks *keystore.KeyStore, acc accounts.Account) ([]byte, error) {
	return ks.SignHash(acc, h[:])
}

func EthHash(b []byte) Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(b), b)
	return HashBytes([]byte(msg))
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

// VerifySigBytes verifies the signature of a byte array given an ethereum address.
func VerifySigBytes(addr common.Address, sig, msg []byte) bool {
	msgHash := EthHash(msg)
	return VerifySig(addr, sig, msgHash[:])
}

// VerifySigBytesDate verifies the signature of a byte array with a date
// appended given an ethereum address.
func VerifySigBytesDate(addr common.Address, sig, msg []byte, date uint64) bool {
	dateBytes := Uint64ToEthBytes(date)
	return VerifySigBytes(addr, sig, append(msg[:], dateBytes...))
}
