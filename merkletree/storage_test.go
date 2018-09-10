package merkletree

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestingStorage(f Fatalable) Storage {
	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		f.Fatal(err)
		return nil
	}
	sto, err := NewLevelDbStorage(dir)
	if err != nil {
		f.Fatal(err)
		return nil
	}
	return sto
}

func TestStorageInsertGet(t *testing.T) {
	sto := newTestingStorage(t)

	h := HashBytes([]byte("a"))

	tx, err := sto.NewTx()
	assert.Nil(t, err)
	tx.Insert(h, normalNodeType, 0, []byte("data"))
	assert.Nil(t, tx.Commit())

	nodeType, indexLength, nodeBytes, err := tx.Get(h)
	assert.Nil(t, err)

	assert.Equal(t, nodeType, byte(normalNodeType))
	assert.Equal(t, indexLength, uint32(0))
	assert.Equal(t, []byte("data"), nodeBytes)
}

func TestStorageWithPrefix(t *testing.T) {

	h1 := HashBytes([]byte{1})

	sto := newTestingStorage(t)
	sto1 := sto.WithPrefix([]byte{1})
	sto2 := sto.WithPrefix([]byte{2})

	// check within tx

	sto1tx, err := sto1.NewTx()
	assert.Nil(t, err)
	sto1tx.Insert(h1, 1, 2, []byte{4, 5, 6})
	typ, len, val, err := sto1tx.Get(h1)
	assert.Nil(t, err)
	assert.Equal(t, typ, uint8(1))
	assert.Equal(t, len, uint32(2))
	assert.Equal(t, val, []byte{4, 5, 6})
	assert.Nil(t, sto1tx.Commit())

	sto2tx, err := sto2.NewTx()
	assert.Nil(t, err)
	sto2tx.Insert(h1, 4, 5, []byte{6, 7})
	typ, len, val, err = sto2tx.Get(h1)
	assert.Nil(t, err)
	assert.Equal(t, typ, uint8(4))
	assert.Equal(t, len, uint32(5))
	assert.Equal(t, val, []byte{6, 7})
	assert.Nil(t, sto2tx.Commit())

	// check outside tx

	typ, len, val, err = sto1.Get(h1)
	assert.Nil(t, err)
	assert.Equal(t, typ, uint8(1))
	assert.Equal(t, len, uint32(2))
	assert.Equal(t, val, []byte{4, 5, 6})

	typ, len, val, err = sto2.Get(h1)
	assert.Nil(t, err)
	assert.Equal(t, typ, uint8(4))
	assert.Equal(t, len, uint32(5))
	assert.Equal(t, val, []byte{6, 7})
}
