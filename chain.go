package core

import (
	"fmt"

	"github.com/iden3/go-iden3-core/v2/w3c"
)

// ChainID is alias for int32 that represents ChainID
type ChainID int32

// ChainIDs Object containing chain IDs for various blockchains and networks.
var chainIDs = map[string]ChainID{
	"eth":            1,
	"eth:main":       1,
	"eth:goerli":     5,
	"eth:sepolia":    11155111,
	"polygon":        137,
	"polygon:main":   137,
	"polygon:mumbai": 80001,
	"zkevm":          1101,
	"zkevm:main":     1101,
	"zkevm:test":     1442,
}

// ChainIDfromDID returns chain name from w3c.DID
func ChainIDfromDID(did w3c.DID) (ChainID, error) {
	// TODO: fix for networks like eth / polygon / zkevm

	id, err := IDFromDID(did)
	if err != nil {
		return 0, err
	}

	blockchain, err := BlockchainFromID(id)
	if err != nil {
		return 0, err
	}

	networkID, err := NetworkIDFromID(id)
	if err != nil {
		return 0, err
	}

	chainID, ok := chainIDs[fmt.Sprintf("%s:%s", blockchain, networkID)]
	if !ok {
		return 0, fmt.Errorf("chainID not found for %s:%s", blockchain, networkID)
	}

	return chainID, nil
}

// RegisterChainID registers chainID for blockchain and network
func RegisterChainID(blockchain Blockchain, network NetworkID, chainID int) error {
	if _, ok := blockchains[blockchain]; !ok {
		return fmt.Errorf("blockchain not registered: %s", blockchain)
	}

	if _, ok := networks[network]; !ok {
		return fmt.Errorf("network not registered: %s", network)
	}

	k := fmt.Sprintf("%s:%s", blockchain, network)
	if _, ok := chainIDs[k]; ok {
		return fmt.Errorf("chainID already registered for %s:%s", blockchain, network)
	}
	chainIDs[k] = ChainID(chainID)

	return nil
}

// GetChainID returns chainID for blockchain and network
func GetChainID(blockchain Blockchain, network NetworkID) (ChainID, error) {
	k := fmt.Sprintf("%s:%s", blockchain, network)
	if _, ok := chainIDs[k]; !ok {
		return 0, fmt.Errorf("chainID not registered for %s:%s", blockchain, network)
	}

	return chainIDs[k], nil
}
