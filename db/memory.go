package db

import (
	"bytes"
	"sort"
)

type MemoryStorage struct {
	prefix []byte
	kv     kvMap
}

type MemoryStorageTx struct {
	s  *MemoryStorage
	kv kvMap
}

func NewMemoryStorage() *MemoryStorage {
	kvmap := make(kvMap)
	return &MemoryStorage{[]byte{}, kvmap}
}

func (l *MemoryStorage) Info() string {
	return "in-memory"
}

func (m *MemoryStorage) WithPrefix(prefix []byte) Storage {
	return &MemoryStorage{concat(m.prefix, prefix), m.kv}
}

func (m *MemoryStorage) NewTx() (Tx, error) {
	return &MemoryStorageTx{m, make(kvMap)}, nil
}

// Get retreives a value from a key in the mt.Lvl
func (l *MemoryStorage) Get(key []byte) ([]byte, error) {

	if v, ok := l.kv.Get(concat(l.prefix, key[:])); ok {
		return v, nil
	}
	return nil, ErrNotFound
}

func (tx *MemoryStorageTx) Get(key []byte) ([]byte, error) {

	if v, ok := tx.kv.Get(concat(tx.s.prefix, key)); ok {
		return v, nil
	}
	if v, ok := tx.s.kv.Get(concat(tx.s.prefix, key)); ok {
		return v, nil
	}

	return nil, ErrNotFound
}

func (tx *MemoryStorageTx) Put(k, v []byte) {
	tx.kv.Put(concat(tx.s.prefix, k), v)
}

func (tx *MemoryStorageTx) Commit() error {
	for _, v := range tx.kv {
		tx.s.kv.Put(v.K, v.V)
	}
	tx.kv = nil
	return nil
}

func (tx *MemoryStorageTx) Add(atx Tx) {
	mstx := atx.(*MemoryStorageTx)
	for _, v := range mstx.kv {
		tx.kv.Put(v.K, v.V)
	}
}

func (tx *MemoryStorageTx) Close() {
	tx.kv = nil
}

func (m *MemoryStorage) Close() {
}

func (l *MemoryStorage) List(limit int) ([]KV, error) {

	ret := make([]KV, 0)
	for _, v := range l.kv {

		if len(v.K) < len(l.prefix) || !bytes.Equal(v.K[:len(l.prefix)], l.prefix) {
			continue
		}
		localkey := v.K[len(l.prefix):]
		ret = append(ret, KV{localkey, v.V})

	}
	sort.SliceStable(ret, func(i, j int) bool { return bytes.Compare(ret[i].K, ret[j].K) < 0 })
	if len(ret) > limit {
		ret = ret[:limit]
	}
	return ret, nil
}
