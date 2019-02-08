package namesrv

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/identitysrv"
	"github.com/iden3/go-iden3/services/signsrv"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	VinculateId(vinculateIdMsg VinculateIdMsg) (*core.ClaimAssignName, error)
	ResolvClaimAssignName(nameid string) (*core.ClaimAssignName, error)
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

// VinculateId creates an adds a ClaimAssignName vinculating a name and an address, into the merkletree
func (ns *ServiceImpl) VinculateId(vinculateIdMsg VinculateIdMsg) (*core.ClaimAssignName, error) {
	// verify vinculateIdMsg.Msg signature with the Operational Key of the identity vinculateIdMsg.IdAddr
	// get the operational key
	fmt.Println(vinculateIdMsg.IdAddr)
	identity, err := ns.identitysrv.Get(vinculateIdMsg.IdAddr)
	if err != nil {
		fmt.Println("aaa")
		return &core.ClaimAssignName{}, err
	}
	opkey := identity.Operational

	if !utils.VerifySigEthMsg(opkey, vinculateIdMsg.Signature, vinculateIdMsg.Bytes()) {
		return &core.ClaimAssignName{}, errors.New("signature can not be verified")
	}

	// add ClaimAssignName to merkle tree
	assignNameClaim := core.NewClaimAssignName(vinculateIdMsg.Name, vinculateIdMsg.IdAddr)
	err = ns.claimsrv.AddClaimAssignName(*assignNameClaim)
	if err != nil {
		return &core.ClaimAssignName{}, err
	}

	return assignNameClaim, nil
}

// ResolvClaimAssignName returns the ClaimAssignName from the merkletree, given a nameid and a namespace
func (ns *ServiceImpl) ResolvClaimAssignName(nameid string) (*core.ClaimAssignName, error) {
	// get name and domain
	s := strings.Split(nameid, "@")
	if len(s) != 2 {
		return &core.ClaimAssignName{}, fmt.Errorf("Invalid nameid %v", s)
	}
	name := s[0]
	// domain := s[1]

	// build the ClaimAssignName Partial with the given data of the Index
	claimPartial := core.NewClaimAssignName(name, common.Address{})

	version, err := claimsrv.GetNextVersion(ns.claimsrv.MT(), claimPartial.Entry().HIndex())
	if err != nil {
		return &core.ClaimAssignName{}, err
	}
	claimPartial.Version = version - 1

	// get the complete ClaimAssignName in that merkle tree position
	leafDataInPos, err := ns.claimsrv.MT().GetDataByIndex(claimPartial.Entry().HIndex())
	if err != nil {
		return &core.ClaimAssignName{}, err
	}
	// if bytes.Equal(claimInPosBytes, merkletree.EmptyNodeValue[:]) {
	//         return core.ClaimAssignName{}, errors.New("not found")
	// }
	entry := &merkletree.Entry{
		Data: *leafDataInPos,
	}
	assignNameClaim := core.NewClaimAssignNameFromEntry(entry)
	// assignNameClaim, err := core.ParseClaimAssignNameBytes(claimInPosBytes)
	return assignNameClaim, nil
}
