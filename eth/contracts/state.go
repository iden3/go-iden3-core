// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StateABI is the input ABI used to generate the binding from.
const StateABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"state\",\"type\":\"bytes32\"}],\"name\":\"StateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"}],\"name\":\"getStateDataByBlock\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getStateDataById\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getStateDataByTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"genesisState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"bytes\",\"name\":\"kOpProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"itp\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"sigR8\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sigS\",\"type\":\"bytes32\"}],\"name\":\"initState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"bytes\",\"name\":\"kOpProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"itp\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"sigR8\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sigS\",\"type\":\"bytes32\"}],\"name\":\"setState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// StateFuncSigs maps the 4-byte function signature to its string representation.
var StateFuncSigs = map[string]string{
	"c1056e53": "getState(bytes31)",
	"c68631e1": "getStateDataByBlock(bytes31,uint64)",
	"4cabaefa": "getStateDataById(bytes31)",
	"5710773a": "getStateDataByTime(bytes31,uint64)",
	"9b2d7aac": "initState(bytes32,bytes32,bytes31,bytes,bytes,bytes32,bytes32)",
	"ea1662de": "setState(bytes32,bytes31,bytes,bytes,bytes32,bytes32)",
}

// StateBin is the compiled bytecode used for deploying new contracts.
var StateBin = "0x608060405234801561001057600080fd5b50610e82806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80634cabaefa146100675780635710773a146100b25780639b2d7aac146100e2578063c1056e5314610229578063c68631e11461025c578063ea1662de1461028c575b600080fd5b6100886004803603602081101561007d57600080fd5b503560ff19166103ca565b604080516001600160401b0394851681529290931660208301528183015290519081900360600190f35b610088600480360360408110156100c857600080fd5b50803560ff191690602001356001600160401b031661047a565b610227600480360360e08110156100f857600080fd5b81359160208101359160ff196040830135169190810190608081016060820135600160201b81111561012957600080fd5b82018360208201111561013b57600080fd5b803590602001918460018302840111600160201b8311171561015c57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b8111156101ae57600080fd5b8201836020820111156101c057600080fd5b803590602001918460018302840111600160201b831117156101e157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550508235935050506020013561087a565b005b61024a6004803603602081101561023f57600080fd5b503560ff19166108b0565b60408051918252519081900360200190f35b6100886004803603604081101561027257600080fd5b50803560ff191690602001356001600160401b031661090e565b610227600480360360c08110156102a257600080fd5b81359160ff1960208201351691810190606081016040820135600160201b8111156102cc57600080fd5b8201836020820111156102de57600080fd5b803590602001918460018302840111600160201b831117156102ff57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561035157600080fd5b82018360208201111561036357600080fd5b803590602001918460018302840111600160201b8311171561038457600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505082359350505060200135610b93565b60ff198116600090815260208190526040812054819081906103f6575050600154600091508190610473565b6103fe610e0b565b60ff19851660009081526020819052604090208054600019810190811061042157fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b909204169484018590526001909101549290910182905295509093509150505b9193909250565b600080600042846001600160401b0316106104d1576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff1985166000908152602081905260409020546104f9575050600154600091508190610873565b60ff19851660009081526020819052604081208054600019810190811061051c57fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116915085168110156106105760ff19861660009081526020819052604090208054600019810190811061056d57fe5b6000918252602080832060029092029091015460ff198916835290829052604090912080546001600160401b039092169160001981019081106105ac57fe5b600091825260208083206002929092029091015460ff198a1683529082905260409091208054600160401b9092046001600160401b03169160001981019081106105f257fe5b90600052602060002090600202016001015493509350935050610873565b60ff198616600090815260208190526040812054600019015b8082116108635760ff19881660009081526020819052604090208054600284840104916001600160401b038a16918390811061066157fe5b6000918252602090912060029091020154600160401b90046001600160401b031614156107455760ff19891660009081526020819052604090208054829081106106a757fe5b6000918252602080832060029092029091015460ff198c16835290829052604090912080546001600160401b0390921691839081106106e257fe5b600091825260208083206002929092029091015460ff198d1683529082905260409091208054600160401b9092046001600160401b0316918490811061072457fe5b90600052602060002090600202016001015496509650965050505050610873565b60ff198916600090815260208190526040902080548290811061076457fe5b60009182526020909120600290910201546001600160401b03600160401b90910481169089161180156107dc575060ff19891660009081526020819052604090208054600183019081106107b457fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908916105b156108005760ff19891660009081526020819052604090208054829081106106a757fe5b60ff198916600090815260208190526040902080548290811061081f57fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690891611156108565780600101925061085d565b6001810391505b50610629565b5050600154600094508493509150505b9250925092565b60ff1985166000908152602081905260409020541561089857600080fd5b6108a787878787878787610c6f565b50505050505050565b60ff1981166000908152602081905260408120546108d15750600154610909565b60ff1982166000908152602081905260409020805460001981019081106108f457fe5b90600052602060002090600202016001015490505b919050565b600080600043846001600160401b031610610965576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff19851660009081526020819052604090205461098d575050600154600091508190610873565b60ff1985166000908152602081905260408120805460001981019081106109b057fe5b60009182526020909120600290910201546001600160401b03908116915085168110156109fa5760ff19861660009081526020819052604090208054600019810190811061056d57fe5b60ff198616600090815260208190526040812054600019015b8082116108635760ff19881660009081526020819052604090208054600284840104916001600160401b038a169183908110610a4b57fe5b60009182526020909120600290910201546001600160401b03161415610a8a5760ff19891660009081526020819052604090208054829081106106a757fe5b60ff1989166000908152602081905260409020805482908110610aa957fe5b60009182526020909120600290910201546001600160401b03908116908916118015610b13575060ff1989166000908152602081905260409020805460018301908110610af257fe5b60009182526020909120600290910201546001600160401b03908116908916105b15610b375760ff19891660009081526020819052604090208054829081106106a757fe5b60ff1989166000908152602081905260409020805482908110610b5657fe5b60009182526020909120600290910201546001600160401b039081169089161115610b8657806001019250610b8d565b6001810391505b50610a13565b60ff198516600090815260208190526040902054610bb057600080fd5b610bb8610e0b565b60ff198616600090815260208190526040902080546000198101908110610bdb57fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b909204169484019490945260010154908201529150431415610c605760405162461bcd60e51b8152600401808060200182810382526021815260200180610e2c6021913960400191505060405180910390fd5b6108a787826040015188888888885b610c798785610dfa565b1515600114610c8757600080fd5b610c92868885610e02565b1515600114610ca057600080fd5b6000808660ff191660ff191681526020019081526020016000206040518060600160405280436001600160401b03168152602001426001600160401b0316815260200189815250908060018154018082558091505060019003906000526020600020906002020160009091909190915060008201518160000160006101000a8154816001600160401b0302191690836001600160401b0316021790555060208201518160000160086101000a8154816001600160401b0302191690836001600160401b031602179055506040820151816001015550507fbe8b490a2c0932bd7e7463b73a740738a6a4f6a99ac6e50aa73ddc788673c88d8543428a604051808560ff191660ff19168152602001846001600160401b03166001600160401b03168152602001836001600160401b03166001600160401b0316815260200182815260200194505050505060405180910390a150505050505050565b600192915050565b60019392505050565b60408051606081018252600080825260208201819052918101919091529056fe6e6f206d756c7469706c652073657420696e207468652073616d6520626c6f636ba264697066735822122084eb3fcfb546e07b0b8e4c23bd1c36134dea2916fb1535abd917a1eda8d580c264736f6c63430006010033"

