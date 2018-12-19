package adminsrv

import (
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
	RawDump() string
	ClaimsDump() string
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
func (as *ServiceImpl) RawDump() string {
	var out string
	sto := as.mt.Storage()
	sto.Iterate(func(key, value []byte) {
		out = out + "key: " + common3.BytesToHex(key) + ", value: " + common3.BytesToHex(value) + "\n"
	})
	return out
}

// ClaimsDump returns all the claims key and values from the database
func (as *ServiceImpl) ClaimsDump() string {
	var out string
	sto := as.mt.Storage()
	sto.Iterate(func(key, value []byte) {
		if value[0] == byte(1) { // TODO when the new merkletree version is ready, instead of byte(1) use the type indicator
			out = out + "key: " + common3.BytesToHex(key) + ", value: " + common3.BytesToHex(value) + "\n"
		}
	})
	return out
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
