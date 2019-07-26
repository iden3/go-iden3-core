package discoverysrv

import (
	"encoding/json"
	"fmt"
	"os"

	// "github.com/ethereum/go-ethereum/common"
	// common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-core/core"
	// "github.com/iden3/go-iden3-core/utils"
)

type Trusted struct {
	Relay bool `json:"relay"`
}

type Entity struct {
	Id            core.ID            `json:"-"`
	Name          string             `json:"name"`
	OperationalPk *babyjub.PublicKey `json:"kOpPub"`
	// OperationalAddr common.Address   `json:"kOpAddr"`
	Trusted Trusted `json:"trusted"`
}

type Service struct {
	Entitites map[core.ID]*Entity
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
	for id, identity := range service.Entitites {
		identity.Id = id
	}

	return &service, nil
}

func (ds *Service) GetEntity(id core.ID) (*Entity, error) {
	entity, ok := ds.Entitites[id]
	if !ok {
		return nil, fmt.Errorf("Id %v not found in the internal DB", &id)
	}
	return entity, nil
}
