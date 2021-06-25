package issuer

import (
	"os"
	"testing"
	"time"

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
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-merkletree"
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
	idenStatePending, _ := issuer.idenStatePending()
	assert.Equal(t, &merkletree.HashZero, idenStatePending)

	tx, err := issuer.storage.NewTx()
	require.Nil(t, err)
	idenStateListLen, err := issuer.idenStateList.Length(tx)
	require.Nil(t, err)
	assert.Equal(t, uint32(1), idenStateListLen)
	idenStateLast, _, err := issuer.getIdenStateByIdx(tx, -1)
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
	idenStatePending, _ = issuer.idenStatePending()
	assert.Equal(t, newState, idenStatePending)

	// Sync (not yet on the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, &merkletree.HashZero, issuer.idenStateOnChain())
	idenStatePending, _ = issuer.idenStatePending()
	assert.Equal(t, newState, idenStatePending)

	// Sync (finally in the smart contract)
	idenPubOnChain.Sync()
	blockN += 10
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, newState, issuer.idenStateOnChain())
	idenStatePending, _ = issuer.idenStatePending()
	assert.Equal(t, &merkletree.HashZero, idenStatePending)

	//
	// State Update
	//

	indexBytes, valueBytes = [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x42
	err = issuer.IssueClaim(claims.NewClaimBasic(indexBytes, valueBytes))
	require.Nil(t, err)

	oldState := newState

	// Publishing state update
	err = issuer.PublishState()
	newState, _ = issuer.State()
	require.Nil(t, err)
	assert.Equal(t, oldState, issuer.idenStateOnChain())
	idenStatePending, _ = issuer.idenStatePending()
	assert.Equal(t, newState, idenStatePending)

	// Sync (not yet on the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, oldState, issuer.idenStateOnChain())
	idenStatePending, _ = issuer.idenStatePending()
	assert.Equal(t, newState, idenStatePending)

	// Sync (finally in the smart contract)
	idenPubOnChain.Sync()
	blockN += 10
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, newState, issuer.idenStateOnChain())
	idenStatePending, _ = issuer.idenStatePending()
	assert.Equal(t, &merkletree.HashZero, idenStatePending)
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
	idenStatePending, _ := issuer.idenStatePending()
	assert.Equal(t, &merkletree.HashZero, idenStatePending)

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
	zkFiles := zkutils.NewZkFiles("http://161.35.72.58:9000/circuit-idstate/", "/tmp/iden3/idenstatezk-issuer",
		zkutils.ProvingKeyFormatJSON,
		zkutils.ZkFilesHashes{
			ProvingKey:      "2c72fceb10323d8b274dbd7649a63c1b6a11fff3a1e4cd7f5ec12516f32ec452",
			VerificationKey: "473952ff80aef85403005eb12d1e78a3f66b1cc11e7bd55d6bfe94e0b5577640",
			WitnessCalcWASM: "8eafd9314c4d2664a23bf98a4f42cd0c29984960ae3544747ba5fbd60905c41f",
		}, true)
	if err := zkFiles.LoadAll(); err != nil {
		panic(err)
	}

	var err error
	vk, err = zkFiles.VerificationKey()
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
		Levels: 16,
		Files:  *zkFiles,
	}
	os.Exit(m.Run())
}
