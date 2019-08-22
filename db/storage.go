package db

import (
	"errors"
)

var ErrNotFound = errors.New("key not found")

type KV struct {
	K []byte
	V []byte
}

type Storage interface {
	NewTx() (Tx, error)
	WithPrefix(prefix []byte) Storage
	Get([]byte) ([]byte, error)
	List(int) ([]KV, error)
	Close()
	Info() string
	Iterate(func([]byte, []byte) (bool, error)) error
}

type Tx interface {
	Get([]byte) ([]byte, error)
	Put(k, v []byte)
	Add(Tx)
	Commit() error
	Close()
}
