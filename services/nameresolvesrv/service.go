package nameresolvesrv

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

type Service struct {
	Names map[string]common.Address
}

func New(namesFilePath string) (*Service, error) {
	namesFile, err := os.Open(namesFilePath)
	if err != nil {
		return nil, err
	}
	var service Service
	if err = json.NewDecoder(namesFile).Decode(&service.Names); err != nil {
		return nil, err
	}

	return &service, nil
}

func (ns *Service) Resolve(name string) (*common.Address, error) {
	idAddr, ok := ns.Names[name]
	if !ok {
		return nil, fmt.Errorf("Name %v not found in the internal DB", name)
	}
	return &idAddr, nil
}
