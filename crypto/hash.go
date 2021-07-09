package crypto

import (
	"encoding/hex"
	"fmt"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-iden3-crypto/utils"
	"math/big"

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

// PoseidonHashBytes hashes a msg byte slice by blocks of 31 bytes encoded as
// little-endian
func PoseidonHashBytes(b []byte) (*big.Int, error) {
	n := 31
	bElems := make([]*big.Int, 0, len(b)/n+1)
	for i := 0; i < len(b)/n; i++ {
		v := big.NewInt(int64(0))
		utils.SetBigIntFromLEBytes(v, b[n*i:n*(i+1)])
		bElems = append(bElems, v)

	}
	if len(b)%n != 0 {
		v := big.NewInt(int64(0))
		utils.SetBigIntFromLEBytes(v, b[(len(b)/n)*n:])
		bElems = append(bElems, v)
	}
	h, err := poseidon.Hash(bElems)
	if err != nil {
		return nil, err
	}
	return h, nil
}
