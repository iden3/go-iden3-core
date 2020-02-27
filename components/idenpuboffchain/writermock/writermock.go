package writermock

import (
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/stretchr/testify/mock"
)

type IdenPubOffChainWriteMock struct {
	mock.Mock
}

func New() *IdenPubOffChainWriteMock {
	return &IdenPubOffChainWriteMock{}
}

func (i *IdenPubOffChainWriteMock) Publish(id *core.ID, publicData *idenpuboffchain.PublicData) error {
	args := i.Called(id, publicData)
	return args.Error(0)
}

func (i *IdenPubOffChainWriteMock) Url() string {
	args := i.Called()
	return args.Get(0).(string)
}
