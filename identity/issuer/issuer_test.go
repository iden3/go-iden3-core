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
	idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter) (*Issuer, db.Storage, *keystore.KeyStore) {
	cfg := ConfigDefault
	cfg.GenesisOnly = genesisOnly
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
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
	issuer, storage, keyStore := newIssuer(t, true, nil, nil)

	issuerLoad, err := Load(storage, keyStore, nil, nil, nil)
	require.Nil(t, err)

	assert.Equal(t, issuer.cfg, issuerLoad.cfg)
	assert.Equal(t, issuer.id, issuerLoad.id)
}

func TestIssuerGenesis(t *testing.T) {
	issuer, _, _ := newIssuer(t, true, nil, nil)

	assert.Equal(t, issuer.revocationsTree.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func TestIssuerFull(t *testing.T) {
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain)

	assert.Equal(t, issuer.revocationsTree.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func TestIssuerPublish(t *testing.T) {
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain)

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
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain)

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
	issuer, _, _ := newIssuer(t, false, idenPubOnChain, idenPubOffChain)
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
