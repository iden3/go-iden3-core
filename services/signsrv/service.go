package signsrv

import (
	// "crypto/ecdsa"
	// "encoding/hex"
	"time"

	"github.com/iden3/go-iden3/crypto/babyjub"
	babykeystore "github.com/iden3/go-iden3/keystore"
	"github.com/iden3/go-iden3/utils"
)

type Service struct {
	ks     *babykeystore.KeyStore
	pk     babyjub.PublicKey
	pkComp babyjub.PublicKeyComp
}

// New creates a new signsrv service.
func New(ks *babykeystore.KeyStore, pk babyjub.PublicKey) *Service {
	return &Service{ks, pk, pk.Compress()}
}

func (s *Service) PublicKey() *babyjub.PublicKey {
	return &s.pk
}

func (s *Service) SignEthMsg(msg []byte) (*babyjub.SignatureComp, error) {
	return s.ks.Sign(&s.pkComp, msg)
}

func (s *Service) SignEthMsgDate(msg []byte) (*babyjub.SignatureComp, int64, error) {
	dateInt64 := time.Now().Unix()
	dateBytes := utils.Uint64ToEthBytes(uint64(dateInt64))
	sig, err := s.SignEthMsg(append(msg, dateBytes...))
	return sig, dateInt64, err
}
