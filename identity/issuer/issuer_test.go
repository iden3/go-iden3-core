package issuer

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/iden3/go-circom-prover-verifier/parsers"
	zkparsers "github.com/iden3/go-circom-prover-verifier/parsers"
	"github.com/iden3/go-circom-prover-verifier/prover"
	zktypes "github.com/iden3/go-circom-prover-verifier/types"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	witnesscalc "github.com/iden3/go-circom-witnesscalc"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	idenpuboffchanlocal "github.com/iden3/go-iden3-core/components/idenpuboffchain/local"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	idenpubonchainlocal "github.com/iden3/go-iden3-core/components/idenpubonchain/local"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"
)

var idenPubOnChain *idenpubonchainlocal.IdenPubOnChain
var idenPubOffChain *idenpuboffchanlocal.IdenPubOffChain
var idenStateZkProofConf *IdenStateZkProofConf

var pass = []byte("my passphrase")

func newIssuer(t *testing.T, genesisOnly bool, idenPubOnChain idenpubonchain.IdenPubOnChainer,
	idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter, skIdx byte) (*Issuer, db.Storage, *keystore.KeyStore) {
	cfg := ConfigDefault
	cfg.GenesisOnly = genesisOnly
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	//
	kOp, err := keyStore.NewKey(pass)
	// DBG BEGIN
	// var sk babyjub.PrivateKey
	// sk[0] = skIdx
	// kOp, err := keyStore.ImportKey(sk, pass)
	// DBG END
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	_, err = Create(cfg, kOp, []claims.Claimer{}, storage, keyStore)
	require.Nil(t, err)
	issuer, err := Load(storage, keyStore, idenPubOnChain, idenStateZkProofConf, idenPubOffChainWrite)
	require.Nil(t, err)
	return issuer, storage, keyStore
}

func TestNewLoadIssuer(t *testing.T) {
	issuer, storage, keyStore := newIssuer(t, true, nil, nil, 0)

	issuerLoad, err := Load(storage, keyStore, nil, nil, nil)
	require.Nil(t, err)

	assert.Equal(t, issuer.cfg, issuerLoad.cfg)
	assert.Equal(t, issuer.id, issuerLoad.id)
}

