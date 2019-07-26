package namesrv

import (
	//"errors"
	"fmt"
	"strings"

	//"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/claimsrv"
	"github.com/iden3/go-iden3-core/services/signsrv"
	//"github.com/iden3/go-iden3-core/utils"
)

type Service interface {
	VinculateId(name string, domain string, id core.ID) (*core.ClaimAssignName, error)
	ResolvClaimAssignName(name string) (*core.ClaimAssignName, error)
}

type ServiceImpl struct {
	claimsrv claimsrv.Service
	signer   signsrv.Service
	domain   string
}

func New(claimsrv claimsrv.Service, signer signsrv.Service, domain string) *ServiceImpl {
	return &ServiceImpl{claimsrv, signer, domain}
}

// VinculateId creates an adds a ClaimAssignName vinculating a name and an id, into the merkletree
func (ns *ServiceImpl) VinculateId(name string, domain string,
	id core.ID) (*core.ClaimAssignName, error) {
	if name == "" {
		return nil, fmt.Errorf("Name is empty")
	}
	if strings.Contains(name, "@") {
		return nil, fmt.Errorf("Name contains a '@' character")
	}
	// add ClaimAssignName to merkle tree
	assignNameClaim := core.NewClaimAssignName(fmt.Sprintf("%v@%v", name, domain), id)
	if err := ns.claimsrv.AddClaim(assignNameClaim); err != nil {
		return nil, err
	}
	return assignNameClaim, nil
}

// ResolvClaimAssignName returns the ClaimAssignName from the merkletree, given a name
func (ns *ServiceImpl) ResolvClaimAssignName(name string) (*core.ClaimAssignName, error) {
	// build the ClaimAssignName Partial with the given data of the Index
	claimPartial := core.NewClaimAssignName(name, core.ID{})

	version, err := claimsrv.GetNextVersion(ns.claimsrv.MT(), claimPartial.Entry().HIndex())
	if err != nil {
		return nil, err
	}
	claimPartial.Version = version - 1

	// get the complete ClaimAssignName in that merkle tree position
	leafDataInPos, err := ns.claimsrv.MT().GetDataByIndex(claimPartial.Entry().HIndex())
	if err != nil {
		return nil, err
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
