package merkletree

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

// Hex returns a hex string from the Hash type
func (hash Hash) Hex() string {
	r := "0x"
	h := hex.EncodeToString(hash[:])
	r = r + h
	return r
}

// Bytes returns a byte array from a Hash
func (hash Hash) Bytes() []byte {
	return hash[:]
}

// HashBytes performs a Keccak256 hash over the bytes
func HashBytes(b []byte) (hash Hash) {
	h := crypto.Keccak256(b)
	copy(hash[:], h)
	return hash
}

// getPath returns the binary path, from the leaf to the root
func getPath(numLevels int, hi Hash) []bool {

	path := []bool{}
	for bitno := numLevels - 2; bitno >= 0; bitno-- {
		path = append(path, testbitmap(hi[:], uint(bitno)))
	}
	return path
}

func comparePaths(b1 []bool, b2 []bool) int {
	for i := len(b1) - 1; i >= 0; i-- {
		if b1[i] != b2[i] {
			return i
		}
	}
	return -1
}

func getEmptiesBetweenIAndPosHash(mt *MerkleTree, iPos int, posHash int) []Hash {
	var sibl []Hash
	for i := iPos; i >= posHash; i-- {
		sibl = append(sibl, EmptyNodeValue)
	}
	return sibl
}

func setbitmap(bitmap []byte, bitno uint) {
	bitmap[uint(len(bitmap))-bitno/8-1] |= 1 << (bitno % 8)
}
func testbitmap(bitmap []byte, bitno uint) bool {
	return bitmap[uint(len(bitmap))-bitno/8-1]&(1<<(bitno%8)) > 0
}
