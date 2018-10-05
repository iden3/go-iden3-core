package db

import (
	"crypto/sha256"
)

type kvMap map[[sha256.Size]byte]KV

func (m kvMap) Get(k []byte) ([]byte, bool) {
	v, ok := m[sha256.Sum256(k)]
	return v.V, ok
}
func (m kvMap) Put(k, v []byte) {
	m[sha256.Sum256(k)] = KV{k, v}
}
