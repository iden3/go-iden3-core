package issuer

import (
	"testing"

	"github.com/iden3/go-iden3-core/db"
	"github.com/stretchr/testify/require"
)

func TestUniqueNonceGen(t *testing.T) {
	storage := db.NewMemoryStorage()
	nonceGen := NewUniqueNonceGen(NewStorageValue([]byte("nonceIdx")))
	tx, err := storage.NewTx()
	require.Nil(t, err)
	nonceGen.Init(tx)

	n0, err := nonceGen.Next(tx)
	require.Nil(t, err)
	require.Equal(t, uint32(0), n0)
	require.Nil(t, err)
	n1, err := nonceGen.Next(tx)
	require.Nil(t, err)
	require.Equal(t, uint32(1), n1)
	n2, err := nonceGen.Next(tx)
	require.Nil(t, err)
	require.Equal(t, uint32(2), n2)
	err = tx.Commit()
	require.Nil(t, err)
}
