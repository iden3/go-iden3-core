package holder

import (
	"fmt"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

var (
	ErrRevokedClaim = fmt.Errorf("Revocation nonce exists in the Revocation Tree.  The claim is revoked.")
)

var ConfigDefault = Config{Config: issuer.ConfigDefault}

type Config struct {
	issuer.Config
}

func init() {
	ConfigDefault.GenesisOnly = true
}

// Holder is an identity that holds claims.  It is an extension of an Issuer.
type Holder struct {
	*issuer.Issuer
	idenPubOffChainReader idenpuboffchain.IdenPubOffChainReader
	idenPubOnChain        idenpubonchain.IdenPubOnChainer
}

// Create a new Holder, calling the internal Issuer.New().
func Create(cfg Config, kOpComp *babyjub.PublicKeyComp, extraGenesisClaims []claims.Claimer,
	storage db.Storage, keyStore *keystore.KeyStore) (*core.ID, error) {
	id, err := issuer.Create(cfg.Config, kOpComp, extraGenesisClaims, storage, keyStore)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// New creates a Holder by loading a previously created Holder (with New, and calling the internal Issuer.Load().
func Load(storage db.Storage, keyStore *keystore.KeyStore,
	idenPubOnChain idenpubonchain.IdenPubOnChainer,
	idenPubOffChainWriter idenpuboffchain.IdenPubOffChainWriter,
	idenPubOffChainReader idenpuboffchain.IdenPubOffChainReader) (*Holder, error) {
	is, err := issuer.Load(storage, keyStore, idenPubOnChain, idenPubOffChainWriter)
	if err != nil {
		return nil, err
	}
	return &Holder{
		Issuer:                is,
		idenPubOffChainReader: idenPubOffChainReader,
		idenPubOnChain:        idenPubOnChain,
	}, nil
}

// HolderGetCredentialValidity gets a Credential of Validity from a Credential
// of Existence.  This requires a request to the Issuer IdenStatePubOffChain.
func (h *Holder) HolderGetCredentialValidity(
	credExist *proof.CredentialExistence) (*proof.CredentialValidity, error) {
	idenStateData, err := h.idenPubOnChain.GetState(credExist.Id)
	if err != nil {
		return nil, err
	}
	publicData, err := h.idenPubOffChainReader.GetPublicData(credExist.IdenPubUrl, credExist.Id, idenStateData.IdenState)
	if err != nil {
		return nil, err
	}
	var claimMetadata claims.Metadata
	claimMetadata.Unmarshal(credExist.Claim)
	// NOTE: Once we add versions, this will require some changes that need to be thought properly!
	revLeaf := claims.NewLeafRevocationsTree(claimMetadata.RevNonce, 0xffffffff).Entry()
	revLeafHi, err := revLeaf.HIndex()
	if err != nil {
		return nil, err
	}
	mtpNotNonce, err := publicData.RevocationsTree.GenerateProof(revLeafHi, nil)
	if err != nil {
		return nil, err
	}
	if mtpNotNonce.Existence {
		return nil, ErrRevokedClaim
	}
	return &proof.CredentialValidity{
		CredentialExistence: *credExist,
		IdenStateData:       *idenStateData,
		MtpNotNonce:         mtpNotNonce,
		ClaimsTreeRoot:      publicData.ClaimsTreeRoot,
		RootsTreeRoot:       publicData.RootsTree.RootKey(),
	}, nil
}

// HolderImportCredentialExistence imports a received Credential of Existence into the ClaimsDB.
// func (h *Holder) HolderImportCredentialExistence(credExist *proof.ProofClaim) error {
// 	return fmt.Errorf("TODO: Implement ClaimDB")
// }
