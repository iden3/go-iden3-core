package db

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
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
	// Export(io.WriteCloser) error
	// Import(io.ReadCloser) error
}

type Tx interface {
	Get([]byte) ([]byte, error)
	Put(k, v []byte)
	Add(Tx)
	Commit() error
	Close()
}

// Export is thread safe as long as the Storage.Iterate implementation is thread safe.
func Export(db Storage, out io.WriteCloser) error {
	var buf [10]byte
	if err := db.Iterate(
		func(key, value []byte) (bool, error) {
			keyLenLen := binary.PutUvarint(buf[:], uint64(len(key)))
			if _, err := out.Write(buf[:keyLenLen]); err != nil {
				return false, err
			}
			valueLenLen := binary.PutUvarint(buf[:], uint64(len(value)))
			if _, err := out.Write(buf[:valueLenLen]); err != nil {
				return false, err
			}
			if _, err := out.Write(key); err != nil {
				return false, err
			}
			if _, err := out.Write(value); err != nil {
				return false, err
			}
			return true, nil
		},
	); err != nil {
		return err
	}
	return out.Close()
}

// setCap makes sure that buf has capacity at least n
func setCap(buf []byte, n int) []byte {
	if cap(buf) < n {
		buf = make([]byte, n)
	}
	return buf[:n]
}

// Import is not thread safe.
func Import(db Storage, in io.Reader) error {
	var buf []byte
	tx, err := db.NewTx()
	if err != nil {
		return err
	}
	bufIn := bufio.NewReader(in)
	for {
		keyLen, err := binary.ReadUvarint(bufIn)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		valueLen, err := binary.ReadUvarint(bufIn)
		if err != nil {
			return err
		}
		// use setCap
		buf = setCap(buf, int(keyLen+valueLen))
		key, value := buf[:keyLen], buf[keyLen:keyLen+valueLen]
		if _, err := bufIn.Read(key); err != nil {
			return err
		}
		if _, err := bufIn.Read(value); err != nil {
			return err
		}
		tx.Put(key, value)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
