package nameresolversrv

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iden3/go-iden3/core"
)

type Service struct {
	Names map[string]core.ID
}

func New(namesFilePath string) (*Service, error) {
	namesFile, err := os.Open(namesFilePath)
	if err != nil {
		return nil, err
	}
	var service Service
	fmt.Println(namesFile)
	if err = json.NewDecoder(namesFile).Decode(&service.Names); err != nil {
		return nil, err
	}

	return &service, nil
}

func (ns *Service) Resolve(name string) (*core.ID, error) {
	id, ok := ns.Names[name]
	if !ok {
		return nil, fmt.Errorf("Name %v not found in the internal DB", name)
	}
	return &id, nil
}
