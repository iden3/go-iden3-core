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

func (l *MemoryStorage) Iterate(f func([]byte, []byte) (bool, error)) error {
	kvs := make([]KV, 0)
	for _, v := range l.kv {
		if len(v.K) < len(l.prefix) || !bytes.Equal(v.K[:len(l.prefix)], l.prefix) {
			continue
		}
		localkey := v.K[len(l.prefix):]
		kvs = append(kvs, KV{localkey, v.V})

	}
	sort.SliceStable(kvs, func(i, j int) bool { return bytes.Compare(kvs[i].K, kvs[j].K) < 0 })

	for _, kv := range kvs {
		if cont, err := f(kv.K, kv.V); err != nil {
			return err
		} else if !cont {
			break
		}
	}
	return nil
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
	ret := []KV{}
	err := l.Iterate(func(key []byte, value []byte) (bool, error) {
		ret = append(ret, KV{clone(key), clone(value)})
		if len(ret) == limit {
			return false, nil
		}
		return true, nil
	})
	return ret, err
}
