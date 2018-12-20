package utils

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

// Hash used in this tree, is the [32]byte keccak()
type Hash [32]byte

// hashBytes performs a Keccak256 hash over the bytes
func HashBytes(b []byte) (hash Hash) {
	h := crypto.Keccak256(b)
	copy(hash[:], h)
	return hash
}

// Bytes returns a byte array from a Hash
func (hash Hash) Bytes() []byte {
	return hash[:]
}

// Hex returns a hex string from the Hash type
func (hash Hash) Hex() string {
	r := "0x"
	h := hex.EncodeToString(hash[:])
	r = r + h
	return r
}
