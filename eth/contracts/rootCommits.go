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

// Iden3HelpersBin is the compiled bytecode used for deploying new contracts.
const Iden3HelpersBin = `0x608060405234801561001057600080fd5b506040516020806102578339810180604052602081101561003057600080fd5b5051600080546001600160a01b039092166001600160a01b03199092169190911790556101f5806100626000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806301b0452c1461003b578063ad05a8d214610104575b600080fd5b6100e86004803603604081101561005157600080fd5b8135919081019060408101602082013564010000000081111561007357600080fd5b82018360208201111561008557600080fd5b803590602001918460018302840111640100000000831117156100a757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610140945050505050565b604080516001600160a01b039092168252519081900360200190f35b6101256004803603602081101561011a57600080fd5b503560ff19166101bf565b6040805164ffffffffff199092168252519081900360200190f35b602081810151604080840151606080860151835160008082528188018087528a905291821a81860181905292810186905260808101849052935190959293919260019260a080820193601f1981019281900390910190855afa1580156101aa573d6000803e3d6000fd5b5050604051601f190151979650505050505050565b60ff191660101b9056fea165627a7a723058203f04d449ff149208f068a236dbf08a86beaf0911a4a2ec3a76851acfcda694230029`

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
const MemoryBin = `0x604c6023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea165627a7a72305820626fcc786ff29a40122c10d81ba654920bb8af0e29d42517e3aa7fca56e6f3700029`

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

