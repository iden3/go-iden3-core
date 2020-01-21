package mock

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/mock"
)

type IdenStateWriteMock struct {
	mock.Mock
}

func New() *IdenStateWriteMock {
	return &IdenStateWriteMock{}
}

func (m *IdenStateWriteMock) Start() {

}

func (m *IdenStateWriteMock) StopAndJoin() {

}

func (m *IdenStateWriteMock) GetRoot(id *core.ID) (*core.RootData, error) {
	args := m.Called(id)
	return args.Get(0).(*core.RootData), args.Error(1)
}

func (m *IdenStateWriteMock) SetRoot(hash merkletree.Hash) {
	m.Called(hash)
}