func TestIssuerGenesis(t *testing.T) {
	issuer, _, _ := newIssuer(t, true, nil, nil, 1)

	assert.Equal(t, issuer.revocationsTree.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func TestIssuerFull(t *testing.T) {
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain, 2)

	assert.Equal(t, issuer.revocationsTree.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func TestIssuerPublish(t *testing.T) {
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain, 3)

	assert.Equal(t, &merkletree.HashZero, issuer.idenStateOnChain())
	assert.Equal(t, &merkletree.HashZero, issuer.idenStatePending())

	tx, err := issuer.storage.NewTx()
	require.Nil(t, err)
	idenStateListLen, err := issuer.idenStateList.Length(tx)
	require.Nil(t, err)
	assert.Equal(t, uint32(1), idenStateListLen)
	idenStateLast, _, err := issuer.getIdenStateByIdx(tx, idenStateListLen-1)
	assert.Nil(t, err)
	genesisState, _ := issuer.state()
	assert.Equal(t, idenStateLast, genesisState)

	// If state hasn't changed, PublisState does nothing
	err = issuer.PublishState()
	require.Nil(t, err)

	//
	// State Init
	//

	indexBytes, valueBytes := [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	err = issuer.IssueClaim(claims.NewClaimBasic(indexBytes, valueBytes))
	require.Nil(t, err)

	// Publishing state for the first time
	err = issuer.PublishState()
	require.Nil(t, err)
	assert.Equal(t, &merkletree.HashZero, issuer.idenStateOnChain())
	newState, _ := issuer.State()
	assert.Equal(t, newState, issuer.idenStatePending())

	// Sync (not yet on the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, &merkletree.HashZero, issuer.idenStateOnChain())
	assert.Equal(t, newState, issuer.idenStatePending())

	// Sync (finally in the smart contract)
	idenPubOnChain.Sync()
	blockN += 10
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, newState, issuer.idenStateOnChain())
	assert.Equal(t, &merkletree.HashZero, issuer.idenStatePending())

	//
	// State Update
	//

	indexBytes, valueBytes = [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x42
	err = issuer.IssueClaim(claims.NewClaimBasic(indexBytes, valueBytes))
	require.Nil(t, err)

	oldState := newState
	newState, _ = issuer.State()

	// Publishing state update
	err = issuer.PublishState()
	require.Nil(t, err)
	assert.Equal(t, oldState, issuer.idenStateOnChain())
	assert.Equal(t, newState, issuer.idenStatePending())

	// Sync (not yet on the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, oldState, issuer.idenStateOnChain())
	assert.Equal(t, newState, issuer.idenStatePending())

	// Sync (finally in the smart contract)
	idenPubOnChain.Sync()
	blockN += 10
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, newState, issuer.idenStateOnChain())
	assert.Equal(t, &merkletree.HashZero, issuer.idenStatePending())
}

func TestIssuerCredential(t *testing.T) {
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain, 4)

	// Issue a Claim
	indexBytes, valueBytes := [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x42
	claim0 := claims.NewClaimBasic(indexBytes, valueBytes)

	err := issuer.IssueClaim(claim0)
	require.Nil(t, err)

	credExist, err := issuer.GenCredentialExistence(claim0)
	assert.Nil(t, credExist)
	assert.Equal(t, ErrIdenStateOnChainZero, err)

	err = issuer.PublishState()
	require.Nil(t, err)

	idenPubOnChain.Sync()
	blockN += 10

	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	newState, _ := issuer.State()
	assert.Equal(t, newState, issuer.idenStateOnChain())
	assert.Equal(t, &merkletree.HashZero, issuer.idenStatePending())

	_, err = issuer.GenCredentialExistence(claim0)
	assert.Nil(t, err)

	// Issue another claim
	indexBytes, valueBytes = [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x81
	claim1 := claims.NewClaimBasic(indexBytes, valueBytes)

	err = issuer.IssueClaim(claim1)
	require.Nil(t, err)

	_, err = issuer.GenCredentialExistence(claim1)
	assert.Equal(t, ErrClaimNotYetInOnChainState, err)
}

func TestIssuerGenZkProofIdenStateUpdate(t *testing.T) {
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain, 5)
	var oldIdState, newIdState merkletree.Hash
	oldIdState[0] = 41
	newIdState[0] = 42
	// 1
	proof, err := issuer.GenZkProofIdenStateUpdate(&oldIdState, &newIdState)
	assert.Nil(t, err)

	// Verify zk proof
	v := verifier.Verify(vk, &proof.Proof, proof.PubSignals)
	assert.True(t, v)
}

// Problematic private key 6b04bf6475e16ad38baec9e722f2e6e7fc1da68aa3764687ff349e46c600d4b3
func testFail1(t *testing.T) {
	// BAD
	// _sk, _ := hex.DecodeString("6b04bf6475e16ad38baec9e722f2e6e7fc1da68aa3764687ff349e46c600d4b3")
	//_sk, _ := hex.DecodeString("6b04bf6475e16ad38baec9e722f2e6e7fc1da68aa3764687ff349e46c600d4b3")
	// 	_sk, _ := hex.DecodeString("9fa7ac2fcdf8636b7620a7314f7465609ea3cd5e371cb852dd3e0c713f54d01f")
	_sk, _ := hex.DecodeString("f338c8c4c325fb5fedeba44f8321d5f92537082c401d517d6ac39a5e40a8e631")
	assert.NotNil(t, _sk)
	var sk babyjub.PrivateKey
	copy(sk[:], _sk)

	cfg := ConfigDefault
	cfg.GenesisOnly = false
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.ImportKey(sk, pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	_, err = Create(cfg, kOp, []claims.Claimer{}, storage, keyStore)
	require.Nil(t, err)
	issuer, err := Load(storage, keyStore, idenPubOnChain, idenStateZkProofConf, idenPubOffChain)
	require.Nil(t, err)

	var oldIdState, newIdState merkletree.Hash
	oldIdState[0] = 41
	newIdState[0] = 42
	// 1
	proof, err := issuer.GenZkProofIdenStateUpdate(&oldIdState, &newIdState)
	assert.Nil(t, err)

	// Verify zk proof
	v := verifier.Verify(vk, &proof.Proof, proof.PubSignals)
	assert.True(t, v)
}

// Good private key, bad newState
func testFail2(t *testing.T) {

	// provingKeyJson, err := ioutil.ReadFile("/tmp/iden3/idenstatezk/proving_key.json")
	// require.Nil(t, err)
	// pk, err := parsers.ParsePk(provingKeyJson)
	// require.Nil(t, err)

	// // sk, _ := hex.DecodeString("2b7b9e65ce76234253abe93020067f6da43b3c2d338c9db87e810c36e5e428f8")
	// // GOOD
	// inputs0 := []witnesscalc.Input{
	// 	witnesscalc.Input{"id", str2bigInt("166174423361963799283253196798464647356443473882845164262571944143767404544")},
	// 	witnesscalc.Input{"oldIdState", str2bigInt("18061852250807542465459324171099679699732328224155010629972281397013214682936")},
	// 	witnesscalc.Input{"userPrivateKey", str2bigInt("6469972784802433160191522894288801153442857113505826083467538673986104889677")},
	// 	witnesscalc.Input{"siblings", []interface{}{new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}},
	// 	witnesscalc.Input{"claimsTreeRoot", str2bigInt("10478935728107947013312566021594556414247698779325036158514933104325301457854")},
	// 	witnesscalc.Input{"newIdState", str2bigInt("9940728637857898226593740028824531568924032722412161416083830975519935587360")},
	// }
	// // pubSignals0 := []*big.Int{
	// // 	str2bigInt("166174423361963799283253196798464647356443473882845164262571944143767404544"),
	// // 	str2bigInt("18061852250807542465459324171099679699732328224155010629972281397013214682936"),
	// // 	str2bigInt("9940728637857898226593740028824531568924032722412161416083830975519935587360"),
	// // }

	// wit0, err := zkutils.CalculateWitness("/tmp/iden3/idenstatezk/circuit.wasm", inputs0)
	// require.Nil(t, err)

	// proof0, pubSignals0, err := prover.GenerateProof(pk, wit0)
	// require.Nil(t, err)
	// PrintProof(proof0)
	// PrintPubSignals(pubSignals0)

	// v0 := verifier.Verify(vk, proof0, pubSignals0)
	// assert.True(t, v0)
	// fmt.Println("### OK 0")

	// // BAD
	// inputs1 := []witnesscalc.Input{
	// 	witnesscalc.Input{"id", str2bigInt("166174423361963799283253196798464647356443473882845164262571944143767404544")},
	// 	witnesscalc.Input{"oldIdState", str2bigInt("9940728637857898226593740028824531568924032722412161416083830975519935587360")},
	// 	witnesscalc.Input{"userPrivateKey", str2bigInt("6469972784802433160191522894288801153442857113505826083467538673986104889677")},
	// 	witnesscalc.Input{"siblings", []interface{}{new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}},
	// 	witnesscalc.Input{"claimsTreeRoot", str2bigInt("10478935728107947013312566021594556414247698779325036158514933104325301457854")},
	// 	witnesscalc.Input{"newIdState", str2bigInt("248088787701355955143148753006927777750089509977341830846802956300562259590")},
	// }
	// // pubSignals1 := []*big.Int{
	// // 	str2bigInt("166174423361963799283253196798464647356443473882845164262571944143767404544"),
	// // 	str2bigInt("9940728637857898226593740028824531568924032722412161416083830975519935587360"),
	// // 	str2bigInt("10651284305116482887729604437381231998942314694780986256266820269511507500678"),
	// // }

	// wit1, err := zkutils.CalculateWitness("/tmp/iden3/idenstatezk/circuit.wasm", inputs1)
	// require.Nil(t, err)

	// proof1, pubSignals1, err := prover.GenerateProof(pk, wit1)
	// require.Nil(t, err)
	// PrintProof(proof1)
	// PrintPubSignals(pubSignals1)

	// v1 := verifier.Verify(vk, proof1, pubSignals1)
	// assert.True(t, v1)
	// fmt.Println("### OK 1")

	proofJSON := `{"pi_a":["4182483818622962789457604543489145876976636145955445021927382241489522733674","7429345183110352880961027029084203791153470339532918565052444558221068844946","1"],"pi_b":[["16221422688188661682643450899040461754094990633044510049566061940437024717567","3746175388990977192742822154495081109714896872864757039381142289970052317519"],["20967708563755093197064464978289461555751883844848212319516420860952633812939","8292377492820166206360383623275321786802679847150981278961633605797422836503"],["1","0"]],"pi_c":["11715940235232946139783773862937799956089406192864943913592975124424179637913","7028706879621311966236363294892168449788339681539833530960032155943506155322","1"],"protocol":"groth"}`
	pubSignals1a := []*big.Int{
		str2bigInt("166174423361963799283253196798464647356443473882845164262571944143767404544"),
		str2bigInt("9940728637857898226593740028824531568924032722412161416083830975519935587360"),
		str2bigInt("10651284305116482887729604437381231998942314694780986256266820269511507500678"),
	}

	proof1a, err := zkparsers.ParseProof([]byte(proofJSON))
	require.Nil(t, err)

	v1a := verifier.Verify(vk, proof1a, pubSignals1a)
	assert.True(t, v1a)
	fmt.Println("### OK 1a")
}

func str2bigInt(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 0)
	if !ok {
		panic("bad big int string")
	}
	return v
}

func TestFail4(t *testing.T) {

	inputs := []witnesscalc.Input{
		witnesscalc.Input{"id", str2bigInt("74298297829111562926187599078727627559142484638934671394407492439729438720")},
		witnesscalc.Input{"oldIdState", str2bigInt("16750471626935575499752646689084686834854638468319988234800834523078625853271")},
		witnesscalc.Input{"userPrivateKey", str2bigInt("4466212258221602502691917127056870703087944399848567709970335594201931334125")},
		witnesscalc.Input{"siblings", []interface{}{new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}},
		witnesscalc.Input{"claimsTreeRoot", str2bigInt("9343032743785627766292049090761560376294067387689658337057045244881419074409")},
		// witnesscalc.Input{"newIdState", str2bigInt("257429961926807294422348145094196353778746940119888650397072999436360291293")},
		witnesscalc.Input{"newIdState", str2bigInt("257429961926807294422348145094196353778746940119888650397072999436360291293")},
	}

	wit, err := zkutils.CalculateWitness("/tmp/iden3/idenstatezk/circuit.wasm", inputs)
	require.Nil(t, err)
	fmt.Println("~~~ Wit 0:", wit[0])
	fmt.Println("~~~ Wit 1:", wit[1])
	fmt.Println("~~~ Wit 2:", wit[2])
	fmt.Println("~~~ Wit 3:", wit[3])
	fmt.Println("~~~ Wit 4:", wit[4])

	provingKeyJson, err := ioutil.ReadFile("/tmp/iden3/idenstatezk/proving_key.json")
	require.Nil(t, err)
	pk, err := parsers.ParsePk(provingKeyJson)
	require.Nil(t, err)
	proof0, pubSignals0, err := prover.GenerateProof(pk, wit)
	require.Nil(t, err)
	PrintProof(proof0)
	PrintPubSignals(pubSignals0)

	// fmt.Println("~~~ Input:", hex.EncodeToString(inputs[len(inputs)-1].Value.(*big.Int).Bytes()))
	// fmt.Println("~~~ Wit  :", hex.EncodeToString(wit[3].Bytes()))
	fmt.Println("~~~ Input:", inputs[len(inputs)-1].Value.(*big.Int))
	fmt.Println("~~~ Wit  :", wit[3])
	fmt.Println("~~~ PubSi:", pubSignals0[len(pubSignals0)-1])

	v := verifier.Verify(vk, proof0, pubSignals0)
	assert.True(t, v)
}

func testFail3(t *testing.T) {
	// GOOD
	// inputs := []witnesscalc.Input{
	// 	witnesscalc.Input{"id", str2bigInt("210345594388067481897037608142723337231630498852282180796063929411563290624")},
	// 	witnesscalc.Input{"oldIdState", str2bigInt("41")},
	// 	witnesscalc.Input{"userPrivateKey", str2bigInt("3991346357692901872276868545384750811175206482315550203289194924898393435079")},
	// 	witnesscalc.Input{"siblings", []interface{}{new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}},
	// 	witnesscalc.Input{"claimsTreeRoot", str2bigInt("3271425843806034297235996875862237425526109822426577697813893574173248821884")},
	// 	witnesscalc.Input{"newIdState", str2bigInt("42")},
	// }

	// BAD
	//inputs := []witnesscalc.Input{
	//	witnesscalc.Input{"id", str2bigInt("325175891201904061770219815708117754716107445234074270409340186911740723200")},
	//	witnesscalc.Input{"oldIdState", str2bigInt("1880270691508214773256930953745800337960378579419980864662765820435726350915")},
	//	witnesscalc.Input{"userPrivateKey", str2bigInt("4606569494897889584070207882588197822892560892629729378653253101357531598856")},
	//	witnesscalc.Input{"siblings", []interface{}{new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}},
	//	witnesscalc.Input{"claimsTreeRoot", str2bigInt("121541882510195987597965007205101834241828157974971046026151697698584672909")},
	//	witnesscalc.Input{"newIdState", str2bigInt("12002223098192248627320135710707280782046802862792472435575438852402516963184")},
	//}
	sk, _ := hex.DecodeString("6b04bf6475e16ad38baec9e722f2e6e7fc1da68aa3764687ff349e46c600d4b3")
	inputs := []witnesscalc.Input{
		witnesscalc.Input{"id", str2bigInt("418819213778301409762533179606089431647905360492196824360893336916713865216")},
		witnesscalc.Input{"oldIdState", str2bigInt("41")},
		witnesscalc.Input{"userPrivateKey", str2bigInt("6797712463363996089658188893540155103242810209169027195324641517139573685419")},
		witnesscalc.Input{"siblings", []interface{}{new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}},
		witnesscalc.Input{"claimsTreeRoot", str2bigInt("231766116244639713778938132818311683805082353426148533056698258597560449034")},
		witnesscalc.Input{"newIdState", str2bigInt("42")},
	}

	require.NotNil(t, inputs)
	require.NotNil(t, sk)

	/////
	wit, err := zkutils.CalculateWitness("/tmp/iden3/idenstatezk/circuit.wasm", inputs)
	require.Nil(t, err)

	provingKeyJson, err := ioutil.ReadFile("/tmp/iden3/idenstatezk/proving_key.json")
	require.Nil(t, err)
	pk, err := parsers.ParsePk(provingKeyJson)
	require.Nil(t, err)
	proof0, pubSignals0, err := prover.GenerateProof(pk, wit)
	require.Nil(t, err)
	PrintProof(proof0)
	PrintPubSignals(pubSignals0)

	v := verifier.Verify(vk, proof0, pubSignals0)
	assert.True(t, v)
	fmt.Println("### OK 0")

	return
	/////

	// GOOD
	// pubSignals := []*big.Int{
	// 	str2bigInt("210345594388067481897037608142723337231630498852282180796063929411563290624"),
	// 	str2bigInt("41"),
	// 	str2bigInt("42"),
	// }

	// BAD
	// pubSignals := []*big.Int{
	// 	str2bigInt("325175891201904061770219815708117754716107445234074270409340186911740723200"),
	// 	str2bigInt("1880270691508214773256930953745800337960378579419980864662765820435726350915"),
	// 	str2bigInt("12002223098192248627320135710707280782046802862792472435575438852402516963184"),
	// }
	pubSignals := []*big.Int{
		str2bigInt("418819213778301409762533179606089431647905360492196824360893336916713865216"),
		str2bigInt("41"),
		str2bigInt("42"),
	}

	// GOOD
	// proofJSON := `{"pi_a":["17121680560078803949092571371457517692394117524984333832297266994314891453981","5182587154084293109628994627039528999621720381658302077755004006059561006773","1"],"pi_b":[["18283156164936766209016307943511972181583059224027264002405807951443749051393","1614455403515250685551181318284419058529501446417380794985805297587539949280"],["11133574583507279696795301069098283376348897330090789375820568304616692077363","18296317170826843116315986665349444498466742529253113139653812312450257784864"],["1","0"]],"pi_c":["2188868170530467947384242437891911634909822584402855510318077715648918670780","5047743681687443948747324124477361890830782512035992997889344688091933305051","1"],"protocol":"groth"}`

	// BAD
	// proofJSON := `{"pi_a":["20778664224953409296976945287939109438694820507452424811099375225359937285952","21667851703947503639706576639963504172539372684430553878144517158772144899386","1"],"pi_b":[["5310656177980863263118505533331700673453892513357177418427765501229179266469","14644187988562369431563004050015817393837345313282511727522766370057669141108"],["13680652776805440525070853446099906696336640566860843313858633058268781485299","17777307777058307587373036730576330839421390796407149585066261101053048787617"],["1","0"]],"pi_c":["11519688400344297020347306976710054392753352273290082278933763596213211972432","7670943716371955179636215412041439324843742940549923991627388743806904707703","1"],"protocol":"groth"}`
	proofJSON := `{"pi_a":["5773561996726672532338383446183069536085501342879381318437116432684173391503","19520986992324673188106123860336056927322025025097881960515877690858456442476","1"],"pi_b":[["17912346158151027474659052414293338615579781182946687983504822338718604693958","9895127303343255076319043073971445149763880632532533855217528814588911982914"],["1459421304848879625692008551998525872743201647133539345715141766737151966221","18715333805700597881192990920442164339964764583874237017410230677371818391822"],["1","0"]],"pi_c":["3879886212009124699496862200212851347309811270564850549303198092792065645059","12810428327202120439902136418825623769996203526143152982875117752745190331190","1"],"protocol":"groth"}`

	proof, err := zkparsers.ParseProof([]byte(proofJSON))
	require.Nil(t, err)

	// Verify zk proof
	v = verifier.Verify(vk, proof, pubSignals)
	assert.True(t, v)
}

// func TestIssuerGenZkProofIdenStateUpdate(t *testing.T) {
// 	// N := 8
// 	N := 1
// 	var wg sync.WaitGroup
// 	wg.Add(N)
// 	for i := 0; i < N; i++ {
// 		go func() {
// 			issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain, 6)
// 			var oldIdState, newIdState merkletree.Hash
// 			oldIdState[0] = 41
// 			newIdState[0] = 42
// 			// 1
// 			proof, err := issuer.GenZkProofIdenStateUpdate(&oldIdState, &newIdState)
// 			assert.Nil(t, err)
//
// 			// Verify zk proof
// 			v := verifier.Verify(vk, &proof.Proof, proof.PubSignals)
// 			assert.True(t, v)
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }

var vk *zktypes.Vk
var blockN uint64

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	err := GetIdenStateZKFiles("http://161.35.72.58:9000/circuit1/")
	if err != nil {
		panic(err)
	}
	vkJSON, err := ioutil.ReadFile("/tmp/iden3/idenstatezk/verification_key.json")
	if err != nil {
		panic(err)
	}
	vk, err = parsers.ParseVk(vkJSON)
	if err != nil {
		panic(err)
	}
	idenPubOnChain = idenpubonchainlocal.New(
		func() time.Time {
			return time.Now()
		},
		func() uint64 {
			blockN += 1
			return blockN
		},
		vk,
	)
	idenPubOffChain = idenpuboffchanlocal.NewIdenPubOffChain("http://foo.bar")
	idenStateZkProofConf = &IdenStateZkProofConf{
		Levels:              16,
		PathWitnessCalcWASM: "/tmp/iden3/idenstatezk/circuit.wasm",
		PathProvingKey:      "/tmp/iden3/idenstatezk/proving_key.json",
		PathVerifyingKey:    "/tmp/iden3/idenstatezk/verification_key.json",
		CacheProvingKey:     true,
	}
	os.Exit(m.Run())
}
