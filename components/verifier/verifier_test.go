package verifier

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	idenpuboffchanlocal "github.com/iden3/go-iden3-core/components/idenpuboffchain/local"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/merkletree"

	idenpubonchainlocal "github.com/iden3/go-iden3-core/components/idenpubonchain/local"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/holder"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var blockN uint64
var blockTs int64

var idenPubOffChain *idenpuboffchanlocal.IdenPubOffChain
var idenPubOnChain *idenpubonchainlocal.IdenPubOnChain

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

func newIssuer(t *testing.T, idenPubOnChain idenpubonchain.IdenPubOnChainer,
	idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter) (*issuer.Issuer, db.Storage, *keystore.KeyStore) {
	cfg := issuer.ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	is, err := issuer.New(cfg, kOp, []claims.Claimer{}, storage, keyStore, idenPubOnChain, idenPubOffChainWrite)
	require.Nil(t, err)
	return is, storage, keyStore
}

func TestVerifyCredentialExistence(t *testing.T) {
	indexBytes, valueBytes := [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x42
	claim := claims.NewClaimBasic(indexBytes, valueBytes)

	is, _, _ := newIssuer(t, idenPubOnChain, idenPubOffChain)
	err := is.IssueClaim(claim)
	require.Nil(t, err)

	// Publishing state for the first time
	blockTs, blockN = 105000, 12
	err = is.PublishState()
	require.Nil(t, err)
	idenPubOnChain.Sync()

	err = is.SyncIdenStatePublic()
	require.Nil(t, err)

	credExist, err := is.GenCredentialExistence(claim)
	require.Nil(t, err)

	verifier := NewWithTimeNow(idenPubOnChain, func() time.Time {
		return time.Unix(blockTs, 0)
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
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.IdenStateData.BlockN = 01
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad BlockTs
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	credExistBad.IdenStateData.BlockTs = 02
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)

	// Cred Exist has bad Claim
	credExistBad = &proof.CredentialExistence{}
	Copy(credExistBad, credExist)
	indexBytes, valueBytes = [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x88
	credExistBad.Claim = claims.NewClaimBasic(indexBytes, valueBytes).Entry()
	require.NotEqual(t, credExist, credExistBad)
	err = verifier.VerifyCredentialExistence(credExistBad)
	assert.NotNil(t, err)
}

func newHolder(t *testing.T, idenPubOnChain idenpubonchain.IdenPubOnChainer,
	idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter,
	idenPubOffChainRead idenpuboffchain.IdenPubOffChainReader) (*holder.Holder, db.Storage, *keystore.KeyStore) {
	cfg := holder.ConfigDefault
	storage := db.NewMemoryStorage()
	ksStorage := keystore.MemStorage([]byte{})
	keyStore, err := keystore.NewKeyStore(&ksStorage, keystore.LightKeyStoreParams)
	require.Nil(t, err)
	kOp, err := keyStore.NewKey(pass)
	require.Nil(t, err)
	err = keyStore.UnlockKey(kOp, pass)
	require.Nil(t, err)
	ho, err := holder.New(cfg, kOp, []claims.Claimer{}, storage, keyStore, idenPubOnChain, idenPubOffChainWrite, idenPubOffChainRead)
	require.Nil(t, err)
	return ho, storage, keyStore
}

func TestVerifyCredentialValidity(t *testing.T) {
	verifier := NewWithTimeNow(idenPubOnChain, func() time.Time {
		return time.Unix(blockTs, 0)
	})

	ho, _, _ := newHolder(t, idenPubOnChain, nil, idenPubOffChain)

	//
	// {Ts: 100, BlockN: 12} -> claim1 is added
	//
	blockTs, blockN = 100, 12

	// ISSUER: Publish state first time with claim1

	indexBytes, valueBytes := [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x42
	claim1 := claims.NewClaimBasic(indexBytes, valueBytes)

	is, _, _ := newIssuer(t, idenPubOnChain, idenPubOffChain)
	err := is.IssueClaim(claim1)
	require.Nil(t, err)

	// Publishing state for the first time
	err = is.PublishState()
	require.Nil(t, err)
	idenPubOnChain.Sync()

	err = is.SyncIdenStatePublic()
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
	blockTs, blockN = 200, 13

	// ISSUER: Publish state a second time with another claim2, claim3

	indexBytes, valueBytes = [claims.IndexSlotLen]byte{}, [claims.ValueSlotLen]byte{}
	indexBytes[0] = 0x48
	claim2 := claims.NewClaimBasic(indexBytes, valueBytes)

	err = is.IssueClaim(claim2)
	require.Nil(t, err)

	// claim3 is a claim with expiration at T=350

	header := claims.ClaimHeader{
		Type:       *claims.NewClaimTypeNum(9999),
		Dest:       claims.ClaimRecipSelf,
		Expiration: true,
		Version:    false,
	}
	metadata := claims.NewMetadata(header)
	metadata.Expiration = 350
	var entry merkletree.Entry
	metadata.Marshal(&entry)
	claim3 := claims.NewClaimGeneric(&entry)

	err = is.IssueClaim(claim3)
	require.Nil(t, err)

	err = is.PublishState()
	require.Nil(t, err)
	idenPubOnChain.Sync()

	err = is.SyncIdenStatePublic()
	require.Nil(t, err)

	credExistClaim2, err := is.GenCredentialExistence(claim2)
	require.Nil(t, err)
	credExistClaim3, err := is.GenCredentialExistence(claim3)
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
	blockTs, blockN = 300, 14

	// ISSUER: Publish state a third time revoking claim1

	err = is.RevokeClaim(claim1)
	require.Nil(t, err)

	err = is.PublishState()
	require.Nil(t, err)
	idenPubOnChain.Sync()

	err = is.SyncIdenStatePublic()
	require.Nil(t, err)

	err = verifier.VerifyCredentialExistence(credExistClaim1)
	assert.Nil(t, err)

	// HOLDER + VERIFIER

	_, err = ho.HolderGetCredentialValidity(credExistClaim1)
	assert.Equal(t, holder.ErrRevokedClaim, err)

	credValidClaim2t3, err := ho.HolderGetCredentialValidity(credExistClaim2)
	assert.Nil(t, err)

	credValidClaim3t3, err := ho.HolderGetCredentialValidity(credExistClaim3)
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

	// C3T3 has not expired at T=300 (expiration=350)
	err = verifier.VerifyCredentialValidity(credValidClaim3t3, 1000*time.Second)
	assert.Nil(t, err)

	//
	// {Ts: 400, BlockN: --}
	//
	blockTs, blockN = 400, 15

	// C3T3 has expired at T=400 (expiration=350)
	err = verifier.VerifyCredentialValidity(credValidClaim3t3, 1000*time.Second)
	assert.Equal(t, ErrClaimExpired, err)
}

func TestMain(m *testing.M) {
	idenPubOnChain = idenpubonchainlocal.New(
		func() time.Time {
			return time.Unix(blockTs, 0)
		},
		func() uint64 {
			return blockN
		},
	)
	idenPubOffChain = idenpuboffchanlocal.NewIdenPubOffChain("http://foo.bar")
	os.Exit(m.Run())
}
