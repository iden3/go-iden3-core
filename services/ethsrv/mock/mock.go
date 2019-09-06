package mock

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/mock"
)

type EthServiceMock struct {
	mock.Mock
}

func New() *EthServiceMock {
	return &EthServiceMock{}
}

func (m *EthServiceMock) GetRoot(id *core.ID) (*core.RootData, error) {
	args := m.Called(id)
	return args.Get(0).(*core.RootData), args.Error(1)
}

func (m *EthServiceMock) GetRootByBlock(id *core.ID, blockN uint64) (merkletree.Hash, error) {
	args := m.Called(id, blockN)
	return args.Get(0).(merkletree.Hash), args.Error(1)
}
func (m *EthServiceMock) GetRootByTime(id *core.ID, blockTimestamp int64) (merkletree.Hash, error) {
	args := m.Called(id, blockTimestamp)
	return args.Get(0).(merkletree.Hash), args.Error(1)
}
func (m *EthServiceMock) VerifyProofClaim(pc *core.ProofClaim) (bool, error) {
	args := m.Called(pc)
	return args.Get(0).(bool), args.Error(1)
}

func (m *EthServiceMock) Client() *eth.Client2 {
	args := m.Called()
	return args.Get(0).(*eth.Client2)
}
