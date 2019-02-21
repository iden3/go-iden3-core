package signsrv

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	SignHash(h utils.Hash) (*utils.SignatureEthMsg, error)
}

type ServiceImpl struct {
	ks  *keystore.KeyStore
	acc accounts.Account
}

// New creates a new signsrv service.
func New(ks *keystore.KeyStore, acc accounts.Account) *ServiceImpl {
	return &ServiceImpl{ks, acc}
}

// SignHash signs a hash.
func (s *ServiceImpl) SignHash(h utils.Hash) (*utils.SignatureEthMsg, error) {
	sigBytes, err := s.ks.SignHash(s.acc, h[:])
	if err != nil {
		return nil, err
	}
	sig := &utils.SignatureEthMsg{}
	copy(sig[:], sigBytes)
	return sig, nil
}

// SignBytes signs a byte array.
func SignBytes(s Service, data []byte) (*utils.SignatureEthMsg, error) {
	h := utils.EthHash(data[:])
	return s.SignHash(h)
}

// SignBytesDate signs a byte array appended by the current time and returns
// the signature and the time used in the signature.
func SignBytesDate(s Service, data []byte) (*utils.SignatureEthMsg, int64, error) {
	dateUint64 := uint64(time.Now().Unix())
	dateBytes := utils.Uint64ToEthBytes(dateUint64)
	h := utils.EthHash(append(data[:], dateBytes...))
	sig, err := s.SignHash(h)
	if err != nil {
		return nil, 0, err
	}
	return sig, int64(dateUint64), nil
}
