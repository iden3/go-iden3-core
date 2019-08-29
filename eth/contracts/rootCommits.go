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

// Iden3HelpersABI is the input ABI used to generate the binding from.
const Iden3HelpersABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\"},{\"name\":\"rsv\",\"type\":\"bytes\"}],\"name\":\"checkSig\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getRootFromId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes27\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_mimcContractAddr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// Iden3HelpersFuncSigs maps the 4-byte function signature to its string representation.
var Iden3HelpersFuncSigs = map[string]string{
	"01b0452c": "checkSig(bytes32,bytes)",
	"ad05a8d2": "getRootFromId(bytes31)",
}

// Iden3HelpersBin is the compiled bytecode used for deploying new contracts.
var Iden3HelpersBin = "0x608060405234801561001057600080fd5b506040516102653803806102658339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b0319909216919091179055610200806100656000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806301b0452c1461003b578063ad05a8d214610104575b600080fd5b6100e86004803603604081101561005157600080fd5b8135919081019060408101602082013564010000000081111561007357600080fd5b82018360208201111561008557600080fd5b803590602001918460018302840111640100000000831117156100a757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610140945050505050565b604080516001600160a01b039092168252519081900360200190f35b6101256004803603602081101561011a57600080fd5b503560ff19166101bf565b6040805164ffffffffff199092168252519081900360200190f35b602081810151604080840151606080860151835160008082528188018087528a905291821a81860181905292810186905260808101849052935190959293919260019260a080820193601f1981019281900390910190855afa1580156101aa573d6000803e3d6000fd5b5050604051601f190151979650505050505050565b60101b62ffffff19169056fea265627a7a723058200c3b202c57cc47fc64f71c113cc87cd8f9f3e1c4c225db3df887d89dc077f82664736f6c634300050a0032"

