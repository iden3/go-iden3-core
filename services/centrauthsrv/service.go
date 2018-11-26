package centrauthsrv

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/utils"
)

// Auth validates that the given AuthMsg data matches the requirments
func Auth(authMsg AuthMsg) error {
	err := VerifyChallengeTimestamp(authMsg.Challenge)
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
	proofOfKSign, err := authMsg.KSignProof.Unhex()
	if err != nil {
		return err
	}

	// TODO get the Relay address, now it's hardcoded, will be getted from the counterfactual contract of the Relay
	addrBytes, err := common3.HexToBytes("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	if err != nil {
		return err
	}
	relayAddr := common.BytesToAddress(addrBytes)
	if !claimsrv.CheckProofOfClaim(relayAddr, proofOfKSign, 140) { //TODO send the address of the Relay, to check the signature of proofOfKSign
		return errors.New("ProofOfKSign can not be verified")
	}

	// verify the Signature of the Challenge with the KSign
	msgHash := utils.EthHash([]byte(authMsg.Challenge))
	sigBytes[64] -= 27
	verified := utils.VerifySig(ksign, sigBytes, msgHash[:])
	if !verified {
		return errors.New("signature of challenge can not be verified")
	}

	return nil
}
