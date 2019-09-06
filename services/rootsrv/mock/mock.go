package mock

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/stretchr/testify/mock"
)

type RootServiceMock struct {
	mock.Mock
}

func New() *RootServiceMock {
	return &RootServiceMock{}
}

func (m *RootServiceMock) Start() {

}

func (m *RootServiceMock) StopAndJoin() {

}

func (m *RootServiceMock) GetRoot(id *core.ID) (*core.RootData, error) {
	args := m.Called(id)
	return args.Get(0).(*core.RootData), args.Error(1)
}

func (m *RootServiceMock) SetRoot(hash merkletree.Hash) {
	m.Called(hash)
}
