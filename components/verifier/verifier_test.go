package verifier

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	idenpubonchain "github.com/iden3/go-iden3-core/components/idenpubonchain/mock"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func newIssuer(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock) (*issuer.Issuer, db.Storage, *keystore.KeyStore) {
	cfg := issuer.ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	is, err := issuer.New(cfg, kOp, []merkletree.Entrier{}, storage, keyStore, idenPubOnChain)
	require.Nil(t, err)
	return is, storage, keyStore
}

func mockInitState(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock, is *issuer.Issuer, genesisState *merkletree.Hash) (*types.Transaction, *merkletree.Hash) {
	var ethTx types.Transaction
	newState, _ := is.State()
	sig, err := is.SignBinary(issuer.SigPrefixSetState, append(genesisState[:], newState[:]...))
	require.Nil(t, err)
	idenPubOnChain.On("InitState", is.ID(), genesisState, newState, []byte(nil), []byte(nil), sig).Return(&ethTx, nil).Once()
	return &ethTx, newState
}

func newIssuerIssuedClaim(t *testing.T, idenPubOnChain *idenpubonchain.IdenPubOnChainMock, claim merkletree.Entrier) *issuer.Issuer {
	is, _, _ := newIssuer(t, idenPubOnChain)
	genesisState, _ := is.State()
	err := is.IssueClaim(claim)
	require.Nil(t, err)

	_, newState := mockInitState(t, idenPubOnChain, is, genesisState)

	// Publishing state for the first time
	err = is.PublishState()
	require.Nil(t, err)

	blockN := uint64(12)
	blockTs := int64(105000)
	idenPubOnChain.On("GetState", is.ID()).Return(&proof.IdenStateData{IdenState: newState, BlockN: blockN, BlockTs: blockTs}, nil)
	idenPubOnChain.On("GetStateByBlock", is.ID(), blockN).Return(newState, nil)

	err = is.SyncIdenStatePublic()
	require.Nil(t, err)

	return is
}

func TestVerifyCredentialExistence(t *testing.T) {
	idenPubOnChain := idenpubonchain.New()
	indexBytes, dataBytes := [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	indexBytes[0] = 0x42
	claim := claims.NewClaimBasic(indexBytes, dataBytes, 0)
	is := newIssuerIssuedClaim(t, idenPubOnChain, claim)

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
	idenPubOnChain.On("GetStateByBlock", credExistBad.Id, credExistBad.IdenStateData.BlockN).Return(&merkletree.HashZero, nil)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad RootsRoot
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.RootsRoot[0] = 0x00
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad IdenState
	credExistBad = &proof.CredentialExistence{}
	//copier.Copy(credExistBad, credExist)
	Copy(credExistBad, credExist)
	credExistBad.IdenStateData.IdenState[1] = 0x00
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// TODO: Update once smart contract returns BlockTs and BlockN every time
	// Cred Exist has bad BlockN
	// stateDataOnChain := is.StateDataOnChain()
	idenPubOnChain.On("GetStateByBlock", is.ID(), uint64(01)).Return(&merkletree.HashZero, nil)
	credExistBad = &proof.CredentialExistence{}
	copier.Copy(credExistBad, credExist)
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
	copier.Copy(credExistBad, credExist)
	indexBytes, dataBytes = [claims.IndexSlotBytes]byte{}, [claims.DataSlotBytes]byte{}
	indexBytes[0] = 0x88
	credExistBad.Claim = claims.NewClaimBasic(indexBytes, dataBytes, 0).Entry()
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

}
