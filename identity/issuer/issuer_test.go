package issuer

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/iden3/go-circom-prover-verifier/parsers"
	zktypes "github.com/iden3/go-circom-prover-verifier/types"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	idenpuboffchanlocal "github.com/iden3/go-iden3-core/components/idenpuboffchain/local"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	idenpubonchainlocal "github.com/iden3/go-iden3-core/components/idenpubonchain/local"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
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

/*
func TestDebug(t *testing.T) {
	str2bigInt := func(s string) *big.Int {
		v, ok := new(big.Int).SetString(s, 0)
		if !ok {
			panic("bad big int string")
		}
		return v
	}
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
	inputs := []witnesscalc.Input{
		witnesscalc.Input{"id", str2bigInt("325175891201904061770219815708117754716107445234074270409340186911740723200")},
		witnesscalc.Input{"oldIdState", str2bigInt("1880270691508214773256930953745800337960378579419980864662765820435726350915")},
		witnesscalc.Input{"userPrivateKey", str2bigInt("4606569494897889584070207882588197822892560892629729378653253101357531598856")},
		witnesscalc.Input{"siblings", []interface{}{new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}},
		witnesscalc.Input{"claimsTreeRoot", str2bigInt("121541882510195987597965007205101834241828157974971046026151697698584672909")},
		witnesscalc.Input{"newIdState", str2bigInt("12002223098192248627320135710707280782046802862792472435575438852402516963184")},
	}

	require.NotNil(t, inputs)

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
	/////

	// GOOD
	// pubSignals := []*big.Int{
	// 	str2bigInt("210345594388067481897037608142723337231630498852282180796063929411563290624"),
	// 	str2bigInt("41"),
	// 	str2bigInt("42"),
	// }

	// BAD
	pubSignals := []*big.Int{
		str2bigInt("325175891201904061770219815708117754716107445234074270409340186911740723200"),
		str2bigInt("1880270691508214773256930953745800337960378579419980864662765820435726350915"),
		str2bigInt("12002223098192248627320135710707280782046802862792472435575438852402516963184"),
	}

	// GOOD
	// proofJSON := `{"pi_a":["17121680560078803949092571371457517692394117524984333832297266994314891453981","5182587154084293109628994627039528999621720381658302077755004006059561006773","1"],"pi_b":[["18283156164936766209016307943511972181583059224027264002405807951443749051393","1614455403515250685551181318284419058529501446417380794985805297587539949280"],["11133574583507279696795301069098283376348897330090789375820568304616692077363","18296317170826843116315986665349444498466742529253113139653812312450257784864"],["1","0"]],"pi_c":["2188868170530467947384242437891911634909822584402855510318077715648918670780","5047743681687443948747324124477361890830782512035992997889344688091933305051","1"],"protocol":"groth"}`

	// BAD
	proofJSON := `{"pi_a":["20778664224953409296976945287939109438694820507452424811099375225359937285952","21667851703947503639706576639963504172539372684430553878144517158772144899386","1"],"pi_b":[["5310656177980863263118505533331700673453892513357177418427765501229179266469","14644187988562369431563004050015817393837345313282511727522766370057669141108"],["13680652776805440525070853446099906696336640566860843313858633058268781485299","17777307777058307587373036730576330839421390796407149585066261101053048787617"],["1","0"]],"pi_c":["11519688400344297020347306976710054392753352273290082278933763596213211972432","7670943716371955179636215412041439324843742940549923991627388743806904707703","1"],"protocol":"groth"}`

	proof, err := zkparsers.ParseProof([]byte(proofJSON))
	require.Nil(t, err)

	// Verify zk proof
	v = verifier.Verify(vk, proof, pubSignals)
	assert.True(t, v)
}
*/

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
