package signsrv

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/iden3/go-iden3/merkletree"
)

type Service interface {
	SignHash(h merkletree.Hash) ([]byte, error)
}

type ServiceImpl struct {
	ks  *keystore.KeyStore
	acc accounts.Account
}

func New(ks *keystore.KeyStore, acc accounts.Account) *ServiceImpl {
	return &ServiceImpl{ks, acc}
}

func (s *ServiceImpl) SignHash(h merkletree.Hash) ([]byte, error) {
	return s.ks.SignHash(s.acc, h[:])
}
