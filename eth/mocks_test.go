package eth

import (
	"math/big"

	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) NetworkID() (*big.Int, error) {
	args := m.Called()
	return args.Get(0).(*big.Int), args.Error(1)
}
func (m *ClientMock) BalanceInfo() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *ClientMock) SendTransactionSync(to *common.Address, value *big.Int, gasLimit uint64, calldata []byte) (*types.Transaction, *types.Receipt, error) {
	args := m.Called(to, value, gasLimit, calldata)
	return args.Get(0).(*types.Transaction), args.Get(1).(*types.Receipt), args.Error(2)
}
func (m *ClientMock) Call(to *common.Address, value *big.Int, calldata []byte) ([]byte, error) {
	args := m.Called(to, value, calldata)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *ClientMock) Sign(data ...[]byte) ([3][32]byte, error) {
	args := m.Called(data)
	return args.Get(0).([3][32]byte), args.Error(1)
}
func (m *ClientMock) CodeAt(account common.Address) ([]byte, error) {
	args := m.Called(account)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *ClientMock) SendRawTxSync(rawtx []byte) (*types.Receipt, error) {
	args := m.Called(rawtx)
	return args.Get(0).(*types.Receipt), args.Error(1)
}
