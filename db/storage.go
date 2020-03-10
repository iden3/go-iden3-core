package db

import (
	"errors"
)

var ErrNotFound = errors.New("key not found")

type KV struct {
	K []byte
	V []byte
}

type StorageLegacy interface {
	NewTx() (TxLegacy, error)
	WithPrefix(prefix []byte) StorageLegacy
	Get([]byte) ([]byte, error)
	List(int) ([]KV, error)
	Close()
	Info() string
	Iterate(func([]byte, []byte) (bool, error)) error
}

type TxLegacy interface {
	Get([]byte) ([]byte, error)
	Put(k, v []byte)
	Add(TxLegacy)
	Commit() error
	Close()
}

type Storage interface {
	NewTx() (Tx, error)
	WithPrefix(prefix []byte) Storage
	Get([]byte) ([]byte, error)
	Close()
	Iterate(func([]byte, []byte) (bool, error)) error
}

type Tx interface {
	Get([]byte) ([]byte, error)
	Put(k, v []byte)
	Delete(k []byte) error
	Commit() error
	Close()
}
