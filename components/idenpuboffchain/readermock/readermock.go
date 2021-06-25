package readermock

import (
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-merkletree"
	"github.com/stretchr/testify/mock"
)

type IdenPubOffChainReadMock struct {
	mock.Mock
}

func New() *IdenPubOffChainReadMock {
	return &IdenPubOffChainReadMock{}
}

func (i *IdenPubOffChainReadMock) GetPublicData(idenPubUrl string, id *core.ID, idenState *merkletree.Hash) (*idenpuboffchain.PublicData, error) {
	args := i.Called(idenPubUrl, id, idenState)
	return args.Get(0).(*idenpuboffchain.PublicData), args.Error(1)
}
