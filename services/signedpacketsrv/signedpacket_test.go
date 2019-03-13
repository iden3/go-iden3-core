package signedpacketsrv

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
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/discoverysrv"
	"github.com/iden3/go-iden3/services/nameresolvesrv"
	// "github.com/iden3/go-iden3/utils"
)

const debug = false

const passphrase = "secret"

const relaySkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
const relayIdAddrHex = "0x0123456789abcdef0123456789abcdef01234567"

//const kSignSkHex = "7517685f1693593d3263460200ed903370c2318e8ba4b9bb5727acae55c32b3d"
const kSignSkHex = "0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4"

//const idAddrHex = "0x970e8128ab834e8eac17ab8e3812f010678cf791"
const idAddrHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"

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

const namesFileContent = `
{
  "iden3.io": "0x0123456789abcdef0123456789abcdef01234567"
}
`

const identitiesFileContent = `
{
  "0x0123456789abcdef0123456789abcdef01234567": {
    "name": "iden3-test-relay",
    "kOpAddr": "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c",
    "kOpPub": "0x036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd4d59",
    "trusted": { "relay": true }
  }
}
`

var namesFilePath string
var identitiesFilePath string

var signedPacketSrv *Service

var dbDir string
var keyStoreDir string
var keyStore *keystore.KeyStore
var relaySk *ecdsa.PrivateKey
var relayPk *ecdsa.PublicKey
var kSignSk *ecdsa.PrivateKey
var kSignPk *ecdsa.PublicKey

var relayIdAddr common.Address

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
	// _, err = keyStore.ImportECDSA(relaySk, passphrase)
	// if err != nil {
	// 	panic(err)
	// }
	// keyStore.Unlock(accounts.Account{Address: crypto.PubkeyToAddress(*relayPk)}, passphrase)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("relayPk:", common3.HexEncode(crypto.CompressPubkey(relayPk)))

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

	//relayAddr = crypto.PubkeyToAddress(*relayPk)

	namesFile, err := ioutil.TempFile("", "go-iden3-test-namesFile")
	if err != nil {
		panic(err)
	}
	namesFile.WriteString(namesFileContent)
	namesFilePath = namesFile.Name()
	namesFile.Close()

	nameResolveSrv, err := nameresolvesrv.New(namesFilePath)
	if err != nil {
		panic(err)
	}
	// common3.HexDecodeInto(relayIdAddr[:], []byte(relayIdAddrHex))

	// nameResolveSrv := nameresolvesrv.Service{
	// 	Names: map[string]common.Address{"iden3.io": relayIdAddr},
	// }

	identitiesFile, err := ioutil.TempFile("", "go-iden3-test-identitiesFile")
	if err != nil {
		panic(err)
	}
	identitiesFile.WriteString(identitiesFileContent)
	identitiesFilePath = identitiesFile.Name()
	identitiesFile.Close()

	discoverySrv, err := discoverysrv.New(identitiesFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("DBG: Checkpoint A")

	// var iden3TestOpAddr common.Address
	// common3.HexDecodeInto(iden3TestOpAddr[:], []byte("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c"))
	// discoverySrv := discoverysrv.Service{
	// 	Identities: map[common.Address]*discoverysrv.Identity{
	// 		relayIdAddr: &discoverysrv.Identity{
	// 			Name:            "iden3-test-relay",
	// 			OperationalPk:   &utils.PublicKey{PublicKey: *relayPk},
	// 			OperationalAddr: iden3TestOpAddr,
	// 			Trusted:         discoverysrv.Trusted{Relay: true},
	// 		},
	// 	},
	// }

	signedPacketSrv = New(discoverySrv, nameResolveSrv)
}

func teardown() {
	os.RemoveAll(keyStoreDir)
	os.RemoveAll(dbDir)
	os.Remove(namesFilePath)
	os.Remove(identitiesFilePath)
}

func TestSignedPacket(t *testing.T) {
	setup()
	defer teardown()

	t.Run("SignPacketV01", testSignPacketV01)
	t.Run("SignGenericSigV01", testSignGenericSigV01)
	t.Run("SignIdenAssertV01", testSignIdenAssertV01)
	t.Run("MarshalUnmarshal", testMarshalUnmarshal)

}

