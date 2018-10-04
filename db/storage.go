package db

import "errors"

var ErrNotFound = errors.New("key not found")

type Storage interface {
	NewTx() (Tx, error)
	WithPrefix(prefix []byte) Storage
	Get([]byte) ([]byte, error)
	Close()
	Info() string
}

type Tx interface {
	Get([]byte) ([]byte, error)
	Put(k, v []byte)
	Commit() error
	Close()
}
