package issuer

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	idenpubonchain "github.com/iden3/go-iden3-core/components/idenpubonchain/mock"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var pass = []byte("my passphrase")

func newIssuer(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock) (*Issuer, db.Storage, *keystore.KeyStore) {
	cfg := ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	issuer, err := New(cfg, kOp, []merkletree.Entrier{}, storage, keyStore, idenPubOnChain)
	require.Nil(t, err)
	return issuer, storage, keyStore
}

func TestNewLoadIssuer(t *testing.T) {
	issuer, storage, keyStore := newIssuer(t, nil)

	issuerLoad, err := Load(storage, keyStore, nil)
	require.Nil(t, err)

	assert.Equal(t, issuer.cfg, issuerLoad.cfg)
	assert.Equal(t, issuer.id, issuerLoad.id)
}

func TestIssuerGenesis(t *testing.T) {
	issuer, _, _ := newIssuer(t, nil)

	assert.Equal(t, issuer.revocationsMt.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func TestIssuerFull(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	issuer, _, _ := newIssuer(t, idenPubOnChain)

	assert.Equal(t, issuer.revocationsMt.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func TestIssuerPublish(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	issuer, _, _ := newIssuer(t, idenPubOnChain)

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

	indexBytes, dataBytes := [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	err = issuer.IssueClaim(claims.NewClaimBasic(indexBytes, dataBytes, 0))
	require.Nil(t, err)

	var ethTx types.Transaction
	newState, _ := issuer.state()
	require.Nil(t, err)
	sig, err := issuer.SignBinary(sigPrefixSetState, append(genesisState[:], newState[:]...))
	require.Nil(t, err)
	idenPubOnChain.On("InitState", issuer.id, genesisState, newState, []byte(nil), []byte(nil), sig).Return(&ethTx, nil).Once()

	// Publishing state for the first time
	err = issuer.PublishState()
	require.Nil(t, err)
	assert.Equal(t, &merkletree.HashZero, issuer.idenStateOnChain())
	assert.Equal(t, newState, issuer.idenStatePending())

	idenPubOnChain.On("GetState", issuer.id).Return(&proof.IdenStateData{IdenState: &merkletree.HashZero}, nil).Once()

	// Sync (not yet on the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, &merkletree.HashZero, issuer.idenStateOnChain())
	assert.Equal(t, newState, issuer.idenStatePending())

	idenPubOnChain.On("GetState", issuer.id).Return(&proof.IdenStateData{IdenState: newState}, nil).Once()

	// Sync (finally in the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, newState, issuer.idenStateOnChain())
	assert.Equal(t, &merkletree.HashZero, issuer.idenStatePending())

	//
	// State Update
	//

	indexBytes, dataBytes = [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	indexBytes[0] = 0x42
	err = issuer.IssueClaim(claims.NewClaimBasic(indexBytes, dataBytes, 0))
	require.Nil(t, err)

	oldState := newState
	newState, _ = issuer.state()
	require.Nil(t, err)
	sig, err = issuer.SignBinary(sigPrefixSetState, append(oldState[:], newState[:]...))
	require.Nil(t, err)
	idenPubOnChain.On("SetState", issuer.id, newState, []byte(nil), []byte(nil), sig).Return(&ethTx, nil).Once()

	// Publishing state update
	err = issuer.PublishState()
	require.Nil(t, err)
	assert.Equal(t, oldState, issuer.idenStateOnChain())
	assert.Equal(t, newState, issuer.idenStatePending())

	idenPubOnChain.On("GetState", issuer.id).Return(&proof.IdenStateData{IdenState: oldState}, nil).Once()

	// Sync (not yet on the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, oldState, issuer.idenStateOnChain())
	assert.Equal(t, newState, issuer.idenStatePending())

	idenPubOnChain.On("GetState", issuer.id).Return(&proof.IdenStateData{IdenState: newState}, nil).Once()

	// Sync (finally in the smart contract)
	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
	assert.Equal(t, newState, issuer.idenStateOnChain())
	assert.Equal(t, &merkletree.HashZero, issuer.idenStatePending())
}
