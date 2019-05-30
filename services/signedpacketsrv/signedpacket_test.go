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

	// "github.com/ethereum/go-ethereum/accounts"
	//"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3/crypto/babyjub"
	babykeystore "github.com/iden3/go-iden3/keystore"
	"github.com/stretchr/testify/assert"
	//"github.com/iden3/go-iden3/merkletree"

	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/discoverysrv"
	"github.com/iden3/go-iden3/services/nameresolversrv"
	"github.com/iden3/go-iden3/services/signsrv"
	// "github.com/iden3/go-iden3/utils"
)

const debug = false

const passphrase = "secret"

const relaySkHex = "4be5471a938bdf3606888472878baace4a6a64e14a153adf9a1333969e4e573c"

const kSignSkHex = "9b3260823e7b07dd26ef357ccfed23c10bcef1c85940baa3d02bbf29461bbbbe"

//const idAddrHex = "0x970e8128ab834e8eac17ab8e3812f010678cf791"
// const idAddrHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
const idHex = "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"

// root 0x1d9d41171c4b621ff279e2acb84d8ab45612fef53e37225bdf67e8ad761c3922
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
    "date": 1559037174,
    "signature": "44060200ba32ea144ffe346242a6bfe530f3767cedcb039cca068fce14563300b358ca2eaacfac04a8332f21bcd74231334d7e7ef5b27fbab762e97e8f2de705",
    "signer": "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2fps"
}
`

const ethName = "testName@iden3.eth"

//o0x0d1fadf720af488d10aaa4bdaf6a8d1163ad30b19624082b0e4403934ab57ff3
const proofAssignNameJSON = `
{
    "proofs":
    [
        {
            "mtp0": "0x0001000000000000000000000000000000000000000000000000000000000001066996bfceb028398017cf44ef9e6aab2a13412b7dc9ee00d90a305cb97ae78e",
            "mtp1": "0x0301000000000000000000000000000000000000000000000000000000000001301a38afef4b2e600259190ab68ee7dddd194766d4ba90c93f49e8310a1b5cba20b468baa588865efc5df741e0a48569aa1171143a8627f425fff0d4fa7803c7022a1e2c3a59747c79b0cddee114e3bfb2d24777281ed568b364d43a6eea33a8",
            "root": "0x1c5d63fcb41321f5648ec038d852a345a4c08434896fac8dbcc2bde1d8541015",
            "aux": null
        }
    ],
    "leaf": "0x00000000000000000000000000000000000000000000000000000000000000000000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f400178118069763dbe18ad9c512b09b4f9a9b7ae14c4ead00200ceabdcbac85950000000000000000000000000000000000000000000000000000000000000003",
    "date": 1559037174,
    "signature": "6b9be136e74c5a0da8996e3456efab8bf61d0ce87ed726529919a5874b62928a826be4e3d0953ace275272112a502cb8da5e33e8351e81aec8d416cb95a00a01",
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

var namesFilePath string
var entititesFilePath string

var signedPacketVerifier *SignedPacketVerifier
var signedPacketSigner *SignedPacketSigner

var dbDir string

// var keyStoreDir string
var keyStore *babykeystore.KeyStore
var relaySk babyjub.PrivateKey
var relayPkComp *babyjub.PublicKeyComp
var relayPk *babyjub.PublicKey

var kSignSk babyjub.PrivateKey
var kSignPkComp *babyjub.PublicKeyComp
var kSignPk *babyjub.PublicKey

var proofKSign core.ProofClaim

var relayId core.ID

var id core.ID

func genPrivateKey() {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(crypto.FromECDSA(privateKeyECDSA)))
}

