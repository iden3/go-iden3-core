package db

type MemoryStorage struct {
	prefix []byte
	kv     kvMap
}

type MemoryStorageTx struct {
	s  *MemoryStorage
	kv kvMap
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{[]byte{}, make(kvMap)}
}

func (l *MemoryStorage) Info() string {
	return "in-memory"
}

func (m *MemoryStorage) WithPrefix(prefix []byte) Storage {
	return &MemoryStorage{append(m.prefix, prefix...), m.kv}
}

func (m *MemoryStorage) NewTx() (Tx, error) {
	return &MemoryStorageTx{m, make(kvMap)}, nil
}

// Get retreives a value from a key in the mt.Lvl
func (l *MemoryStorage) Get(key []byte) ([]byte, error) {

	if v, ok := l.kv.Get(append(l.prefix, key[:]...)); ok {
		return v, nil
	}
	return nil, ErrNotFound
}

func (tx *MemoryStorageTx) Get(key []byte) ([]byte, error) {

	if v, ok := tx.kv.Get(key); ok {
		return v, nil
	}
	if v, ok := tx.s.kv.Get(append(tx.s.prefix, key...)); ok {
		return v, nil
	}

	return nil, ErrNotFound
}

func (tx *MemoryStorageTx) Put(k, v []byte) {
	tx.kv.Put(k, v)
}

func (tx *MemoryStorageTx) Commit() error {
	for _, v := range tx.kv {
		tx.s.kv.Put(append(tx.s.prefix, v.k...), v.v)
	}
	tx.kv = nil
	return nil
}

func (tx *MemoryStorageTx) Close() {
	tx.kv = nil
}

func (m *MemoryStorage) Close() {
}
