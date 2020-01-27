package idensigner

import (
	// "crypto/ecdsa"
	// "encoding/hex"
	// "fmt"
	"time"

	"github.com/iden3/go-iden3-core/common"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

type IdenSigner struct {
	ks     *babykeystore.KeyStore
	pk     babyjub.PublicKey
	pkComp babyjub.PublicKeyComp
}

// New creates a new IdenSigner service.
func New(ks *babykeystore.KeyStore, pk babyjub.PublicKey) *IdenSigner {
	return &IdenSigner{ks, pk, pk.Compress()}
}

func (s *IdenSigner) PublicKey() *babyjub.PublicKey {
	return &s.pk
}

func (s *IdenSigner) SignEthMsg(msg []byte) (*babyjub.SignatureComp, error) {
	// sig, err := s.ks.Sign(&s.pkComp, msg)
	// fmt.Println("publicKey", s.pk.String())
	// fmt.Println("signature", sig.String())
	// fmt.Println("msg", hex.EncodeToString(msg))
	// return sig, err
	return s.ks.SignRaw(&s.pkComp, msg)
}

func (s *IdenSigner) SignEthMsgDate(msg []byte) (*babyjub.SignatureComp, int64, error) {
	dateInt64 := time.Now().Unix()
	dateBytes := common.Uint64ToEthBytes(uint64(dateInt64))
	sig, err := s.SignEthMsg(append(msg, dateBytes...))
	return sig, dateInt64, err
}
