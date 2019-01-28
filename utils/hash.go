package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// Hash used in this tree, is the [32]byte keccak()
type Hash [32]byte

// hashBytes performs a Keccak256 hash over the bytes
func HashBytes(b ...[]byte) (hash Hash) {
	h := crypto.Keccak256(b...)
	copy(hash[:], h)
	return hash
}

// Hex returns a hex string from the Hash type
func (hash Hash) Hex() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(hash[:]))
}