// DeployState deploys a new Ethereum contract, binding an instance of State to it.
func DeployState(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *State, error) {
	parsed, err := abi.JSON(strings.NewReader(StateABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StateBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &State{StateCaller: StateCaller{contract: contract}, StateTransactor: StateTransactor{contract: contract}, StateFilterer: StateFilterer{contract: contract}}, nil
}

// State is an auto generated Go binding around an Ethereum contract.
type State struct {
	StateCaller     // Read-only binding to the contract
	StateTransactor // Write-only binding to the contract
	StateFilterer   // Log filterer for contract events
}

// StateCaller is an auto generated read-only Go binding around an Ethereum contract.
type StateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StateSession struct {
	Contract     *State            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StateCallerSession struct {
	Contract *StateCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StateTransactorSession struct {
	Contract     *StateTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateRaw is an auto generated low-level Go binding around an Ethereum contract.
type StateRaw struct {
	Contract *State // Generic contract binding to access the raw methods on
}

// StateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StateCallerRaw struct {
	Contract *StateCaller // Generic read-only contract binding to access the raw methods on
}

// StateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StateTransactorRaw struct {
	Contract *StateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewState creates a new instance of State, bound to a specific deployed contract.
func NewState(address common.Address, backend bind.ContractBackend) (*State, error) {
	contract, err := bindState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &State{StateCaller: StateCaller{contract: contract}, StateTransactor: StateTransactor{contract: contract}, StateFilterer: StateFilterer{contract: contract}}, nil
}

// NewStateCaller creates a new read-only instance of State, bound to a specific deployed contract.
func NewStateCaller(address common.Address, caller bind.ContractCaller) (*StateCaller, error) {
	contract, err := bindState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StateCaller{contract: contract}, nil
}

// NewStateTransactor creates a new write-only instance of State, bound to a specific deployed contract.
func NewStateTransactor(address common.Address, transactor bind.ContractTransactor) (*StateTransactor, error) {
	contract, err := bindState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StateTransactor{contract: contract}, nil
}

// NewStateFilterer creates a new log filterer instance of State, bound to a specific deployed contract.
func NewStateFilterer(address common.Address, filterer bind.ContractFilterer) (*StateFilterer, error) {
	contract, err := bindState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StateFilterer{contract: contract}, nil
}

// bindState binds a generic wrapper to an already deployed contract.
func bindState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StateABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_State *StateRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _State.Contract.StateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_State *StateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.Contract.StateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_State *StateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _State.Contract.StateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_State *StateCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _State.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_State *StateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _State.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_State *StateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _State.Contract.contract.Transact(opts, method, params...)
}

// GetState is a free data retrieval call binding the contract method 0xc1056e53.
//
// Solidity: function getState(bytes31 id) constant returns(bytes32)
func (_State *StateCaller) GetState(opts *bind.CallOpts, id [31]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _State.contract.Call(opts, out, "getState", id)
	return *ret0, err
}

// GetState is a free data retrieval call binding the contract method 0xc1056e53.
//
// Solidity: function getState(bytes31 id) constant returns(bytes32)
func (_State *StateSession) GetState(id [31]byte) ([32]byte, error) {
	return _State.Contract.GetState(&_State.CallOpts, id)
}

// GetState is a free data retrieval call binding the contract method 0xc1056e53.
//
// Solidity: function getState(bytes31 id) constant returns(bytes32)
func (_State *StateCallerSession) GetState(id [31]byte) ([32]byte, error) {
	return _State.Contract.GetState(&_State.CallOpts, id)
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xc68631e1.
//
// Solidity: function getStateDataByBlock(bytes31 id, uint64 blockN) constant returns(uint64, uint64, bytes32)
func (_State *StateCaller) GetStateDataByBlock(opts *bind.CallOpts, id [31]byte, blockN uint64) (uint64, uint64, [32]byte, error) {
	var (
		ret0 = new(uint64)
		ret1 = new(uint64)
		ret2 = new([32]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _State.contract.Call(opts, out, "getStateDataByBlock", id, blockN)
	return *ret0, *ret1, *ret2, err
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xc68631e1.
//
// Solidity: function getStateDataByBlock(bytes31 id, uint64 blockN) constant returns(uint64, uint64, bytes32)
func (_State *StateSession) GetStateDataByBlock(id [31]byte, blockN uint64) (uint64, uint64, [32]byte, error) {
	return _State.Contract.GetStateDataByBlock(&_State.CallOpts, id, blockN)
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xc68631e1.
//
// Solidity: function getStateDataByBlock(bytes31 id, uint64 blockN) constant returns(uint64, uint64, bytes32)
func (_State *StateCallerSession) GetStateDataByBlock(id [31]byte, blockN uint64) (uint64, uint64, [32]byte, error) {
	return _State.Contract.GetStateDataByBlock(&_State.CallOpts, id, blockN)
}

// GetStateDataById is a free data retrieval call binding the contract method 0x4cabaefa.
//
// Solidity: function getStateDataById(bytes31 id) constant returns(uint64, uint64, bytes32)
func (_State *StateCaller) GetStateDataById(opts *bind.CallOpts, id [31]byte) (uint64, uint64, [32]byte, error) {
	var (
		ret0 = new(uint64)
		ret1 = new(uint64)
		ret2 = new([32]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _State.contract.Call(opts, out, "getStateDataById", id)
	return *ret0, *ret1, *ret2, err
}

// GetStateDataById is a free data retrieval call binding the contract method 0x4cabaefa.
//
// Solidity: function getStateDataById(bytes31 id) constant returns(uint64, uint64, bytes32)
func (_State *StateSession) GetStateDataById(id [31]byte) (uint64, uint64, [32]byte, error) {
	return _State.Contract.GetStateDataById(&_State.CallOpts, id)
}

// GetStateDataById is a free data retrieval call binding the contract method 0x4cabaefa.
//
// Solidity: function getStateDataById(bytes31 id) constant returns(uint64, uint64, bytes32)
func (_State *StateCallerSession) GetStateDataById(id [31]byte) (uint64, uint64, [32]byte, error) {
	return _State.Contract.GetStateDataById(&_State.CallOpts, id)
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x5710773a.
//
// Solidity: function getStateDataByTime(bytes31 id, uint64 timestamp) constant returns(uint64, uint64, bytes32)
func (_State *StateCaller) GetStateDataByTime(opts *bind.CallOpts, id [31]byte, timestamp uint64) (uint64, uint64, [32]byte, error) {
	var (
		ret0 = new(uint64)
		ret1 = new(uint64)
		ret2 = new([32]byte)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _State.contract.Call(opts, out, "getStateDataByTime", id, timestamp)
	return *ret0, *ret1, *ret2, err
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x5710773a.
//
// Solidity: function getStateDataByTime(bytes31 id, uint64 timestamp) constant returns(uint64, uint64, bytes32)
func (_State *StateSession) GetStateDataByTime(id [31]byte, timestamp uint64) (uint64, uint64, [32]byte, error) {
	return _State.Contract.GetStateDataByTime(&_State.CallOpts, id, timestamp)
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x5710773a.
//
// Solidity: function getStateDataByTime(bytes31 id, uint64 timestamp) constant returns(uint64, uint64, bytes32)
func (_State *StateCallerSession) GetStateDataByTime(id [31]byte, timestamp uint64) (uint64, uint64, [32]byte, error) {
	return _State.Contract.GetStateDataByTime(&_State.CallOpts, id, timestamp)
}

// InitState is a paid mutator transaction binding the contract method 0x9b2d7aac.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sigR8, bytes32 sigS) returns()
func (_State *StateTransactor) InitState(opts *bind.TransactOpts, newState [32]byte, genesisState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sigR8 [32]byte, sigS [32]byte) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "initState", newState, genesisState, id, kOpProof, itp, sigR8, sigS)
}

// InitState is a paid mutator transaction binding the contract method 0x9b2d7aac.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sigR8, bytes32 sigS) returns()
func (_State *StateSession) InitState(newState [32]byte, genesisState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sigR8 [32]byte, sigS [32]byte) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, kOpProof, itp, sigR8, sigS)
}

// InitState is a paid mutator transaction binding the contract method 0x9b2d7aac.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sigR8, bytes32 sigS) returns()
func (_State *StateTransactorSession) InitState(newState [32]byte, genesisState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sigR8 [32]byte, sigS [32]byte) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, kOpProof, itp, sigR8, sigS)
}

// SetState is a paid mutator transaction binding the contract method 0xea1662de.
//
// Solidity: function setState(bytes32 newState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sigR8, bytes32 sigS) returns()
func (_State *StateTransactor) SetState(opts *bind.TransactOpts, newState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sigR8 [32]byte, sigS [32]byte) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "setState", newState, id, kOpProof, itp, sigR8, sigS)
}

// SetState is a paid mutator transaction binding the contract method 0xea1662de.
//
// Solidity: function setState(bytes32 newState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sigR8, bytes32 sigS) returns()
func (_State *StateSession) SetState(newState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sigR8 [32]byte, sigS [32]byte) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, kOpProof, itp, sigR8, sigS)
}

// SetState is a paid mutator transaction binding the contract method 0xea1662de.
//
// Solidity: function setState(bytes32 newState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sigR8, bytes32 sigS) returns()
func (_State *StateTransactorSession) SetState(newState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sigR8 [32]byte, sigS [32]byte) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, kOpProof, itp, sigR8, sigS)
}

// StateStateUpdatedIterator is returned from FilterStateUpdated and is used to iterate over the raw logs and unpacked data for StateUpdated events raised by the State contract.
type StateStateUpdatedIterator struct {
	Event *StateStateUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StateStateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateStateUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StateStateUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StateStateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateStateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateStateUpdated represents a StateUpdated event raised by the State contract.
type StateStateUpdated struct {
	Id        [31]byte
	BlockN    uint64
	Timestamp uint64
	State     [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateUpdated is a free log retrieval operation binding the contract event 0xbe8b490a2c0932bd7e7463b73a740738a6a4f6a99ac6e50aa73ddc788673c88d.
//
// Solidity: event StateUpdated(bytes31 id, uint64 blockN, uint64 timestamp, bytes32 state)
func (_State *StateFilterer) FilterStateUpdated(opts *bind.FilterOpts) (*StateStateUpdatedIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return &StateStateUpdatedIterator{contract: _State.contract, event: "StateUpdated", logs: logs, sub: sub}, nil
}

// WatchStateUpdated is a free log subscription operation binding the contract event 0xbe8b490a2c0932bd7e7463b73a740738a6a4f6a99ac6e50aa73ddc788673c88d.
//
// Solidity: event StateUpdated(bytes31 id, uint64 blockN, uint64 timestamp, bytes32 state)
func (_State *StateFilterer) WatchStateUpdated(opts *bind.WatchOpts, sink chan<- *StateStateUpdated) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateStateUpdated)
				if err := _State.contract.UnpackLog(event, "StateUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStateUpdated is a log parse operation binding the contract event 0xbe8b490a2c0932bd7e7463b73a740738a6a4f6a99ac6e50aa73ddc788673c88d.
//
// Solidity: event StateUpdated(bytes31 id, uint64 blockN, uint64 timestamp, bytes32 state)
func (_State *StateFilterer) ParseStateUpdated(log types.Log) (*StateStateUpdated, error) {
	event := new(StateStateUpdated)
	if err := _State.contract.UnpackLog(event, "StateUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
