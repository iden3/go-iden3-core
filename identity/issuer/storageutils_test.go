package issuer

import (
	"testing"

	"github.com/iden3/go-iden3-core/db"
	"github.com/stretchr/testify/require"
)

func TestStorageValue(t *testing.T) {
	storage := db.NewMemoryStorage()
	sv := NewStorageValue([]byte("key"))

	tx, err := storage.NewTx()
	require.Nil(t, err)
	v0 := uint32(42)
	sv.Set(tx, 42)
	err = tx.Commit()
	require.Nil(t, err)

	tx, err = storage.NewTx()
	require.Nil(t, err)
	v1, err := sv.Get(tx)
	require.Nil(t, err)
	require.Equal(t, v0, v1)
	tx.Close()
}

func TestStorageList(t *testing.T) {
	storage := db.NewMemoryStorage()
	sl := NewStorageList([]byte("list:"))

	type Entry struct {
		Value uint32
	}

	type KV struct {
		Key   []byte
		Value Entry
	}

	entries := []KV{
		KV{Key: []byte("zero"), Value: Entry{Value: 0}},
		KV{Key: []byte("one"), Value: Entry{Value: 1}},
		KV{Key: []byte("two"), Value: Entry{Value: 2}},
		KV{Key: []byte("three"), Value: Entry{Value: 3}},
	}

	tx, err := storage.NewTx()
	require.Nil(t, err)
	sl.Init(tx)
	for _, kv := range entries {
		require.Nil(t, sl.Append(tx, kv.Key, kv.Value))
	}
	err = tx.Commit()
	require.Nil(t, err)

	tx, err = storage.NewTx()
	require.Nil(t, err)
	sl.Init(tx)
	for _, kv := range entries {
		var value Entry
		err := sl.Get(tx, kv.Key, &value)
		require.Nil(t, err)
		require.Equal(t, kv.Value, value)
	}
	tx.Close()

	tx, err = storage.NewTx()
	require.Nil(t, err)
	sl.Init(tx)
	for idx, kv := range entries {
		var value Entry
		key, err := sl.GetByIdx(tx, uint32(idx), &value)
		require.Nil(t, err)
		require.Equal(t, kv.Value, value)
		require.Equal(t, kv.Key, key)
	}
	tx.Close()
}
