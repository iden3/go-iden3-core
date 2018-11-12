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
	"github.com/iden3/go-iden3/services/identitysrv"
	"github.com/iden3/go-iden3/services/signsrv"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	VinculateID(vinculateIDMsg VinculateIDMsg) (core.AssignNameClaim, error)
	ResolvAssignNameClaim(nameid string) (core.AssignNameClaim, error)
}

type ServiceImpl struct {
	claimsrv    claimsrv.Service
	identitysrv identitysrv.Service
	signer      signsrv.Service
	domain      string
}

func New(claimsrv claimsrv.Service, identitysrv identitysrv.Service, signer signsrv.Service, domain string) *ServiceImpl {
	return &ServiceImpl{claimsrv, identitysrv, signer, domain}
}

// VinculateID creates an adds a AssignNameClaim vinculating a name and an address, into the merkletree
func (ns *ServiceImpl) VinculateID(vinculateIDMsg VinculateIDMsg) (core.AssignNameClaim, error) {
	// verify vinculateIDMsg.Msg signature with the Operational Key of the identity vinculateIDMsg.EthID
	// get the operational key
	identity, err := ns.identitysrv.Get(vinculateIDMsg.EthID)
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	opkey := identity.Operational

	sigBytes, err := common3.HexToBytes(vinculateIDMsg.Signature)
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	msgHash := vinculateIDMsg.MsgHash()
	sigBytes[64] -= 27
	verified := utils.VerifySig(opkey, sigBytes, msgHash[:])
	if !verified {
		return core.AssignNameClaim{}, errors.New("signature can not be verified")
	}

	// add AssignNameClaim to merkle tree
	nameHash := merkletree.HashBytes([]byte(vinculateIDMsg.Name))
	domainHash := merkletree.HashBytes([]byte(ns.domain))
	assignNameClaim := core.NewAssignNameClaim(nameHash, domainHash, vinculateIDMsg.EthID)
	err = ns.claimsrv.AddAssignNameClaim(assignNameClaim)
	if err != nil {
		return core.AssignNameClaim{}, err
	}

	return assignNameClaim, nil
}

// ResolvAssignNameClaim returns the AssignNameClaim from the merkletree, given a nameid and a namespace
func (ns *ServiceImpl) ResolvAssignNameClaim(nameid string) (core.AssignNameClaim, error) {
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
	claimPartial := core.NewAssignNameClaim(nameHash, domainHash, common.Address{})

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
