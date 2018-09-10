// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package web3srv

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// IdentitiesContractABI is the input ABI used to generate the binding from.
const IdentitiesContractABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"getRoot\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"},{\"name\":\"_blockN\",\"type\":\"uint64\"}],\"name\":\"getRootByBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"getRootBlockN\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"getRootTimestamp\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"},{\"name\":\"_timestamp\",\"type\":\"uint64\"}],\"name\":\"getRootByTime\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_root\",\"type\":\"bytes32\"}],\"name\":\"setRoot\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IdentitiesContractBin is the compiled bytecode used for deploying new contracts.
const IdentitiesContractBin = `0x608060405234801561001057600080fd5b50610669806100206000396000f3006080604052600436106100775763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663079cf76e811461007c57806318c9990b146100af57806336025294146100dd5780633cf204951461011b578063ab481d231461013c578063dab5f3401461016a575b600080fd5b34801561008857600080fd5b5061009d600160a060020a0360043516610184565b60408051918252519081900360200190f35b3480156100bb57600080fd5b5061009d600160a060020a036004351667ffffffffffffffff602435166101c5565b3480156100e957600080fd5b506100fe600160a060020a036004351661036c565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561012757600080fd5b506100fe600160a060020a03600435166103b5565b34801561014857600080fd5b5061009d600160a060020a036004351667ffffffffffffffff6024351661040a565b34801561017657600080fd5b506101826004356105a7565b005b600160a060020a0381166000908152602081905260408120805460001981019081106101ac57fe5b9060005260206000209060020201600101549050919050565b60008080804367ffffffffffffffff86161061024257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6572724e6f467574757265416c6c6f7765640000000000000000000000000000604482015290519081900360640190fd5b600160a060020a0386166000908152602081905260408120549093506000190191505b8183116103635750600160a060020a038516600090815260208190526040902080546002848401049167ffffffffffffffff871691839081106102a457fe5b600091825260209091206002909102015467ffffffffffffffff16141561030257600160a060020a03861660009081526020819052604090208054829081106102e957fe5b9060005260206000209060020201600101549350610363565b600160a060020a038616600090815260208190526040902080548290811061032657fe5b600091825260209091206002909102015467ffffffffffffffff90811690861611156103575780600101925061035e565b6001810391505b610265565b50505092915050565b600160a060020a03811660009081526020819052604081208054600019810190811061039457fe5b600091825260209091206002909102015467ffffffffffffffff1692915050565b600160a060020a0381166000908152602081905260408120805460001981019081106103dd57fe5b600091825260209091206002909102015468010000000000000000900467ffffffffffffffff1692915050565b60008080804267ffffffffffffffff86161061048757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6572724e6f467574757265416c6c6f7765640000000000000000000000000000604482015290519081900360640190fd5b600160a060020a0386166000908152602081905260408120549093506000190191505b8183116103635750600160a060020a038516600090815260208190526040902080546002848401049167ffffffffffffffff871691839081106104e957fe5b600091825260209091206002909102015468010000000000000000900467ffffffffffffffff16141561053a57600160a060020a03861660009081526020819052604090208054829081106102e957fe5b600160a060020a038616600090815260208190526040902080548290811061055e57fe5b600091825260209091206002909102015467ffffffffffffffff680100000000000000009091048116908616111561059b578060010192506105a2565b6001810391505b6104aa565b33600090815260208181526040808320815160608101835267ffffffffffffffff438116825242811682860190815293820196875282546001818101855593875294909520905160029094020180549251851668010000000000000000026fffffffffffffffff0000000000000000199490951667ffffffffffffffff19909316929092179290921692909217825591519101555600a165627a7a7230582052b01c82f035d62ae30fec1ef562512a30d83fd6bcd64b1b6652b9f9c69e4cb40029`

// DeployIdentitiesContract deploys a new Ethereum contract, binding an instance of IdentitiesContract to it.
func DeployIdentitiesContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IdentitiesContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IdentitiesContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IdentitiesContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IdentitiesContract{IdentitiesContractCaller: IdentitiesContractCaller{contract: contract}, IdentitiesContractTransactor: IdentitiesContractTransactor{contract: contract}, IdentitiesContractFilterer: IdentitiesContractFilterer{contract: contract}}, nil
}

// IdentitiesContract is an auto generated Go binding around an Ethereum contract.
type IdentitiesContract struct {
	IdentitiesContractCaller     // Read-only binding to the contract
	IdentitiesContractTransactor // Write-only binding to the contract
	IdentitiesContractFilterer   // Log filterer for contract events
}

// IdentitiesContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type IdentitiesContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentitiesContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IdentitiesContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentitiesContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IdentitiesContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IdentitiesContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IdentitiesContractSession struct {
	Contract     *IdentitiesContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IdentitiesContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IdentitiesContractCallerSession struct {
	Contract *IdentitiesContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// IdentitiesContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IdentitiesContractTransactorSession struct {
	Contract     *IdentitiesContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IdentitiesContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type IdentitiesContractRaw struct {
	Contract *IdentitiesContract // Generic contract binding to access the raw methods on
}

// IdentitiesContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IdentitiesContractCallerRaw struct {
	Contract *IdentitiesContractCaller // Generic read-only contract binding to access the raw methods on
}

// IdentitiesContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IdentitiesContractTransactorRaw struct {
	Contract *IdentitiesContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIdentitiesContract creates a new instance of IdentitiesContract, bound to a specific deployed contract.
func NewIdentitiesContract(address common.Address, backend bind.ContractBackend) (*IdentitiesContract, error) {
	contract, err := bindIdentitiesContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IdentitiesContract{IdentitiesContractCaller: IdentitiesContractCaller{contract: contract}, IdentitiesContractTransactor: IdentitiesContractTransactor{contract: contract}, IdentitiesContractFilterer: IdentitiesContractFilterer{contract: contract}}, nil
}

// NewIdentitiesContractCaller creates a new read-only instance of IdentitiesContract, bound to a specific deployed contract.
func NewIdentitiesContractCaller(address common.Address, caller bind.ContractCaller) (*IdentitiesContractCaller, error) {
	contract, err := bindIdentitiesContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IdentitiesContractCaller{contract: contract}, nil
}

// NewIdentitiesContractTransactor creates a new write-only instance of IdentitiesContract, bound to a specific deployed contract.
func NewIdentitiesContractTransactor(address common.Address, transactor bind.ContractTransactor) (*IdentitiesContractTransactor, error) {
	contract, err := bindIdentitiesContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IdentitiesContractTransactor{contract: contract}, nil
}

// NewIdentitiesContractFilterer creates a new log filterer instance of IdentitiesContract, bound to a specific deployed contract.
func NewIdentitiesContractFilterer(address common.Address, filterer bind.ContractFilterer) (*IdentitiesContractFilterer, error) {
	contract, err := bindIdentitiesContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IdentitiesContractFilterer{contract: contract}, nil
}

// bindIdentitiesContract binds a generic wrapper to an already deployed contract.
func bindIdentitiesContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IdentitiesContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentitiesContract *IdentitiesContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IdentitiesContract.Contract.IdentitiesContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentitiesContract *IdentitiesContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentitiesContract.Contract.IdentitiesContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentitiesContract *IdentitiesContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentitiesContract.Contract.IdentitiesContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IdentitiesContract *IdentitiesContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IdentitiesContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IdentitiesContract *IdentitiesContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IdentitiesContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IdentitiesContract *IdentitiesContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IdentitiesContract.Contract.contract.Transact(opts, method, params...)
}

// GetRoot is a free data retrieval call binding the contract method 0x079cf76e.
//
// Solidity: function getRoot(_address address) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractCaller) GetRoot(opts *bind.CallOpts, _address common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _IdentitiesContract.contract.Call(opts, out, "getRoot", _address)
	return *ret0, err
}

// GetRoot is a free data retrieval call binding the contract method 0x079cf76e.
//
// Solidity: function getRoot(_address address) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractSession) GetRoot(_address common.Address) ([32]byte, error) {
	return _IdentitiesContract.Contract.GetRoot(&_IdentitiesContract.CallOpts, _address)
}

// GetRoot is a free data retrieval call binding the contract method 0x079cf76e.
//
// Solidity: function getRoot(_address address) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractCallerSession) GetRoot(_address common.Address) ([32]byte, error) {
	return _IdentitiesContract.Contract.GetRoot(&_IdentitiesContract.CallOpts, _address)
}

// GetRootBlockN is a free data retrieval call binding the contract method 0x36025294.
//
// Solidity: function getRootBlockN(_address address) constant returns(uint64)
func (_IdentitiesContract *IdentitiesContractCaller) GetRootBlockN(opts *bind.CallOpts, _address common.Address) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _IdentitiesContract.contract.Call(opts, out, "getRootBlockN", _address)
	return *ret0, err
}

// GetRootBlockN is a free data retrieval call binding the contract method 0x36025294.
//
// Solidity: function getRootBlockN(_address address) constant returns(uint64)
func (_IdentitiesContract *IdentitiesContractSession) GetRootBlockN(_address common.Address) (uint64, error) {
	return _IdentitiesContract.Contract.GetRootBlockN(&_IdentitiesContract.CallOpts, _address)
}

// GetRootBlockN is a free data retrieval call binding the contract method 0x36025294.
//
// Solidity: function getRootBlockN(_address address) constant returns(uint64)
func (_IdentitiesContract *IdentitiesContractCallerSession) GetRootBlockN(_address common.Address) (uint64, error) {
	return _IdentitiesContract.Contract.GetRootBlockN(&_IdentitiesContract.CallOpts, _address)
}

