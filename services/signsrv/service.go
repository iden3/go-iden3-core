package signsrv

import (
	"time"

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

func SignBytes(s Service, data []byte) ([]byte, error) {
	h := utils.EthHash(data[:])
	return s.SignHash(h)
}

func SignBytesDate(s Service, data []byte) ([]byte, uint64, error) {
	dateUint64 := uint64(time.Now().Unix())
	dateBytes := utils.Uint64ToEthBytes(dateUint64)
	h := utils.EthHash(append(data[:], dateBytes...))
	sig, err := s.SignHash(h)
	if err != nil {
		return nil, 0, err
	}
	return sig, dateUint64, nil
}
