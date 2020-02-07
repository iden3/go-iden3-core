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
const StateABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"state\",\"type\":\"bytes32\"}],\"name\":\"StateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"}],\"name\":\"getStateByBlock\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getStateByTime\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getStateDataById\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getStateFromId\",\"outputs\":[{\"internalType\":\"bytes27\",\"name\":\"\",\"type\":\"bytes27\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"genesisState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"bytes\",\"name\":\"kOpProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"itp\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"sig\",\"type\":\"bytes32\"}],\"name\":\"initState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"bytes\",\"name\":\"kOpProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"itp\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"sig\",\"type\":\"bytes32\"}],\"name\":\"setState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// StateFuncSigs maps the 4-byte function signature to its string representation.
var StateFuncSigs = map[string]string{
	"c1056e53": "getState(bytes31)",
	"b812400d": "getStateByBlock(bytes31,uint64)",
	"415e20cc": "getStateByTime(bytes31,uint64)",
	"4cabaefa": "getStateDataById(bytes31)",
	"b2c7f012": "getStateFromId(bytes31)",
	"402894c4": "initState(bytes32,bytes32,bytes31,bytes,bytes,bytes32)",
	"6d94975b": "setState(bytes32,bytes31,bytes,bytes,bytes32)",
}

// StateBin is the compiled bytecode used for deploying new contracts.
var StateBin = "0x608060405234801561001057600080fd5b50610dbd806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80636d94975b1161005b5780636d94975b14610253578063b2c7f0121461038e578063b812400d146103ca578063c1056e53146103fa5761007d565b8063402894c414610082578063415e20cc146101c65780634cabaefa14610208575b600080fd5b6101c4600480360360c081101561009857600080fd5b81359160208101359160ff196040830135169190810190608081016060820135600160201b8111156100c957600080fd5b8201836020820111156100db57600080fd5b803590602001918460018302840111600160201b831117156100fc57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561014e57600080fd5b82018360208201111561016057600080fd5b803590602001918460018302840111600160201b8311171561018157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550509135925061041b915050565b005b6101f6600480360360408110156101dc57600080fd5b50803560ff191690602001356001600160401b031661044f565b60408051918252519081900360200190f35b6102296004803603602081101561021e57600080fd5b503560ff1916610733565b604080516001600160401b0394851681529290931660208301528183015290519081900360600190f35b6101c4600480360360a081101561026957600080fd5b81359160ff1960208201351691810190606081016040820135600160201b81111561029357600080fd5b8201836020820111156102a557600080fd5b803590602001918460018302840111600160201b831117156102c657600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561031857600080fd5b82018360208201111561032a57600080fd5b803590602001918460018302840111600160201b8311171561034b57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955050913592506107e3915050565b6103af600480360360208110156103a457600080fd5b503560ff19166108c2565b6040805164ffffffffff199092168252519081900360200190f35b6101f6600480360360408110156103e057600080fd5b50803560ff191690602001356001600160401b03166108d3565b6101f66004803603602081101561041057600080fd5b503560ff1916610b4e565b60ff1984166000908152602081905260409020541561043957600080fd5b610447868686868686610bab565b505050505050565b600042826001600160401b0316106104a3576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff1983166000908152602081905260409020546104c4575060015461072d565b60ff1983166000908152602081905260408120805460001981019081106104e757fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116915083168110156105525760ff19841660009081526020819052604090208054600019810190811061053857fe5b90600052602060002090600202016001015491505061072d565b60ff198416600090815260208190526040812054600019015b8082116107245760ff19861660009081526020819052604090208054600284840104916001600160401b03881691839081106105a357fe5b6000918252602090912060029091020154600160401b90046001600160401b031614156106065760ff19871660009081526020819052604090208054829081106105e957fe5b90600052602060002090600202016001015494505050505061072d565b60ff198716600090815260208190526040902080548290811061062557fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690871611801561069d575060ff198716600090815260208190526040902080546001830190811061067557fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908716105b156106c15760ff19871660009081526020819052604090208054829081106105e957fe5b60ff19871660009081526020819052604090208054829081106106e057fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690871611156107175780600101925061071e565b6001810391505b5061056b565b60015493505050505b92915050565b60ff1981166000908152602081905260408120548190819061075f5750506001546000915081906107dc565b610767610d46565b60ff19851660009081526020819052604090208054600019810190811061078a57fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b909204169484018590526001909101549290910182905295509093509150505b9193909250565b60ff19841660009081526020819052604090205461080057600080fd5b610808610d46565b60ff19851660009081526020819052604090208054600019810190811061082b57fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b9092041694840194909452600101549082015291504314156108b05760405162461bcd60e51b8152600401808060200182810382526021815260200180610d676021913960400191505060405180910390fd5b61044786826040015187878787610bab565b62ffffff19601082901b165b919050565b600043826001600160401b031610610927576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff198316600090815260208190526040902054610948575060015461072d565b60ff19831660009081526020819052604081208054600019810190811061096b57fe5b60009182526020909120600290910201546001600160401b03908116915083168110156109b55760ff19841660009081526020819052604090208054600019810190811061053857fe5b60ff198416600090815260208190526040812054600019015b8082116107245760ff19861660009081526020819052604090208054600284840104916001600160401b0388169183908110610a0657fe5b60009182526020909120600290910201546001600160401b03161415610a455760ff19871660009081526020819052604090208054829081106105e957fe5b60ff1987166000908152602081905260409020805482908110610a6457fe5b60009182526020909120600290910201546001600160401b03908116908716118015610ace575060ff1987166000908152602081905260409020805460018301908110610aad57fe5b60009182526020909120600290910201546001600160401b03908116908716105b15610af25760ff19871660009081526020819052604090208054829081106105e957fe5b60ff1987166000908152602081905260409020805482908110610b1157fe5b60009182526020909120600290910201546001600160401b039081169087161115610b4157806001019250610b48565b6001810391505b506109ce565b60ff198116600090815260208190526040812054610b6f57506001546108ce565b60ff198216600090815260208190526040902080546000198101908110610b9257fe5b9060005260206000209060020201600101549050919050565b610bb58684610d35565b1515600114610bc357600080fd5b610bce858784610d3d565b1515600114610bdc57600080fd5b6000808560ff191660ff191681526020019081526020016000206040518060600160405280436001600160401b03168152602001426001600160401b0316815260200188815250908060018154018082558091505060019003906000526020600020906002020160009091909190915060008201518160000160006101000a8154816001600160401b0302191690836001600160401b0316021790555060208201518160000160086101000a8154816001600160401b0302191690836001600160401b031602179055506040820151816001015550507fbe8b490a2c0932bd7e7463b73a740738a6a4f6a99ac6e50aa73ddc788673c88d84434289604051808560ff191660ff19168152602001846001600160401b03166001600160401b03168152602001836001600160401b03166001600160401b0316815260200182815260200194505050505060405180910390a1505050505050565b600192915050565b60019392505050565b60408051606081018252600080825260208201819052918101919091529056fe6e6f206d756c7469706c652073657420696e207468652073616d6520626c6f636ba26469706673582212202517a0cec26968359c4645b918aa6addeef1f55bbb3deccd62ced579a07ad6d064736f6c63430006010033"

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

