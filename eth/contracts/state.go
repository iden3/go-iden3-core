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

// PairingABI is the input ABI used to generate the binding from.
const PairingABI = "[]"

// PairingBin is the compiled bytecode used for deploying new contracts.
var PairingBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212202cbffd420089035d717bf24840a24ee6f3a7cabaadb24c64e5de3d71f84217f364736f6c63430006060033"

// DeployPairing deploys a new Ethereum contract, binding an instance of Pairing to it.
func DeployPairing(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Pairing, error) {
	parsed, err := abi.JSON(strings.NewReader(PairingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PairingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Pairing{PairingCaller: PairingCaller{contract: contract}, PairingTransactor: PairingTransactor{contract: contract}, PairingFilterer: PairingFilterer{contract: contract}}, nil
}

// Pairing is an auto generated Go binding around an Ethereum contract.
type Pairing struct {
	PairingCaller     // Read-only binding to the contract
	PairingTransactor // Write-only binding to the contract
	PairingFilterer   // Log filterer for contract events
}

// PairingCaller is an auto generated read-only Go binding around an Ethereum contract.
type PairingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PairingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PairingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PairingSession struct {
	Contract     *Pairing          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PairingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PairingCallerSession struct {
	Contract *PairingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// PairingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PairingTransactorSession struct {
	Contract     *PairingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PairingRaw is an auto generated low-level Go binding around an Ethereum contract.
type PairingRaw struct {
	Contract *Pairing // Generic contract binding to access the raw methods on
}

// PairingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PairingCallerRaw struct {
	Contract *PairingCaller // Generic read-only contract binding to access the raw methods on
}

// PairingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PairingTransactorRaw struct {
	Contract *PairingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPairing creates a new instance of Pairing, bound to a specific deployed contract.
func NewPairing(address common.Address, backend bind.ContractBackend) (*Pairing, error) {
	contract, err := bindPairing(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pairing{PairingCaller: PairingCaller{contract: contract}, PairingTransactor: PairingTransactor{contract: contract}, PairingFilterer: PairingFilterer{contract: contract}}, nil
}

// NewPairingCaller creates a new read-only instance of Pairing, bound to a specific deployed contract.
func NewPairingCaller(address common.Address, caller bind.ContractCaller) (*PairingCaller, error) {
	contract, err := bindPairing(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PairingCaller{contract: contract}, nil
}

// NewPairingTransactor creates a new write-only instance of Pairing, bound to a specific deployed contract.
func NewPairingTransactor(address common.Address, transactor bind.ContractTransactor) (*PairingTransactor, error) {
	contract, err := bindPairing(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PairingTransactor{contract: contract}, nil
}

// NewPairingFilterer creates a new log filterer instance of Pairing, bound to a specific deployed contract.
func NewPairingFilterer(address common.Address, filterer bind.ContractFilterer) (*PairingFilterer, error) {
	contract, err := bindPairing(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PairingFilterer{contract: contract}, nil
}

// bindPairing binds a generic wrapper to an already deployed contract.
func bindPairing(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PairingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pairing *PairingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Pairing.Contract.PairingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pairing *PairingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pairing.Contract.PairingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pairing *PairingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pairing.Contract.PairingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pairing *PairingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Pairing.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pairing *PairingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pairing.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pairing *PairingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pairing.Contract.contract.Transact(opts, method, params...)
}

// StateABI is the input ABI used to generate the binding from.
const StateABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_verifierContractAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"StateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"}],\"name\":\"getStateDataByBlock\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getStateDataById\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getStateDataByTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"genesisState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"name\":\"initState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"name\":\"setState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// StateFuncSigs maps the 4-byte function signature to its string representation.
var StateFuncSigs = map[string]string{
	"44c9af28": "getState(uint256)",
	"d8dcd971": "getStateDataByBlock(uint256,uint64)",
	"c8d1e53e": "getStateDataById(uint256)",
	"0281bec2": "getStateDataByTime(uint256,uint64)",
	"307ec167": "initState(uint256,uint256,uint256,uint256[2],uint256[2][2],uint256[2])",
	"befb5b77": "setState(uint256,uint256,uint256[2],uint256[2][2],uint256[2])",
}

// StateBin is the compiled bytecode used for deploying new contracts.
var StateBin = "0x608060405234801561001057600080fd5b50604051610f34380380610f348339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b0319909216919091179055610ecf806100656000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80630281bec214610067578063307ec167146100bd57806344c9af28146101a0578063befb5b77146101cf578063c8d1e53e146102ab578063d8dcd971146102c8575b600080fd5b6100936004803603604081101561007d57600080fd5b50803590602001356001600160401b03166102f4565b604080516001600160401b0394851681529290931660208301528183015290519081900360600190f35b61019e60048036036101608110156100d457600080fd5b60408051808201825283359360208101359383820135939082019260a08301916060840190600290839083908082843760009201829052506040805180820190915293969594608081019493509150600290835b828210156101665760408051808201825290808402860190600290839083908082843760009201919091525050508152600190910190602001610128565b505060408051808201825293969594818101949350915060029083908390808284376000920191909152509194506106be9350505050565b005b6101bd600480360360208110156101b657600080fd5b50356106ed565b60408051918252519081900360200190f35b61019e60048036036101408110156101e657600080fd5b6040805180820182528335936020810135938101929091608083019180840190600290839083908082843760009201829052506040805180820190915293969594608081019493509150600290835b828210156102735760408051808201825290808402860190600290839083908082843760009201919091525050508152600190910190602001610235565b505060408051808201825293969594818101949350915060029083908390808284376000920191909152509194506107419350505050565b610093600480360360208110156102c157600080fd5b5035610816565b610093600480360360408110156102de57600080fd5b50803590602001356001600160401b03166108bc565b600080600042846001600160401b03161061034b576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60008581526001602052604090205461036e5750506002546000915081906106b7565b60008581526001602052604081208054600019810190811061038c57fe5b60009182526020909120600290910201546001600160401b03600160401b909104811691508516811015610475576000868152600160205260409020805460001981019081106103d857fe5b600091825260208083206002909202909101548883526001909152604090912080546001600160401b0390921691600019810190811061041457fe5b6000918252602080832060029290920290910154898352600190915260409091208054600160401b9092046001600160401b031691600019810190811061045757fe5b906000526020600020906002020160010154935093509350506106b7565b600086815260016020526040812054600019015b8082116106a75760008881526001602052604090208054600284840104916001600160401b038a1691839081106104bc57fe5b6000918252602090912060029091020154600160401b90046001600160401b031614156105995760008981526001602052604090208054829081106104fd57fe5b600091825260208083206002909202909101548b83526001909152604090912080546001600160401b03909216918390811061053557fe5b906000526020600020906002020160000160089054906101000a90046001600160401b0316600160008c8152602001908152602001600020838154811061057857fe5b906000526020600020906002020160010154965096509650505050506106b7565b60008981526001602052604090208054829081106105b357fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690891611801561062a5750600160008a8152602001908152602001600020816001018154811061060257fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908916105b156106495760008981526001602052604090208054829081106104fd57fe5b600089815260016020526040902080548290811061066357fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908916111561069a578060010192506106a1565b6001810391505b50610489565b5050600254600094508493509150505b9250925092565b600084815260016020526040902054156106d757600080fd5b6106e5868686868686610b13565b505050505050565b600081815260016020526040812054610709575060025461073c565b60008281526001602052604090208054600019810190811061072757fe5b90600052602060002090600202016001015490505b919050565b60008481526001602052604090205461075957600080fd5b610761610e0e565b60008581526001602052604090208054600019810190811061077f57fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b9092041694840194909452600101549082015291504314156108045760405162461bcd60e51b8152600401808060200182810382526021815260200180610e796021913960400191505060405180910390fd5b6106e586826040015187878787610b13565b6000818152600160205260408120548190819061083d5750506002546000915081906108b5565b610845610e0e565b60008581526001602052604090208054600019810190811061086357fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b909204169484018590526001909101549290910182905295509093509150505b9193909250565b600080600043846001600160401b031610610913576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b6000858152600160205260409020546109365750506002546000915081906106b7565b60008581526001602052604081208054600019810190811061095457fe5b60009182526020909120600290910201546001600160401b0390811691508516811015610999576000868152600160205260409020805460001981019081106103d857fe5b600086815260016020526040812054600019015b8082116106a75760008881526001602052604090208054600284840104916001600160401b038a1691839081106109e057fe5b60009182526020909120600290910201546001600160401b03161415610a1a5760008981526001602052604090208054829081106104fd57fe5b6000898152600160205260409020805482908110610a3457fe5b60009182526020909120600290910201546001600160401b03908116908916118015610a9d5750600160008a81526020019081526020016000208160010181548110610a7c57fe5b60009182526020909120600290910201546001600160401b03908116908916105b15610abc5760008981526001602052604090208054829081106104fd57fe5b6000898152600160205260409020805482908110610ad657fe5b60009182526020909120600290910201546001600160401b039081169089161115610b0657806001019250610b0d565b6001810391505b506109ad565b610b1b610e2e565b60405180606001604052808681526020018781526020018881525090506000809054906101000a90046001600160a01b03166001600160a01b03166311479fea858585856040518563ffffffff1660e01b81526004018085600260200280838360005b83811015610b96578181015183820152602001610b7e565b505050509050018460026000925b81841015610be45760208402830151604080838360005b83811015610bd3578181015183820152602001610bbb565b505050509050019260010192610ba4565b9250505083600260200280838360005b83811015610c0c578181015183820152602001610bf4565b5050505090500182600360200280838360005b83811015610c37578181015183820152602001610c1f565b5050505090500194505050505060206040518083038186803b158015610c5c57600080fd5b505afa158015610c70573d6000803e3d6000fd5b505050506040513d6020811015610c8657600080fd5b5051610cc35760405162461bcd60e51b815260040180806020018281038252602c815260200180610e4d602c913960400191505060405180910390fd5b600160008681526020019081526020016000206040518060600160405280436001600160401b03168152602001426001600160401b0316815260200189815250908060018154018082558091505060019003906000526020600020906002020160009091909190915060008201518160000160006101000a8154816001600160401b0302191690836001600160401b0316021790555060208201518160000160086101000a8154816001600160401b0302191690836001600160401b031602179055506040820151816001015550507f81c6f328b24014ef550c34a433275b52f3a8a0f32aa871adec069ab526a023908543428a60405180858152602001846001600160401b03166001600160401b03168152602001836001600160401b03166001600160401b0316815260200182815260200194505050505060405180910390a150505050505050565b604080516060810182526000808252602082018190529181019190915290565b6040518060600160405280600390602082028036833750919291505056fe7a6b50726f6f6620696453746174652075706461746520636f756c64206e6f742062652076657269666965646e6f206d756c7469706c652073657420696e207468652073616d6520626c6f636ba26469706673582212207b680bc0868359424cba050f5a4b295828a98db567d2a99cf119b0ffb533625964736f6c63430006060033"

// DeployState deploys a new Ethereum contract, binding an instance of State to it.
func DeployState(auth *bind.TransactOpts, backend bind.ContractBackend, _verifierContractAddr common.Address) (common.Address, *types.Transaction, *State, error) {
	parsed, err := abi.JSON(strings.NewReader(StateABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StateBin), backend, _verifierContractAddr)
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

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(uint256 id) view returns(uint256)
func (_State *StateCaller) GetState(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _State.contract.Call(opts, out, "getState", id)
	return *ret0, err
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(uint256 id) view returns(uint256)
func (_State *StateSession) GetState(id *big.Int) (*big.Int, error) {
	return _State.Contract.GetState(&_State.CallOpts, id)
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(uint256 id) view returns(uint256)
func (_State *StateCallerSession) GetState(id *big.Int) (*big.Int, error) {
	return _State.Contract.GetState(&_State.CallOpts, id)
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xd8dcd971.
//
// Solidity: function getStateDataByBlock(uint256 id, uint64 blockN) view returns(uint64, uint64, uint256)
func (_State *StateCaller) GetStateDataByBlock(opts *bind.CallOpts, id *big.Int, blockN uint64) (uint64, uint64, *big.Int, error) {
	var (
		ret0 = new(uint64)
		ret1 = new(uint64)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _State.contract.Call(opts, out, "getStateDataByBlock", id, blockN)
	return *ret0, *ret1, *ret2, err
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xd8dcd971.
//
// Solidity: function getStateDataByBlock(uint256 id, uint64 blockN) view returns(uint64, uint64, uint256)
func (_State *StateSession) GetStateDataByBlock(id *big.Int, blockN uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByBlock(&_State.CallOpts, id, blockN)
}

// GetStateDataByBlock is a free data retrieval call binding the contract method 0xd8dcd971.
//
// Solidity: function getStateDataByBlock(uint256 id, uint64 blockN) view returns(uint64, uint64, uint256)
func (_State *StateCallerSession) GetStateDataByBlock(id *big.Int, blockN uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByBlock(&_State.CallOpts, id, blockN)
}

// GetStateDataById is a free data retrieval call binding the contract method 0xc8d1e53e.
//
// Solidity: function getStateDataById(uint256 id) view returns(uint64, uint64, uint256)
func (_State *StateCaller) GetStateDataById(opts *bind.CallOpts, id *big.Int) (uint64, uint64, *big.Int, error) {
	var (
		ret0 = new(uint64)
		ret1 = new(uint64)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _State.contract.Call(opts, out, "getStateDataById", id)
	return *ret0, *ret1, *ret2, err
}

// GetStateDataById is a free data retrieval call binding the contract method 0xc8d1e53e.
//
// Solidity: function getStateDataById(uint256 id) view returns(uint64, uint64, uint256)
func (_State *StateSession) GetStateDataById(id *big.Int) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataById(&_State.CallOpts, id)
}

// GetStateDataById is a free data retrieval call binding the contract method 0xc8d1e53e.
//
// Solidity: function getStateDataById(uint256 id) view returns(uint64, uint64, uint256)
func (_State *StateCallerSession) GetStateDataById(id *big.Int) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataById(&_State.CallOpts, id)
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x0281bec2.
//
// Solidity: function getStateDataByTime(uint256 id, uint64 timestamp) view returns(uint64, uint64, uint256)
func (_State *StateCaller) GetStateDataByTime(opts *bind.CallOpts, id *big.Int, timestamp uint64) (uint64, uint64, *big.Int, error) {
	var (
		ret0 = new(uint64)
		ret1 = new(uint64)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _State.contract.Call(opts, out, "getStateDataByTime", id, timestamp)
	return *ret0, *ret1, *ret2, err
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x0281bec2.
//
// Solidity: function getStateDataByTime(uint256 id, uint64 timestamp) view returns(uint64, uint64, uint256)
func (_State *StateSession) GetStateDataByTime(id *big.Int, timestamp uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByTime(&_State.CallOpts, id, timestamp)
}

// GetStateDataByTime is a free data retrieval call binding the contract method 0x0281bec2.
//
// Solidity: function getStateDataByTime(uint256 id, uint64 timestamp) view returns(uint64, uint64, uint256)
func (_State *StateCallerSession) GetStateDataByTime(id *big.Int, timestamp uint64) (uint64, uint64, *big.Int, error) {
	return _State.Contract.GetStateDataByTime(&_State.CallOpts, id, timestamp)
}

// InitState is a paid mutator transaction binding the contract method 0x307ec167.
//
// Solidity: function initState(uint256 newState, uint256 genesisState, uint256 id, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactor) InitState(opts *bind.TransactOpts, newState *big.Int, genesisState *big.Int, id *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "initState", newState, genesisState, id, a, b, c)
}

// InitState is a paid mutator transaction binding the contract method 0x307ec167.
//
// Solidity: function initState(uint256 newState, uint256 genesisState, uint256 id, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateSession) InitState(newState *big.Int, genesisState *big.Int, id *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, a, b, c)
}

// InitState is a paid mutator transaction binding the contract method 0x307ec167.
//
// Solidity: function initState(uint256 newState, uint256 genesisState, uint256 id, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactorSession) InitState(newState *big.Int, genesisState *big.Int, id *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, a, b, c)
}

// SetState is a paid mutator transaction binding the contract method 0xbefb5b77.
//
// Solidity: function setState(uint256 newState, uint256 id, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactor) SetState(opts *bind.TransactOpts, newState *big.Int, id *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "setState", newState, id, a, b, c)
}

// SetState is a paid mutator transaction binding the contract method 0xbefb5b77.
//
// Solidity: function setState(uint256 newState, uint256 id, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateSession) SetState(newState *big.Int, id *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, a, b, c)
}

// SetState is a paid mutator transaction binding the contract method 0xbefb5b77.
//
// Solidity: function setState(uint256 newState, uint256 id, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_State *StateTransactorSession) SetState(newState *big.Int, id *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, a, b, c)
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
	Id        *big.Int
	BlockN    uint64
	Timestamp uint64
	State     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateUpdated is a free log retrieval operation binding the contract event 0x81c6f328b24014ef550c34a433275b52f3a8a0f32aa871adec069ab526a02390.
//
// Solidity: event StateUpdated(uint256 id, uint64 blockN, uint64 timestamp, uint256 state)
func (_State *StateFilterer) FilterStateUpdated(opts *bind.FilterOpts) (*StateStateUpdatedIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "StateUpdated")
	if err != nil {
		return nil, err
	}
	return &StateStateUpdatedIterator{contract: _State.contract, event: "StateUpdated", logs: logs, sub: sub}, nil
}

// WatchStateUpdated is a free log subscription operation binding the contract event 0x81c6f328b24014ef550c34a433275b52f3a8a0f32aa871adec069ab526a02390.
//
// Solidity: event StateUpdated(uint256 id, uint64 blockN, uint64 timestamp, uint256 state)
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

// ParseStateUpdated is a log parse operation binding the contract event 0x81c6f328b24014ef550c34a433275b52f3a8a0f32aa871adec069ab526a02390.
//
// Solidity: event StateUpdated(uint256 id, uint64 blockN, uint64 timestamp, uint256 state)
func (_State *StateFilterer) ParseStateUpdated(log types.Log) (*StateStateUpdated, error) {
	event := new(StateStateUpdated)
	if err := _State.contract.UnpackLog(event, "StateUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// VerifierABI is the input ABI used to generate the binding from.
const VerifierABI = "[{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[3]\",\"name\":\"input\",\"type\":\"uint256[3]\"}],\"name\":\"verifyProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"r\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// VerifierFuncSigs maps the 4-byte function signature to its string representation.
var VerifierFuncSigs = map[string]string{
	"11479fea": "verifyProof(uint256[2],uint256[2][2],uint256[2],uint256[3])",
}

// VerifierBin is the compiled bytecode used for deploying new contracts.
var VerifierBin = "0x608060405234801561001057600080fd5b50610f80806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806311479fea14610030575b600080fd5b61012c600480360361016081101561004757600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201829052506040805180820190915293969594608081019493509150600290835b828210156100c8576040805180820182529080840286019060029083908390808284376000920191909152505050815260019091019060200161008a565b5050604080518082018252939695948181019493509150600290839083908082843760009201919091525050604080516060818101909252929594938181019392509060039083908390808284376000920191909152509194506101409350505050565b604080519115158252519081900360200190f35b600061014a610e1f565b60408051808201825287518152602080890151818301529083528151608080820184528851518285019081528951840151606080850191909152908352845180860186528a85018051518252518501518186015283850152858401929092528351808501855288518152888401518185015285850152835160038082529181019094529092918201838036833701905050905060005b6003811015610219578481600381106101f557fe5b602002015182828151811061020657fe5b60209081029190910101526001016101e0565b506102248183610242565b6102335760019250505061023a565b6000925050505b949350505050565b60007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000161026d610e51565b61027561041a565b90508060800151518551600101146102c9576040805162461bcd60e51b81526020600482015260126024820152711d995c9a599a595c8b5898590b5a5b9c1d5d60721b604482015290519081900360640190fd5b6102d1610e98565b50604080518082019091526000808252602082018190525b86518110156103a357838782815181106102ff57fe5b602002602001015110610359576040805162461bcd60e51b815260206004820152601f60248201527f76657269666965722d6774652d736e61726b2d7363616c61722d6669656c6400604482015290519081900360640190fd5b610399826103948560800151846001018151811061037357fe5b60200260200101518a858151811061038757fe5b602002602001015161087e565b610913565b91506001016102e9565b506103c68183608001516000815181106103b957fe5b6020026020010151610913565b90506103fc6103d886600001516109a4565b8660200151846000015185602001518587604001518b604001518960600151610a30565b61040c5760019350505050610414565b600093505050505b92915050565b610422610e51565b6040805180820182527f1d1b564fe591b27c73b7a8b30af2bbe5280bd51d8580782806ca9cb9f35c7a9e81527e649d50eaadb8c1b56cb2794f9cc2c32704c08f2d276c0f80eaa1fe3ed63b726020808301919091529083528151608080820184527f2751f7664296eeff00c361c64cb0af1359295b3ba2cce285abfd87a5bb470ad78285019081527f25615130bb2ce70726ac5f1946fd009be3878d8683d56e575dd1658942c22c96606080850191909152908352845180860186527f22cc692dcbbec5325482fe335891395fe1a8f9f680a7e7cfdd0ba59aa6e3f36f81527f0c1693385daf2f1b7352afb8d8e8b9feedc40291537888d21278af242fd38cc2818601528385015285840192909252835180820185527f1ef928cd82271708275bdab4372d3ca067244ee8643619e633d8926761fc18398186019081527f2182344ca47325eba00c7a557b9be01909231e30b6a10eba808b21db10eb7a91828501528152845180860186527f26841b7038ddb849c5ade8ea4d99c8b1f9a67de31098136b37acc7078e35f6a381527f23706c84e9dd669cb75eb8c82dbac37c26f9a6cfaa80688eb4f9e7e942192777818601528185015285850152835190810184527f16bf1568da1fa5218b06cf6fc94877db15fe12e0f937677d0ba162627c4a93e48185019081527f2e77e0e7653db8924db6a05cc706f4da39e3c14d35fd329ed01a956ffd950f0c828401528152835180850185527f22ab01ae2ad480a7486d871cb6cbb34d7a43ff427e906cb1ece7f960639298da81527f28c0e4829e660b83eb4e3c897e0ceaacb76dd1ef0dd33ba4fa367718b8e0fb418185015281840152908401528151600480825260a08201909352919082015b6106a4610e98565b81526020019060019003908161069c57505060808201908152604080518082019091527f234aec3ebf46b30d274c50e769ad687677637517aef79baad44d4a654529007781527f036042b092fc73a49d8f767efa214669e8f870abc3e75b8622c67d6b92b45b3d60208201529051805160009061071d57fe5b602002602001018190525060405180604001604052807f24d059e9d9d3a62e1909d7b8beb4f714345f2b20c1b80f91940153e5259838aa81526020017f09336d9a5db25c97c9365950b1997166b600d1bb721a04ad02a3afc258c77830815250816080015160018151811061078e57fe5b602002602001018190525060405180604001604052807f21f254afe0780856e350a98c520dd4aa87b31713f768991f1d5dd90a7ef7e42e81526020017f29ab32045695753ef60d12202ed968839a6539920acac3b629aefd4311184cf981525081608001516002815181106107ff57fe5b602002602001018190525060405180604001604052807f2a5caeb40e8fbf6179e3c183941f870b457b4192608b59afe242ed52eae4a64681526020017f2f59d6e51f92dd378ac99a9c069a266e831e5159a9cdad06653aeb3da7ba14aa815250816080015160038151811061087057fe5b602002602001018190525090565b610886610e98565b61088e610eb2565b835181526020808501519082015260408101839052600060608360808460076107d05a03fa90508080156108c1576108c3565bfe5b508061090b576040805162461bcd60e51b81526020600482015260126024820152711c185a5c9a5b99cb5b5d5b0b59985a5b195960721b604482015290519081900360640190fd5b505092915050565b61091b610e98565b610923610ed0565b8351815260208085015181830152835160408301528301516060808301919091526000908360c08460066107d05a03fa90508080156108c157508061090b576040805162461bcd60e51b81526020600482015260126024820152711c185a5c9a5b99cb5859190b59985a5b195960721b604482015290519081900360640190fd5b6109ac610e98565b81517f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47901580156109df57506020830151155b156109ff5750506040805180820190915260008082526020820152610a2b565b60405180604001604052808460000151815260200182856020015181610a2157fe5b0683038152509150505b919050565b60408051600480825260a0820190925260009160609190816020015b610a54610e98565b815260200190600190039081610a4c57505060408051600480825260a0820190925291925060609190602082015b610a8a610eee565b815260200190600190039081610a825790505090508a82600081518110610aad57fe5b60200260200101819052508882600181518110610ac657fe5b60200260200101819052508682600281518110610adf57fe5b60200260200101819052508482600381518110610af857fe5b60200260200101819052508981600081518110610b1157fe5b60200260200101819052508781600181518110610b2a57fe5b60200260200101819052508581600281518110610b4357fe5b60200260200101819052508381600381518110610b5c57fe5b6020026020010181905250610b718282610b80565b9b9a5050505050505050505050565b60008151835114610bd1576040805162461bcd60e51b81526020600482015260166024820152751c185a5c9a5b99cb5b195b99dd1a1ccb59985a5b195960521b604482015290519081900360640190fd5b82516006810260608167ffffffffffffffff81118015610bf057600080fd5b50604051908082528060200260200182016040528015610c1a578160200160208202803683370190505b50905060005b83811015610d9f57868181518110610c3457fe5b602002602001015160000151828260060260000181518110610c5257fe5b602002602001018181525050868181518110610c6a57fe5b602002602001015160200151828260060260010181518110610c8857fe5b602002602001018181525050858181518110610ca057fe5b602090810291909101015151518251839060026006850201908110610cc157fe5b602002602001018181525050858181518110610cd957fe5b60209081029190910101515160016020020151828260060260030181518110610cfe57fe5b602002602001018181525050858181518110610d1657fe5b602002602001015160200151600060028110610d2e57fe5b6020020151828260060260040181518110610d4557fe5b602002602001018181525050858181518110610d5d57fe5b602002602001015160200151600160028110610d7557fe5b6020020151828260060260050181518110610d8c57fe5b6020908102919091010152600101610c20565b50610da8610f0e565b6000602082602086026020860160086107d05a03fa90508080156108c1575080610e11576040805162461bcd60e51b81526020600482015260156024820152741c185a5c9a5b99cb5bdc18dbd9194b59985a5b1959605a1b604482015290519081900360640190fd5b505115159695505050505050565b6040518060600160405280610e32610e98565b8152602001610e3f610eee565b8152602001610e4c610e98565b905290565b6040518060a00160405280610e64610e98565b8152602001610e71610eee565b8152602001610e7e610eee565b8152602001610e8b610eee565b8152602001606081525090565b604051806040016040528060008152602001600081525090565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280610f01610f2c565b8152602001610e4c610f2c565b60405180602001604052806001906020820280368337509192915050565b6040518060400160405280600290602082028036833750919291505056fea26469706673582212200359cfa79bd869b1553f6d9aceb5b7b32ecd302e9154fdcf08a0962fdc795b6264736f6c63430006060033"

// DeployVerifier deploys a new Ethereum contract, binding an instance of Verifier to it.
func DeployVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Verifier, error) {
	parsed, err := abi.JSON(strings.NewReader(VerifierABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// Verifier is an auto generated Go binding around an Ethereum contract.
type Verifier struct {
	VerifierCaller     // Read-only binding to the contract
	VerifierTransactor // Write-only binding to the contract
	VerifierFilterer   // Log filterer for contract events
}

// VerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierSession struct {
	Contract     *Verifier         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierCallerSession struct {
	Contract *VerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// VerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierTransactorSession struct {
	Contract     *VerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// VerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierRaw struct {
	Contract *Verifier // Generic contract binding to access the raw methods on
}

// VerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierCallerRaw struct {
	Contract *VerifierCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierTransactorRaw struct {
	Contract *VerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifier creates a new instance of Verifier, bound to a specific deployed contract.
func NewVerifier(address common.Address, backend bind.ContractBackend) (*Verifier, error) {
	contract, err := bindVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Verifier{VerifierCaller: VerifierCaller{contract: contract}, VerifierTransactor: VerifierTransactor{contract: contract}, VerifierFilterer: VerifierFilterer{contract: contract}}, nil
}

// NewVerifierCaller creates a new read-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierCaller(address common.Address, caller bind.ContractCaller) (*VerifierCaller, error) {
	contract, err := bindVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierCaller{contract: contract}, nil
}

// NewVerifierTransactor creates a new write-only instance of Verifier, bound to a specific deployed contract.
func NewVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierTransactor, error) {
	contract, err := bindVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierTransactor{contract: contract}, nil
}

// NewVerifierFilterer creates a new log filterer instance of Verifier, bound to a specific deployed contract.
func NewVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierFilterer, error) {
	contract, err := bindVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierFilterer{contract: contract}, nil
}

// bindVerifier binds a generic wrapper to an already deployed contract.
func bindVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VerifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.VerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.VerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Verifier *VerifierCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Verifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Verifier *VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Verifier *VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Verifier.Contract.contract.Transact(opts, method, params...)
}

// VerifyProof is a free data retrieval call binding the contract method 0x11479fea.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[3] input) view returns(bool r)
func (_Verifier *VerifierCaller) VerifyProof(opts *bind.CallOpts, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [3]*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Verifier.contract.Call(opts, out, "verifyProof", a, b, c, input)
	return *ret0, err
}

// VerifyProof is a free data retrieval call binding the contract method 0x11479fea.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[3] input) view returns(bool r)
func (_Verifier *VerifierSession) VerifyProof(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [3]*big.Int) (bool, error) {
	return _Verifier.Contract.VerifyProof(&_Verifier.CallOpts, a, b, c, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x11479fea.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[3] input) view returns(bool r)
func (_Verifier *VerifierCallerSession) VerifyProof(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [3]*big.Int) (bool, error) {
	return _Verifier.Contract.VerifyProof(&_Verifier.CallOpts, a, b, c, input)
}
