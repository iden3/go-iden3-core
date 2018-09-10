package endpoint

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/id/config"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/name"
	"github.com/iden3/go-iden3/services/web3"
	"github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
)

var testPrivHex = "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"

func performRequest(r http.Handler, method, path, data string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

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
	config.MustRead("..", "config")
	// MerkleTree leveldb
	var err error
	mt, err = newTestingMerkle(140)
	if err != nil {
		return err
	}

	// Ethereum
	err = web3srv.Open(config.C.Geth.URL, config.C.Relay.PrivK)
	if err != nil {
		return err
	}
	return nil
}

func TestVinculateID(t *testing.T) {
	err := initializeEnvironment()
	if err != nil {
		t.Errorf(err.Error())
	}
	r := gin.Default()
	r.POST("/vinculateid", handleVinculateID)

	testPrivK, err := crypto.HexToECDSA(testPrivHex)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	var vinculateIDMsg namesrv.VinculateIDMsg
	vinculateIDMsg.Msg.Name = "johndoe"
	vinculateIDMsg.Msg.RawIdentityTx.KSignOperational_p = "0xKSign_p"
	vinculateIDMsg.Msg.RawIdentityTx.KRecovery_p = "0xKRecovery_p"
	vinculateIDMsg.Msg.RawIdentityTx.KRevocation_p = "0xKRevocation_p"
	vinculateIDMsg.Msg.EthID = testAddr.Hex()
	msgHash := vinculateIDMsg.MsgHash()
	sig, err := utils.Sign(msgHash, testPrivK)
	assert.Nil(t, err)
	vinculateIDMsg.MsgSignature = hexutil.Encode(sig)
	json, err := json.Marshal(vinculateIDMsg)
	if err != nil {
		t.Errorf(err.Error())
	}
	w := performRequest(r, "POST", "/vinculateid", string(json))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleClaimNameResolv(t *testing.T) {
	r := gin.Default()
	r.GET("/identities/resolv/:nameid", handleAssignNameClaimResolv)

	w := performRequest(r, "GET", "/identities/resolv/johndoe@iden3.io", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	// err := json.Unmarshal([]byte(w.Body.String()), &response)
	json.Unmarshal([]byte(w.Body.String()), &response)
	ethIDHex, exists := response["ethID"]
	assert.True(t, exists)
	ethID := common.HexToAddress(ethIDHex)
	testPrivK, err := crypto.HexToECDSA(testPrivHex)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)
	if !bytes.Equal(testAddr.Bytes(), ethID.Bytes()) {
		t.Errorf("EthID not equal to the expected Address")
	}
}