// GetStateByBlock is a free data retrieval call binding the contract method 0xb812400d.
//
// Solidity: function getStateByBlock(bytes31 id, uint64 blockN) constant returns(bytes32)
func (_State *StateCaller) GetStateByBlock(opts *bind.CallOpts, id [31]byte, blockN uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _State.contract.Call(opts, out, "getStateByBlock", id, blockN)
	return *ret0, err
}

// GetStateByBlock is a free data retrieval call binding the contract method 0xb812400d.
//
// Solidity: function getStateByBlock(bytes31 id, uint64 blockN) constant returns(bytes32)
func (_State *StateSession) GetStateByBlock(id [31]byte, blockN uint64) ([32]byte, error) {
	return _State.Contract.GetStateByBlock(&_State.CallOpts, id, blockN)
}

// GetStateByBlock is a free data retrieval call binding the contract method 0xb812400d.
//
// Solidity: function getStateByBlock(bytes31 id, uint64 blockN) constant returns(bytes32)
func (_State *StateCallerSession) GetStateByBlock(id [31]byte, blockN uint64) ([32]byte, error) {
	return _State.Contract.GetStateByBlock(&_State.CallOpts, id, blockN)
}

// GetStateByTime is a free data retrieval call binding the contract method 0x415e20cc.
//
// Solidity: function getStateByTime(bytes31 id, uint64 timestamp) constant returns(bytes32)
func (_State *StateCaller) GetStateByTime(opts *bind.CallOpts, id [31]byte, timestamp uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _State.contract.Call(opts, out, "getStateByTime", id, timestamp)
	return *ret0, err
}

