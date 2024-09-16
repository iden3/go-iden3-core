package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/iden3/go-iden3-core/v2/w3c"
)

// ChainID is alias for int32 that represents ChainID
type ChainID int32

type chainIDKey struct {
	blockchain Blockchain
	networkID  NetworkID
}

var ErrChainIDNotRegistered = errors.New("chainID is not registered")

var chainIDsLock sync.RWMutex

// chainIDs Object containing chain IDs for various blockchains and networks.
// It can be modified using RegisterChainID public function. So it is guarded
// by chainIDsLock mutex.
var chainIDs = map[chainIDKey]ChainID{
	{Ethereum, Main}:    1,
	{Ethereum, Goerli}:  5,
	{Ethereum, Sepolia}: 11155111,
	{Polygon, Main}:     137,
	{Polygon, Mumbai}:   80001,
	{Polygon, Amoy}:     80002,
	{Polygon, Zkevm}:    1101,
	{Polygon, Cardona}:  2442,
	{Privado, Main}:     21000,
	{Privado, Test}:     21001,
	{Linea, Main}:       59144,
	{Linea, Sepolia}:    59141,
}

// ChainIDfromDID returns chain name from w3c.DID
func ChainIDfromDID(did w3c.DID) (ChainID, error) {
	id, err := IDFromDID(did)
	if err != nil {
		return 0, err
	}

	return ChainIDfromID(id)
}

// ChainIDfromID(id ID) returns chain name from ID
func ChainIDfromID(id ID) (ChainID, error) {
	blockchain, err := BlockchainFromID(id)
	if err != nil {
		return 0, err
	}

	networkID, err := NetworkIDFromID(id)
	if err != nil {
		return 0, err
	}

	return GetChainID(blockchain, networkID)
}

// RegisterChainID registers chainID for blockchain and network
func RegisterChainID(blockchain Blockchain, network NetworkID, chainID int) error {
	chainIDsLock.Lock()
	defer chainIDsLock.Unlock()

	k := chainIDKey{
		blockchain: blockchain,
		networkID:  network,
	}
	existingChainID, ok := chainIDs[k]
	if ok && existingChainID == ChainID(chainID) {
		return nil
	}

	for _, v := range chainIDs {
		if v == ChainID(chainID) {
			return fmt.Errorf(`can't register chain id %d for '%v:%v' because it's already registered for another chain id`,
				chainID, k.blockchain, k.networkID)
		}
	}

	chainIDs[k] = ChainID(chainID)

	return nil
}

// GetChainID returns chainID for blockchain and network
func GetChainID(blockchain Blockchain, network NetworkID) (ChainID, error) {
	chainIDsLock.RLock()
	defer chainIDsLock.RUnlock()

	k := chainIDKey{
		blockchain: blockchain,
		networkID:  network,
	}
	if _, ok := chainIDs[k]; !ok {
		return 0, fmt.Errorf("%w for %s:%s", ErrChainIDNotRegistered, blockchain,
			network)
	}

	return chainIDs[k], nil
}

// NetworkByChainID returns blockchain and networkID for registered chain ID.
// Or ErrChainIDNotRegistered error if chainID is not registered.
func NetworkByChainID(chainID ChainID) (Blockchain, NetworkID, error) {
	chainIDsLock.RLock()
	defer chainIDsLock.RUnlock()

	for k, v := range chainIDs {
		if v == chainID {
			return k.blockchain, k.networkID, nil
		}
	}
	return NoChain, NoNetwork, ErrChainIDNotRegistered
}