func setup() {
	//genPrivateKey()

	pass := []byte("my passphrase")
	storage := babykeystore.MemStorage([]byte{})
	keyStore, err := babykeystore.NewKeyStore(&storage, babykeystore.LightKeyStoreParams)
	if err != nil {
		panic(err)
	}

	if _, err := hex.Decode(relaySk[:], []byte(relaySkHex)); err != nil {
		panic(err)
	}
	if relayPkComp, err = keyStore.ImportKey(relaySk, pass); err != nil {
		panic(err)
	}
	if err := keyStore.UnlockKey(relayPkComp, pass); err != nil {
		panic(err)
	}
	if relayPk, err = relayPkComp.Decompress(); err != nil {
		panic(err)
	}

	if _, err := hex.Decode(kSignSk[:], []byte(kSignSkHex)); err != nil {
		panic(err)
	}
	if kSignPkComp, err = keyStore.ImportKey(kSignSk, pass); err != nil {
		panic(err)
	}
	if err := keyStore.UnlockKey(kSignPkComp, pass); err != nil {
		panic(err)
	}
	if kSignPk, err = kSignPkComp.Decompress(); err != nil {
		panic(err)
	}

	// common3.HexDecodeInto(idAddr[:], []byte(idAddrHex))
	id, err = core.IDFromString(idHex)
	if err != nil {
		panic(err)
	}

	namesFile, err := ioutil.TempFile("", "go-iden3-test-namesFile")
	if err != nil {
		panic(err)
	}
	namesFile.WriteString(namesFileContent)
	namesFilePath = namesFile.Name()
	namesFile.Close()

	nameResolverSrv, err := nameresolversrv.New(namesFilePath)
	if err != nil {
		panic(err)
	}

	entititesFile, err := ioutil.TempFile("", "go-iden3-test-entititesFile")
	if err != nil {
		panic(err)
	}
	entititesFile.WriteString(entititesFileContent)
	entititesFilePath = entititesFile.Name()
	entititesFile.Close()

	discoverySrv, err := discoverysrv.New(entititesFilePath)
	if err != nil {
		panic(err)
	}

	signSrv := signsrv.New(keyStore, *kSignPk)

	signedPacketVerifier = NewSignedPacketVerifier(discoverySrv, nameResolverSrv)

	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	signedPacketSigner = NewSignedPacketSigner(*signSrv, proofKSign, id)
}

func teardown() {
	os.RemoveAll(dbDir)
	os.Remove(namesFilePath)
	os.Remove(entititesFilePath)
}

func TestSignedPacket(t *testing.T) {
	setup()
	defer teardown()

	t.Run("SignPacketV01", testSignPacketV01)
	t.Run("SignGenericSigV01", testSignGenericSigV01)
	t.Run("SignIdenAssertV01Name", testSignIdenAssertV01Name)
	t.Run("SignIdenAssertV01NoName", testSignIdenAssertV01NoName)
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
	signedPacket, err := signedPacketSigner.
		NewSignPacketV02(600, GENERICSIGV01, nil, form)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	err = signedPacketVerifier.VerifySignedPacket(signedPacket)
	assert.Nil(t, err)
}

func testSignGenericSigV01(t *testing.T) {
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	form := map[string]string{"foo": "baz"}
	signedPacket, err := signedPacketSigner.
		NewSignGenericSigV01(600, form)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	err = signedPacketVerifier.VerifySignedPacketGeneric(signedPacket)
	assert.Nil(t, err)
}

func benchmarkSignGenericSigV01(b *testing.B) {
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	form := map[string]string{"foo": "baz"}

	for n := 0; n < b.N; n++ {
		signedPacketSigner.
			NewSignGenericSigV01(600, form)
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
	signedPacket, err := signedPacketSigner.
		NewSignGenericSigV01(600, form)
	assert.Nil(b, err)

	for n := 0; n < b.N; n++ {
		signedPacketVerifier.VerifySignedPacketGeneric(signedPacket)
	}
}

func testSignIdenAssertV01Name(t *testing.T) {
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
	signedPacket, err := signedPacketSigner.
		NewSignIdenAssertV01(requestIdenAssert,
			&IdenAssertForm{EthName: ethName, ProofAssignName: &proofAssignName}, 600)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	// Login Server
	idenAssertResult, err := signedPacketVerifier.
		VerifySignedPacketIdenAssert(signedPacket, nonceDb, "example.com")
	assert.Nil(t, err)
	if debug {
		fmt.Println(idenAssertResult)
	}
}

func testSignIdenAssertV01NoName(t *testing.T) {
	// Login Server
	nonceDb := core.NewNonceDb()
	requestIdenAssert := NewRequestIdenAssert(nonceDb, "example.com", 60)

	// Client
	var proofKSign core.ProofClaim
	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	signedPacket, err := signedPacketSigner.
		NewSignIdenAssertV01(requestIdenAssert, nil,
			600)
	assert.Nil(t, err)
	signedPacketStr, err := signedPacket.Marshal()
	assert.Nil(t, err)
	if debug {
		fmt.Println(signedPacketStr)
	}

	// Login Server
	idenAssertResult, err := signedPacketVerifier.
		VerifySignedPacketIdenAssert(signedPacket, nonceDb, "example.com")
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
	signedPacket, err := signedPacketSigner.
		NewSignPacketV02(600,
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

	signedPacket3, err := signedPacketSigner.
		NewSignPacketV02(600, "invalid", nil, nil)
	assert.Nil(t, err)
	signedPacketStr2, err := signedPacket3.Marshal()
	assert.Nil(t, err)
	var signedPacket4 SignedPacket
	// "invalid" is not a valid signed packet type, so unmarshal must error
	err = signedPacket4.Unmarshal(signedPacketStr2)
	assert.Error(t, err)
}
