package verifier

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	zktypes "github.com/iden3/go-circom-prover-verifier/types"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-merkletree-sql"
)

var (
	ErrIdenStateOnChainDoesntMatch    = fmt.Errorf("IdenState on chain doesn't match the one in the credential")
	ErrMtpNonExistence                = fmt.Errorf("The Merkle Tree Proof is of non-existence")
	ErrMtpExistence                   = fmt.Errorf("The Merkle Tree Proof is of existence")
	ErrCalculatedIdenStateDoesntMatch = fmt.Errorf("Calculated IdenState doesn't match the one in the credential")
	ErrClaimExpired                   = fmt.Errorf("Expired claim")
	ErrFailedVerifyZkProofCredential  = fmt.Errorf("failed verifing generated zk proof of credential")
)

// Verifier allows verifying claims in three forms: credential of existence,
// credential of validity, and zero knowledge proof of valid credential
// ownership.
type Verifier struct {
	idenPubOnChain idenpubonchain.IdenPubOnChainer
	timeNow        func() time.Time
}

// NewWithTimeNow creates a verifier that uses the real time to validate freshness of claims.
func New(idenPubOnChain idenpubonchain.IdenPubOnChainer) *Verifier {
	return &Verifier{
		idenPubOnChain: idenPubOnChain,
		timeNow: func() time.Time {
			return time.Now()
		},
	}
}

// NewWithTimeNow creates a verifier that uses the time returned by `timeNow`
// to validate freshness of claims.  Mainly used for testing (by faking the
// current time).
func NewWithTimeNow(idenPubOnChain idenpubonchain.IdenPubOnChainer, timeNow func() time.Time) *Verifier {
	return &Verifier{
		idenPubOnChain: idenPubOnChain,
		timeNow:        timeNow,
	}
}

// VerifyCredentialExistence verifies a credential of existence.  That is, that
// the claim was issued by a particular identity.
func (v *Verifier) VerifyCredentialExistence(credExist *proof.CredentialExistence) error {
	if !credExist.MtpClaim.Existence {
		return ErrMtpNonExistence
	}
	// Verify that the idenState is built from claims merkle tree where the
	// claim exists.
	hi, hv, err := credExist.Claim.HiHv()
	if err != nil {
		return err
	}
	claimsRoot, err := merkletree.RootFromProof(credExist.MtpClaim, hi.BigInt(), hv.BigInt())
	if err != nil {
		return err
	}
	idenState := core.IdenState(claimsRoot, credExist.RevocationsTreeRoot, credExist.RootsTreeRoot)
	if !idenState.Equals(credExist.IdenStateData.IdenState) {
		return ErrCalculatedIdenStateDoesntMatch
	}

	// Verify that the IdenStateData from the existence credential is in the smart contract.
	idenStateDataOnChain, err := v.idenPubOnChain.GetStateByBlock(credExist.Id, credExist.IdenStateData.BlockN)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(idenStateDataOnChain, &credExist.IdenStateData) {
		return ErrIdenStateOnChainDoesntMatch
	}
	return nil
}

// validateFreshness is a helper function that validates that the passed
// `idenState` is not older than `freshness`, or that it's the most recent one.
// The link between `idenState` and `blockTs` is not checked here.
func (v *Verifier) validateFreshness(id *core.ID, idenState *merkletree.Hash, blockTs int64, freshness time.Duration) error {
	// if now minus freshness is not a time before the validity credential
	// IdenState block ts, it means that the validity credential IdenState
	// may be too old!  This will be the case except for when the validity
	// credential IdenState is the last idenstate on chain.
	now := v.timeNow()
	timeOldestAccepted := now.Add(-freshness)
	credentialTimestamp := time.Unix(blockTs, 0)
	if !timeOldestAccepted.Before(credentialTimestamp) {
		// Check if the last IdenState matches with the validity
		// credential IdenState.
		idenStateDataLast, err := v.idenPubOnChain.GetState(id)
		if err != nil {
			return err
		}
		if !idenStateDataLast.IdenState.Equals(idenState) {
			return fmt.Errorf("Outdated validity credential.  validity credential IdenState timestamp is %v"+
				" Accepting IdenState only after timestamp %v", credentialTimestamp, timeOldestAccepted)
		}
	}
	return nil
}

