package db

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func levelDbStorage(t *testing.T) Storage {
	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		t.Fatal(err)
		return nil
	}
	sto, err := NewLevelDbStorage(dir, false)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return sto
}

func testReturnKnownErrIfNotExists(t *testing.T, sto Storage) {
	k := []byte("key")

	tx, err := sto.NewTx()
	assert.Nil(t, err)
	_, err = tx.Get(k)
	assert.EqualError(t, err, ErrNotFound.Error())
}

func testStorageInsertGet(t *testing.T, sto Storage) {
	key := []byte("key")
	value := []byte("data")

	tx, err := sto.NewTx()
	assert.Nil(t, err)
	tx.Put(key, value)
	v, err := tx.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, v)
	assert.Nil(t, tx.Commit())

	tx, err = sto.NewTx()
	assert.Nil(t, err)
	v, err = tx.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, v)
}

func testStorageWithPrefix(t *testing.T, sto Storage) {

	k := []byte{9}

	sto1 := sto.WithPrefix([]byte{1})
	sto2 := sto.WithPrefix([]byte{2})

	// check within tx

	sto1tx, err := sto1.NewTx()
	assert.Nil(t, err)
	sto1tx.Put(k, []byte{4, 5, 6})
	v1, err := sto1tx.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, v1, []byte{4, 5, 6})
	assert.Nil(t, sto1tx.Commit())

	sto2tx, err := sto2.NewTx()
	assert.Nil(t, err)
	sto2tx.Put(k, []byte{8, 9})
	v2, err := sto2tx.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, v2, []byte{8, 9})
	assert.Nil(t, sto2tx.Commit())

	// check outside tx

	v1, err = sto1.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, v1, []byte{4, 5, 6})

	v2, err = sto2.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, v2, []byte{8, 9})
}

func TestLevelDb(t *testing.T) {
	testReturnKnownErrIfNotExists(t, levelDbStorage(t))
	testStorageInsertGet(t, levelDbStorage(t))
	testStorageWithPrefix(t, levelDbStorage(t))
}

func TestMemory(t *testing.T) {
	testReturnKnownErrIfNotExists(t, NewMemoryStorage())
	testStorageInsertGet(t, NewMemoryStorage())
	testStorageWithPrefix(t, NewMemoryStorage())
}
