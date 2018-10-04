package db

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

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
	if err = ldb.Put([]byte("init"), []byte{1}, nil); err != nil {
		return nil, err
	}
	return &LevelDbStorage{ldb, []byte{}}, nil
}

type storageInfo struct {
	KeyCount int
}

func (l *LevelDbStorage) Info() string {

	keycount := 0
	iter := l.ldb.NewIterator(nil, nil)
	for iter.Next() {
		keycount++
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return err.Error()
	}
	json, _ := json.MarshalIndent(
		storageInfo{keycount},
		"", "  ",
	)
	return string(json)
}

func (l *LevelDbStorage) WithPrefix(prefix []byte) Storage {
	return &LevelDbStorage{l.ldb, append(l.prefix, prefix...)}
}

func (l *LevelDbStorage) NewTx() (Tx, error) {
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

// Get retreives a value from a key in the mt.Lvl
func (l *LevelDbStorageTx) Get(key []byte) ([]byte, error) {
	var err error

	if value, ok := l.cache.Get(key); ok {
		return value, nil
	}

	value, err := l.ldb.Get(append(l.prefix, key[:]...), nil)
	if err == errors.ErrNotFound {
		return nil, ErrNotFound
	}

	return value, err
}

// Insert saves a key:value into the mt.Lvl
func (l *LevelDbStorageTx) Put(k, v []byte) {
	l.cache.Put(k, v)
}

func (l *LevelDbStorageTx) Commit() error {

	var batch leveldb.Batch
	for _, v := range l.cache {
		batch.Put(append(l.prefix, v.k[:]...), v.v)
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
