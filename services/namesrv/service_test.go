package namesrv

/*
import (
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3-core/cmd/id/config"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/web3"
	"github.com/iden3/go-iden3-core/utils"
	"github.com/stretchr/testify/assert"
)

const (
	testPrivKHex = "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"
	gethURL      = "https://ropsten.infura.io/TFnR8BWJlqZOKxHHZNcs"
)

var mt *merkletree.MerkleTree

func newTestingMerkle(numLevels int) (*merkletree.MerkleTree, error) {
	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}
	sto, err := merkletree.NewLevelDbStorage(dir)
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}
	mt, err := merkletree.New(sto, numLevels)
	return mt, err
}
func initializeEnvironment() error {
	// initialize
	config.MustRead("../../cmd/relay", "config")
	// MerkleTree leveldb
	var err error
	mt, err = newTestingMerkle(140)
	if err != nil {
		return err
	}
	// Ethereum
	err = web3srv.Open(gethURL, testPrivKHex)
	if err != nil {
		return err
	}
	return nil
}
func TestVinculateID(t *testing.T) {
	initializeEnvironment()
	testPrivK, err := crypto.HexToECDSA(config.C.Server.PrivK)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	var vinculateIDMsg VinculateIDMsg
	vinculateIDMsg.Msg.Name = "johndoe"
	vinculateIDMsg.Msg.RawIdentityTx.KSignOperational_p = "0xKSign_p"
	vinculateIDMsg.Msg.RawIdentityTx.KRecovery_p = "0xKRecovery_p"
	vinculateIDMsg.Msg.RawIdentityTx.KRevocation_p = "0xKRevocation_p"
	vinculateIDMsg.Msg.EthAddr = testAddr.Hex()
	msgHash := vinculateIDMsg.MsgHash()

	sig, err := utils.Sign(msgHash, testPrivK)
	assert.Nil(t, err)
	vinculateIDMsg.MsgSignature = common3.HexEncode(sig)
	assert.Equal(t, "0x87da22a43b63bda4fa77f65f966677161a1e5dc6af65f71eb84b628d485c881d", vinculateIDMsg.MsgHash().Hex())

	assignNameClaim, err := VinculateID(mt, vinculateIDMsg, config.C.ContractsAddress.Identities, testPrivK)
	assert.Nil(t, err)
	assert.Equal(t, "0xaa5a987f33c986616fe49c291448ddde00ef80b555fdeed35fe4d87b3f7445e7", assignNameClaim.Ht().Hex())
}
*/
