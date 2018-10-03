package centrauthsrv

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

func Auth(authMsg AuthMsg) error {
	err := VerifyTimestamp(authMsg.Challenge)
	if err != nil {
		return err
	}
	// addr := common.HexToAddress(authMsg.Address)
	ksign := common.HexToAddress(authMsg.KSign)

	sigBytes, err := common3.HexToBytes(authMsg.Signature)
	if err != nil {
		return err
	}

	// check the authMsg.KSignProof
	ksignProof, err := authMsg.KSignProof.Unhex()
	if err != nil {
		return err
	}
	verified := merkletree.CheckProof(ksignProof.ClaimProof.Root, ksignProof.ClaimProof.Proof, ksignProof.ClaimProof.Hi, merkletree.HashBytes(ksignProof.ClaimProof.Leaf), 140)
	if !verified {
		return errors.New("ksignProof.ClaimProof failed")
	}
	verified = merkletree.CheckProof(ksignProof.SetRootClaimProof.Root, ksignProof.SetRootClaimProof.Proof, ksignProof.SetRootClaimProof.Hi, merkletree.HashBytes(ksignProof.SetRootClaimProof.Leaf), 140)
	if !verified {
		return errors.New("ksignProof.SetRootClaimProof failed")
	}
	verified = merkletree.CheckProof(ksignProof.ClaimNonRevocationProof.Root, ksignProof.ClaimNonRevocationProof.Proof, ksignProof.ClaimNonRevocationProof.Hi, merkletree.EmptyNodeValue, 140)
	if !verified {
		return errors.New("ksignProof.ClaimNonRevocationProof failed")
	}
	verified = merkletree.CheckProof(ksignProof.SetRootClaimNonRevocationProof.Root, ksignProof.SetRootClaimNonRevocationProof.Proof, ksignProof.SetRootClaimNonRevocationProof.Hi, merkletree.EmptyNodeValue, 140)
	if !verified {
		return errors.New("ksignProof.SetRootClaimNonRevocationProof failed")
	}

	// verify the Signature of the Challenge with the KSign
	msgHash := utils.EthHash([]byte(authMsg.Challenge))
	sigBytes[64] -= 27
	verified = utils.VerifySig(ksign, sigBytes, msgHash[:])
	if !verified {
		return errors.New("signature error")
	}

	//

	return nil
}
