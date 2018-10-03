package namesrv

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/cmd/id/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/utils"
)

// VinculateID creates an adds a AssignNameClaim vinculating a name and an address, into the merkletree
func VinculateID(cs claimsrv.Service, vinculateIDMsg VinculateIDMsg) (core.AssignNameClaim, error) {
	// TODO calculate EthID from the AssignNameClaim.RawIdentityTx
	ethID := common.HexToAddress(vinculateIDMsg.Msg.EthID) // tmp

	// verify vinculateIDMsg.Msg signature with EthID
	sigBytes, err := common3.HexToBytes(vinculateIDMsg.MsgSignature)
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	msgHash := vinculateIDMsg.MsgHash()
	verified := utils.VerifySig(ethID, sigBytes, msgHash[:])
	if !verified {
		return core.AssignNameClaim{}, errors.New("signature can not be verified")
	}

	// add AssignNameClaim to merkle tree
	nameHash := merkletree.HashBytes([]byte(vinculateIDMsg.Msg.Name))
	domainHash := merkletree.HashBytes([]byte(config.C.Domain))
	assignNameClaim := core.NewAssignNameClaim(config.C.Namespace, nameHash, domainHash, ethID)

	err = cs.AddAssignNameClaim(assignNameClaim)
	return assignNameClaim, err
}
