package namesrv

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
	ResolvAssignNameClaim(nameid, namespace string) (core.AssignNameClaim, error)
}

type ServiceImpl struct {
	claimsrv  claimsrv.Service
	signer    signsrv.Service
	domain    string
	namespace string
}

func New(rootsrv rootsrv.Service, claimsrv claimsrv.Service, signer signsrv.Service, domain string, namespace string) *ServiceImpl {
	return &ServiceImpl{claimsrv, signer, domain, namespace}
}

// VinculateID creates an adds a AssignNameClaim vinculating a name and an address, into the merkletree
func (ns *ServiceImpl) VinculateID(vinculateIDMsg VinculateIDMsg) (core.AssignNameClaim, error) {
	// verify vinculateIDMsg.Msg signature with EthID
	sigBytes, err := common3.HexToBytes(vinculateIDMsg.Signature)
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	msgHash := vinculateIDMsg.MsgHash()
	sigBytes[64] -= 27
	verified := utils.VerifySig(vinculateIDMsg.EthID, sigBytes, msgHash[:])
	if !verified {
		return core.AssignNameClaim{}, errors.New("signature can not be verified")
	}

	// add AssignNameClaim to merkle tree
	nameHash := merkletree.HashBytes([]byte(vinculateIDMsg.Name))
	domainHash := merkletree.HashBytes([]byte(ns.domain))
	assignNameClaim := core.NewAssignNameClaim(ns.namespace, nameHash, domainHash, vinculateIDMsg.EthID)
	err = ns.claimsrv.AddAssignNameClaim(assignNameClaim)
	if err != nil {
		return core.AssignNameClaim{}, err
	}

	return assignNameClaim, nil
}

// ResolvAssignNameClaim returns the AssignNameClaim from the merkletree, given a nameid and a namespace
func (ns *ServiceImpl) ResolvAssignNameClaim(nameid, namespace string) (core.AssignNameClaim, error) {
	// get name and domain
	s := strings.Split(nameid, "@")
	if len(s) != 2 {
		return core.AssignNameClaim{}, fmt.Errorf("Invalid nameid %v", s)
	}
	name := s[0]
	domain := s[1]

	// build the AssignNameClaim Partial with the given data of the Index
	nameHash := merkletree.HashBytes([]byte(name))
	domainHash := merkletree.HashBytes([]byte(domain))
	claimPartial := core.NewAssignNameClaim(ns.namespace, nameHash, domainHash, common.Address{})

	version, err := claimsrv.GetNextVersion(ns.claimsrv.MT(), claimPartial.Hi())
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	claimPartial.BaseIndex.Version = version - 1

	// get the complete AssignNameClaim in that merkle tree position
	claimInPosBytes, err := ns.claimsrv.MT().GetValueInPos(claimPartial.Hi())
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	if bytes.Equal(claimInPosBytes, merkletree.EmptyNodeValue[:]) {
		return core.AssignNameClaim{}, errors.New("not found")
	}
	assignNameClaim, err := core.ParseAssignNameClaimBytes(claimInPosBytes)
	return assignNameClaim, err
}
