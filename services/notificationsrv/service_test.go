// Run integration tests with:
// TEST=int go test -v -count=1 ./... -run=TestInt

package notificationsrv

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/signedpacketsrv"
	"github.com/iden3/go-iden3/services/signsrv"
	"github.com/stretchr/testify/assert"
)

var integration bool

func init() {
	if os.Getenv("TEST") == "int" {
		integration = true
	}
}

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

const urlNotificationService = "http://127.0.0.1:10000/api/unstable"

const passphrase = "secret"
const relaySkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
const idAddrHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
const sendIdAddrHex = "0xdcde41e52633bcf03c68248b54fc48875acc978f"
const keySignSkHex = "0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4"

var proofKSign core.ProofClaim
var keyStoreDir string
var keyStore *keystore.KeyStore
var idAddr core.ID
var sendIdAddr core.ID

var service *Service

func setup() {
	var err error
	keyStoreDir, err = ioutil.TempDir("", "go-iden3-test-keystore")
	if err != nil {
		panic(err)
	}
	keyStore = keystore.NewKeyStore(keyStoreDir, 2, 1)

	keySignSk, err := crypto.HexToECDSA(keySignSkHex)
	if err != nil {
		panic(err)
	}
	keySignPk := keySignSk.Public().(*ecdsa.PublicKey)
	if _, err = keyStore.ImportECDSA(keySignSk, passphrase); err != nil {
		panic(err)
	}
	account := accounts.Account{Address: crypto.PubkeyToAddress(*keySignPk)}
	if err = keyStore.Unlock(account, passphrase); err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	common3.HexDecodeInto(idAddr[:], []byte(idAddrHex))
	common3.HexDecodeInto(sendIdAddr[:], []byte(sendIdAddrHex))
	signer, err := signsrv.New(keyStore, account)
	if err != nil {
		panic(err)
	}
	signedPacketSigner := signedpacketsrv.NewSignedPacketSigner(signer, proofKSign, idAddr)
	service = New(urlNotificationService, signedPacketSigner)
}

func teardown() {
	os.RemoveAll(keyStoreDir)
}

func TestIntNotificationService(t *testing.T) {
	if !integration {
		t.Skip()
	}
	setup()
	defer teardown()

	t.Run("SendNotification", testSendNotification)
}

func testSendNotification(t *testing.T) {
	// Send notification with proofClaim
	notification := NewMsgProofClaim(&proofKSign)
	err := service.SendNotification(notification, sendIdAddr)
	assert.Nil(t, err)
	// Send notification with text
	notification = NewMsgTxt("notificationText")
	err = service.SendNotification(notification, sendIdAddr)
	assert.Nil(t, err)
}
