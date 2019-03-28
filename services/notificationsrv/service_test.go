// +build integration

package notificationsrv

import (
	"github.com/stretchr/testify/assert"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/signedpacketsrv"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
)

const urlNotificationService = "http://127.0.0.1:10000/api/unstable"

const passphrase = "secret"
const relaySkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
const idAddrHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
const keySignSkHex = "0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4"

const namesFileContent = `
{
  "iden3.io": "0x0123456789abcdef0123456789abcdef01234567"
}
`

const entititesFileContent = `
{
  "0x0123456789abcdef0123456789abcdef01234567": {
    "name": "iden3-test-relay",
    "kOpAddr": "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c",
    "kOpPub": "0x036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd4d59",
    "trusted": { "relay": true }
  }
}
`
const proofKSignJSON = `
{
    "proofs": [
      {
        "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "mtp1": "0x030000000000000000000000000000000000000000000000000000000000000028f8267fb21e8ce0cdd9888a6e532764eb8d52dd6c1e354157c78b7ea281ce801541a6b5aa9bf7d9be3d5cb0bcc7cacbca26242016a0feebfc19c90f2224baed",
        "root": "0x1d9d41171c4b621ff279e2acb84d8ab45612fef53e37225bdf67e8ad761c3922",
        "aux": {
          "version": 0,
          "era": 0,
          "idAddr": "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
        }
      },
      {
        "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "mtp1": "0x030000000000000000000000000000000000000000000000000000000000000021c9cceb8a61605050a029cecda7e36eeaffef11910778b3e5ea32f79659cee125451237d9133b0f5c1386b7b822f382cb14c5fff612a913956ef5436fb6208a",
        "root": "0x0b109c78f2679a405cc0b5a8d999129ee4429e86a986b8001e0a4df61c359690",
        "aux": null
      }
    ],
    "leaf": "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003c2e48632c87932663beff7a1f6deb692cc61b041262ae8f310203d0f5ff50000000000000000000000000000000000007833000000000000000000000004",
    "date": 1551434314,
    "signature": "0x8af874853697b3893c69ba992637676cacaa8fe84195d113f2d4a276e8002af265eb9bbc45766c1908fffc4d1cedab2f8c12a6ccced6f50147ed61b9f7c8d1b71b",
  "signer": "0x0123456789abcdef0123456789abcdef01234567"
}
`
const ethName = "testName@iden3.io"

const proofAssignNameJSON = `
{
  "proofs": [
    {
      "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "mtp1": "0x0300000000000000000000000000000000000000000000000000000000000000212b74b502c20a92dc8ce4e479fde418e0fa2d0da00fdffa239aed393efdccc9116f9ac9335ac96ec780b18bf4b82ab0b39e661e4818116226ffb2d2acb9ece2",
      "root": "0x0d1fadf720af488d10aaa4bdaf6a8d1163ad30b19624082b0e4403934ab57ff3",
      "aux": null
    }
  ],
  "leaf": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000308eff1357e7b5881c00ae22463b0f69a0d58adb0063ea0983a784a474e012c2ce392b45296419d9b57f91533c579a691db028f30000000000000000000000000000000000000000000000000000000000000003",
  "date": 1551367178,
  "signature": "0x3d350047f9f1464d4d57ff67e73a86689df4f358697a56ca68fc4ce68bf9a69a10ccd66c43b3c45757c19c7fcdc1d10c1d7e38eb66bb647571cdc56bc85f0b031c",
  "signer": "0x0123456789abcdef0123456789abcdef01234567"
}
`

var service *Service

var namesFilePath string
var entititesFilePath string

var proofKSign core.ProofClaim
var proofAssignName core.ProofClaim

var keyStoreDir string
var keyStore *keystore.KeyStore
var relaySk *ecdsa.PrivateKey
var relayPk *ecdsa.PublicKey
var keySignSk *ecdsa.PrivateKey
var keySignPk *ecdsa.PublicKey

func setup() {
	//genPrivateKey()
	var err error
	keyStoreDir, err = ioutil.TempDir("", "go-iden3-test-keystore")
	if err != nil {
		panic(err)
	}
	keyStore = keystore.NewKeyStore(keyStoreDir, 2, 1)
	relaySk, err = crypto.HexToECDSA(relaySkHex)
	if err != nil {
		panic(err)
	}
	relayPk = relaySk.Public().(*ecdsa.PublicKey)

	keySignSk, err = crypto.HexToECDSA(keySignSkHex)
	if err != nil {
		panic(err)
	}
	keySignPk = keySignSk.Public().(*ecdsa.PublicKey)
	_, err = keyStore.ImportECDSA(keySignSk, passphrase)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(proofAssignNameJSON), &proofAssignName); err != nil {
		panic(err)
	}
	err = keyStore.Unlock(accounts.Account{Address: crypto.PubkeyToAddress(*keySignPk)}, passphrase)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}

	var idAddr common.Address
	common3.HexDecodeInto(idAddr[:], []byte(idAddrHex))
	service = New(urlNotificationService, idAddr)
}

func teardown() {
	os.RemoveAll(keyStoreDir)
}

func TestNotificationService(t *testing.T) {
	setup()
	defer teardown()

	t.Run("Login", testLogin)
	t.Run("SendNotification", testSendNotification)
}

func testLogin(t *testing.T) {
	err := service.Login(keyStore, keySignPk, proofKSign, &signedpacketsrv.IdenAssertForm{EthName: ethName, ProofAssignName: &proofAssignName})
	if err != nil {
		fmt.Println(err)	
	}
	assert.Nil(t, err)
	assert.NotEqual(t, service.Token, "", "Token has not been received")
	fmt.Println("Login correctly to notification server")
}

func testSendNotification(t *testing.T) {
	notification := Notification{ Type: "notif.claim.v01", Data: proofKSign }
  // send notification with proofClaim
	err := service.SendNotification(&notification)
	if err != nil {
		fmt.Println(err)	
	}
	assert.Nil(t, err)
	// send notification with text
	notification = Notification{ Type: "notif.txt.v01", Data: "notificationText" }
	err = service.SendNotification(&notification)
	if err != nil {
		fmt.Println(err)	
	}
	assert.Nil(t, err)
}
