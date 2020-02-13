package verifier

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain/readermock"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain/writermock"
	idenpubonchain "github.com/iden3/go-iden3-core/components/idenpubonchain/mock"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/holder"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var idenpuboffchaindata = map[merkletree.Hash]*idenpuboffchain.PublicData{}

var pass = []byte("my passphrase")

func Copy(dst interface{}, src interface{}) {
	srcJSON, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(srcJSON, dst); err != nil {
		panic(err)
	}
}

func newIssuer(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	idenPubOffChainWrite *writermock.IdenPubOffChainWriteMock) (*issuer.Issuer, db.Storage, *keystore.KeyStore) {
	cfg := issuer.ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	is, err := issuer.New(cfg, kOp, []merkletree.Entrier{}, storage, keyStore, idenPubOnChain, idenPubOffChainWrite)
	require.Nil(t, err)
	return is, storage, keyStore
}

func mockInitState(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	idenPubOffChainWrite *writermock.IdenPubOffChainWriteMock,
	is *issuer.Issuer, genesisState *merkletree.Hash) (*types.Transaction, *merkletree.Hash) {
	var ethTx types.Transaction
	newState, _ := is.State()
	sig, err := is.SignBinary(issuer.SigPrefixSetState, append(genesisState[:], newState[:]...))
	require.Nil(t, err)
	idenPubOnChain.On("InitState", is.ID(), genesisState, newState, []byte(nil), []byte(nil), sig).Return(&ethTx, nil).Once()
	idenPubOffChainWrite.On("Publish", mock.AnythingOfType("*idenpuboffchain.PublicData")).Return().Run(func(args mock.Arguments) {
		publicData := args.Get(0).(*idenpuboffchain.PublicData)
		idenpuboffchaindata[*publicData.IdenState] = publicData
	}).Return(nil)
	idenPubOffChainWrite.On("Url").Return("https://foo.bar")
	return &ethTx, newState
}

func mockSetState(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock, is *issuer.Issuer, oldState *merkletree.Hash) (*types.Transaction, *merkletree.Hash) {
	var ethTx types.Transaction
	newState, _ := is.State()
	sig, err := is.SignBinary(issuer.SigPrefixSetState, append(oldState[:], newState[:]...))
	require.Nil(t, err)
	idenPubOnChain.On("SetState", is.ID(), newState, []byte(nil), []byte(nil), sig).Return(&ethTx, nil).Once()
	return &ethTx, newState
}

func _newIssuerIssuedClaim(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	idenPubOffChainWrite *writermock.IdenPubOffChainWriteMock,
	claim merkletree.Entrier) (*issuer.Issuer, *merkletree.Hash) {
	is, _, _ := newIssuer(t, idenPubOnChain, idenPubOffChainWrite)
	genesisState, _ := is.State()
	err := is.IssueClaim(claim)
	require.Nil(t, err)

	_, newState := mockInitState(t, idenPubOnChain, idenPubOffChainWrite, is, genesisState)

	// Publishing state for the first time
	err = is.PublishState()
	require.Nil(t, err)

	return is, newState
}

func newIssuerIssuedClaim(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	idenPubOffChainWrite *writermock.IdenPubOffChainWriteMock,
	claim merkletree.Entrier) *issuer.Issuer {
	is, newState := _newIssuerIssuedClaim(t, idenPubOnChain, idenPubOffChainWrite, claim)

	blockN := uint64(12)
	blockTs := int64(105000)
	idenPubOnChain.On("GetState", is.ID()).Return(&proof.IdenStateData{IdenState: newState, BlockN: blockN, BlockTs: blockTs}, nil)
	idenPubOnChain.On("GetStateByBlock", is.ID(), blockN).
		Return(&proof.IdenStateData{BlockN: blockN, BlockTs: blockTs, IdenState: newState}, nil)

	err := is.SyncIdenStatePublic()
	require.Nil(t, err)

	return is
}

