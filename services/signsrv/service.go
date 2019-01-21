package signsrv

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	SignHash(h utils.Hash) ([]byte, error)
}

type ServiceImpl struct {
	ks  *keystore.KeyStore
	acc accounts.Account
}

func New(ks *keystore.KeyStore, acc accounts.Account) *ServiceImpl {
	return &ServiceImpl{ks, acc}
}

func (s *ServiceImpl) SignHash(h utils.Hash) ([]byte, error) {
	return s.ks.SignHash(s.acc, h[:])
}
