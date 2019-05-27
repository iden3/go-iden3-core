// Run integration tests with:
// TEST=int go test -v -count=1 ./... -run=TestInt

package notificationsrv

import (
	"crypto/ecdsa"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	// common3 "github.com/iden3/go-iden3/common"
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
      "mtp0": "0x00020000000000000000000000000000000000000000000000000000000000022910a6fba42851f8282e0266c887e09db4fd84975a76a6c6ce468651683d2346",
      "mtp1": "0x010100000000000000000000000000000000000000000000000000000000000125024058dff8730e7c283b2eb8b1553f32b5db48b2dc3499f1f610591b7cb5ab",
      "root": "0x2ad101bc0d0e1b1efa9e74d03f017f531016e1b77c7cd5f514c864e8f4f22f90",
      "aux": {
        "version": 0,
        "era": 0,
        "id": "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"
      }
    },
    {
      "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "mtp1": "0x030000000000000000000000000000000000000000000000000000000000000020b468baa588865efc5df741e0a48569aa1171143a8627f425fff0d4fa7803c701bfeaf3af8775cbd1884bde8bec9762d167dcd4c77b3eafa13f938364b89772",
      "root": "0x1b7a0d2cdea1bd692f8fae6fafa774a4dc8fe28be8f81e464d259a603079e4c5",
      "aux": null
    }
  ],
  "leaf": "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003c2e48632c87932663beff7a1f6deb692cc61b041262ae8f310203d0f5ff50000000000000000000000000000000000007833000000000000000000000004",
  "date": 1558091875,
  "signature": "0x99314ccf4e79472f55019ce348c7f367bc1a8a508bd43c972aa586d0f8bf198c53aea3b40a86bf9a3fba5c8cdadc3bb01c94b3432f9ccce4493ad01ef46dffbb1c",
  "signer": "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoWij"
}
`

const namesFileContent = `
{
  "iden3.io": "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2fps"
}
`

const entititesFileContent = `
{
  "11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoWij": {
    "name": "iden3-test-relay2",
    "kOpAddr": "0x7633bc9012f924100fae50d6dda7162b0bba720d",
    "kOpPub": "0x036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd4d59",
    "trusted": { "relay": true }
  },
  "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2fps": {
    "name": "iden3-test-relay3",
    "kOpAddr": "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c",
    "kOpPub": "0x036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd4d59",
    "trusted": { "relay": true }
  }
}
`

const urlNotificationService = "http://127.0.0.1:10000/api/unstable"

const passphrase = "secret"

// const relaySkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"

// const idHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
const idB58 = "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"

// const sendIdHex = "0xdcde41e52633bcf03c68248b54fc48875acc978f"
const sendIdB58 = "1pquYVpccuB491VyD3rEwhqJXUiKGJonbdxcWorpz"

// const keySignSkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
const keySignSkHex = "0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4"

var proofKSign core.ProofClaim
var keyStoreDir string
var keyStore *keystore.KeyStore
var id core.ID
var sendId core.ID

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
	if id, err = core.IDFromString(idB58); err != nil {
		panic(err)
	}

	if sendId, err = core.IDFromString(sendIdB58); err != nil {
		panic(err)
	}
	signer, err := signsrv.New(keyStore, account)
	if err != nil {
		panic(err)
	}
	signedPacketSigner := signedpacketsrv.NewSignedPacketSigner(signer, proofKSign, id)
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
	err := service.SendNotification(notification, sendId)
	assert.Nil(t, err)
	// Send notification with text
	notification = NewMsgTxt("notificationText")
	err = service.SendNotification(notification, sendId)
	assert.Nil(t, err)
}
