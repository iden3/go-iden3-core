package db

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// LEGACY

type LevelDbStorage struct {
	ldb    *leveldb.DB
	prefix []byte
}

type LevelDbStorageTx struct {
	*LevelDbStorage
	cache kvMap
}

func NewLevelDbStorage(path string, errorIfMissing bool) (*LevelDbStorage, error) {
	o := &opt.Options{
		ErrorIfMissing: errorIfMissing,
	}
	ldb, err := leveldb.OpenFile(path, o)
	if err != nil {
		return nil, err
	}
	return &LevelDbStorage{ldb, []byte{}}, nil
}

type storageInfo struct {
	KeyCount   int
	ClaimCount int
}

func (l *LevelDbStorage) Info() string {
	snapshot, err := l.ldb.GetSnapshot()
	if err != nil {
		return err.Error()
	}

	keycount := 0
	claimcount := 0
	iter := snapshot.NewIterator(nil, nil)
	for iter.Next() {
		if iter.Value()[0] == byte(1) { // TODO when the new merkletree version is ready, instead of byte(1) use the type indicator
			claimcount++
		}

		keycount++
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return err.Error()
	}
	json, _ := json.MarshalIndent(
		storageInfo{
			KeyCount:   keycount,
			ClaimCount: claimcount,
		},
		"", "  ",
	)
	return string(json)
}

func (l *LevelDbStorage) WithPrefix(prefix []byte) StorageLegacy {
	return &LevelDbStorage{l.ldb, concat(l.prefix, prefix)}
}

func (l *LevelDbStorage) NewTx() (TxLegacy, error) {
	return &LevelDbStorageTx{l, make(kvMap)}, nil
}

// Get retreives a value from a key in the mt.Lvl
func (l *LevelDbStorage) Get(key []byte) ([]byte, error) {
	v, err := l.ldb.Get(append(l.prefix, key[:]...), nil)
	if err == errors.ErrNotFound {
		return nil, ErrNotFound
	}
	return v, err
}

func (l *LevelDbStorage) Iterate(f func([]byte, []byte) (bool, error)) error {
	// FIXME: Use the prefix!
	snapshot, err := l.ldb.GetSnapshot()
	if err != nil {
		return err
	}
	iter := snapshot.NewIterator(util.BytesPrefix(l.prefix), nil)
	defer iter.Release()
	for iter.Next() {
		localKey := iter.Key()[len(l.prefix):]
		if cont, err := f(localKey, iter.Value()); err != nil {
			return err
		} else if !cont {
			break
		}
	}
	iter.Release()
	return iter.Error()
}

// Get retreives a value from a key in the mt.Lvl
func (l *LevelDbStorageTx) Get(key []byte) ([]byte, error) {
	var err error

	fullkey := concat(l.prefix, key)

	if value, ok := l.cache.Get(fullkey); ok {
		return value, nil
	}

	value, err := l.ldb.Get(fullkey, nil)
	if err == errors.ErrNotFound {
		return nil, ErrNotFound
	}

	return value, err
}

// Insert saves a key:value into the mt.Lvl
func (tx *LevelDbStorageTx) Put(k, v []byte) {
	tx.cache.Put(concat(tx.prefix, k[:]), v)
}

func (tx *LevelDbStorageTx) Add(atx TxLegacy) {
	ldbtx := atx.(*LevelDbStorageTx)
	for _, v := range ldbtx.cache {
		tx.cache.Put(v.K, v.V)
	}
}

func (l *LevelDbStorageTx) Commit() error {

	var batch leveldb.Batch
	for _, v := range l.cache {
		batch.Put(v.K, v.V)
	}

	l.cache = nil
	return l.ldb.Write(&batch, nil)
}

func (l *LevelDbStorageTx) Close() {
	l.cache = nil
}

func (l *LevelDbStorage) Close() {
	if err := l.ldb.Close(); err != nil {
		panic(err)
	}
	log.Info("Database closed")
}

func (l *LevelDbStorage) LevelDB() *leveldb.DB {
	return l.ldb
}

func (l *LevelDbStorage) List(limit int) ([]KV, error) {
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
