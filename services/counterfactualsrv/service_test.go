package counterfactualsrv

import (
	"crypto/ecdsa"
	"encoding/hex"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var debug = false

func TestPack(t *testing.T) {

	kb, _ := hex.DecodeString("3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000ee602447b5a75cf4f25367f5d199b860844d10c40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000259e9d8000000000967a7600")
	kr, _ := hex.DecodeString("0c7fbb73b49a62b75c44cc0b8559a67af866bcd942fa3bc1e7888d43e2f186f2")
	kep, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	knep, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000017a0ec823c79c6d1756a29edbf52eb228a69c5435ead519eb96cdb2412927b865")
	rb, _ := hex.DecodeString("3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000d79ae0a65e7dd29db1eac700368e693de09610b80c7fbb73b49a62b75c44cc0b8559a67af866bcd942fa3bc1e7888d43e2f186f2")
	rr, _ := hex.DecodeString("a392bc7458973721c1266b2ac65db038a87bb6ad2e822c2509298803e9941119")
	rep, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	rnep, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000016602097464f2c4a8f7854f1c29a7671a85d5aa670dbbe04a65f9d9c50a70626d")

	sigdate := int64(1539438904)
	signature, _ := hex.DecodeString("fc6ce7ce736ec3dc88c9f5f8d54d9d4d91dd7ee8f0b9d2beeed578b35104c7d6762ddf9c6b4a022aceba1ee09016b4edee8c697c5c677eaf0857eb290eee72441c")

	expected, _ := hex.DecodeString("00a43cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000ee602447b5a75cf4f25367f5d199b860844d10c40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000259e9d8000000000967a76000c7fbb73b49a62b75c44cc0b8559a67af866bcd942fa3bc1e7888d43e2f186f200200000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000017a0ec823c79c6d1756a29edbf52eb228a69c5435ead519eb96cdb2412927b86500743cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000d79ae0a65e7dd29db1eac700368e693de09610b80c7fbb73b49a62b75c44cc0b8559a67af866bcd942fa3bc1e7888d43e2f186f2a392bc7458973721c1266b2ac65db038a87bb6ad2e822c2509298803e994111900200000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000016602097464f2c4a8f7854f1c29a7671a85d5aa670dbbe04a65f9d9c50a70626d000000005bc1f9380041fc6ce7ce736ec3dc88c9f5f8d54d9d4d91dd7ee8f0b9d2beeed578b35104c7d6762ddf9c6b4a022aceba1ee09016b4edee8c697c5c677eaf0857eb290eee72441c")

	actual := packAuth(kb, kr, kep, knep, rb, rr, rep, rnep, sigdate, signature)

	assert.Equal(t, hex.EncodeToString(expected), hex.EncodeToString(actual))

}

var relaySecKey *ecdsa.PrivateKey
var relayPubKey *ecdsa.PublicKey
var relayKOpAddr common.Address
var relayIdAddr core.ID

type RootServiceMock struct {
	mock.Mock
}

func (m *RootServiceMock) Start() {

}

func (m *RootServiceMock) StopAndJoin() {

}

func (m *RootServiceMock) GetRoot(addr core.ID) (merkletree.Hash, error) {
	args := m.Called(addr)
	return args.Get(0).(merkletree.Hash), args.Error(1)
}

func (m *RootServiceMock) SetRoot(hash merkletree.Hash) {
	m.Called(hash)
	return
}

type SignServiceMock struct {
	mock.Mock
}

func (m *SignServiceMock) SignEthMsg(msg []byte) (*utils.SignatureEthMsg, error) {
	h := utils.EthHash(msg)
	sig, err := crypto.Sign(h[:], relaySecKey)
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	sigEthMsg := &utils.SignatureEthMsg{}
	copy(sigEthMsg[:], sig)
	return sigEthMsg, nil
}

func (m *SignServiceMock) SignEthMsgDate(msg []byte) (*utils.SignatureEthMsg, int64, error) {
	dateInt64 := time.Now().Unix()
	dateBytes := utils.Uint64ToEthBytes(uint64(dateInt64))
	sig, err := m.SignEthMsg(append(msg, dateBytes...))
	return sig, dateInt64, err
}

func (m *SignServiceMock) PublicKey() *ecdsa.PublicKey {
	return relayPubKey
}

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