// GetRootByBlock is a free data retrieval call binding the contract method 0x18c9990b.
//
// Solidity: function getRootByBlock(_address address, _blockN uint64) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractCaller) GetRootByBlock(opts *bind.CallOpts, _address common.Address, _blockN uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _IdentitiesContract.contract.Call(opts, out, "getRootByBlock", _address, _blockN)
	return *ret0, err
}

// GetRootByBlock is a free data retrieval call binding the contract method 0x18c9990b.
//
// Solidity: function getRootByBlock(_address address, _blockN uint64) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractSession) GetRootByBlock(_address common.Address, _blockN uint64) ([32]byte, error) {
	return _IdentitiesContract.Contract.GetRootByBlock(&_IdentitiesContract.CallOpts, _address, _blockN)
}

// GetRootByBlock is a free data retrieval call binding the contract method 0x18c9990b.
//
// Solidity: function getRootByBlock(_address address, _blockN uint64) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractCallerSession) GetRootByBlock(_address common.Address, _blockN uint64) ([32]byte, error) {
	return _IdentitiesContract.Contract.GetRootByBlock(&_IdentitiesContract.CallOpts, _address, _blockN)
}

// GetRootByTime is a free data retrieval call binding the contract method 0xab481d23.
//
// Solidity: function getRootByTime(_address address, _timestamp uint64) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractCaller) GetRootByTime(opts *bind.CallOpts, _address common.Address, _timestamp uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _IdentitiesContract.contract.Call(opts, out, "getRootByTime", _address, _timestamp)
	return *ret0, err
}

// GetRootByTime is a free data retrieval call binding the contract method 0xab481d23.
//
// Solidity: function getRootByTime(_address address, _timestamp uint64) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractSession) GetRootByTime(_address common.Address, _timestamp uint64) ([32]byte, error) {
	return _IdentitiesContract.Contract.GetRootByTime(&_IdentitiesContract.CallOpts, _address, _timestamp)
}

// GetRootByTime is a free data retrieval call binding the contract method 0xab481d23.
//
// Solidity: function getRootByTime(_address address, _timestamp uint64) constant returns(bytes32)
func (_IdentitiesContract *IdentitiesContractCallerSession) GetRootByTime(_address common.Address, _timestamp uint64) ([32]byte, error) {
	return _IdentitiesContract.Contract.GetRootByTime(&_IdentitiesContract.CallOpts, _address, _timestamp)
}

// GetRootTimestamp is a free data retrieval call binding the contract method 0x3cf20495.
//
// Solidity: function getRootTimestamp(_address address) constant returns(uint64)
func (_IdentitiesContract *IdentitiesContractCaller) GetRootTimestamp(opts *bind.CallOpts, _address common.Address) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _IdentitiesContract.contract.Call(opts, out, "getRootTimestamp", _address)
	return *ret0, err
}

// GetRootTimestamp is a free data retrieval call binding the contract method 0x3cf20495.
//
// Solidity: function getRootTimestamp(_address address) constant returns(uint64)
func (_IdentitiesContract *IdentitiesContractSession) GetRootTimestamp(_address common.Address) (uint64, error) {
	return _IdentitiesContract.Contract.GetRootTimestamp(&_IdentitiesContract.CallOpts, _address)
}

// GetRootTimestamp is a free data retrieval call binding the contract method 0x3cf20495.
//
// Solidity: function getRootTimestamp(_address address) constant returns(uint64)
func (_IdentitiesContract *IdentitiesContractCallerSession) GetRootTimestamp(_address common.Address) (uint64, error) {
	return _IdentitiesContract.Contract.GetRootTimestamp(&_IdentitiesContract.CallOpts, _address)
}

// SetRoot is a paid mutator transaction binding the contract method 0xdab5f340.
//
// Solidity: function setRoot(_root bytes32) returns()
func (_IdentitiesContract *IdentitiesContractTransactor) SetRoot(opts *bind.TransactOpts, _root [32]byte) (*types.Transaction, error) {
	return _IdentitiesContract.contract.Transact(opts, "setRoot", _root)
}

// SetRoot is a paid mutator transaction binding the contract method 0xdab5f340.
//
// Solidity: function setRoot(_root bytes32) returns()
func (_IdentitiesContract *IdentitiesContractSession) SetRoot(_root [32]byte) (*types.Transaction, error) {
	return _IdentitiesContract.Contract.SetRoot(&_IdentitiesContract.TransactOpts, _root)
}

// SetRoot is a paid mutator transaction binding the contract method 0xdab5f340.
//
// Solidity: function setRoot(_root bytes32) returns()
func (_IdentitiesContract *IdentitiesContractTransactorSession) SetRoot(_root [32]byte) (*types.Transaction, error) {
	return _IdentitiesContract.Contract.SetRoot(&_IdentitiesContract.TransactOpts, _root)
}
