package adminsrv

import (
	"fmt"
	"math/big"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	merkletree "github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/claimsrv"
	"github.com/iden3/go-iden3-core/services/rootsrv"
	"github.com/iden3/go-iden3-crypto/mimc7"
)

type Service interface {
	Info(id *core.ID) map[string]string
	RawDump(c *gin.Context)
	RawImport(raw map[string]string) (int, error)
	ClaimsDump() map[string]string
	Mimc7(data []*big.Int) (*big.Int, error)
	// AddClaimBasic(indexSlot [400 / 8]byte, dataSlot [496 / 8]byte) (*core.ProofClaim, error)
}

type ServiceImpl struct {
	mt       *merkletree.MerkleTree
	rootsrv  rootsrv.Service
	claimsrv claimsrv.Service
}

func New(mt *merkletree.MerkleTree, rootsrv rootsrv.Service, claimsrv claimsrv.Service) *ServiceImpl {
	return &ServiceImpl{mt, rootsrv, claimsrv}
}

// Info returns the info overview of the Relay
func (as *ServiceImpl) Info(id *core.ID) map[string]string {
	o := make(map[string]string)
	o["db"] = as.mt.Storage().Info()
	o["root"] = as.mt.RootKey().Hex()

	root, err := as.claimsrv.RootSrv().GetRoot(id)
	if err != nil {
		o["root_contract"] = "error getting root from contract"
	} else {
		o["root_contract"] = root.Root.Hex()
		o["root_block_number"] = fmt.Sprint(root.BlockN)
		o["root_block_timestamp"] = fmt.Sprint(root.BlockTimestamp)
	}

	return o
}

// RawDump returns all the key and values from the database
func (as *ServiceImpl) RawDump(c *gin.Context) {
	// var out string
	sto := as.mt.Storage()
	sto.Iterate(func(key, value []byte) (bool, error) {
		c.String(200, common3.HexEncode(key)+", "+common3.HexEncode(value)+"\n")
		return true, nil
	})
	return
}

// RawImport imports the key and values from the RawDump() to the database
func (as *ServiceImpl) RawImport(raw map[string]string) (int, error) {
	count := 0

	tx, err := as.mt.Storage().NewTx()
	if err != nil {
		return count, err
	}

	defer func() {
		if err == nil {
			tx.Commit()
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
func (as *ServiceImpl) ClaimsDump() map[string]string {
	data := make(map[string]string)
	sto := as.mt.Storage()
	sto.Iterate(func(key, value []byte) (bool, error) {
		if value[0] == byte(merkletree.NodeTypeLeaf) {
			data[common3.HexEncode(key)] = common3.HexEncode(value)
		}
		return true, nil
	})
	return data
}

// Mimc7 performs the MIMC7 hash over a given data
func (as *ServiceImpl) Mimc7(data []*big.Int) (*big.Int, error) {
	helement, err := mimc7.Hash(data, nil)
	return (*big.Int)(helement), err

}

// DECPRECATED
// func (as *ServiceImpl) AddClaimBasic(indexSlot [400 / 8]byte, dataSlot [496 / 8]byte) (*core.ProofClaim, error) {
// 	// TODO check if indexSlot and dataSlot fit inside R element
// 	// var indexSlot [400 / 8]byte
// 	// var dataSlot [496 / 8]byte
// 	// copy(indexSlot[:], indexData[:400/8])
// 	// copy(dataSlot[:], data[:496/8])
// 	claim := core.NewClaimBasic(indexSlot, dataSlot)
//
// 	err := as.mt.Add(claim.Entry())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// update Relay Root in Smart Contract
// 	as.rootsrv.SetRoot(*as.mt.RootKey())
//
// 	proofClaim, err := as.claimsrv.GetClaimProofByHi(claim.Entry().HIndex())
// 	if err != nil {
// 		fmt.Println("err", err.Error())
// 		return nil, err
// 	}
// 	return proofClaim, nil
// }