// DeployIden3Helpers deploys a new Ethereum contract, binding an instance of Iden3Helpers to it.
func DeployIden3Helpers(auth *bind.TransactOpts, backend bind.ContractBackend, _mimcContractAddr common.Address) (common.Address, *types.Transaction, *Iden3Helpers, error) {
	parsed, err := abi.JSON(strings.NewReader(Iden3HelpersABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Iden3HelpersBin), backend, _mimcContractAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Iden3Helpers{Iden3HelpersCaller: Iden3HelpersCaller{contract: contract}, Iden3HelpersTransactor: Iden3HelpersTransactor{contract: contract}, Iden3HelpersFilterer: Iden3HelpersFilterer{contract: contract}}, nil
}

// Iden3Helpers is an auto generated Go binding around an Ethereum contract.
type Iden3Helpers struct {
	Iden3HelpersCaller     // Read-only binding to the contract
	Iden3HelpersTransactor // Write-only binding to the contract
	Iden3HelpersFilterer   // Log filterer for contract events
}

// Iden3HelpersCaller is an auto generated read-only Go binding around an Ethereum contract.
type Iden3HelpersCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Iden3HelpersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Iden3HelpersTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Iden3HelpersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Iden3HelpersFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Iden3HelpersSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Iden3HelpersSession struct {
	Contract     *Iden3Helpers     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Iden3HelpersCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Iden3HelpersCallerSession struct {
	Contract *Iden3HelpersCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// Iden3HelpersTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Iden3HelpersTransactorSession struct {
	Contract     *Iden3HelpersTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// Iden3HelpersRaw is an auto generated low-level Go binding around an Ethereum contract.
type Iden3HelpersRaw struct {
	Contract *Iden3Helpers // Generic contract binding to access the raw methods on
}

// Iden3HelpersCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Iden3HelpersCallerRaw struct {
	Contract *Iden3HelpersCaller // Generic read-only contract binding to access the raw methods on
}

// Iden3HelpersTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Iden3HelpersTransactorRaw struct {
	Contract *Iden3HelpersTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIden3Helpers creates a new instance of Iden3Helpers, bound to a specific deployed contract.
func NewIden3Helpers(address common.Address, backend bind.ContractBackend) (*Iden3Helpers, error) {
	contract, err := bindIden3Helpers(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Iden3Helpers{Iden3HelpersCaller: Iden3HelpersCaller{contract: contract}, Iden3HelpersTransactor: Iden3HelpersTransactor{contract: contract}, Iden3HelpersFilterer: Iden3HelpersFilterer{contract: contract}}, nil
}

// NewIden3HelpersCaller creates a new read-only instance of Iden3Helpers, bound to a specific deployed contract.
func NewIden3HelpersCaller(address common.Address, caller bind.ContractCaller) (*Iden3HelpersCaller, error) {
	contract, err := bindIden3Helpers(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Iden3HelpersCaller{contract: contract}, nil
}

// NewIden3HelpersTransactor creates a new write-only instance of Iden3Helpers, bound to a specific deployed contract.
func NewIden3HelpersTransactor(address common.Address, transactor bind.ContractTransactor) (*Iden3HelpersTransactor, error) {
	contract, err := bindIden3Helpers(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Iden3HelpersTransactor{contract: contract}, nil
}

// NewIden3HelpersFilterer creates a new log filterer instance of Iden3Helpers, bound to a specific deployed contract.
func NewIden3HelpersFilterer(address common.Address, filterer bind.ContractFilterer) (*Iden3HelpersFilterer, error) {
	contract, err := bindIden3Helpers(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Iden3HelpersFilterer{contract: contract}, nil
}

// bindIden3Helpers binds a generic wrapper to an already deployed contract.
func bindIden3Helpers(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Iden3HelpersABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Iden3Helpers *Iden3HelpersRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Iden3Helpers.Contract.Iden3HelpersCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Iden3Helpers *Iden3HelpersRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Iden3Helpers.Contract.Iden3HelpersTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Iden3Helpers *Iden3HelpersRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Iden3Helpers.Contract.Iden3HelpersTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Iden3Helpers *Iden3HelpersCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Iden3Helpers.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Iden3Helpers *Iden3HelpersTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Iden3Helpers.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Iden3Helpers *Iden3HelpersTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Iden3Helpers.Contract.contract.Transact(opts, method, params...)
}

// CheckSig is a free data retrieval call binding the contract method 0x01b0452c.
//
// Solidity: function checkSig(bytes32 msgHash, bytes rsv) constant returns(address)
func (_Iden3Helpers *Iden3HelpersCaller) CheckSig(opts *bind.CallOpts, msgHash [32]byte, rsv []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Iden3Helpers.contract.Call(opts, out, "checkSig", msgHash, rsv)
	return *ret0, err
}

// CheckSig is a free data retrieval call binding the contract method 0x01b0452c.
//
// Solidity: function checkSig(bytes32 msgHash, bytes rsv) constant returns(address)
func (_Iden3Helpers *Iden3HelpersSession) CheckSig(msgHash [32]byte, rsv []byte) (common.Address, error) {
	return _Iden3Helpers.Contract.CheckSig(&_Iden3Helpers.CallOpts, msgHash, rsv)
}

// CheckSig is a free data retrieval call binding the contract method 0x01b0452c.
//
// Solidity: function checkSig(bytes32 msgHash, bytes rsv) constant returns(address)
func (_Iden3Helpers *Iden3HelpersCallerSession) CheckSig(msgHash [32]byte, rsv []byte) (common.Address, error) {
	return _Iden3Helpers.Contract.CheckSig(&_Iden3Helpers.CallOpts, msgHash, rsv)
}

// GetRootFromId is a free data retrieval call binding the contract method 0xad05a8d2.
//
// Solidity: function getRootFromId(bytes31 id) constant returns(bytes27)
func (_Iden3Helpers *Iden3HelpersCaller) GetRootFromId(opts *bind.CallOpts, id [31]byte) ([27]byte, error) {
	var (
		ret0 = new([27]byte)
	)
	out := ret0
	err := _Iden3Helpers.contract.Call(opts, out, "getRootFromId", id)
	return *ret0, err
}

// GetRootFromId is a free data retrieval call binding the contract method 0xad05a8d2.
//
// Solidity: function getRootFromId(bytes31 id) constant returns(bytes27)
func (_Iden3Helpers *Iden3HelpersSession) GetRootFromId(id [31]byte) ([27]byte, error) {
	return _Iden3Helpers.Contract.GetRootFromId(&_Iden3Helpers.CallOpts, id)
}

// GetRootFromId is a free data retrieval call binding the contract method 0xad05a8d2.
//
// Solidity: function getRootFromId(bytes31 id) constant returns(bytes27)
func (_Iden3Helpers *Iden3HelpersCallerSession) GetRootFromId(id [31]byte) ([27]byte, error) {
	return _Iden3Helpers.Contract.GetRootFromId(&_Iden3Helpers.CallOpts, id)
}

// MemoryABI is the input ABI used to generate the binding from.
const MemoryABI = "[]"

// MemoryBin is the compiled bytecode used for deploying new contracts.
var MemoryBin = "0x60556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a7230582045f8113dc5e8da564c8bdc6bea4dc57e886c6c74e19ac40d37f3e9c4e8d44b8c64736f6c634300050a0032"

// DeployMemory deploys a new Ethereum contract, binding an instance of Memory to it.
func DeployMemory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Memory, error) {
	parsed, err := abi.JSON(strings.NewReader(MemoryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MemoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Memory{MemoryCaller: MemoryCaller{contract: contract}, MemoryTransactor: MemoryTransactor{contract: contract}, MemoryFilterer: MemoryFilterer{contract: contract}}, nil
}

// Memory is an auto generated Go binding around an Ethereum contract.
type Memory struct {
	MemoryCaller     // Read-only binding to the contract
	MemoryTransactor // Write-only binding to the contract
	MemoryFilterer   // Log filterer for contract events
}

// MemoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type MemoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MemoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MemoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MemoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MemoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MemorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MemorySession struct {
	Contract     *Memory           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MemoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MemoryCallerSession struct {
	Contract *MemoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MemoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MemoryTransactorSession struct {
	Contract     *MemoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MemoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type MemoryRaw struct {
	Contract *Memory // Generic contract binding to access the raw methods on
}

// MemoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MemoryCallerRaw struct {
	Contract *MemoryCaller // Generic read-only contract binding to access the raw methods on
}

// MemoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MemoryTransactorRaw struct {
	Contract *MemoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMemory creates a new instance of Memory, bound to a specific deployed contract.
func NewMemory(address common.Address, backend bind.ContractBackend) (*Memory, error) {
	contract, err := bindMemory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Memory{MemoryCaller: MemoryCaller{contract: contract}, MemoryTransactor: MemoryTransactor{contract: contract}, MemoryFilterer: MemoryFilterer{contract: contract}}, nil
}

// NewMemoryCaller creates a new read-only instance of Memory, bound to a specific deployed contract.
func NewMemoryCaller(address common.Address, caller bind.ContractCaller) (*MemoryCaller, error) {
	contract, err := bindMemory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MemoryCaller{contract: contract}, nil
}

// NewMemoryTransactor creates a new write-only instance of Memory, bound to a specific deployed contract.
func NewMemoryTransactor(address common.Address, transactor bind.ContractTransactor) (*MemoryTransactor, error) {
	contract, err := bindMemory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MemoryTransactor{contract: contract}, nil
}

// NewMemoryFilterer creates a new log filterer instance of Memory, bound to a specific deployed contract.
func NewMemoryFilterer(address common.Address, filterer bind.ContractFilterer) (*MemoryFilterer, error) {
	contract, err := bindMemory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MemoryFilterer{contract: contract}, nil
}

// bindMemory binds a generic wrapper to an already deployed contract.
func bindMemory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MemoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Memory *MemoryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Memory.Contract.MemoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Memory *MemoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Memory.Contract.MemoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Memory *MemoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Memory.Contract.MemoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Memory *MemoryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Memory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Memory *MemoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Memory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Memory *MemoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Memory.Contract.contract.Transact(opts, method, params...)
}

// MimcUnitABI is the input ABI used to generate the binding from.
const MimcUnitABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"MiMCpe7\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// MimcUnitFuncSigs maps the 4-byte function signature to its string representation.
var MimcUnitFuncSigs = map[string]string{
	"d15ca109": "MiMCpe7(uint256,uint256)",
}

// MimcUnitBin is the compiled bytecode used for deploying new contracts.
var MimcUnitBin = "0x6080604052348015600f57600080fd5b50609c8061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063d15ca10914602d575b600080fd5b604d60048036036040811015604157600080fd5b5080359060200135605f565b60408051918252519081900360200190f35b60009291505056fea265627a7a72305820117562ad7b874f2f4dc1ff83201e81c18d5a43ed1ba319c83a293d24f8cb5ffb64736f6c634300050a0032"

// DeployMimcUnit deploys a new Ethereum contract, binding an instance of MimcUnit to it.
func DeployMimcUnit(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MimcUnit, error) {
	parsed, err := abi.JSON(strings.NewReader(MimcUnitABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MimcUnitBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MimcUnit{MimcUnitCaller: MimcUnitCaller{contract: contract}, MimcUnitTransactor: MimcUnitTransactor{contract: contract}, MimcUnitFilterer: MimcUnitFilterer{contract: contract}}, nil
}

// MimcUnit is an auto generated Go binding around an Ethereum contract.
type MimcUnit struct {
	MimcUnitCaller     // Read-only binding to the contract
	MimcUnitTransactor // Write-only binding to the contract
	MimcUnitFilterer   // Log filterer for contract events
}

// MimcUnitCaller is an auto generated read-only Go binding around an Ethereum contract.
type MimcUnitCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MimcUnitTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MimcUnitTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MimcUnitFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MimcUnitFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MimcUnitSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MimcUnitSession struct {
	Contract     *MimcUnit         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MimcUnitCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MimcUnitCallerSession struct {
	Contract *MimcUnitCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MimcUnitTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MimcUnitTransactorSession struct {
	Contract     *MimcUnitTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MimcUnitRaw is an auto generated low-level Go binding around an Ethereum contract.
type MimcUnitRaw struct {
	Contract *MimcUnit // Generic contract binding to access the raw methods on
}

// MimcUnitCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MimcUnitCallerRaw struct {
	Contract *MimcUnitCaller // Generic read-only contract binding to access the raw methods on
}

// MimcUnitTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MimcUnitTransactorRaw struct {
	Contract *MimcUnitTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMimcUnit creates a new instance of MimcUnit, bound to a specific deployed contract.
func NewMimcUnit(address common.Address, backend bind.ContractBackend) (*MimcUnit, error) {
	contract, err := bindMimcUnit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MimcUnit{MimcUnitCaller: MimcUnitCaller{contract: contract}, MimcUnitTransactor: MimcUnitTransactor{contract: contract}, MimcUnitFilterer: MimcUnitFilterer{contract: contract}}, nil
}

// NewMimcUnitCaller creates a new read-only instance of MimcUnit, bound to a specific deployed contract.
func NewMimcUnitCaller(address common.Address, caller bind.ContractCaller) (*MimcUnitCaller, error) {
	contract, err := bindMimcUnit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MimcUnitCaller{contract: contract}, nil
}

// NewMimcUnitTransactor creates a new write-only instance of MimcUnit, bound to a specific deployed contract.
func NewMimcUnitTransactor(address common.Address, transactor bind.ContractTransactor) (*MimcUnitTransactor, error) {
	contract, err := bindMimcUnit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MimcUnitTransactor{contract: contract}, nil
}

// NewMimcUnitFilterer creates a new log filterer instance of MimcUnit, bound to a specific deployed contract.
func NewMimcUnitFilterer(address common.Address, filterer bind.ContractFilterer) (*MimcUnitFilterer, error) {
	contract, err := bindMimcUnit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MimcUnitFilterer{contract: contract}, nil
}

// bindMimcUnit binds a generic wrapper to an already deployed contract.
func bindMimcUnit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MimcUnitABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MimcUnit *MimcUnitRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MimcUnit.Contract.MimcUnitCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MimcUnit *MimcUnitRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MimcUnit.Contract.MimcUnitTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MimcUnit *MimcUnitRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MimcUnit.Contract.MimcUnitTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MimcUnit *MimcUnitCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MimcUnit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MimcUnit *MimcUnitTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MimcUnit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MimcUnit *MimcUnitTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MimcUnit.Contract.contract.Transact(opts, method, params...)
}

// MiMCpe7 is a free data retrieval call binding the contract method 0xd15ca109.
//
// Solidity: function MiMCpe7(uint256 , uint256 ) constant returns(uint256)
func (_MimcUnit *MimcUnitCaller) MiMCpe7(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MimcUnit.contract.Call(opts, out, "MiMCpe7", arg0, arg1)
	return *ret0, err
}

// MiMCpe7 is a free data retrieval call binding the contract method 0xd15ca109.
//
// Solidity: function MiMCpe7(uint256 , uint256 ) constant returns(uint256)
func (_MimcUnit *MimcUnitSession) MiMCpe7(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MimcUnit.Contract.MiMCpe7(&_MimcUnit.CallOpts, arg0, arg1)
}

// MiMCpe7 is a free data retrieval call binding the contract method 0xd15ca109.
//
// Solidity: function MiMCpe7(uint256 , uint256 ) constant returns(uint256)
func (_MimcUnit *MimcUnitCallerSession) MiMCpe7(arg0 *big.Int, arg1 *big.Int) (*big.Int, error) {
	return _MimcUnit.Contract.MiMCpe7(&_MimcUnit.CallOpts, arg0, arg1)
}

// RootCommitsABI is the input ABI used to generate the binding from.
const RootCommitsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\"},{\"name\":\"rsv\",\"type\":\"bytes\"}],\"name\":\"checkSig\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"},{\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getRootByTime\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getRootFromId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes27\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getRootDataById\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"},{\"name\":\"\",\"type\":\"uint64\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"},{\"name\":\"blockN\",\"type\":\"uint64\"}],\"name\":\"getRootByBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newRoot\",\"type\":\"bytes32\"},{\"name\":\"id\",\"type\":\"bytes31\"},{\"name\":\"mtp\",\"type\":\"bytes\"}],\"name\":\"setRoot\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getRoot\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_mimcContractAddr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes31\"},{\"indexed\":false,\"name\":\"blockN\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootUpdated\",\"type\":\"event\"}]"

// RootCommitsFuncSigs maps the 4-byte function signature to its string representation.
var RootCommitsFuncSigs = map[string]string{
	"01b0452c": "checkSig(bytes32,bytes)",
	"fead90d7": "getRoot(bytes31)",
	"b816ff6f": "getRootByBlock(bytes31,uint64)",
	"4175dae5": "getRootByTime(bytes31,uint64)",
	"b5818699": "getRootDataById(bytes31)",
	"ad05a8d2": "getRootFromId(bytes31)",
	"e0681acd": "setRoot(bytes32,bytes31,bytes)",
}

// RootCommitsBin is the compiled bytecode used for deploying new contracts.
var RootCommitsBin = "0x608060405234801561001057600080fd5b5060405161133f38038061133f8339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b03199092169190911790556112da806100656000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063b58186991161005b578063b5818699146101c9578063b816ff6f14610214578063e0681acd14610244578063fead90d7146102fc5761007d565b806301b0452c146100825780634175dae51461014b578063ad05a8d21461018d575b600080fd5b61012f6004803603604081101561009857600080fd5b813591908101906040810160208201356401000000008111156100ba57600080fd5b8201836020820111156100cc57600080fd5b803590602001918460018302840111640100000000831117156100ee57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061031d945050505050565b604080516001600160a01b039092168252519081900360200190f35b61017b6004803603604081101561016157600080fd5b50803560ff191690602001356001600160401b031661039d565b60408051918252519081900360200190f35b6101ae600480360360208110156101a357600080fd5b503560ff1916610682565b6040805164ffffffffff199092168252519081900360200190f35b6101ea600480360360208110156101df57600080fd5b503560ff1916610693565b604080516001600160401b0394851681529290931660208301528183015290519081900360600190f35b61017b6004803603604081101561022a57600080fd5b50803560ff191690602001356001600160401b0316610743565b6102fa6004803603606081101561025a57600080fd5b81359160ff196020820135169181019060608101604082013564010000000081111561028557600080fd5b82018360208201111561029757600080fd5b803590602001918460018302840111640100000000831117156102b957600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506109c1945050505050565b005b61017b6004803603602081101561031257600080fd5b503560ff1916610c5e565b602081810151604080840151606080860151835160008082528188018087528a905291821a81860181905292810186905260808101849052935190959293919260019260a080820193601f1981019281900390910190855afa158015610387573d6000803e3d6000fd5b5050506020604051035193505050505b92915050565b600042826001600160401b0316106103f1576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff1983166000908152600160205260409020546104125750600254610397565b60ff19831660009081526001602052604081208054600019810190811061043557fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116915083168110156104a05760ff19841660009081526001602052604090208054600019810190811061048657fe5b906000526020600020906002020160010154915050610397565b60ff198416600090815260016020526040812054600019015b8082116106755760ff19861660009081526001602052604090208054600284840104916001600160401b03881691839081106104f157fe5b6000918252602090912060029091020154600160401b90046001600160401b031614156105545760ff198716600090815260016020526040902080548290811061053757fe5b906000526020600020906002020160010154945050505050610397565b60ff198716600090815260016020526040902080548290811061057357fe5b60009182526020909120600290910201546001600160401b03600160401b90910481169087161180156105ee575060ff19871660009081526001602081905260409091208054909183019081106105c657fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908716105b156106125760ff198716600090815260016020526040902080548290811061053757fe5b60ff198716600090815260016020526040902080548290811061063157fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690871611156106685780600101925061066f565b6001810391505b506104b9565b5050600254949350505050565b62ffffff19601082901b165b919050565b60ff198116600090815260016020526040812054819081906106bf57505060025460009150819061073c565b6106c76111ee565b60ff1985166000908152600160205260409020805460001981019081106106ea57fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b909204169484018590526001909101549290910182905295509093509150505b9193909250565b600043826001600160401b031610610797576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff1983166000908152600160205260409020546107b85750600254610397565b60ff1983166000908152600160205260408120805460001981019081106107db57fe5b60009182526020909120600290910201546001600160401b03908116915083168110156108255760ff19841660009081526001602052604090208054600019810190811061048657fe5b60ff198416600090815260016020526040812054600019015b8082116106755760ff19861660009081526001602052604090208054600284840104916001600160401b038816918390811061087657fe5b60009182526020909120600290910201546001600160401b031614156108b55760ff198716600090815260016020526040902080548290811061053757fe5b60ff19871660009081526001602052604090208054829081106108d457fe5b60009182526020909120600290910201546001600160401b03908116908716118015610941575060ff198716600090815260016020819052604090912080549091830190811061092057fe5b60009182526020909120600290910201546001600160401b03908116908716105b156109655760ff198716600090815260016020526040902080548290811061053757fe5b60ff198716600090815260016020526040902080548290811061098457fe5b60009182526020909120600290910201546001600160401b0390811690871611156109b4578060010192506109bb565b6001810391505b5061083e565b6109c961120e565b6109d233610cbb565b90506109dc61120e565b6109e582610cf9565b90506000806109f383610d51565b90925090506000610a0387610682565b9050610a1381878585608c610d8a565b1515600114610a69576040805162461bcd60e51b815260206004820152601b60248201527f4d65726b6c6520747265652070726f6f66206e6f742076616c69640000000000604482015290519081900360640190fd5b60ff19871660009081526001602052604090205415610aff5760ff198716600090815260016020526040902080544391906000198101908110610aa857fe5b60009182526020909120600290910201546001600160401b03161415610aff5760405162461bcd60e51b81526004018080602001828103825260218152602001806112856021913960400191505060405180910390fd5b600160008860ff191660ff191681526020019081526020016000206040518060600160405280436001600160401b03168152602001426001600160401b031681526020018a8152509080600181540180825580915050906001820390600052602060002090600202016000909192909190915060008201518160000160006101000a8154816001600160401b0302191690836001600160401b0316021790555060208201518160000160086101000a8154816001600160401b0302191690836001600160401b03160217905550604082015181600101555050507fc26526f78af19d7325e758deea00a04b0db26e055a576222266dfe241755a57f8743428b604051808560ff191660ff19168152602001846001600160401b03166001600160401b03168152602001836001600160401b03166001600160401b0316815260200182815260200194505050505060405180910390a15050505050505050565b60ff198116600090815260016020526040812054610c7f575060025461068e565b60ff198216600090815260016020526040902080546000198101908110610ca257fe5b9060005260206000209060020201600101549050919050565b610cc361120e565b600960c01b815260006020820152606091821b6bffffffffffffffffffffffff19166040820152600360e01b9181019190915290565b610d0161120e565b815160c01c6060828101828152602085015160a01c6bffffffff000000000000000016909217909152604080840151821c8184018181529290940151901c63ffffffff60a01b1690921790915290565b600080610d6e836040015160001c846060015160001c6000610f8e565b83516020850151919350610d83916000610f8e565b9050915091565b6000610d94611235565b610d9d86611008565b9050600081600001518015610db3575081602001515b15610e1c576000805b81158015610dc957508581105b15610df1578088901c600116818560800151901c6001161860ff169150600181019050610dbc565b81610e03576000945050505050610f85565b610e1784608001518560a001516001610f8e565b925050505b6060826040015160ff16604051908082528060200260200182016040528015610e4f578160200160208202803883390190505b509050600080805b856040015160ff16811015610ee05760608601516001600160f01b0316811c600190811690811415610ebc5760008360200260400160ff169050808d0151945084868481518110610ea457fe5b60200260200101818152505060018401935050610ed7565b6000858381518110610eca57fe5b6020026020010181815250505b50600101610e57565b508451600090610efb57610ef68a8a6001610f8e565b610efd565b845b9050600080600188604001510360ff1690505b858181518110610f1c57fe5b6020908102919091010151915060018c821c81161480610f4757610f4284846000610f8e565b610f53565b610f5383856000610f8e565b935081610f605750610f6a565b5060001901610f10565b505060281b64ffffffffff19908116908c1614955050505050505b95945050505050565b60408051600280825260608083018452600093909291906020830190803883390190505090508481600081518110610fc257fe5b6020026020010181815250508381600181518110610fdc57fe5b60200260200101818152505082610ffd57610ff8816000611097565b610f85565b610f85816001611097565b611010611235565b61101861126a565b611021836110aa565b9050600061102e826110d3565b60018082168114855281811c8116146020850152905061104d826110d3565b60ff16604084015261105e826110e5565b6001600160f01b031660608401526020830151156110905783518401601f198101519051608085019190915260a08401525b5050919050565b60006110a383836110f7565b9392505050565b6110b261126a565b50604080518082019091526020828101825282518301810190820152919050565b80518051600190910190915260f81c90565b80518051601e90910190915260101c90565b6000817f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001825b85518110156111e457600054865183916001600160a01b03169063d15ca1099089908590811061114957fe5b6020026020010151866040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b15801561118d57600080fd5b505afa1580156111a1573d6000803e3d6000fd5b505050506040513d60208110156111b757600080fd5b505187518890849081106111c757fe5b6020026020010151850101816111d957fe5b06925060010161111d565b5090949350505050565b604080516060810182526000808252602082018190529181019190915290565b60408051608081018252600080825260208201819052918101829052606081019190915290565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a081019190915290565b60405180604001604052806000815260200160008152509056fe6e6f206d756c7469706c652073657420696e207468652073616d6520626c6f636ba265627a7a72305820d8243def2ead5842a12c6ae7a62d205ca5771d3a723dd96a0ffa60ae521f1e2764736f6c634300050a0032"

// DeployRootCommits deploys a new Ethereum contract, binding an instance of RootCommits to it.
func DeployRootCommits(auth *bind.TransactOpts, backend bind.ContractBackend, _mimcContractAddr common.Address) (common.Address, *types.Transaction, *RootCommits, error) {
	parsed, err := abi.JSON(strings.NewReader(RootCommitsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RootCommitsBin), backend, _mimcContractAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RootCommits{RootCommitsCaller: RootCommitsCaller{contract: contract}, RootCommitsTransactor: RootCommitsTransactor{contract: contract}, RootCommitsFilterer: RootCommitsFilterer{contract: contract}}, nil
}

// RootCommits is an auto generated Go binding around an Ethereum contract.
type RootCommits struct {
	RootCommitsCaller     // Read-only binding to the contract
	RootCommitsTransactor // Write-only binding to the contract
	RootCommitsFilterer   // Log filterer for contract events
}

// RootCommitsCaller is an auto generated read-only Go binding around an Ethereum contract.
type RootCommitsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootCommitsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RootCommitsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootCommitsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RootCommitsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootCommitsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RootCommitsSession struct {
	Contract     *RootCommits      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RootCommitsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RootCommitsCallerSession struct {
	Contract *RootCommitsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// RootCommitsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RootCommitsTransactorSession struct {
	Contract     *RootCommitsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// RootCommitsRaw is an auto generated low-level Go binding around an Ethereum contract.
type RootCommitsRaw struct {
	Contract *RootCommits // Generic contract binding to access the raw methods on
}

// RootCommitsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RootCommitsCallerRaw struct {
	Contract *RootCommitsCaller // Generic read-only contract binding to access the raw methods on
}

// RootCommitsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RootCommitsTransactorRaw struct {
	Contract *RootCommitsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRootCommits creates a new instance of RootCommits, bound to a specific deployed contract.
func NewRootCommits(address common.Address, backend bind.ContractBackend) (*RootCommits, error) {
	contract, err := bindRootCommits(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RootCommits{RootCommitsCaller: RootCommitsCaller{contract: contract}, RootCommitsTransactor: RootCommitsTransactor{contract: contract}, RootCommitsFilterer: RootCommitsFilterer{contract: contract}}, nil
}

// NewRootCommitsCaller creates a new read-only instance of RootCommits, bound to a specific deployed contract.
func NewRootCommitsCaller(address common.Address, caller bind.ContractCaller) (*RootCommitsCaller, error) {
	contract, err := bindRootCommits(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RootCommitsCaller{contract: contract}, nil
}

// NewRootCommitsTransactor creates a new write-only instance of RootCommits, bound to a specific deployed contract.
func NewRootCommitsTransactor(address common.Address, transactor bind.ContractTransactor) (*RootCommitsTransactor, error) {
	contract, err := bindRootCommits(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RootCommitsTransactor{contract: contract}, nil
}

// NewRootCommitsFilterer creates a new log filterer instance of RootCommits, bound to a specific deployed contract.
func NewRootCommitsFilterer(address common.Address, filterer bind.ContractFilterer) (*RootCommitsFilterer, error) {
	contract, err := bindRootCommits(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RootCommitsFilterer{contract: contract}, nil
}

// bindRootCommits binds a generic wrapper to an already deployed contract.
func bindRootCommits(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RootCommitsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootCommits *RootCommitsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootCommits.Contract.RootCommitsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootCommits *RootCommitsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootCommits.Contract.RootCommitsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootCommits *RootCommitsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootCommits.Contract.RootCommitsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootCommits *RootCommitsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootCommits.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootCommits *RootCommitsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootCommits.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootCommits *RootCommitsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootCommits.Contract.contract.Transact(opts, method, params...)
}

// CheckSig is a free data retrieval call binding the contract method 0x01b0452c.
//
// Solidity: function checkSig(bytes32 msgHash, bytes rsv) constant returns(address)
func (_RootCommits *RootCommitsCaller) CheckSig(opts *bind.CallOpts, msgHash [32]byte, rsv []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RootCommits.contract.Call(opts, out, "checkSig", msgHash, rsv)
	return *ret0, err
}

// CheckSig is a free data retrieval call binding the contract method 0x01b0452c.
//
// Solidity: function checkSig(bytes32 msgHash, bytes rsv) constant returns(address)
func (_RootCommits *RootCommitsSession) CheckSig(msgHash [32]byte, rsv []byte) (common.Address, error) {
	return _RootCommits.Contract.CheckSig(&_RootCommits.CallOpts, msgHash, rsv)
}

// CheckSig is a free data retrieval call binding the contract method 0x01b0452c.
//
// Solidity: function checkSig(bytes32 msgHash, bytes rsv) constant returns(address)
func (_RootCommits *RootCommitsCallerSession) CheckSig(msgHash [32]byte, rsv []byte) (common.Address, error) {
	return _RootCommits.Contract.CheckSig(&_RootCommits.CallOpts, msgHash, rsv)
}

// GetRoot is a free data retrieval call binding the contract method 0xfead90d7.
//
// Solidity: function getRoot(bytes31 id) constant returns(bytes32)
func (_RootCommits *RootCommitsCaller) GetRoot(opts *bind.CallOpts, id [31]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RootCommits.contract.Call(opts, out, "getRoot", id)
	return *ret0, err
}

// GetRoot is a free data retrieval call binding the contract method 0xfead90d7.
//
// Solidity: function getRoot(bytes31 id) constant returns(bytes32)
func (_RootCommits *RootCommitsSession) GetRoot(id [31]byte) ([32]byte, error) {
	return _RootCommits.Contract.GetRoot(&_RootCommits.CallOpts, id)
}

// GetRoot is a free data retrieval call binding the contract method 0xfead90d7.
//
// Solidity: function getRoot(bytes31 id) constant returns(bytes32)
func (_RootCommits *RootCommitsCallerSession) GetRoot(id [31]byte) ([32]byte, error) {
	return _RootCommits.Contract.GetRoot(&_RootCommits.CallOpts, id)
}

// GetRootByBlock is a free data retrieval call binding the contract method 0xb816ff6f.
//
// Solidity: function getRootByBlock(bytes31 id, uint64 blockN) constant returns(bytes32)
func (_RootCommits *RootCommitsCaller) GetRootByBlock(opts *bind.CallOpts, id [31]byte, blockN uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RootCommits.contract.Call(opts, out, "getRootByBlock", id, blockN)
	return *ret0, err
}

// GetRootByBlock is a free data retrieval call binding the contract method 0xb816ff6f.
//
// Solidity: function getRootByBlock(bytes31 id, uint64 blockN) constant returns(bytes32)
func (_RootCommits *RootCommitsSession) GetRootByBlock(id [31]byte, blockN uint64) ([32]byte, error) {
	return _RootCommits.Contract.GetRootByBlock(&_RootCommits.CallOpts, id, blockN)
}

// GetRootByBlock is a free data retrieval call binding the contract method 0xb816ff6f.
//
// Solidity: function getRootByBlock(bytes31 id, uint64 blockN) constant returns(bytes32)
func (_RootCommits *RootCommitsCallerSession) GetRootByBlock(id [31]byte, blockN uint64) ([32]byte, error) {
	return _RootCommits.Contract.GetRootByBlock(&_RootCommits.CallOpts, id, blockN)
}

// GetRootByTime is a free data retrieval call binding the contract method 0x4175dae5.
//
// Solidity: function getRootByTime(bytes31 id, uint64 timestamp) constant returns(bytes32)
func (_RootCommits *RootCommitsCaller) GetRootByTime(opts *bind.CallOpts, id [31]byte, timestamp uint64) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RootCommits.contract.Call(opts, out, "getRootByTime", id, timestamp)
	return *ret0, err
}

// GetRootByTime is a free data retrieval call binding the contract method 0x4175dae5.
//
// Solidity: function getRootByTime(bytes31 id, uint64 timestamp) constant returns(bytes32)
func (_RootCommits *RootCommitsSession) GetRootByTime(id [31]byte, timestamp uint64) ([32]byte, error) {
	return _RootCommits.Contract.GetRootByTime(&_RootCommits.CallOpts, id, timestamp)
}

// GetRootByTime is a free data retrieval call binding the contract method 0x4175dae5.
//
// Solidity: function getRootByTime(bytes31 id, uint64 timestamp) constant returns(bytes32)
func (_RootCommits *RootCommitsCallerSession) GetRootByTime(id [31]byte, timestamp uint64) ([32]byte, error) {
	return _RootCommits.Contract.GetRootByTime(&_RootCommits.CallOpts, id, timestamp)
}

// GetRootDataById is a free data retrieval call binding the contract method 0xb5818699.
//
// Solidity: function getRootDataById(bytes31 id) constant returns(uint64, uint64, bytes32)
func (_RootCommits *RootCommitsCaller) GetRootDataById(opts *bind.CallOpts, id [31]byte) (uint64, uint64, [32]byte, error) {
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
	err := _RootCommits.contract.Call(opts, out, "getRootDataById", id)
	return *ret0, *ret1, *ret2, err
}

// GetRootDataById is a free data retrieval call binding the contract method 0xb5818699.
//
// Solidity: function getRootDataById(bytes31 id) constant returns(uint64, uint64, bytes32)
func (_RootCommits *RootCommitsSession) GetRootDataById(id [31]byte) (uint64, uint64, [32]byte, error) {
	return _RootCommits.Contract.GetRootDataById(&_RootCommits.CallOpts, id)
}

// GetRootDataById is a free data retrieval call binding the contract method 0xb5818699.
//
// Solidity: function getRootDataById(bytes31 id) constant returns(uint64, uint64, bytes32)
func (_RootCommits *RootCommitsCallerSession) GetRootDataById(id [31]byte) (uint64, uint64, [32]byte, error) {
	return _RootCommits.Contract.GetRootDataById(&_RootCommits.CallOpts, id)
}

// GetRootFromId is a free data retrieval call binding the contract method 0xad05a8d2.
//
// Solidity: function getRootFromId(bytes31 id) constant returns(bytes27)
func (_RootCommits *RootCommitsCaller) GetRootFromId(opts *bind.CallOpts, id [31]byte) ([27]byte, error) {
	var (
		ret0 = new([27]byte)
	)
	out := ret0
	err := _RootCommits.contract.Call(opts, out, "getRootFromId", id)
	return *ret0, err
}

// GetRootFromId is a free data retrieval call binding the contract method 0xad05a8d2.
//
// Solidity: function getRootFromId(bytes31 id) constant returns(bytes27)
func (_RootCommits *RootCommitsSession) GetRootFromId(id [31]byte) ([27]byte, error) {
	return _RootCommits.Contract.GetRootFromId(&_RootCommits.CallOpts, id)
}

// GetRootFromId is a free data retrieval call binding the contract method 0xad05a8d2.
//
// Solidity: function getRootFromId(bytes31 id) constant returns(bytes27)
func (_RootCommits *RootCommitsCallerSession) GetRootFromId(id [31]byte) ([27]byte, error) {
	return _RootCommits.Contract.GetRootFromId(&_RootCommits.CallOpts, id)
}

// SetRoot is a paid mutator transaction binding the contract method 0xe0681acd.
//
// Solidity: function setRoot(bytes32 newRoot, bytes31 id, bytes mtp) returns()
func (_RootCommits *RootCommitsTransactor) SetRoot(opts *bind.TransactOpts, newRoot [32]byte, id [31]byte, mtp []byte) (*types.Transaction, error) {
	return _RootCommits.contract.Transact(opts, "setRoot", newRoot, id, mtp)
}

// SetRoot is a paid mutator transaction binding the contract method 0xe0681acd.
//
// Solidity: function setRoot(bytes32 newRoot, bytes31 id, bytes mtp) returns()
func (_RootCommits *RootCommitsSession) SetRoot(newRoot [32]byte, id [31]byte, mtp []byte) (*types.Transaction, error) {
	return _RootCommits.Contract.SetRoot(&_RootCommits.TransactOpts, newRoot, id, mtp)
}

// SetRoot is a paid mutator transaction binding the contract method 0xe0681acd.
//
// Solidity: function setRoot(bytes32 newRoot, bytes31 id, bytes mtp) returns()
func (_RootCommits *RootCommitsTransactorSession) SetRoot(newRoot [32]byte, id [31]byte, mtp []byte) (*types.Transaction, error) {
	return _RootCommits.Contract.SetRoot(&_RootCommits.TransactOpts, newRoot, id, mtp)
}

// RootCommitsRootUpdatedIterator is returned from FilterRootUpdated and is used to iterate over the raw logs and unpacked data for RootUpdated events raised by the RootCommits contract.
type RootCommitsRootUpdatedIterator struct {
	Event *RootCommitsRootUpdated // Event containing the contract specifics and raw log

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
func (it *RootCommitsRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootCommitsRootUpdated)
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
		it.Event = new(RootCommitsRootUpdated)
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
func (it *RootCommitsRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootCommitsRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootCommitsRootUpdated represents a RootUpdated event raised by the RootCommits contract.
type RootCommitsRootUpdated struct {
	Id        [31]byte
	BlockN    uint64
	Timestamp uint64
	Root      [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRootUpdated is a free log retrieval operation binding the contract event 0xc26526f78af19d7325e758deea00a04b0db26e055a576222266dfe241755a57f.
//
// Solidity: event RootUpdated(bytes31 id, uint64 blockN, uint64 timestamp, bytes32 root)
func (_RootCommits *RootCommitsFilterer) FilterRootUpdated(opts *bind.FilterOpts) (*RootCommitsRootUpdatedIterator, error) {

	logs, sub, err := _RootCommits.contract.FilterLogs(opts, "RootUpdated")
	if err != nil {
		return nil, err
	}
	return &RootCommitsRootUpdatedIterator{contract: _RootCommits.contract, event: "RootUpdated", logs: logs, sub: sub}, nil
}

// WatchRootUpdated is a free log subscription operation binding the contract event 0xc26526f78af19d7325e758deea00a04b0db26e055a576222266dfe241755a57f.
//
// Solidity: event RootUpdated(bytes31 id, uint64 blockN, uint64 timestamp, bytes32 root)
func (_RootCommits *RootCommitsFilterer) WatchRootUpdated(opts *bind.WatchOpts, sink chan<- *RootCommitsRootUpdated) (event.Subscription, error) {

	logs, sub, err := _RootCommits.contract.WatchLogs(opts, "RootUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootCommitsRootUpdated)
				if err := _RootCommits.contract.UnpackLog(event, "RootUpdated", log); err != nil {
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

// ParseRootUpdated is a log parse operation binding the contract event 0xc26526f78af19d7325e758deea00a04b0db26e055a576222266dfe241755a57f.
//
// Solidity: event RootUpdated(bytes31 id, uint64 blockN, uint64 timestamp, bytes32 root)
func (_RootCommits *RootCommitsFilterer) ParseRootUpdated(log types.Log) (*RootCommitsRootUpdated, error) {
	event := new(RootCommitsRootUpdated)
	if err := _RootCommits.contract.UnpackLog(event, "RootUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
