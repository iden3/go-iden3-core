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

type Entity struct {
	IdAddr          common.Address   `json:"-"`
	Name            string           `json:"name"`
	OperationalPk   *utils.PublicKey `json:"kOpPub"`
	OperationalAddr common.Address   `json:"kOpAddr"`
	Trusted         Trusted          `json:"trusted"`
}

type Service struct {
	Entitites map[common.Address]*Entity
}

func New(entititesFilePath string) (*Service, error) {
	entititesFile, err := os.Open(entititesFilePath)
	if err != nil {
		return nil, err
	}

	var service Service
	if err = json.NewDecoder(entititesFile).Decode(&service.Entitites); err != nil {
		return nil, err
	}
	for idAddr, identity := range service.Entitites {
		identity.IdAddr = idAddr
	}

	return &service, nil
}

func (ds *Service) GetEntity(idAddr common.Address) (*Entity, error) {
	id, ok := ds.Entitites[idAddr]
	if !ok {
		return nil, fmt.Errorf("IdAddr %v not found in the internal DB", common3.HexEncode(idAddr[:]))
	}
	return id, nil
}