func BenchmarkSignedPacket(b *testing.B) {
	setup()
	defer teardown()

	b.Run("SignGenericSigV01", benchmarkSignGenericSigV01)
	b.Run("VerifySignGenericSigV01", benchmarkVerifySignGenericSigV01)
}

func testSignPacketV01(t *testing.T) {
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	if debug {
		fmt.Println(&proofKSign)
	}
	form := map[string]string{"foo": "baz"}
	signedPacket, err := NewSignPacketV01(keyStore, idAddr, kSignPk, proofKSign, 600,
		GENERICSIGV01, nil, form)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	err = signedPacketSrv.VerifySignedPacket(signedPacket)
	assert.Nil(t, err)
}

func testSignGenericSigV01(t *testing.T) {
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	form := map[string]string{"foo": "baz"}
	signedPacket, err := NewSignGenericSigV01(keyStore, idAddr, kSignPk, proofKSign, 600, form)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	err = signedPacketSrv.VerifySignedPacketGeneric(signedPacket)
	assert.Nil(t, err)
}

func benchmarkSignGenericSigV01(b *testing.B) {
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	form := map[string]string{"foo": "baz"}

	for n := 0; n < b.N; n++ {
		NewSignGenericSigV01(keyStore, idAddr, kSignPk, proofKSign, 600, form)
	}
}

// VerifySignedPacketGeneric is a bit slow right now.  The bottleneck resides
// in the VerifyProofClaim, in particular in the calculation of the mimc7.Hash.
// There is room for optimization in the mimc7.Hash
func benchmarkVerifySignGenericSigV01(b *testing.B) {
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	form := map[string]string{"foo": "baz"}
	signedPacket, err := NewSignGenericSigV01(keyStore, idAddr, kSignPk, proofKSign, 600, form)
	assert.Nil(b, err)

	for n := 0; n < b.N; n++ {
		signedPacketSrv.VerifySignedPacketGeneric(signedPacket)
	}
}

func testSignIdenAssertV01(t *testing.T) {
	// Login Server
	nonceDb := core.NewNonceDb()
	requestIdenAssert := NewRequestIdenAssert(nonceDb, "example.com", 60)

	// Client
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	var proofAssignName core.ProofClaim
	if err := json.Unmarshal([]byte(proofAssignNameJSON), &proofAssignName); err != nil {
		panic(err)
	}
	signedPacket, err := NewSignIdenAssertV01(requestIdenAssert, ethName, &proofAssignName,
		keyStore, idAddr, kSignPk, proofKSign, 600)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	// Login Server
	idenAssertResult, err := signedPacketSrv.VerifySignedPacketIdenAssert(signedPacket, nonceDb, "example.com")
	assert.Nil(t, err)
	if debug {
		fmt.Println(idenAssertResult)
	}
}

func testMarshalUnmarshal(t *testing.T) {
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}

	form := map[string]string{"foo": "baz", "bar": "biz"}
	signedPacket, err := NewSignPacketV01(keyStore, idAddr, kSignPk, proofKSign, 600,
		GENERICSIGV01, nil, form)
	assert.Nil(t, err)
	if debug {
		fmt.Println("\nSignedPacket:")
		fmt.Printf("Data: %#v\n", signedPacket.Payload.Data)
		fmt.Printf("Form: %#v\n", signedPacket.Payload.Form)
	}
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	//if debug {
	//	fmt.Println("\nMarshal:")
	//	fmt.Println(signedPacketStr)
	//}
	var signedPacket2 SignedPacket
	err = signedPacket2.Unmarshal(signedPacketStr)
	assert.Nil(t, err)
	if debug {
		fmt.Println("\nUnmarshal:")
		fmt.Printf("Data: %#v\n", signedPacket2.Payload.Data)
		fmt.Printf("Form: %#v\n", signedPacket2.Payload.Form)
	}

	signedPacket3, err := NewSignPacketV01(keyStore, idAddr, kSignPk, proofKSign, 600,
		"invalid", nil, nil)
	assert.Nil(t, err)
	signedPacketStr2, err := signedPacket3.Marshal()
	assert.Nil(t, err)
	var signedPacket4 SignedPacket
	// "invalid" is not a valid signed packet type, so unmarshal must error
	err = signedPacket4.Unmarshal(signedPacketStr2)
	assert.Error(t, err)
}
