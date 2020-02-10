package idenstatereader

// TODO: Rename this to IdenStatePubOnchain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/eth/contracts"
	"github.com/iden3/go-iden3-core/merkletree"
)

// IdenStateReader is an interface that gives access to the IdenStates Smart Contract.
type IdenStateReader interface {
	GetState(id *core.ID) (*proof.IdenStateData, error)
	GetStateByBlock(id *core.ID, blockN uint64) (merkletree.Hash, error)
	GetStateByTime(id *core.ID, blockTimestamp int64) (merkletree.Hash, error)
	SetState(id *core.ID, newState *merkletree.Hash, kOpProof []byte, stateTransitionProof []byte, signature *merkletree.Hash) (*types.Transaction, error)
	InitState(id *core.ID, genesisState *merkletree.Hash, newState *merkletree.Hash, kOpProof []byte, stateTransitionProof []byte, signature *merkletree.Hash) (*types.Transaction, error)
	// VerifyProofClaim(pc *proof.ProofClaim) (bool, error)
}

// ContractAddresses are the list of Smart Contract addresses used for the on chain identity state data.
type ContractAddresses struct {
	IdenStates common.Address
}

// IdenStateRead is the regular implementation of IdenStateReader
type IdenStateRead struct {
	client    *eth.Client2
	addresses ContractAddresses
}

// New creates a new IdenStateRead
func New(client *eth.Client2, addresses ContractAddresses) *IdenStateRead {
	return &IdenStateRead{
		client:    client,
		addresses: addresses,
	}
}

// GetState returns the Identity State of the given ID from the IdenStates Smart Contract.
func (s *IdenStateRead) GetState(id *core.ID) (*proof.IdenStateData, error) {
	var idenState [32]byte
	var blockN uint64
	var blockTS uint64
	err := s.client.Call(func(c *ethclient.Client) error {
		idenStates, err := contracts.NewState(s.addresses.IdenStates, c)
		if err != nil {
			return err
		}
		blockN, blockTS, idenState, err = idenStates.GetStateDataById(nil, *id)
		return err
	})
	return &proof.IdenStateData{
		BlockN:    blockN,
		BlockTs:   int64(blockTS),
		IdenState: (*merkletree.Hash)(&idenState),
	}, err
}

// GetState returns the Identity State of the given ID closest to the blockN
// from the IdenStates Smart Contract.
func (s *IdenStateRead) GetStateByBlock(id *core.ID, blockN uint64) (merkletree.Hash, error) {
	var idenState [32]byte
	err := s.client.Call(func(c *ethclient.Client) error {
		idenStates, err := contracts.NewState(s.addresses.IdenStates, c)
		if err != nil {
			return err
		}
		idenState, err = idenStates.GetStateByBlock(nil, *id, blockN)
		return err
	})
	return merkletree.Hash(idenState), err
}

// GetState returns the Identity State of the given ID closest to the blockTimeStamp
// from the IdenStates Smart Contract.
func (s *IdenStateRead) GetStateByTime(id *core.ID, blockTimeStamp int64) (merkletree.Hash, error) {
	var idenState [32]byte
	err := s.client.Call(func(c *ethclient.Client) error {
		idenStates, err := contracts.NewState(s.addresses.IdenStates, c)
		if err != nil {
			return err
		}
		idenState, err = idenStates.GetStateByTime(nil, *id, uint64(blockTimeStamp))
		return err
	})
	return merkletree.Hash(idenState), err
}

// SetState updates the Identity State of the given ID in the IdenStates Smart Contract.
func (s *IdenStateRead) SetState(id *core.ID, newState *merkletree.Hash, kOpProof []byte, stateTransitionProof []byte, signature *merkletree.Hash) (*types.Transaction, error) {
	if tx, err := s.client.CallAuth(
		func(c *ethclient.Client, auth *bind.TransactOpts) (*types.Transaction, error) {
			idenStates, err := contracts.NewState(s.addresses.IdenStates, c)
			if err != nil {
				return nil, err
			}
			return idenStates.SetState(auth, *newState, *id, kOpProof, stateTransitionProof, *signature)
		},
	); err != nil {
		return nil, fmt.Errorf("Failed setting identity state in the Smart Contract (setState): %w", err)
	} else {
		return tx, nil
	}
}

// InitState initializes the first Identity State of the given ID in the IdenStates Smart Contract.
func (s *IdenStateRead) InitState(id *core.ID, genesisState *merkletree.Hash, newState *merkletree.Hash, kOpProof []byte, stateTransitionProof []byte, signature *merkletree.Hash) (*types.Transaction, error) {
	if tx, err := s.client.CallAuth(
		func(c *ethclient.Client, auth *bind.TransactOpts) (*types.Transaction, error) {
			idenStates, err := contracts.NewState(s.addresses.IdenStates, c)
			if err != nil {
				return nil, err
			}
			return idenStates.InitState(auth, *newState, *genesisState, *id, kOpProof, stateTransitionProof, *signature)
		},
	); err != nil {
		return nil, fmt.Errorf("Failed initalizating identity state in the Smart Contract (initState): %w", err)
	} else {
		return tx, nil
	}
}

// Should this really be here?
// func (s *IdenStateRead) VerifyProofClaim(pc *proof.ProofClaim) (bool, error) {
// 	if ok, err := pc.Verify(pc.Proof.Root); !ok {
// 		return false, err
// 	}
// 	id, blockN, blockTime := pc.PublishedData()
// 	rootByBlock, err := s.GetStateByBlock(id, blockN)
// 	if err != nil {
// 		return false, err
// 	}
// 	rootByTime, err := s.GetStateByTime(id, blockTime)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	if !pc.Proof.Root.Equals(&rootByBlock) {
// 		return false, fmt.Errorf("ProofClaim Root doesn't match the one " +
// 			"from the smart contract queried by (id, blockN)")
// 	}
// 	if !pc.Proof.Root.Equals(&rootByTime) {
// 		return false, fmt.Errorf("ProofClaim Root doesn't match the one " +
// 			"from the smart contract queried by (id, blockTime)")
// 	}
// 	return true, nil
// }

// func (s *IdenStateRead) Client() *eth.Client2 {
// 	return s.client
// }
