package holder

import (
	"fmt"
	"math/big"
	"time"

	"github.com/iden3/go-circom-prover-verifier/prover"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	witnesscalc "github.com/iden3/go-circom-witnesscalc"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/issuer"
	"github.com/iden3/go-iden3-core/keystore"
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

const (
	ErrStrTooManySiblings = "number of mtp siblings in %v (%v) is higher than requested levels (%v)"
)

var (
	ErrRevokedClaim                   = fmt.Errorf("revocation nonce exists in the Revocation Tree.  The claim is revoked.")
	ErrRootNotFound                   = fmt.Errorf("claims tree root not found in roots tree.")
	ErrFailedVerifyZkProofCredential  = fmt.Errorf("failed verifing generated zk proof of credential")
	ErrCalculatedIdenStateDoesntMatch = fmt.Errorf("Calculated IdenState from public data doesn't match the one queried")
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

// CredentialValidityAux contains the data used in a validity proof.
type CredentialValidityAux struct {
	IdenStateData  *proof.IdenStateData
	MtpNotNonce    *merkletree.Proof
	ClaimsTreeRoot *merkletree.Hash
	RevTreeRoot    *merkletree.Hash
	RootsTreeRoot  *merkletree.Hash
	PublicData     *idenpuboffchain.PublicData
}

// HolderGetCredentialValidityData is a helper function to get the data used in
// a validity proof from a credential existence proof.
func (h *Holder) HolderGetCredentialValidityData(
	credExist *proof.CredentialExistence) (*CredentialValidityAux, error) {
	idenStateData, err := h.idenPubOnChain.GetState(credExist.Id)
	if err != nil {
		return nil, err
	}
	log.WithField("state", idenStateData.IdenState).Debug("Holder.idenPubOnChain.GetState()")
	publicData, err := h.idenPubOffChainReader.GetPublicData(credExist.IdenPubUrl, credExist.Id, idenStateData.IdenState)
	if err != nil {
		return nil, err
	}

	// Verify that the returned public data is consistent with the queried IdenState
	idenState := core.IdenState(publicData.ClaimsTreeRoot, publicData.RevocationsTree.Root(),
		publicData.RootsTree.Root())
	if !idenState.Equals(idenStateData.IdenState) {
		return nil, ErrCalculatedIdenStateDoesntMatch
	}

	var claimMetadata claims.Metadata
	claimMetadata.Unmarshal(credExist.Claim)
	// NOTE: Once we add versions, this will require some changes that need to be thought properly!
	revLeaf := claims.NewLeafRevocationsTree(claimMetadata.RevNonce, 0xffffffff).Entry()
	revLeafHi, err := revLeaf.HIndex()
	if err != nil {
		return nil, err
	}
	mtpNotNonce, _, err := publicData.RevocationsTree.GenerateProof(revLeafHi, nil)
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
		RevTreeRoot:    publicData.RevocationsTree.Root(),
		RootsTreeRoot:  publicData.RootsTree.Root(),
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

// CredentialProofInputs are all the iinputs for the credential ownership proof
// `credential.circom`.
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

// HolderGetCredentialProofInputs generates the inputs for the credential
// ownership proof `credential.circom`.
func (h *Holder) HolderGetCredentialProofInputs(
	idOwnershipGenesisInputs *issuer.IdOwnershipGenesisInputs,
	credExist *proof.CredentialExistence,
	credValidData *CredentialValidityAux,
	issuerLevels int) (*CredentialProofInputs, error) {
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
	mtpRoot, _, err := credValidData.PublicData.RootsTree.GenerateProof(rootLeafHi.BigInt(), nil)
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

// ZkProofCredOut is the data output of a generated credential zkp,
// and contains the inputs required for verification of a credential zkp.
type ZkProofCredOut struct {
	ZkProofOut      zkutils.ZkProofOut
	IssuerID        *core.ID
	IdenStateBlockN uint64
}

// HolderGenZkProofCredential generates a zkp of a credential.  This function
// prepares all the inputs of the `credential.circom` circuit and removes the
// "claim" input.  The `addInputs` function allows adding circuit inputs as
// necessary (for example, inputs used to build the claim).
func (h *Holder) HolderGenZkProofCredential(
	credExist *proof.CredentialExistence,
	addInputs func(inputs map[string]interface{}) error,
	idOwnershipLevels, issuerLevels int,
	zkFiles *zkutils.ZkFiles) (*ZkProofCredOut, error) {

	pk, err := zkFiles.ProvingKey()
	if err != nil {
		return nil, fmt.Errorf("error loading zk pk: %w", err)
	}
	vk, err := zkFiles.VerificationKey()
	if err != nil {
		return nil, fmt.Errorf("error loading zk vk: %w", err)
	}
	witnessCalcWASM, err := zkFiles.WitnessCalcWASM()
	if err != nil {
		return nil, fmt.Errorf("error loading zk witnessCalc WASM: %w", err)
	}

	idOwnershipInputs, err := h.GenIdOwnershipGenesisInputs(idOwnershipLevels)
	if err != nil {
		return nil, err
	}
	credValidData, err := h.HolderGetCredentialValidityData(credExist)
	if err != nil {
		return nil, err
	}
	credProofInputs, err := h.HolderGetCredentialProofInputs(idOwnershipInputs,
		credExist, credValidData, issuerLevels)
	if err != nil {
		return nil, err
	}

	var inputs map[string]interface{}
	if err := mapstructure.Decode(credProofInputs, &inputs); err != nil {
		return nil, err
	}
	delete(inputs, "claim")
	if err := addInputs(inputs); err != nil {
		return nil, err
	}

	wit, err := witnesscalc.CalculateWitnessBinWASM(witnessCalcWASM, inputs)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	proof, pubSignals, err := prover.GenerateProof(pk, wit)
	if err != nil {
		return nil, err
	}
	// Verify zk proof
	if !verifier.Verify(vk, proof, pubSignals) {
		return nil, ErrFailedVerifyZkProofCredential
	}

	log.WithField("elapsed", time.Since(start)).Debug("Proof generated")
	return &ZkProofCredOut{
		ZkProofOut:      zkutils.ZkProofOut{Proof: *proof, PubSignals: pubSignals},
		IssuerID:        credExist.Id,
		IdenStateBlockN: credValidData.IdenStateData.BlockN,
	}, nil
}

// HolderImportCredentialExistence imports a received Credential of Existence into the ClaimsDB.
// func (h *Holder) HolderImportCredentialExistence(credExist *proof.ProofClaim) error {
// 	return fmt.Errorf("TODO: Implement ClaimDB")
// }
