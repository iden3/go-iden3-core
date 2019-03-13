package discoverysrv

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/utils"
)

type Trusted struct {
	Relay bool `json:"relay"`
}

type Identity struct {
	IdAddr          common.Address   `json:"-"`
	Name            string           `json:"name"`
	OperationalPk   *utils.PublicKey `json:"kOpPub"`
	OperationalAddr common.Address   `json:"kOpAddr"`
	Trusted         Trusted          `json:"trusted"`
}

type Service struct {
	Identities map[common.Address]*Identity
}

func New(identitiesFilePath string) (*Service, error) {
	identitiesFile, err := os.Open(identitiesFilePath)
	if err != nil {
		return nil, err
	}

	var service Service
	if err = json.NewDecoder(identitiesFile).Decode(&service.Identities); err != nil {
		return nil, err
	}
	for idAddr, identity := range service.Identities {
		identity.IdAddr = idAddr
	}

	return &service, nil
}

func (ds *Service) GetIdentity(idAddr common.Address) (*Identity, error) {
	id, ok := ds.Identities[idAddr]
	if !ok {
		return nil, fmt.Errorf("IdAddr %v not found in the internal DB", common3.HexEncode(idAddr[:]))
	}
	return id, nil
}