// MimcUnitBin is the compiled bytecode used for deploying new contracts.
const MimcUnitBin = `0x6080604052348015600f57600080fd5b5060938061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063d15ca10914602d575b600080fd5b604d60048036036040811015604157600080fd5b5080359060200135605f565b60408051918252519081900360200190f35b60009291505056fea165627a7a72305820570b24afe47c62fe91b60d5395f9f41d29754bc80e63dc18c6742d587428d54b0029`

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
const RootCommitsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\"},{\"name\":\"rsv\",\"type\":\"bytes\"}],\"name\":\"checkSig\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"},{\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getRootByTime\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getRootFromId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes27\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"},{\"name\":\"blockN\",\"type\":\"uint64\"}],\"name\":\"getRootByBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newRoot\",\"type\":\"bytes32\"},{\"name\":\"id\",\"type\":\"bytes31\"},{\"name\":\"mtp\",\"type\":\"bytes\"}],\"name\":\"setRoot\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getRoot\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_mimcContractAddr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes31\"},{\"indexed\":false,\"name\":\"blockN\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootUpdated\",\"type\":\"event\"}]"

// RootCommitsBin is the compiled bytecode used for deploying new contracts.
const RootCommitsBin = `0x608060405234801561001057600080fd5b506040516020806112258339810180604052602081101561003057600080fd5b5051600080546001600160a01b039092166001600160a01b03199092169190911790556111c3806100626000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806301b0452c146100675780634175dae514610130578063ad05a8d214610172578063b816ff6f146101ae578063e0681acd146101de578063fead90d714610296575b600080fd5b6101146004803603604081101561007d57600080fd5b8135919081019060408101602082013564010000000081111561009f57600080fd5b8201836020820111156100b157600080fd5b803590602001918460018302840111640100000000831117156100d357600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506102b7945050505050565b604080516001600160a01b039092168252519081900360200190f35b6101606004803603604081101561014657600080fd5b50803560ff191690602001356001600160401b0316610337565b60408051918252519081900360200190f35b6101936004803603602081101561018857600080fd5b503560ff1916610622565b6040805164ffffffffff199092168252519081900360200190f35b610160600480360360408110156101c457600080fd5b50803560ff191690602001356001600160401b0316610630565b610294600480360360608110156101f457600080fd5b81359160ff196020820135169181019060608101604082013564010000000081111561021f57600080fd5b82018360208201111561023157600080fd5b8035906020019184600183028401116401000000008311171561025357600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506108b4945050505050565b005b610160600480360360208110156102ac57600080fd5b503560ff1916610b57565b602081810151604080840151606080860151835160008082528188018087528a905291821a81860181905292810186905260808101849052935190959293919260019260a080820193601f1981019281900390910190855afa158015610321573d6000803e3d6000fd5b5050506020604051035193505050505b92915050565b600042826001600160401b0316106103915760408051600160e51b62461bcd0281526020600482015260126024820152600160721b71195c9c939bd19d5d1d5c99505b1b1bddd95902604482015290519081900360640190fd5b60ff1983166000908152600160205260409020546103b25750600254610331565b60ff1983166000908152600160205260408120805460001981019081106103d557fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116915083168110156104405760ff19841660009081526001602052604090208054600019810190811061042657fe5b906000526020600020906002020160010154915050610331565b60ff198416600090815260016020526040812054600019015b8082116106155760ff19861660009081526001602052604090208054600284840104916001600160401b038816918390811061049157fe5b6000918252602090912060029091020154600160401b90046001600160401b031614156104f45760ff19871660009081526001602052604090208054829081106104d757fe5b906000526020600020906002020160010154945050505050610331565b60ff198716600090815260016020526040902080548290811061051357fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690871611801561058e575060ff198716600090815260016020819052604090912080549091830190811061056657fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908716105b156105b25760ff19871660009081526001602052604090208054829081106104d757fe5b60ff19871660009081526001602052604090208054829081106105d157fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690871611156106085780600101925061060f565b6001810391505b50610459565b5050600254949350505050565b60ff19811660101b5b919050565b600043826001600160401b03161061068a5760408051600160e51b62461bcd0281526020600482015260126024820152600160721b71195c9c939bd19d5d1d5c99505b1b1bddd95902604482015290519081900360640190fd5b60ff1983166000908152600160205260409020546106ab5750600254610331565b60ff1983166000908152600160205260408120805460001981019081106106ce57fe5b60009182526020909120600290910201546001600160401b03908116915083168110156107185760ff19841660009081526001602052604090208054600019810190811061042657fe5b60ff198416600090815260016020526040812054600019015b8082116106155760ff19861660009081526001602052604090208054600284840104916001600160401b038816918390811061076957fe5b60009182526020909120600290910201546001600160401b031614156107a85760ff19871660009081526001602052604090208054829081106104d757fe5b60ff19871660009081526001602052604090208054829081106107c757fe5b60009182526020909120600290910201546001600160401b03908116908716118015610834575060ff198716600090815260016020819052604090912080549091830190811061081357fe5b60009182526020909120600290910201546001600160401b03908116908716105b156108585760ff19871660009081526001602052604090208054829081106104d757fe5b60ff198716600090815260016020526040902080548290811061087757fe5b60009182526020909120600290910201546001600160401b0390811690871611156108a7578060010192506108ae565b6001810391505b50610731565b6108bc611100565b6108c533610bb4565b90506108cf611100565b6108d882610bf8565b90506000806108e683610c60565b909250905060006108f687610622565b905061090681878585608c610c99565b151560011461095f5760408051600160e51b62461bcd02815260206004820152601b60248201527f4d65726b6c6520747265652070726f6f66206e6f742076616c69640000000000604482015290519081900360640190fd5b60ff198716600090815260016020526040902054156109f85760ff19871660009081526001602052604090208054439190600019810190811061099e57fe5b60009182526020909120600290910201546001600160401b031614156109f857604051600160e51b62461bcd0281526004018080602001828103825260218152602001806111776021913960400191505060405180910390fd5b600160008860ff191660ff191681526020019081526020016000206040518060600160405280436001600160401b03168152602001426001600160401b031681526020018a8152509080600181540180825580915050906001820390600052602060002090600202016000909192909190915060008201518160000160006101000a8154816001600160401b0302191690836001600160401b0316021790555060208201518160000160086101000a8154816001600160401b0302191690836001600160401b03160217905550604082015181600101555050507fc26526f78af19d7325e758deea00a04b0db26e055a576222266dfe241755a57f8743428b604051808560ff191660ff19168152602001846001600160401b03166001600160401b03168152602001836001600160401b03166001600160401b0316815260200182815260200194505050505060405180910390a15050505050505050565b60ff198116600090815260016020526040812054610b78575060025461062b565b60ff198216600090815260016020526040902080546000198101908110610b9b57fe5b9060005260206000209060020201600101549050919050565b610bbc611100565b600160c01b600902815260006020820152606091821b6bffffffffffffffffffffffff19166040820152600160e01b6003029181019190915290565b610c00611100565b81516001600160c01b03191660c01c606080830182815260208501516001600160e01b031990811660a01c90931790526040808501516bffffffffffffffffffffffff1916821c81850181815292860151909316901c9091179052919050565b600080610c7d836040015160001c846060015160001c6000610e9d565b83516020850151919350610c92916000610e9d565b9050915091565b6000610ca3611127565b610cac86610f17565b9050600081600001518015610cc2575081602001515b15610d2b576000805b81158015610cd857508581105b15610d00578088901c600116818560800151901c6001161860ff169150600181019050610ccb565b81610d12576000945050505050610e94565b610d2684608001518560a001516001610e9d565b925050505b6060826040015160ff16604051908082528060200260200182016040528015610d5e578160200160208202803883390190505b509050600080805b856040015160ff16811015610def5760608601516001600160f01b0316811c600190811690811415610dcb5760008360200260400160ff169050808d0151945084868481518110610db357fe5b60200260200101818152505060018401935050610de6565b6000858381518110610dd957fe5b6020026020010181815250505b50600101610d66565b508451600090610e0a57610e058a8a6001610e9d565b610e0c565b845b9050600080600188604001510360ff1690505b858181518110610e2b57fe5b6020908102919091010151915060018c821c81161480610e5657610e5184846000610e9d565b610e62565b610e6283856000610e9d565b935081610e6f5750610e79565b5060001901610e1f565b505060281b64ffffffffff19908116908c1614955050505050505b95945050505050565b60408051600280825260608083018452600093909291906020830190803883390190505090508481600081518110610ed157fe5b6020026020010181815250508381600181518110610eeb57fe5b60200260200101818152505082610f0c57610f07816000610fa9565b610e94565b610e94816001610fa9565b610f1f611127565b610f2761115c565b610f3083610fbc565b90506000610f3d82610fe5565b60018082168114855260ff8216811c81161460208501529050610f5f82610fe5565b60ff166040840152610f7082610ff7565b6001600160f01b03166060840152602083015115610fa25783518401601f198101519051608085019190915260a08401525b5050919050565b6000610fb58383611009565b9392505050565b610fc461115c565b50604080518082019091526020828101825282518301810190820152919050565b80518051600190910190915260f81c90565b80518051601e90910190915260101c90565b6000817f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001825b85518110156110f657600054865183916001600160a01b03169063d15ca1099089908590811061105b57fe5b6020026020010151866040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b15801561109f57600080fd5b505afa1580156110b3573d6000803e3d6000fd5b505050506040513d60208110156110c957600080fd5b505187518890849081106110d957fe5b6020026020010151850101816110eb57fe5b06925060010161102f565b5090949350505050565b60408051608081018252600080825260208201819052918101829052606081019190915290565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a081019190915290565b60405180604001604052806000815260200160008152509056fe6e6f206d756c7469706c652073657420696e207468652073616d6520626c6f636ba165627a7a7230582001b029c4014fe5fd337e556dc81472ac563b45ca231bb10aedd19d998c9aafd60029`

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