// GetStateByTime is a free data retrieval call binding the contract method 0x415e20cc.
//
// Solidity: function getStateByTime(bytes31 id, uint64 timestamp) constant returns(bytes32)
func (_State *StateSession) GetStateByTime(id [31]byte, timestamp uint64) ([32]byte, error) {
	return _State.Contract.GetStateByTime(&_State.CallOpts, id, timestamp)
}

// GetStateByTime is a free data retrieval call binding the contract method 0x415e20cc.
//
// Solidity: function getStateByTime(bytes31 id, uint64 timestamp) constant returns(bytes32)
func (_State *StateCallerSession) GetStateByTime(id [31]byte, timestamp uint64) ([32]byte, error) {
	return _State.Contract.GetStateByTime(&_State.CallOpts, id, timestamp)
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

// GetStateFromId is a free data retrieval call binding the contract method 0xb2c7f012.
//
// Solidity: function getStateFromId(bytes31 id) constant returns(bytes27)
func (_State *StateCaller) GetStateFromId(opts *bind.CallOpts, id [31]byte) ([27]byte, error) {
	var (
		ret0 = new([27]byte)
	)
	out := ret0
	err := _State.contract.Call(opts, out, "getStateFromId", id)
	return *ret0, err
}

// GetStateFromId is a free data retrieval call binding the contract method 0xb2c7f012.
//
// Solidity: function getStateFromId(bytes31 id) constant returns(bytes27)
func (_State *StateSession) GetStateFromId(id [31]byte) ([27]byte, error) {
	return _State.Contract.GetStateFromId(&_State.CallOpts, id)
}

// GetStateFromId is a free data retrieval call binding the contract method 0xb2c7f012.
//
// Solidity: function getStateFromId(bytes31 id) constant returns(bytes27)
func (_State *StateCallerSession) GetStateFromId(id [31]byte) ([27]byte, error) {
	return _State.Contract.GetStateFromId(&_State.CallOpts, id)
}

// InitState is a paid mutator transaction binding the contract method 0x402894c4.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sig) returns()
func (_State *StateTransactor) InitState(opts *bind.TransactOpts, newState [32]byte, genesisState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sig [32]byte) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "initState", newState, genesisState, id, kOpProof, itp, sig)
}

// InitState is a paid mutator transaction binding the contract method 0x402894c4.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sig) returns()
func (_State *StateSession) InitState(newState [32]byte, genesisState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sig [32]byte) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, kOpProof, itp, sig)
}

// InitState is a paid mutator transaction binding the contract method 0x402894c4.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sig) returns()
func (_State *StateTransactorSession) InitState(newState [32]byte, genesisState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sig [32]byte) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, kOpProof, itp, sig)
}

// SetState is a paid mutator transaction binding the contract method 0x6d94975b.
//
// Solidity: function setState(bytes32 newState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sig) returns()
func (_State *StateTransactor) SetState(opts *bind.TransactOpts, newState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sig [32]byte) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "setState", newState, id, kOpProof, itp, sig)
}

// SetState is a paid mutator transaction binding the contract method 0x6d94975b.
//
// Solidity: function setState(bytes32 newState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sig) returns()
func (_State *StateSession) SetState(newState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sig [32]byte) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, kOpProof, itp, sig)
}

// SetState is a paid mutator transaction binding the contract method 0x6d94975b.
//
// Solidity: function setState(bytes32 newState, bytes31 id, bytes kOpProof, bytes itp, bytes32 sig) returns()
func (_State *StateTransactorSession) SetState(newState [32]byte, id [31]byte, kOpProof []byte, itp []byte, sig [32]byte) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, kOpProof, itp, sig)
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
