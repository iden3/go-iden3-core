package adminsrv

import (
	common3 "github.com/iden3/go-iden3/common"
	merkletree "github.com/iden3/go-iden3/merkletree"
)

type Service interface {
	RawDump() string
	ClaimsDump() string
}

type ServiceImpl struct {
	mt *merkletree.MerkleTree
}

func New(mt *merkletree.MerkleTree) *ServiceImpl {
	return &ServiceImpl{mt}
}

func (as *ServiceImpl) RawDump() string {
	var out string
	sto := as.mt.Storage()
	sto.Iterate(func(key, value []byte) {
		out = out + "key: " + common3.BytesToHex(key) + ", value: " + common3.BytesToHex(value)
	})
	return out
}

func (as *ServiceImpl) ClaimsDump() string {
	var out string
	sto := as.mt.Storage()
	sto.Iterate(func(key, value []byte) {
		if value[0] == byte(1) { // TODO when the new merkletree version is ready, instead of byte(1) use the type indicator
			out = out + "key: " + common3.BytesToHex(key) + ", value: " + common3.BytesToHex(value)
		}
	})
	return out
}
