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
	"github.com/ethereum/go-ethereum/crypto"
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

const relaySkHex = "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
const relayIdAddrHex = "0x0123456789abcdef0123456789abcdef01234567"

//const kSignSkHex = "7517685f1693593d3263460200ed903370c2318e8ba4b9bb5727acae55c32b3d"
// const kSignSkHex = "0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4"
const kSignSkHex = "7517685f1693593d3263460200ed903370c2318e8ba4b9bb5727acae55c32b3d"

//const idAddrHex = "0x970e8128ab834e8eac17ab8e3812f010678cf791"
// const idAddrHex = "0x308eff1357e7b5881c00ae22463b0f69a0d58adb"
const idAddrHex = "1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY"

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
      "idAddr": "1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY"
    }
  }, {
    "mtp0": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "mtp1": "0x03000000000000000000000000000000000000000000000000000000000000000253fb07b4de48eca430e41cbf73dc1808716b1d3cb18f3c9c787ccdf4b8aae108c63afb2b908f0d1c79a88c212d88e95ecdf4987c15706848a10b124ebc1389",
    "root": "0x1e037ddc9e134f3e3d444241213cc289ded7209e145e4b0a511a6a002b08541c",
    "aux": null
  }],
  "leaf": "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034a377d7da59e11540b28034f08f10c35de698b99484b8efdafa2bd8ec476000000000000000000000000000000000000e231000000000000000000000004",
  "date": 1557841171,
  "signature": "0xb29b1c6bf6e8801cde2a2c0e4bde4e90f385fa29da45f42a971784ee97f5e3bc6bddaf4d6abc2a3167c1aa0c3238abe64a0dc54a4904d87562dcec4f5dac0d731b",
  "signer": "1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY"
}
`

const ethName = "testName@iden3.io"

//o0x0d1fadf720af488d10aaa4bdaf6a8d1163ad30b19624082b0e4403934ab57ff3
const proofAssignNameJSON = `
{
  "proofs": [{
    "mtp0": "0x00020000000000000000000000000000000000000000000000000000000000021e037ddc9e134f3e3d444241213cc289ded7209e145e4b0a511a6a002b08541c",
    "mtp1": "0x030200000000000000000000000000000000000000000000000000000000000224ad7dcfa903830b90183014644590b4a5b8a76454d82dd14b56a2e2dc69b5ce0253fb07b4de48eca430e41cbf73dc1808716b1d3cb18f3c9c787ccdf4b8aae108c63afb2b908f0d1c79a88c212d88e95ecdf4987c15706848a10b124ebc1389",
    "root": "0x1c33a5125b07e90fce3e9698f0c9c510dbd9b2638605fa9f5cc83b81aa990649",
    "aux": null
  }],
  "leaf": "0x000000000000000000000000000000000000000000000000000000000000000000000407be6b1c3fe8ca2e03bf7ed1f29917b8e2cd56e8dcd401d65ea0e6796f0063ea0983a784a474e012c2ce392b45296419d9b57f91533c579a691db028f30000000000000000000000000000000000000000000000000000000000000003",
  "date": 1557841171,
  "signature": "0xdd5a999be0f9c915aa30592936c00051892c5facdc0dce68556a81aa3062148616bb56eea33419defefc300d0e97c8e573a720e00803cbc1cb1f0df072d11d451c",
  "signer": "1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY"
}
`

const namesFileContent = `
{
  "iden3.io": "1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY"
}
`

const entititesFileContent = `
{
  "1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY": {
    "name": "iden3-test-relay",
    "kOpAddr": "0xe0fbce58cfaa72812103f003adce3f284fe5fc7c",
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
var keyStoreDir string
var keyStore *keystore.KeyStore
var relaySk *ecdsa.PrivateKey
var relayPk *ecdsa.PublicKey
var kSignSk *ecdsa.PrivateKey
var kSignPk *ecdsa.PublicKey

var proofKSign core.ProofClaim

var relayIdAddr core.ID

var idAddr core.ID

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
	if _, err = keyStore.ImportECDSA(kSignSk, passphrase); err != nil {
		panic(err)
	}
	if err = keyStore.Unlock(accounts.Account{Address: crypto.PubkeyToAddress(*kSignPk)},
		passphrase); err != nil {
		panic(err)
	}

	// common3.HexDecodeInto(idAddr[:], []byte(idAddrHex))
	idAddr, err = core.IDFromString(idAddrHex)
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

	signSrv, err := signsrv.New(keyStore, accounts.Account{Address: crypto.PubkeyToAddress(*kSignPk)})
	if err != nil {
		panic(err)
	}

	signedPacketVerifier = NewSignedPacketVerifier(discoverySrv, nameResolverSrv)

	if err := json.Unmarshal([]byte(proofKSignJSON), &proofKSign); err != nil {
		panic(err)
	}
	signedPacketSigner = NewSignedPacketSigner(signSrv, proofKSign, idAddr)
}

func teardown() {
	os.RemoveAll(keyStoreDir)
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
