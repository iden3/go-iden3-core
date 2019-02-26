package core

import (
	"crypto/ecdsa"
	//"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	//"github.com/iden3/go-iden3/merkletree"
	common3 "github.com/iden3/go-iden3/common"
)

const passphrase = "secret"
const relaySkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
const kSignSkHex = "7517685f1693593d3263460200ed903370c2318e8ba4b9bb5727acae55c32b3d"
const idAddrHex = "0x970e8128ab834e8eac17ab8e3812f010678cf791"

const proofKSignJSON = `
{
  "proofs": [
    {
      "mtp0": "0x00010000000000000000000000000000000000000000000000000000000000012e2e10e151e3d54e45854ed2bc267783207a5319a99dc0517d6da32554813c37",
      "mtp1": "0x03010000000000000000000000000000000000000000000000000000000000012e2e10e151e3d54e45854ed2bc267783207a5319a99dc0517d6da32554813c3716a1806ea95b318e00e4304970d89cc3e38c743e073ae5ef8c3781afc2df03721541a6b5aa9bf7d9be3d5cb0bcc7cacbca26242016a0feebfc19c90f2224baed",
      "root": "0x1ddf0aaf15b0e59cd69dd85daa6b36d20c2ee094dc4c1937a598ab16d217be82",
      "aux": {
        "version": 0,
        "era": 0,
        "idAddr": "0x970e8128ab834e8eac17ab8e3812f010678cf791"
      }
    },
    {
      "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "mtp1": "0x03000000000000000000000000000000000000000000000000000000000000002a34ebf086363c2ed843b862794f1ff12fbebb508f2946c64a4ceddf4528016204422c79b3367ba451e340cf98ff204c5866dfa7766b889c8d3176f683e17d40",
      "root": "0x073abe9560e2c3c9cdadb0a62b5afd2be82dacd9f6fae045b9d0638444ea540b",
      "aux": null
    }
  ],
  "leaf": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034a377d7da59e11540b28034f08f10c35de698b99484b8efdafa2bd8ec476000000000000000000000000000000000000e231000000000000000000000004",
  "date": 1550672244,
  "signature": "0x549756d70737fb08120ecbe3d51a1bf7faea2033e831c218005283249d2a076312a98c4208e49ddc2856335633a13e2e0268f1022578ce700f18b1c42fa5ae661b"
}
`

var dbDir string
var keyStoreDir string
var keyStore *keystore.KeyStore
var relaySk *ecdsa.PrivateKey
var relayPk *ecdsa.PublicKey
var kSignSk *ecdsa.PrivateKey
var kSignPk *ecdsa.PublicKey

var idAddr common.Address

func genPrivateKey() {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(crypto.FromECDSA(privateKeyECDSA)))
}

func setup() {
	//genPrivateKey()
	var err error
	keyStoreDir, err = ioutil.TempDir("/tmp/", "go-iden3-test-keystore")
	if err != nil {
		panic(err)
	}
	keyStore = keystore.NewKeyStore(keyStoreDir, 2, 1)
	relaySk, err = crypto.HexToECDSA(relaySkHex)
	if err != nil {
		panic(err)
	}
	relayPk = relaySk.Public().(*ecdsa.PublicKey)
	_, err = keyStore.ImportECDSA(relaySk, passphrase)
	if err != nil {
		panic(err)
	}
	keyStore.Unlock(accounts.Account{Address: crypto.PubkeyToAddress(*relayPk)}, passphrase)
	if err != nil {
		panic(err)
	}

	kSignSk, err = crypto.HexToECDSA(kSignSkHex)
	if err != nil {
		panic(err)
	}
	kSignPk = kSignSk.Public().(*ecdsa.PublicKey)
	_, err = keyStore.ImportECDSA(kSignSk, passphrase)
	if err != nil {
		panic(err)
	}
	err = keyStore.Unlock(accounts.Account{Address: crypto.PubkeyToAddress(*kSignPk)}, passphrase)
	if err != nil {
		panic(err)
	}

	common3.HexDecodeInto(idAddr[:], []byte(idAddrHex))

	relayAddr = crypto.PubkeyToAddress(*relayPk)
}

func teardown() {
	if err := os.RemoveAll(keyStoreDir); err != nil {
		panic(err)
	}
	if err := os.RemoveAll(dbDir); err != nil {
		panic(err)
	}
}

func TestSignedPacket(t *testing.T) {
	setup()
	defer teardown()

	t.Run("SignPacketV01", testSignPacketV01)

}

func testSignPacketV01(t *testing.T) {
	var proofKSign ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	if debug {
		fmt.Println(&proofKSign)
	}
	data := map[string]string{"foo": "bar"}
	form := map[string]string{"foo": "baz"}
	signedPacket, err := NewSignPacketV01(keyStore, idAddr, kSignPk, proofKSign, 600,
		"iden3.test", data, form)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	err = VerifySignedPacket(signedPacket)
	assert.Nil(t, err)
}
