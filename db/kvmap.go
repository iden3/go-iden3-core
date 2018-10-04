package db

import "crypto/sha256"

type kvEntry struct {
	k []byte
	v []byte
}
type kvMap map[[sha256.Size]byte]kvEntry

func (m kvMap) Get(k []byte) ([]byte, bool) {
	v, ok := m[sha256.Sum256(k)]
	return v.v, ok
}
func (m kvMap) Put(k, v []byte) {
	m[sha256.Sum256(k)] = kvEntry{k, v}
}
