package adminsrv

import (
	"bytes"
	"fmt"
	"math/big"

	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/crypto/mimc7"
	merkletree "github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/rootsrv"
)

type Service interface {
	Info() map[string]string
	RawDump() map[string]string
	RawImport(raw map[string]string) (int, error)
	ClaimsDump() map[string]string
	Mimc7(data []*big.Int) (*big.Int, error)
	AddGenericClaim(indexData, data []byte) (claimsrv.ProofOfRelayClaim, error)
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
func (as *ServiceImpl) Info() map[string]string {
	o := make(map[string]string)
	o["db"] = as.mt.Storage().Info()
	o["root"] = as.mt.Root().Hex()
	return o
}

// RawDump returns all the key and values from the database
func (as *ServiceImpl) RawDump() map[string]string {
	// var out string
	data := make(map[string]string)
	sto := as.mt.Storage()
	sto.Iterate(func(key, value []byte) {
		// out = out + "key: " + common3.BytesToHex(key) + ", value: " + common3.BytesToHex(value) + "\n"
		data[common3.BytesToHex(key)] = common3.BytesToHex(value)
	})
	return data
}

// RawImport imports the key and values from the RawDump() to the database
func (as *ServiceImpl) RawImport(raw map[string]string) (int, error) {
	fmt.Println("raw", raw)
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
		kBytes, err := common3.HexToBytes(k)
		if err != nil {
			return count, err
		}
		vBytes, err := common3.HexToBytes(v)
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
	sto.Iterate(func(key, value []byte) {
		if value[0] == merkletree.NodeTypeLeaf {
			data[common3.BytesToHex(key)] = common3.BytesToHex(value)
		}
	})
	return data
}

// Mimc7 performs the MIMC7 hash over a given data
func (as *ServiceImpl) Mimc7(data []*big.Int) (*big.Int, error) {
	ielements, err := mimc7.BigIntsToRElems(data)
	if err != nil {
		return &big.Int{}, err
	}
	helement := mimc7.Hash(ielements)
	return (*big.Int)(helement), nil

}

func (as *ServiceImpl) AddGenericClaim(indexData, data []byte) (claimsrv.ProofOfRelayClaim, error) {
	claim := core.NewGenericClaim("iden3.io", "default", indexData, data)

	err := as.mt.Add(claim)
	if err != nil {
		fmt.Println("a")
		return claimsrv.ProofOfRelayClaim{}, err
	}

	// update Relay Root in Smart Contract
	as.rootsrv.SetRoot(as.mt.Root())

	proofOfClaim, err := as.claimsrv.GetRelayClaimByHi(claim.Hi())
	if err != nil {
		fmt.Println("err", err.Error())
		return claimsrv.ProofOfRelayClaim{}, err
	}
	return proofOfClaim, nil
}
