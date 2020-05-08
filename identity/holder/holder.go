package holder

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
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
	idenStateZkProofConf *issuer.IdenStateZkProofConf,
	idenPubOffChainWriter idenpuboffchain.IdenPubOffChainWriter,
	idenPubOffChainReader idenpuboffchain.IdenPubOffChainReader) (*Holder, error) {
	is, err := issuer.Load(storage, keyStore, idenPubOnChain, idenStateZkProofConf, idenPubOffChainWriter)
	if err != nil {
		return nil, err
	}
	return &Holder{
		Issuer:                is,
		idenPubOffChainReader: idenPubOffChainReader,
		idenPubOnChain:        idenPubOnChain,
	}, nil
}

type CredentialValidityAux struct {
	IdenStateData  *proof.IdenStateData
	MtpNotNonce    *merkletree.Proof
	ClaimsTreeRoot *merkletree.Hash
	RevTreeRoot    *merkletree.Hash
	RootsTreeRoot  *merkletree.Hash
}

func (h *Holder) HolderGetCredentialValidityData(
	credExist *proof.CredentialExistence) (*CredentialValidityAux, error) {
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
	return &CredentialValidityAux{
		MtpNotNonce:    mtpNotNonce,
		ClaimsTreeRoot: publicData.ClaimsTreeRoot,
		RevTreeRoot:    publicData.RevocationsTree.RootKey(),
		RootsTreeRoot:  publicData.RootsTree.RootKey(),
	}, nil
}

// HolderGetCredentialValidity gets a Credential of Validity from a Credential
// of Existence.  This requires a request to the Issuer IdenStatePubOffChain.
func (h *Holder) HolderGetCredentialValidity(
	credExist *proof.CredentialExistence) (*proof.CredentialValidity, error) {
	credValidData, err := h.HolderGetCredentialValidityData(credExist)
	if err != nil {
		return nil, err
	}
	return &proof.CredentialValidity{
		CredentialExistence: *credExist,
		IdenStateData:       *credValidData.IdenStateData,
		MtpNotNonce:         credValidData.MtpNotNonce,
		ClaimsTreeRoot:      credValidData.ClaimsTreeRoot,
		RootsTreeRoot:       credValidData.RootsTreeRoot,
	}, nil
}

type CredentialProofInputs struct {
	// A
	Claim [8]*big.Int

	CredExistMtp            []*big.Int
	CredExistClaimsTreeRoot *big.Int

	// D. issuer proof of claim validity
	CredValidMtp            []*big.Int
	CredValidClaimsTreeRoot *big.Int
	CredValidRevTreeRoot    *big.Int
	CredValidRootsTreeRoot  *big.Int

	// E. issuer proof of Root (ExistClaimsTreeRoot)
	CredValidRootMtp []*big.Int

	// F. issuer recent idenState
	IdenState *big.Int
}

func (h *Holder) HolderGetCredentialProofInputs(
	credExist *proof.CredentialExistence, issuerLevels int) (*CredentialProofInputs, error) {
	credValidData, err := h.HolderGetCredentialValidityData(credExist)
	if err != nil {
		return nil, err
	}
	// TODO: Compute ExistClaimTreeRoot
	// TODO: Compute RootMtp

	var claim [8]*big.Int
	for i, elem := range credExist.Claim.Data {
		claim[i] = elem.BigInt()
	}

	return &CredentialProofInputs{
		Claim:                   claim,
		CredExistMtp:            credExist.MtpClaim.AllSiblingsCircom(issuerLevels),
		CredExistClaimsTreeRoot: nil,

		CredValidMtp:            credValidData.MtpNotNonce.AllSiblingsCircom(issuerLevels),
		CredValidClaimsTreeRoot: credValidData.ClaimsTreeRoot.BigInt(),
		CredValidRevTreeRoot:    nil,
		CredValidRootsTreeRoot:  credValidData.RootsTreeRoot.BigInt(),

		CredValidRootMtp: nil,

		IdenState: credValidData.IdenStateData.IdenState.BigInt(),
	}, nil
}

// HolderImportCredentialExistence imports a received Credential of Existence into the ClaimsDB.
// func (h *Holder) HolderImportCredentialExistence(credExist *proof.ProofClaim) error {
// 	return fmt.Errorf("TODO: Implement ClaimDB")
// }
