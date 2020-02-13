package verifier

import (
	"fmt"
	"reflect"
	"time"

	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/merkletree"
)

var (
	ErrIdenStateOnChainDoesntMatch    = fmt.Errorf("IdenState on chain doesn't match the one in the credential")
	ErrMtpNonExistence                = fmt.Errorf("The Merkle Tree Proof is of non-existence")
	ErrMtpExistence                   = fmt.Errorf("The Merkle Tree Proof is of existence")
	ErrCalculatedIdenStateDoesntMatch = fmt.Errorf("Calculated IdenState doesn't match the one in the credential")
)

type Verifier struct {
	idenPubOnChain idenpubonchain.IdenPubOnChainer
	timeNow        func() time.Time
}

func New(idenPubOnChain idenpubonchain.IdenPubOnChainer) *Verifier {
	return &Verifier{
		idenPubOnChain: idenPubOnChain,
		timeNow: func() time.Time {
			return time.Now()
		},
	}
}

func NewWithTimeNow(idenPubOnChain idenpubonchain.IdenPubOnChainer, timeNow func() time.Time) *Verifier {
	return &Verifier{
		idenPubOnChain: idenPubOnChain,
		timeNow:        timeNow,
	}
}

func (v *Verifier) VerifyCredentialExistence(credExist *proof.CredentialExistence) error {
	if !credExist.MtpClaim.Existence {
		return ErrMtpNonExistence
	}
	// Verify that the idenState is built from claims merkle tree where the
	// claim exists.
	claimsRoot, err := merkletree.RootFromProof(credExist.MtpClaim, credExist.Claim.HIndex(), credExist.Claim.HValue())
	if err != nil {
		return err
	}
	idenState := core.IdenState(claimsRoot, credExist.RevocationsTreeRoot, credExist.RootsTreeRoot)
	if !idenState.Equals(credExist.IdenStateData.IdenState) {
		return ErrCalculatedIdenStateDoesntMatch
	}

	// Verify that the IdenStateData from the eistence credential is in the smart contract.
	idenStateDataOnChain, err := v.idenPubOnChain.GetStateByBlock(credExist.Id, credExist.IdenStateData.BlockN)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(idenStateDataOnChain, &credExist.IdenStateData) {
		return ErrIdenStateOnChainDoesntMatch
	}
	return nil
}

func (v *Verifier) VerifyCredentialValidity(credValid *proof.CredentialValidity, freshness time.Duration) error {
	if err := v.VerifyCredentialExistence(&credValid.CredentialExistence); err != nil {
		return err
	}
	if credValid.MtpNotNonce.Existence {
		return ErrMtpExistence
	}
	now := v.timeNow()
	// if now minus freshness is not a time before the validity credential
	// IdenState block ts, it means that the validity credential IdenState
	// may be too old!  This will be the case except for when the validity
	// credential IdenState is the last idenstate on chain.
	timeOldestAccepted := now.Add(-freshness)
	credentialTimestamp := time.Unix(credValid.IdenStateData.BlockTs, 0)
	if !timeOldestAccepted.Before(credentialTimestamp) {
		// Check if the last IdenState matches with the validity
		// credential IdenState.
		idenStateDataLast, err := v.idenPubOnChain.GetState(credValid.CredentialExistence.Id)
		if err != nil {
			return err
		}
		if !idenStateDataLast.IdenState.Equals(credValid.IdenStateData.IdenState) {
			return fmt.Errorf("Outdated validity credential.  validity credential IdenState timestamp is %v"+
				" Accepting IdenState only after timestamp %v", credentialTimestamp, timeOldestAccepted)
		}
	}
	// Verify that the idenState is built from revocations merkle tree
	// where the claim is not revoked (the revocation nonce is not a leaf).
	// NOTE: Once we add versions, this will require some changes that need to be thought properly!
	nonce := claims.GetRevocationNonce(credValid.CredentialExistence.Claim)
	revLeaf := claims.NewLeafRevocationsTree(nonce, 0xffffffff).Entry()
	revocationsTreeRoot, err := merkletree.RootFromProof(credValid.MtpNotNonce, revLeaf.HIndex(), revLeaf.HValue())
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
