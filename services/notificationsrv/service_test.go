// Run integration tests with:
// TEST=int go test -v -count=1 ./... -run=TestInt

package notificationsrv

import (
	// "crypto/ecdsa"
	"encoding/json"
	// "fmt"
	// "io/ioutil"
	"encoding/hex"
	"os"
	"testing"

	"github.com/iden3/go-iden3/crypto/babyjub"
	babykeystore "github.com/iden3/go-iden3/keystore"
	// "github.com/ethereum/go-ethereum/accounts"
	// "github.com/ethereum/go-ethereum/accounts/keystore"
	// "github.com/ethereum/go-ethereum/crypto"
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
    "proofs":
    [
        {
            "mtp0": "0x000100000000000000000000000000000000000000000000000000000000000125024058dff8730e7c283b2eb8b1553f32b5db48b2dc3499f1f610591b7cb5ab",
            "mtp1": "0x0302000000000000000000000000000000000000000000000000000000000003286bbd1d59ecc50d86dbb5ee59e2997d3522d378b0eb70a86fa38e99bc48179d1e7604b4b32e21da52f5f8a0ccf9709e378e033a9c1d458c4d426d57e53f629b2ca6f7a21d09938e1b52786f8b525b19832a84bb59c8ba4de6871728854f60af29af7742f31e4dfe967485d2e10d4f040d3f53236587b7de64717b871e661f84",
            "root": "0x14a946742e18446a877932c0938511bb6df3c77329ccd9c9cab5981212ffff17",
            "aux":
            {
                "version": 0,
                "era": 0,
                "id": "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"
            }
        },
        {
            "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "mtp1": "0x030000000000000000000000000000000000000000000000000000000000000020b468baa588865efc5df741e0a48569aa1171143a8627f425fff0d4fa7803c7022a1e2c3a59747c79b0cddee114e3bfb2d24777281ed568b364d43a6eea33a8",
            "root": "0x066996bfceb028398017cf44ef9e6aab2a13412b7dc9ee00d90a305cb97ae78e",
            "aux": null
        }
    ],
    "leaf": "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002906dcb03d2b068326665e02759eff24d35d40522d9e6efd8e29fb299f67bb1c0000000000000000000000000000000000000001000000000000000000000001",
    "date": 1559555861,
    "signature": "64b9e8f3a0a354b069ce0398ea8b914281bc3c1ebd150dfa7d0cf6db3bcafca914ee304bf2fa99a2f880659f790a1b19a434e68f57baaf181f7a64eb53bd1f02",
    "signer": "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2fps"
}
`

const namesFileContent = `
{
  "iden3.eth": "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2fps"
}
`

const entititesFileContent = `
{
  "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2fps": {
    "name": "iden3-test-relay",
    "kOpPub": "117f0a278b32db7380b078cdb451b509a2ed591664d1bac464e8c35a90646796",
    "trusted": { "relay": true }
  }
}
`

const urlNotificationService = "http://127.0.0.1:10000/api/unstable"

const passphrase = "secret"

const senderSkHex = "9b3260823e7b07dd26ef357ccfed23c10bcef1c85940baa3d02bbf29461bbbbe"

// const idHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
const idB58 = "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"

// const sendIdHex = "0xdcde41e52633bcf03c68248b54fc48875acc978f"
const sendIdB58 = "119PcNGhjJj37xjNSgqn4rmdXeuyPPfycmFqwkWT8"

// const keySignSkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
// const keySignSkHex = "0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4"

var proofKSign core.ProofClaim
var id core.ID
var keyStore *babykeystore.KeyStore
var senderSk babyjub.PrivateKey
var senderPkComp *babyjub.PublicKeyComp
var senderPk *babyjub.PublicKey
var sendId core.ID

var service *Service

func setup() {
	var err error
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	if id, err = core.IDFromString(idB58); err != nil {
		panic(err)
	}

	if sendId, err = core.IDFromString(sendIdB58); err != nil {
		panic(err)
	}

	pass := []byte("my passphrase")
	storage := babykeystore.MemStorage([]byte{})
	keyStore, err := babykeystore.NewKeyStore(&storage, babykeystore.LightKeyStoreParams)
	if err != nil {
		panic(err)
	}

	if _, err := hex.Decode(senderSk[:], []byte(senderSkHex)); err != nil {
		panic(err)
	}
	if senderPkComp, err = keyStore.ImportKey(senderSk, pass); err != nil {
		panic(err)
	}
	if err := keyStore.UnlockKey(senderPkComp, pass); err != nil {
		panic(err)
	}
	if senderPk, err = senderPkComp.Decompress(); err != nil {
		panic(err)
	}

	signSrv := signsrv.New(keyStore, *senderPk)

	signedPacketSigner := signedpacketsrv.NewSignedPacketSigner(*signSrv, proofKSign, id)
	service = New(urlNotificationService, signedPacketSigner)
}

func teardown() {
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
