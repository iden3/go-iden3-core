package issuer

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain/writermock"
	idenpubonchain "github.com/iden3/go-iden3-core/components/idenpubonchain/mock"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var pass = []byte("my passphrase")

func newIssuer(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	idenPubOffChainWrite *writermock.IdenPubOffChainWriteMock) (*Issuer, db.Storage, *keystore.KeyStore) {
	cfg := ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	issuer, err := New(cfg, kOp, []claims.Claimer{}, storage, keyStore, idenPubOnChain, idenPubOffChainWrite)
	require.Nil(t, err)
	return issuer, storage, keyStore
}

func TestNewLoadIssuer(t *testing.T) {
	issuer, storage, keyStore := newIssuer(t, nil, nil)

	issuerLoad, err := Load(storage, keyStore, nil, nil)
	require.Nil(t, err)

	assert.Equal(t, issuer.cfg, issuerLoad.cfg)
	assert.Equal(t, issuer.id, issuerLoad.id)
}

func TestIssuerGenesis(t *testing.T) {
	issuer, _, _ := newIssuer(t, nil, nil)

	assert.Equal(t, issuer.revocationsTree.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func TestIssuerFull(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	idenPubOffChainWrite := writermock.New()
	issuer, _, _ := newIssuer(t, idenPubOnChain, idenPubOffChainWrite)

	assert.Equal(t, issuer.revocationsTree.RootKey(), &merkletree.HashZero)

	idenState, _ := issuer.state()
	assert.Equal(t, core.IdGenesisFromIdenState(idenState), issuer.ID())
}

func mockInitState(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	idenPubOffChainWrite *writermock.IdenPubOffChainWriteMock,
	issuer *Issuer, genesisState *merkletree.Hash) (*types.Transaction, *merkletree.Hash) {
	var ethTx types.Transaction
	newState, _ := issuer.state()
	sig, err := issuer.SignBinary(SigPrefixSetState, append(genesisState[:], newState[:]...))
	require.Nil(t, err)
	idenPubOnChain.On("InitState", issuer.id, genesisState, newState, []byte(nil), []byte(nil), sig).Return(&ethTx, nil).Once()
	idenPubOffChainWrite.On("Publish", mock.AnythingOfType("*idenpuboffchain.PublicData")).Return().Run(func(args mock.Arguments) {
	}).Return(nil)
	idenPubOffChainWrite.On("Url").Return("https://foo.bar")
	return &ethTx, newState
}

func mockSetState(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	issuer *Issuer, oldState *merkletree.Hash) (*types.Transaction, *merkletree.Hash) {
	var ethTx types.Transaction
	newState, _ := issuer.state()
	sig, err := issuer.SignBinary(SigPrefixSetState, append(oldState[:], newState[:]...))
	require.Nil(t, err)
	idenPubOnChain.On("SetState", issuer.id, newState, []byte(nil), []byte(nil), sig).Return(&ethTx, nil).Once()
	return &ethTx, newState
}

func TestIssuerPublish(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	idenPubOffChainWrite := writermock.New()
	issuer, _, _ := newIssuer(t, idenPubOnChain, idenPubOffChainWrite)

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

	_, newState := mockInitState(t, idenPubOnChain, idenPubOffChainWrite, issuer, genesisState)

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

	indexBytes, valueBytes = [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x42
	err = issuer.IssueClaim(claims.NewClaimBasic(indexBytes, valueBytes))
	require.Nil(t, err)

	oldState := newState
	_, newState = mockSetState(t, idenPubOnChain, issuer, oldState)

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

func TestIssuerCredential(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	idenPubOffChainWrite := writermock.New()
	issuer, _, _ := newIssuer(t, idenPubOnChain, idenPubOffChainWrite)
	genesisState, _ := issuer.state()

	// Issue a Claim
	indexBytes, valueBytes := [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x42
	claim0 := claims.NewClaimBasic(indexBytes, valueBytes)

	err := issuer.IssueClaim(claim0)
	require.Nil(t, err)

	credExist, err := issuer.GenCredentialExistence(claim0)
	assert.Nil(t, credExist)
	assert.Equal(t, ErrIdenStateOnChainZero, err)

	_, newState := mockInitState(t, idenPubOnChain, idenPubOffChainWrite, issuer, genesisState)
	err = issuer.PublishState()
	require.Nil(t, err)

	idenPubOnChain.On("GetState", issuer.id).Return(&proof.IdenStateData{IdenState: newState}, nil).Once()

	err = issuer.SyncIdenStatePublic()
	require.Nil(t, err)
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
