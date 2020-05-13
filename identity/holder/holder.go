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

const (
	ErrStrTooManySiblings = "number of mtp siblings in %v (%v) is higher than requested levels (%v)"
)

var (
	ErrRevokedClaim = fmt.Errorf("revocation nonce exists in the Revocation Tree.  The claim is revoked.")
	ErrRootNotFound = fmt.Errorf("claims tree root not found in roots tree.")
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
	PublicData     *idenpuboffchain.PublicData
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
		IdenStateData:  idenStateData,
		MtpNotNonce:    mtpNotNonce,
		ClaimsTreeRoot: publicData.ClaimsTreeRoot,
		RevTreeRoot:    publicData.RevocationsTree.RootKey(),
		RootsTreeRoot:  publicData.RootsTree.RootKey(),
		PublicData:     publicData,
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
	Claim []*big.Int `mapstructure:"claim"`

	// B. holder proof of claimKOp in the genesis
	PrivateKey             *big.Int   `mapstructure:"hoKOpSk"`
	ClaimKOpMtp            []*big.Int `mapstructure:"hoClaimKOpMtp"`
	ClaimKOpClaimsTreeRoot *big.Int   `mapstructure:"hoClaimKOpClaimsTreeRoot"`

	// C. issuer proof of claim existence
	CredExistMtp            []*big.Int `mapstructure:"isProofExistMtp"`
	CredExistClaimsTreeRoot *big.Int   `mapstructure:"isProofExistClaimsTreeRoot"`

	// D. issuer proof of claim validity
	CredValidNotRevMtp      []*big.Int `mapstructure:"isProofValidNotRevMtp"`
	CredValidNotRevMtpNoAux *big.Int   `mapstructure:"isProofValidNotRevMtpNoAux"`
	CredValidNotRevMtpAuxHi *big.Int   `mapstructure:"isProofValidNotRevMtpAuxHi"`
	CredValidNotRevMtpAuxHv *big.Int   `mapstructure:"isProofValidNotRevMtpAuxHv"`
	CredValidClaimsTreeRoot *big.Int   `mapstructure:"isProofValidClaimsTreeRoot"`
	CredValidRevTreeRoot    *big.Int   `mapstructure:"isProofValidRevTreeRoot"`
	CredValidRootsTreeRoot  *big.Int   `mapstructure:"isProofValidRootsTreeRoot"`

	// E. issuer proof of Root (ExistClaimsTreeRoot)
	CredValidRootMtp []*big.Int `mapstructure:"isProofRootMtp"`

	// F. issuer recent idenState
	IdenState *big.Int `mapstructure:"isIdenState"`
}

func (h *Holder) HolderGetCredentialProofInputs(
	idOwnershipGenesisInputs *issuer.IdOwnershipGenesisInputs,
	credExist *proof.CredentialExistence, issuerLevels int) (*CredentialProofInputs, error) {
	credValidData, err := h.HolderGetCredentialValidityData(credExist)
	if err != nil {
		return nil, err
	}
	hi, err := credExist.Claim.HIndex()
	if err != nil {
		return nil, err
	}
	hv, err := credExist.Claim.HValue()
	if err != nil {
		return nil, err
	}
	credExistClaimsTreeRoot, err := merkletree.RootFromProof(credExist.MtpClaim, hi, hv)
	if err != nil {
		return nil, err
	}

	var claimBigInts [8]*big.Int
	for i, elem := range credExist.Claim.Data {
		claimBigInts[i] = elem.BigInt()
	}

	rootLeaf := claims.NewLeafRootsTree(*credExistClaimsTreeRoot).Entry()
	rootLeafHi, err := rootLeaf.HIndex()
	if err != nil {
		return nil, err
	}
	mtpRoot, err := credValidData.PublicData.RootsTree.GenerateProof(rootLeafHi, nil)
	if err != nil {
		return nil, err
	}
	if !mtpRoot.Existence {
		return nil, ErrRootNotFound
	}

	credValidNotRevMtpNoAux := new(big.Int).SetUint64(1) // TODO: Confirm this
	credValidNotRevMtpAuxHi := new(big.Int)
	credValidNotRevMtpAuxHv := new(big.Int)
	if credValidData.MtpNotNonce.NodeAux != nil {
		credValidNotRevMtpNoAux = new(big.Int)
		credValidNotRevMtpAuxHi = credValidData.MtpNotNonce.NodeAux.HIndex.BigInt()
		credValidNotRevMtpAuxHv = credValidData.MtpNotNonce.NodeAux.HValue.BigInt()
	}

	credExistMtp := credExist.MtpClaim.AllSiblingsCircom(issuerLevels)
	if len(credExistMtp) != issuerLevels+1 {
		return nil, fmt.Errorf(ErrStrTooManySiblings, "ClaimTree", len(credExistMtp), issuerLevels+1)
	}
	credValidNotRevMtp := credValidData.MtpNotNonce.AllSiblingsCircom(issuerLevels)
	if len(credValidNotRevMtp) != issuerLevels+1 {
		return nil, fmt.Errorf(ErrStrTooManySiblings, "RevTree", len(credValidNotRevMtp), issuerLevels+1)
	}
	credValidRootMtp := mtpRoot.AllSiblingsCircom(issuerLevels)
	if len(credValidRootMtp) != issuerLevels+1 {
		return nil, fmt.Errorf(ErrStrTooManySiblings, "RootsTree", len(credValidRootMtp), issuerLevels+1)
	}

	return &CredentialProofInputs{
		PrivateKey:             idOwnershipGenesisInputs.PrivateKey,
		ClaimKOpMtp:            idOwnershipGenesisInputs.MtpSiblings,
		ClaimKOpClaimsTreeRoot: idOwnershipGenesisInputs.ClaimsTreeRoot,

		Claim:                   claimBigInts[:],
		CredExistMtp:            credExistMtp,
		CredExistClaimsTreeRoot: credExistClaimsTreeRoot.BigInt(),

		CredValidNotRevMtp:      credValidNotRevMtp,
		CredValidNotRevMtpNoAux: credValidNotRevMtpNoAux,
		CredValidNotRevMtpAuxHi: credValidNotRevMtpAuxHi,
		CredValidNotRevMtpAuxHv: credValidNotRevMtpAuxHv,

		CredValidClaimsTreeRoot: credValidData.ClaimsTreeRoot.BigInt(),
		CredValidRevTreeRoot:    credValidData.RevTreeRoot.BigInt(),
		CredValidRootsTreeRoot:  credValidData.RootsTreeRoot.BigInt(),

		CredValidRootMtp: credValidRootMtp,

		IdenState: credValidData.IdenStateData.IdenState.BigInt(),
	}, nil
}

func (h *Holder) HolderGenZkProofCredential(
	credExist *proof.CredentialExistence,
	addInputs func(inputs map[string]interface{}) error,
	idOwnershipLevels, issuerLevels int,
	circuitWASMPath, provingKeyPath string) (*int, error) {
	return nil, fmt.Errorf("TODO")
}

// HolderImportCredentialExistence imports a received Credential of Existence into the ClaimsDB.
// func (h *Holder) HolderImportCredentialExistence(credExist *proof.ProofClaim) error {
// 	return fmt.Errorf("TODO: Implement ClaimDB")
// }
