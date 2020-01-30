package idenadminutils

import (
	"math/big"

	// "github.com/ethereum/go-ethereum/common"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/components/idenmanager"
	merkletree "github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/idenstatewriter"
	"github.com/iden3/go-iden3-crypto/mimc7"
)

type IdenAdminUtils struct {
	mt              *merkletree.MerkleTree
	idenStateWriter idenstatewriter.IdenStateWriter
	idenManager     *idenmanager.IdenManager
}

func New(mt *merkletree.MerkleTree, idenStateWriter idenstatewriter.IdenStateWriter, idenManager *idenmanager.IdenManager) *IdenAdminUtils {
	return &IdenAdminUtils{mt, idenStateWriter, idenManager}
}

// RawDump returns all the key and values from the database
func (a *IdenAdminUtils) RawDump(f func(key, value string)) {
	// var out string
	sto := a.mt.Storage()
	err := sto.Iterate(func(key, value []byte) (bool, error) {
		f(common3.HexEncode(key), common3.HexEncode(value))
		return true, nil
	})
	if err != nil {
		panic(err)
	}
}

// RawImport imports the key and values from the RawDump() to the database
func (a *IdenAdminUtils) RawImport(raw map[string]string) (int, error) {
	count := 0

	tx, err := a.mt.Storage().NewTx()
	if err != nil {
		return count, err
	}

	defer func() {
		if err == nil {
			if err := tx.Commit(); err != nil {
				tx.Close()
			}

		} else {
			tx.Close()
		}
	}()

	for k, v := range raw {
		kBytes, err := common3.HexDecode(k)
		if err != nil {
			return count, err
		}
		vBytes, err := common3.HexDecode(v)
		if err != nil {
			return count, err
		}
		tx.Put(kBytes, vBytes)
		count++
	}
	return count, nil
}

// ClaimsDump returns all the claims key and values from the database
func (a *IdenAdminUtils) ClaimsDump() map[string]string {
	data := make(map[string]string)
	sto := a.mt.Storage()
	if err := sto.Iterate(func(key, value []byte) (bool, error) {
		if value[0] == byte(merkletree.NodeTypeLeaf) {
			data[common3.HexEncode(key)] = common3.HexEncode(value)
		}
		return true, nil
	}); err != nil {
		panic(err)
	}

	return data
}

// Mimc7 performs the MIMC7 hash over a given data
func (a *IdenAdminUtils) Mimc7(data []*big.Int) (*big.Int, error) {
	helement, err := mimc7.Hash(data, nil)
	return (*big.Int)(helement), err

}