// VerifyCredentialValidity verifies a credential of validity.  That is, that
// the claim was issued by a particular identity, and that it has not been
// revoked passed the freshness time.
func (v *Verifier) VerifyCredentialValidity(credValid *proof.CredentialValidity, freshness time.Duration) error {
	// If the claim has an expiration date, check that it hasn't expired.
	var metadata claims.Metadata
	metadata.Unmarshal(credValid.CredentialExistence.Claim)
	if metadata.Header().Expiration {
		now := v.timeNow()
		if time.Unix(metadata.Expiration, 0).Before(now) {
			return ErrClaimExpired
		}
	}
	if err := v.VerifyCredentialExistence(&credValid.CredentialExistence); err != nil {
		return err
	}
	if credValid.MtpNotNonce.Existence {
		return ErrMtpExistence
	}
	if err := v.validateFreshness(credValid.CredentialExistence.Id,
		credValid.IdenStateData.IdenState,
		credValid.IdenStateData.BlockTs,
		freshness); err != nil {
		return err
	}
	// Verify that the idenState is built from revocations merkle tree
	// where the claim is not revoked (the revocation nonce is not a leaf).
	// NOTE: Once we add versions, this will require some changes that need to be thought properly!
	nonce := claims.GetRevocationNonce(credValid.CredentialExistence.Claim)
	revLeaf := claims.NewLeafRevocationsTree(nonce, 0xffffffff).Entry()
	hi, hv, err := revLeaf.HiHv()
	if err != nil {
		return err
	}
	revocationsTreeRoot, err := merkletree.RootFromProof(credValid.MtpNotNonce, hi.BigInt(), hv.BigInt())
	if err != nil {
		return err
	}
	idenState := core.IdenState(credValid.ClaimsTreeRoot, revocationsTreeRoot, credValid.RootsTreeRoot)
	if !idenState.Equals(credValid.IdenStateData.IdenState) {
		return ErrCalculatedIdenStateDoesntMatch
	}

	// Verify that the IdenStateData from the validity credential is in the smart contract.
	idenStateDataOnChain, err := v.idenPubOnChain.GetStateByBlock(credValid.CredentialExistence.Id, credValid.IdenStateData.BlockN)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(idenStateDataOnChain, &credValid.IdenStateData) {
		return ErrIdenStateOnChainDoesntMatch
	}
	return nil
}

// VerifyZkProofCredential verifies a zkp of a credential. For now expiration
// is not checked.
func (v *Verifier) VerifyZkProofCredential(
	zkProof *zktypes.Proof,
	pubSignals []*big.Int,
	issuerID *core.ID,
	idenStateBlockN uint64,
	zkFiles *zkutils.ZkFiles,
	freshness time.Duration) error {

	vk, err := zkFiles.VerificationKey()
	if err != nil {
		return fmt.Errorf("error loading zk vk: %w", err)
	}

	// Verify the zkp
	if !verifier.Verify(vk, zkProof, pubSignals) {
		return ErrFailedVerifyZkProofCredential
	}

	// Verify that the IdenState used in the proof corresponds to the
	// issuerID at idenStateBlockN in the smart contract.
	idenState := merkletree.NewHashFromBigInt(pubSignals[0])
	idenStateDataOnChain, err := v.idenPubOnChain.GetStateByBlock(issuerID, idenStateBlockN)
	if err != nil {
		return err
	}
	if idenStateDataOnChain.BlockN != idenStateBlockN ||
		!idenStateDataOnChain.IdenState.Equals(idenState) {
		return ErrIdenStateOnChainDoesntMatch
	}

	if err := v.validateFreshness(issuerID,
		idenState,
		idenStateDataOnChain.BlockTs,
		freshness); err != nil {
		return err
	}
	return nil
}
