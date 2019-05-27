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

const relaySkHex = "4406831fa7bb87d8c92fc65f090a6017916bd2197ffca0e1e97933b14e8f5de5"

const kSignSkHex = "4be5471a938bdf3606888472878baace4a6a64e14a153adf9a1333969e4e573c"

//const idAddrHex = "0x970e8128ab834e8eac17ab8e3812f010678cf791"
// const idAddrHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
const idHex = "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"

// root 0x1d9d41171c4b621ff279e2acb84d8ab45612fef53e37225bdf67e8ad761c3922
const proofKSignJSON = `
{
  "proofs": [{
    "mtp0": "0x00010000000000000000000000000000000000000000000000000000000000012910a6fba42851f8282e0266c887e09db4fd84975a76a6c6ce468651683d2346",
    "mtp1": "0x03010000000000000000000000000000000000000000000000000000000000012910a6fba42851f8282e0266c887e09db4fd84975a76a6c6ce468651683d23460df35209acb01f06874be29f52842b53c93512cc9a5c634677db9019acd6289c06d4571fb9634e4bed32e265f91a373a852c476656c5c13b09bc133ac61bc5a6",
    "root": "0x27b5c03fe418e310161239e274c0ee17bb12a69ab951e8c16bb1ff7327f5c194",
    "aux": {
      "version": 0,
      "era": 0,
      "id": "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"
    }
  }, {
    "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "mtp1": "0x030000000000000000000000000000000000000000000000000000000000000020b468baa588865efc5df741e0a48569aa1171143a8627f425fff0d4fa7803c708c63afb2b908f0d1c79a88c212d88e95ecdf4987c15706848a10b124ebc1389",
    "root": "0x01460efcd144ff6e92244cc51c59cb322886599fc2092b9d4b09b5bce79ddfda",
    "aux": null
  }],
  "leaf": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034a377d7da59e11540b28034f08f10c35de698b99484b8efdafa2bd8ec476000000000000000000000000000000000000e231000000000000000000000004",
  "date": 1557843460,
  "signature": "0x59a423c5d6de2dfc9b4e23fef0b23d7300ca47441e8a2cd8d7c4d6be45a86a6716e64a87ad51db043a594116ce8a766dfffe475833703ffe5bf006cd2fb4eff11b",
  "signer": "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"
}
`

const ethName = "testName@iden3.io"

//o0x0d1fadf720af488d10aaa4bdaf6a8d1163ad30b19624082b0e4403934ab57ff3
const proofAssignNameJSON = `
{
  "proofs": [{
    "mtp0": "0x000300000000000000000000000000000000000000000000000000000000000401460efcd144ff6e92244cc51c59cb322886599fc2092b9d4b09b5bce79ddfda",
    "mtp1": "0x010200000000000000000000000000000000000000000000000000000000000226620ae90f5af5a98b0fffc2b2aaf0b1ba00784b5d3731e27d16732a05393a22",
    "root": "0x1820d7d5639e770c1e8ee4583b0bc626c5f61af24d3c358b4cc06faef2802ced",
    "aux": null
  }],
  "leaf": "0x00000000000000000000000000000000000000000000000000000000000000000000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f40063ea0983a784a474e012c2ce392b45296419d9b57f91533c579a691db028f30000000000000000000000000000000000000000000000000000000000000003",
  "date": 1557843460,
  "signature": "0x33a697aad066743f7ba852c93b907691e3d3a53d6dd82e33d808afb6c6714bca4f316bdfc9b3070e049883b6852f0ea6b1162068dcc2fe8347e4131e56c54b0e1b",
  "signer": "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"
}
`

const namesFileContent = `
{
  "iden3.io": "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z"
}
`

const entititesFileContent = `
{
  "1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z": {
    "name": "iden3-test-relay",
    "kOp": "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c",
    "kOpPub": "0x036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd4d59",
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
		NewSignPacketV01(600, GENERICSIGV01, nil, form)
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
		NewSignPacketV01(600,
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
		NewSignPacketV01(600, "invalid", nil, nil)
	assert.Nil(t, err)
	signedPacketStr2, err := signedPacket3.Marshal()
	assert.Nil(t, err)
	var signedPacket4 SignedPacket
	// "invalid" is not a valid signed packet type, so unmarshal must error
	err = signedPacket4.Unmarshal(signedPacketStr2)
	assert.Error(t, err)
}
