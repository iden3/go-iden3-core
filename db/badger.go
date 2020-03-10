package db

import (
	"github.com/dgraph-io/badger/v2"
	log "github.com/sirupsen/logrus"
)

type BadgerTx struct {
	prefix []byte
	bgTx   *badger.Txn
}

// Get retreives a value from a key in the mt.Lvl
func (tx *BadgerTx) Get(key []byte) ([]byte, error) {
	item, err := tx.bgTx.Get(append(tx.prefix, key[:]...))
	if err == badger.ErrKeyNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

// Insert saves a key:value into the mt.Lvl
func (tx *BadgerTx) Put(k, v []byte) {
	err := tx.bgTx.Set(concat(tx.prefix, k[:]), v)
	if err != nil {
		panic(err)
	}
}

func (tx *BadgerTx) Delete(k []byte) error {
	err := tx.bgTx.Delete(concat(tx.prefix, k[:]))
	if err == badger.ErrKeyNotFound {
		return ErrNotFound
	}
	return err
}

func (tx *BadgerTx) Commit() error {
	return tx.bgTx.Commit()
}

func (tx *BadgerTx) Close() {
	tx.bgTx.Discard()
}

type BadgerStorage struct {
	db     *badger.DB
	prefix []byte
}

func NewBadgerStorage(path string, logging bool) (*BadgerStorage, error) {
	opt := badger.DefaultOptions(path)
	if logging {
		opt = opt.WithEventLogging(true)
		opt = opt.WithLogger(log.StandardLogger())
	} else {
		opt = opt.WithEventLogging(false)
	}
	return NewBadgerStorageWithOpt(opt)
}

func NewMemoryStorage() *BadgerStorage {
	s, err := NewBadgerStorageWithOpt(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic(err)
	}
	return s
}

func NewBadgerStorageWithOpt(opt badger.Options) (*BadgerStorage, error) {
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}
	return &BadgerStorage{db, []byte{}}, nil
}

func (s *BadgerStorage) WithPrefix(prefix []byte) Storage {
	return &BadgerStorage{s.db, concat(s.prefix, prefix)}
}

func (s *BadgerStorage) NewTx() (Tx, error) {
	bgTx := s.db.NewTransaction(true)
	return &BadgerTx{prefix: s.prefix, bgTx: bgTx}, nil
}

// Get retreives a value from a key in the mt.Lvl
func (s *BadgerStorage) Get(key []byte) ([]byte, error) {
	var v []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(s.prefix, key[:]...))
		if err == badger.ErrKeyNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		v, err = item.ValueCopy(nil)
		return err
	})
	return v, err
}

func (s *BadgerStorage) Iterate(f func([]byte, []byte) (bool, error)) error {
	return s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(s.prefix); it.ValidForPrefix(s.prefix); it.Next() {
			item := it.Item()
			k := item.Key()[len(s.prefix):]
			var cont bool
			err := item.Value(func(v []byte) error {
				var err error
				cont, err = f(k, v)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
			if !cont {
				break
			}
		}
		return nil
	})
}

func (s *BadgerStorage) Close() {
	if err := s.db.Close(); err != nil {
		panic(err)
	}
}
