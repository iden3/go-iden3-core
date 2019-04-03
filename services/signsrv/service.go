package signsrv

import (
	"crypto/ecdsa"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	SignEthMsg(msg []byte) (*utils.SignatureEthMsg, error)
	SignEthMsgDate(msg []byte) (*utils.SignatureEthMsg, int64, error)
	PublicKey() *ecdsa.PublicKey
}

type ServiceImpl struct {
	ks  *keystore.KeyStore
	acc accounts.Account
	pk  *ecdsa.PublicKey
}

// New creates a new signsrv service.
func New(ks *keystore.KeyStore, acc accounts.Account) (*ServiceImpl, error) {
	pk, err := utils.GetPkFromKeyStore(ks, acc.Address)
	if err != nil {
		return nil, err
	}
	return &ServiceImpl{ks, acc, pk}, nil
}

func (s *ServiceImpl) PublicKey() *ecdsa.PublicKey {
	return s.pk
}

func (s *ServiceImpl) SignEthMsg(msg []byte) (*utils.SignatureEthMsg, error) {
	return utils.SignEthMsg(s.ks, s.acc, msg)
}

func (s *ServiceImpl) SignEthMsgDate(msg []byte) (*utils.SignatureEthMsg, int64, error) {
	dateInt64 := time.Now().Unix()
	dateBytes := utils.Uint64ToEthBytes(uint64(dateInt64))
	sig, err := s.SignEthMsg(append(msg, dateBytes...))
	return sig, dateInt64, err
}
