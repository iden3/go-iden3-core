package endpoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claim"
	"github.com/iden3/go-iden3/services/web3"
	"github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
)

var testPrivHex = "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"

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
func performRequest(r http.Handler, method, path, data string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
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
	err = web3srv.Open(config.C.Geth.URL, config.C.Server.PrivK)
	if err != nil {
		return err
	}
	return nil
}

func TestHandlePostAssignNameClaim(t *testing.T) {
	err := initializeEnvironment()
	if err != nil {
		t.Errorf(err.Error())
	}
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r.POST("/claim/:idaddr", handlePostClaim)

	privK, err := crypto.HexToECDSA(testPrivHex)
	assert.Nil(t, err)
	ethID := crypto.PubkeyToAddress(privK.PublicKey)
	namespaceHash := merkletree.HashBytes([]byte(config.C.Namespace))
	nameHash := merkletree.HashBytes([]byte("johndoe"))
	domainHash := merkletree.HashBytes([]byte(config.C.Domain))
	assignNameClaim := core.NewAssignNameClaim(namespaceHash, nameHash, domainHash, ethID)
	signature, err := utils.Sign(assignNameClaim.Ht(), privK)
	assert.Nil(t, err)
	signatureHex := common3.BytesToHex(signature)
	claimValueMsg := claimsrv.BytesSignedMsg{
		common3.BytesToHex(assignNameClaim.Bytes()),
		signatureHex,
	}
	json, err := json.Marshal(claimValueMsg)
	if err != nil {
		t.Errorf(err.Error())
	}
	w := performRequest(r, "POST", "/claim/"+ethID.Hex(), string(json))
	assert.Equal(t, http.StatusOK, w.Code)
	buf := make([]byte, 1024)

	num, _ := w.Body.Read(buf)
	reqBody := string(buf[0:num])
	fmt.Println(reqBody)
}

func TestHandlePostAuthorizeKSignClaim(t *testing.T) {
	err := initializeEnvironment()
	if err != nil {
		t.Errorf(err.Error())
	}
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r.POST("/claim/:idaddr", handlePostClaim)

	privK, _ := crypto.HexToECDSA(testPrivHex)
	ethID := crypto.PubkeyToAddress(privK.PublicKey)
	authorizeKSignClaim := core.NewAuthorizeKSignClaim("iden3.io", ethID, "app1", "appauthz", 1535208350, 1535208350)

	signature, err := utils.Sign(authorizeKSignClaim.Ht(), privK)
	assert.Nil(t, err)
	signatureHex := common3.BytesToHex(signature)
	claimValueMsg := claimsrv.BytesSignedMsg{
		common3.BytesToHex(authorizeKSignClaim.Bytes()),
		signatureHex,
	}
	json, err := json.Marshal(claimValueMsg)
	if err != nil {
		t.Errorf(err.Error())
	}
	w := performRequest(r, "POST", "/claim/"+ethID.Hex(), string(json))
	assert.Equal(t, http.StatusOK, w.Code)
	// buf := make([]byte, 1024)
	// num, _ := w.Body.Read(buf)
	// reqBody := string(buf[0:num])
	// fmt.Println(reqBody)
}

func TestGetRoot(t *testing.T) {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r.GET("/root", handleGetRoot)

	w := performRequest(r, "GET", "/root", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetIDRoot(t *testing.T) {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r.GET("/claim/:idaddr/root", handleGetIDRoot)
	privK, _ := crypto.HexToECDSA(testPrivHex)
	ethID := crypto.PubkeyToAddress(privK.PublicKey)
	w := performRequest(r, "GET", "/claim/"+ethID.Hex()+"/root", "")
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetClaimByHiThatDontExist(t *testing.T) {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r.GET("/claim/:idaddr/hi/:hi", handleGetClaimByHi)
	ethIDHex := "0x970E8128AB834E8EAC17Ab8E3812F010678CF791"
	hiHex := "0x784adb4a490b9c0521c11298f384bf847881711f1a522a40129d76e3cfc68c9a"
	w := performRequest(r, "GET", "/claim/"+ethIDHex+"/hi/"+hiHex, "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestAddClaimAndGetClaimByHi(t *testing.T) {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r.GET("/claim/:idaddr/hi/:hi", handleGetClaimByHi)

	privK, _ := crypto.HexToECDSA(testPrivHex)
	ethID := crypto.PubkeyToAddress(privK.PublicKey)
	claim := core.NewClaimDefault("namespace.io", "default", []byte("dataasdf"))
	signature, err := utils.Sign(claim.Ht(), privK)
	assert.Nil(t, err)
	signatureHex := common3.BytesToHex(signature)
	claimValueMsg := claimsrv.ClaimValueMsg{
		claim,
		signatureHex,
	}
	_, _, err = claimsrv.AddUserIDClaim(mt, "namespace.io", ethID, claimValueMsg, config.C.ContractsAddress.Identities)
	assert.Nil(t, err)
	hi := claim.Hi()
	w := performRequest(r, "GET", "/claim/"+ethID.Hex()+"/hi/"+hi.Hex(), "")
	assert.Equal(t, http.StatusOK, w.Code)
	// buf := make([]byte, 1024)
	// num, _ := w.Body.Read(buf)
	// reqBody := string(buf[0:num])
	// fmt.Println(reqBody)
}
