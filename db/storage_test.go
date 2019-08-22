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

func testIterate(t *testing.T, sto Storage) {
	r := []KV{}
	lister := func(k []byte, v []byte) (bool, error) {
		r = append(r, KV{clone(k), clone(v)})
		return true, nil
	}

	sto1 := sto.WithPrefix([]byte{1})
	err := sto1.Iterate(lister)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(r))

	sto1tx, _ := sto1.NewTx()
	sto1tx.Put([]byte{1}, []byte{4})
	sto1tx.Put([]byte{2}, []byte{5})
	sto1tx.Put([]byte{3}, []byte{6})
	assert.Nil(t, sto1tx.Commit())

	sto2 := sto.WithPrefix([]byte{2})
	sto2tx, _ := sto2.NewTx()
	sto2tx.Put([]byte{1}, []byte{7})
	sto2tx.Put([]byte{2}, []byte{8})
	sto2tx.Put([]byte{3}, []byte{9})
	assert.Nil(t, sto2tx.Commit())

	r = []KV{}
	err = sto1.Iterate(lister)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(r))
	assert.Equal(t, KV{[]byte{1}, []byte{4}}, r[0])
	assert.Equal(t, KV{[]byte{2}, []byte{5}}, r[1])
	assert.Equal(t, KV{[]byte{3}, []byte{6}}, r[2])

	r = []KV{}
	err = sto2.Iterate(lister)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(r))
	assert.Equal(t, KV{[]byte{1}, []byte{7}}, r[0])
	assert.Equal(t, KV{[]byte{2}, []byte{8}}, r[1])
	assert.Equal(t, KV{[]byte{3}, []byte{9}}, r[2])
}

func testConcatTx(t *testing.T, sto Storage) {
	k := []byte{9}

	sto1 := sto.WithPrefix([]byte{1})
	sto2 := sto.WithPrefix([]byte{2})

	// check within tx

	sto1tx, err := sto1.NewTx()
	sto1tx.Put(k, []byte{4, 5, 6})
	sto2tx, err := sto2.NewTx()
	sto2tx.Put(k, []byte{8, 9})

	sto1tx.Add(sto2tx)
	assert.Nil(t, sto1tx.Commit())

	// check outside tx

	v1, err := sto1.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, v1, []byte{4, 5, 6})

	v2, err := sto2.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, v2, []byte{8, 9})
}

func testList(t *testing.T, sto Storage) {
	sto1 := sto.WithPrefix([]byte{1})
	r1, err := sto1.List(100)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(r1))

	sto1tx, _ := sto1.NewTx()
	sto1tx.Put([]byte{1}, []byte{4})
	sto1tx.Put([]byte{2}, []byte{5})
	sto1tx.Put([]byte{3}, []byte{6})
	assert.Nil(t, sto1tx.Commit())

	sto2 := sto.WithPrefix([]byte{2})
	sto2tx, _ := sto2.NewTx()
	sto2tx.Put([]byte{1}, []byte{7})
	sto2tx.Put([]byte{2}, []byte{8})
	sto2tx.Put([]byte{3}, []byte{9})
	assert.Nil(t, sto2tx.Commit())

	r, err := sto1.List(100)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(r))
	assert.Equal(t, r[0], KV{[]byte{1}, []byte{4}})
	assert.Equal(t, r[1], KV{[]byte{2}, []byte{5}})
	assert.Equal(t, r[2], KV{[]byte{3}, []byte{6}})

	r, err = sto1.List(2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(r))
	assert.Equal(t, r[0], KV{[]byte{1}, []byte{4}})
	assert.Equal(t, r[1], KV{[]byte{2}, []byte{5}})

}

func TestLevelDb(t *testing.T) {
	testReturnKnownErrIfNotExists(t, levelDbStorage(t))
	testStorageInsertGet(t, levelDbStorage(t))
	testStorageWithPrefix(t, levelDbStorage(t))
	testConcatTx(t, levelDbStorage(t))
	testList(t, levelDbStorage(t))
	testIterate(t, levelDbStorage(t))
}

func TestMemory(t *testing.T) {
	testReturnKnownErrIfNotExists(t, NewMemoryStorage())
	testStorageInsertGet(t, NewMemoryStorage())
	testStorageWithPrefix(t, NewMemoryStorage())
	testConcatTx(t, NewMemoryStorage())
	testList(t, NewMemoryStorage())
	testIterate(t, NewMemoryStorage())
}
