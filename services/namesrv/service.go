package namesrv

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/cmd/id/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/rootsrv"
	"github.com/iden3/go-iden3/services/signsrv"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	VinculateID(vinculateIDMsg VinculateIDMsg) (core.AssignNameClaim, error)
}

type ServiceImpl struct {
	mt       *merkletree.MerkleTree
	rootsrv  rootsrv.Service
	claimsrv claimsrv.Service
	signer   signsrv.Service
}

func New(mt *merkletree.MerkleTree, rootsrv rootsrv.Service, claimsrv claimsrv.Service, signer signsrv.Service) *ServiceImpl {
	return &ServiceImpl{mt, rootsrv, claimsrv, signer}
}

// VinculateID creates an adds a AssignNameClaim vinculating a name and an address, into the merkletree
func (ns *ServiceImpl) VinculateID(vinculateIDMsg VinculateIDMsg) (core.AssignNameClaim, error) {
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

	err = ns.claimsrv.AddAssignNameClaim(assignNameClaim)
	return assignNameClaim, err
}
