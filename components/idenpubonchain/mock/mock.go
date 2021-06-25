package mock

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree"
	"github.com/stretchr/testify/mock"
)

type IdenPubOnChainMock struct {
	mock.Mock
}

func New() *IdenPubOnChainMock {
	return &IdenPubOnChainMock{}
}

func (m *IdenPubOnChainMock) GetState(id *core.ID) (*proof.IdenStateData, error) {
	args := m.Called(id)
	return args.Get(0).(*proof.IdenStateData), args.Error(1)
}

func (m *IdenPubOnChainMock) GetStateByBlock(id *core.ID, blockN uint64) (*proof.IdenStateData, error) {
	args := m.Called(id, blockN)
	return args.Get(0).(*proof.IdenStateData), args.Error(1)
}

func (m *IdenPubOnChainMock) GetStateByTime(id *core.ID, blockTimeStamp int64) (*proof.IdenStateData, error) {
	args := m.Called(id, blockTimeStamp)
	return args.Get(0).(*proof.IdenStateData), args.Error(1)
}

func (m *IdenPubOnChainMock) InitState(id *core.ID, genesisState *merkletree.Hash, newState *merkletree.Hash, kOpProof []byte, stateTransitionProof []byte, signature *babyjub.SignatureComp) (*types.Transaction, error) {
	args := m.Called(id, genesisState, newState, kOpProof, stateTransitionProof, signature)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *IdenPubOnChainMock) SetState(id *core.ID, newState *merkletree.Hash, kOpProof []byte, stateTransitionProof []byte, signature *babyjub.SignatureComp) (*types.Transaction, error) {
	args := m.Called(id, newState, kOpProof, stateTransitionProof, signature)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

// func (m *IdenPubOnChainMock) VerifyProofClaim(pc *proof.ProofClaim) (bool, error) {
// 	args := m.Called(pc)
// 	return args.Get(0).(bool), args.Error(1)
// }

// func (m *IdenPubOnChainMock) Client() *eth.Client {
// 	args := m.Called()
// 	return args.Get(0).(*eth.Client)
// }
