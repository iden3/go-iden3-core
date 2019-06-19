package identitysrv

import (
	// "crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"testing"
	// "time"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	babykeystore "github.com/iden3/go-iden3/keystore"
	// "github.com/ethereum/go-ethereum/crypto"
	// common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/crypto/babyjub"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/signsrv"
	// "github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var debug = false

const relaySkHex = "4406831fa7bb87d8c92fc65f090a6017916bd2197ffca0e1e97933b14e8f5de5"

var keyStore *babykeystore.KeyStore
var relaySk babyjub.PrivateKey
var relayPkComp *babyjub.PublicKeyComp
var relayPk *babyjub.PublicKey
var relayId core.ID

type RootServiceMock struct {
	mock.Mock
}

func (m *RootServiceMock) Start() {

}

func (m *RootServiceMock) StopAndJoin() {

}

func (m *RootServiceMock) GetRoot(addr common.Address) (merkletree.Hash, error) {
	args := m.Called(addr)
	return args.Get(0).(merkletree.Hash), args.Error(1)
}

func (m *RootServiceMock) SetRoot(hash merkletree.Hash) {
	m.Called(hash)
	return
}

// type SignServiceMock struct {
// 	mock.Mock
// }
//
// func (m *SignServiceMock) SignEthMsg(msg []byte) (*utils.SignatureEthMsg, error) {
// 	h := utils.EthHash(msg)
// 	sig, err := crypto.Sign(h[:], relaySecKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	sig[64] += 27
// 	sigEthMsg := &utils.SignatureEthMsg{}
// 	copy(sigEthMsg[:], sig)
// 	return sigEthMsg, nil
// }
//
// func (m *SignServiceMock) SignEthMsgDate(msg []byte) (*utils.SignatureEthMsg, int64, error) {
// 	dateInt64 := time.Now().Unix()
// 	dateBytes := utils.Uint64ToEthBytes(uint64(dateInt64))
// 	sig, err := m.SignEthMsg(append(msg, dateBytes...))
// 	return sig, dateInt64, err
// }
//
// func (m *SignServiceMock) PublicKey() *ecdsa.PublicKey {
// 	return relayPubKey
// }

func newTestingMerkle(numLevels int) (*merkletree.MerkleTree, error) {
	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}
	sto, err := db.NewLevelDbStorage(dir, false)
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}

	mt, err := merkletree.NewMerkleTree(sto, numLevels)
	return mt, err
}
func initializeIdService(t *testing.T) *ServiceImpl {
	relayId, err := core.IDFromString("113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf")
	if err != nil {
		t.Error(err)
	}

	// MerkleTree leveldb
	mt, err := newTestingMerkle(140)
	if err != nil {
		t.Error(err)
	}
	// sto := db.NewMemoryStorage()
	rootservicemock := &RootServiceMock{}
	rootservicemock.On("SetRoot", mock.Anything).Return()

	pass := []byte("my passphrase")
	storage := babykeystore.MemStorage([]byte{})
	keyStore, err := babykeystore.NewKeyStore(&storage, babykeystore.LightKeyStoreParams)
	if err != nil {
		panic(err)
	}

	if _, err := hex.Decode(relaySk[:], []byte(relaySkHex)); err != nil {
		panic(err)
	}
	if relayPkComp, err = keyStore.ImportKey(relaySk, pass); err != nil {
		panic(err)
	}
	if err := keyStore.UnlockKey(relayPkComp, pass); err != nil {
		panic(err)
	}
	if relayPk, err = relayPkComp.Decompress(); err != nil {
		panic(err)
	}

	signSrv := signsrv.New(keyStore, *relayPk)

	claimService := claimsrv.New(relayId, mt, rootservicemock, *signSrv)
	idService := New(claimService)

	return idService
}

func TestCreateIdGenesisRandom(t *testing.T) {
	idsrv := initializeIdService(t)

	kOpSk := babyjub.NewRandPrivKey()
	kop := kOpSk.Public()
	if debug {
		fmt.Println("kop", kop)
	}
	kDis := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kReen := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kUpdateRoot := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	id, proofKOp, err := idsrv.CreateIdGenesis(kop, kDis, kReen, kUpdateRoot)
	assert.Nil(t, err)

	id2, _, err := core.CalculateIdGenesis(kop, kDis, kReen, kUpdateRoot)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)

	proofKOpVerified, err := core.VerifyProofClaim(relayPk, proofKOp)
	assert.Nil(t, err)
	assert.True(t, proofKOpVerified)
}

func TestCreateIdGenesisHardcoded(t *testing.T) {
	idsrv := initializeIdService(t)

	kopStr := "0x117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796"
	// krecStr := "0x03f9737be33b5829e3da80160464b2891277dae7d7c23609f9bb34bd4ede397bbf"
	// krevStr := "0x02d2da59d3022b4c1589e4910baa6cbaddd01f95ed198fdc3068d9dc1fb784a9a4"

	var kopComp babyjub.PublicKeyComp
	err := kopComp.UnmarshalText([]byte(kopStr))
	assert.Nil(t, err)
	kopPub, err := kopComp.Decompress()
	assert.Nil(t, err)
	kDis := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kReen := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	kUpdateRoot := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")

	id, proofKOp, err := idsrv.CreateIdGenesis(kopPub, kDis, kReen, kUpdateRoot)
	assert.Nil(t, err)
	if debug {
		fmt.Println("id", id)
		fmt.Println("id (hex)", id.String())
	}
	assert.Equal(t, "117aFcVWPyypFbjCuHRKaAaTV7nN3yT9q6PthJpm96", id.String())

	id2, _, err := core.CalculateIdGenesis(kopPub, kDis, kReen, kUpdateRoot)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)

	proofKOpVerified, err := core.VerifyProofClaim(relayPk, proofKOp)
	assert.Nil(t, err)
	assert.True(t, proofKOpVerified)
}