func TestVerifyCredentialExistence(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	idenPubOffChainWrite := writermock.New()
	indexBytes, dataBytes := [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	indexBytes[0] = 0x42
	claim := claims.NewClaimBasic(indexBytes, dataBytes, 0)
	is := newIssuerIssuedClaim(t, idenPubOnChain, idenPubOffChainWrite, claim)

	credExist, err := is.GenCredentialExistence(claim)
	require.Nil(t, err)

	var now time.Time
	verifier := NewWithTimeNow(idenPubOnChain, func() time.Time {
		return now
	})

	// Good Cred Exist
	err = verifier.VerifyCredentialExistence(credExist)
	assert.Nil(t, err)

	// Cred Exist is proof non existence
	credExistBad := &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.MtpClaim.Existence = false
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad Id
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.Id[4] = 0x00
	credExistBad.Id[5] = 0x00
	credExistBad.Id[6] = 0x00
	require.NotEqual(t, credExist, credExistBad)
	idenPubOnChain.On("GetStateByBlock", credExistBad.Id, credExistBad.IdenStateData.BlockN).
		Return(&proof.IdenStateData{IdenState: &merkletree.HashZero}, nil)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad RootsTreeRoot
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.RootsTreeRoot[0] = 0x00
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad IdenState
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.IdenStateData.IdenState[1] = 0x00
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad BlockN
	idenPubOnChain.On("GetStateByBlock", is.ID(), uint64(01)).
		Return(&proof.IdenStateData{IdenState: &merkletree.HashZero}, nil)
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.IdenStateData.BlockN = 01
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// TODO: Uncomment once smart contract returns BlockTs and BlockN every time
	// Cred Exist has bad BlockTs
	// credExistBad = &proof.CredentialExistence{}
	// copier.Copy(credExistBad, credExist)
	// credExistBad.IdenStateData.BlockTs = 02
	// require.NotEqual(t, credExist, credExistBad)
	// err = verifier.VerifyCredentialExistence(credExistBad)
	// assert.NotNil(t, err)

	// Cred Exist has bad Claim
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	indexBytes, dataBytes = [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	indexBytes[0] = 0x88
	credExistBad.Claim = claims.NewClaimBasic(indexBytes, dataBytes, 0).Entry()
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)
}

func newHolder(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock,
	idenPubOffChainWrite *writermock.IdenPubOffChainWriteMock,
	idenPubOffChainRead *readermock.IdenPubOffChainReadMock) (*holder.Holder, db.Storage, *keystore.KeyStore) {
	cfg := holder.ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	ho, err := holder.New(cfg, kOp, []merkletree.Entrier{}, storage, keyStore, idenPubOnChain, idenPubOffChainWrite, idenPubOffChainRead)
	require.Nil(t, err)
	return ho, storage, keyStore
}

func TestVerifyCredentialValidity(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	idenPubOffChainWrite := writermock.New()

	now := time.Unix(0, 0)
	verifier := NewWithTimeNow(idenPubOnChain, func() time.Time {
		return now
	})

	idenPubOffChainRead := readermock.New()
	mockCall := idenPubOffChainRead.On("GetPublicData", mock.AnythingOfType("string"),
		mock.AnythingOfType("*core.ID"), mock.AnythingOfType("*merkletree.Hash"))
	mockCall.RunFn = func(args mock.Arguments) {
		// url := args.Get(0).(string)
		// id := args.Get(1).(*core.Id)
		idenState := args.Get(2).(*merkletree.Hash)
		publicData, ok := idenpuboffchaindata[*idenState]
		if ok {
			mockCall.ReturnArguments = mock.Arguments{publicData, nil}
		} else {
			mockCall.ReturnArguments = mock.Arguments{nil, fmt.Errorf("No public data found for idenState %v", idenState)}
		}
	}

	ho, _, _ := newHolder(t, idenPubOnChain, nil, idenPubOffChainRead)

	//
	// {Ts: 100, BlockN: 12} -> claim1 is added
	//
	now = time.Unix(100, 0)

	// ISSUER: Publish state first time with claim1

	indexBytes, dataBytes := [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	indexBytes[0] = 0x42
	claim1 := claims.NewClaimBasic(indexBytes, dataBytes, 11)
	is, newState := _newIssuerIssuedClaim(t, idenPubOnChain, idenPubOffChainWrite, claim1)

	blockN := uint64(12)
	blockTs := int64(100)
	idenPubOnChain.On("GetState", is.ID()).Return(&proof.IdenStateData{IdenState: newState, BlockN: blockN, BlockTs: blockTs}, nil).Twice()
	idenPubOnChain.On("GetStateByBlock", is.ID(), blockN).
		Return(&proof.IdenStateData{BlockN: blockN, BlockTs: blockTs, IdenState: newState}, nil)

	err := is.SyncIdenStatePublic()
	require.Nil(t, err)

	credExistClaim1, err := is.GenCredentialExistence(claim1)
	require.Nil(t, err)

	// HOLDER + VERIFIER

	credValidClaim1t1, err := ho.HolderGetCredentialValidity(credExistClaim1)
	require.Nil(t, err)

	err = verifier.VerifyCredentialValidity(credValidClaim1t1, 50*time.Second)
	assert.Nil(t, err)

	//
	// {Ts: 200, BlockN: 13} -> claim2 is added
	//
	now = time.Unix(200, 0)

	// ISSUER: Publish state a second time with another claim2

	indexBytes, dataBytes = [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	indexBytes[0] = 0x48
	claim2 := claims.NewClaimBasic(indexBytes, dataBytes, 22)

	err = is.IssueClaim(claim2)
	require.Nil(t, err)

	_, newState = mockSetState(t, idenPubOnChain, is, newState)
	err = is.PublishState()
	require.Nil(t, err)

	blockN = uint64(13)
	blockTs = int64(200)
	idenPubOnChain.On("GetState", is.ID()).Return(&proof.IdenStateData{IdenState: newState, BlockN: blockN, BlockTs: blockTs}, nil).Times(5)
	idenPubOnChain.On("GetStateByBlock", is.ID(), blockN).
		Return(&proof.IdenStateData{BlockN: blockN, BlockTs: blockTs, IdenState: newState}, nil)

	err = is.SyncIdenStatePublic()
	require.Nil(t, err)

	credExistClaim2, err := is.GenCredentialExistence(claim2)
	require.Nil(t, err)

	// HOLDER + VERIFIER

	credValidClaim1t2, err := ho.HolderGetCredentialValidity(credExistClaim1)
	assert.Nil(t, err)
	assert.NotNil(t, credValidClaim1t2)

	credValidClaim2t2, err := ho.HolderGetCredentialValidity(credExistClaim2)
	assert.Nil(t, err)
	assert.NotNil(t, credValidClaim2t2)

	// Outdated is invalid
	err = verifier.VerifyCredentialValidity(credValidClaim1t1, 50*time.Second)
	assert.Error(t, err)

	// With more freshness time it's valid
	err = verifier.VerifyCredentialValidity(credValidClaim1t1, 150*time.Second)
	assert.Nil(t, err)

	// Recent one is valid
	err = verifier.VerifyCredentialValidity(credValidClaim1t2, 50*time.Second)
	assert.Nil(t, err)

	err = verifier.VerifyCredentialValidity(credValidClaim2t2, 50*time.Second)
	assert.Nil(t, err)

	//
	// {Ts: 300, BlockN: 14} -> claim1 is revoked
	//
	now = time.Unix(300, 0)

	// ISSUER: Publish state a third time revoking claim1

	err = is.RevokeClaim(claim1)
	require.Nil(t, err)

	_, newState = mockSetState(t, idenPubOnChain, is, newState)
	err = is.PublishState()
	require.Nil(t, err)

	blockN = uint64(14)
	blockTs = int64(300)
	idenPubOnChain.On("GetState", is.ID()).Return(&proof.IdenStateData{IdenState: newState, BlockN: blockN, BlockTs: blockTs}, nil)
	idenPubOnChain.On("GetStateByBlock", is.ID(), blockN).
		Return(&proof.IdenStateData{BlockN: blockN, BlockTs: blockTs, IdenState: newState}, nil)

	err = is.SyncIdenStatePublic()
	require.Nil(t, err)

	err = verifier.VerifyCredentialExistence(credExistClaim1)
	assert.Nil(t, err)

	// HOLDER + VERIFIER

	_, err = ho.HolderGetCredentialValidity(credExistClaim1)
	assert.Equal(t, holder.ErrRevokedClaim, err)

	credValidClaim2t3, err := ho.HolderGetCredentialValidity(credExistClaim2)
	assert.Nil(t, err)

	// C1T2 with long freshness is valid
	err = verifier.VerifyCredentialValidity(credValidClaim1t1, 250*time.Second)
	assert.Nil(t, err)

	// C2T2 with mid freshness is valid
	err = verifier.VerifyCredentialValidity(credValidClaim2t2, 150*time.Second)
	assert.Nil(t, err)

	// C2T3 with small freshness is valid
	err = verifier.VerifyCredentialValidity(credValidClaim2t3, 50*time.Second)
	assert.Nil(t, err)

	//
	// {Ts: 400, BlockN: --}
	//
	now = time.Unix(400, 0)
}
