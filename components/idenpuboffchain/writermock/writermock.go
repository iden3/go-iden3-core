package writermock

import (
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/stretchr/testify/mock"
)

type IdenPubOffChainWriteMock struct {
	mock.Mock
}

func New() *IdenPubOffChainWriteMock {
	return &IdenPubOffChainWriteMock{}
}

func (i *IdenPubOffChainWriteMock) Publish(publicData *idenpuboffchain.PublicData) error {
	args := i.Called(publicData)
	return args.Error(0)
}

func (i *IdenPubOffChainWriteMock) Url() string {
	args := i.Called()
	return args.Get(0).(string)
}
