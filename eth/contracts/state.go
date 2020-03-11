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

// BabyJubJubABI is the input ABI used to generate the binding from.
const BabyJubJubABI = "[{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"q\",\"type\":\"uint256[2]\"}],\"name\":\"addition\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"r\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"p\",\"type\":\"uint256[2]\"}],\"name\":\"scalarmul\",\"outputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"r\",\"type\":\"uint256[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// BabyJubJubFuncSigs maps the 4-byte function signature to its string representation.
var BabyJubJubFuncSigs = map[string]string{
	"9c806946": "addition(uint256[2],uint256[2])",
	"4c8632cd": "scalarmul(uint256,uint256[2])",
}

// BabyJubJubBin is the compiled bytecode used for deploying new contracts.
var BabyJubJubBin = "0x6104b9610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100405760003560e01c80634c8632cd146100455780639c806946146100d3575b600080fd5b6100986004803603606081101561005b57600080fd5b6040805180820182528335939283019291606083019190602084019060029083908390808284376000920191909152509194506101479350505050565b6040518082600260200280838360005b838110156100c05781810151838201526020016100a8565b5050505090500191505060405180910390f35b610098600480360360808110156100e957600080fd5b6040805180820182529183019291818301918390600290839083908082843760009201919091525050604080518082018252929594938181019392509060029083908390808284376000920191909152509194506101d69350505050565b61014f610445565b600081526001602082015282610163610445565b50825b811561019d5781600116600114156101855761018283826101d6565b92505b61018f81826101d6565b9050600182901c9150610166565b82516000805160206104648339815191529006835260208301516000805160206104648339815191529006602084015250909392505050565b6101de610445565b81516020840151600091600080516020610464833981519152918291900960208501518651600080516020610464833981519152919009086020808501519086015185518751939450600093600080516020610464833981519152938493909284928391908290620292f8090909096001089050600061026c82600080516020610464833981519152610389565b905060006000805160206104648339815191528285098651885191925060009160008051602061046483398151915291908290620292fc090990506000600080516020610464833981519152828103816020808c0151908d0151090888518a5191925060009160008051602061046483398151915291908290620292f8090990506000600080516020610464833981519152808b600160200201516000805160206104648339815191528e60016020020151860909600080516020610464833981519152036001089050600061035082600080516020610464833981519152610389565b9050600060008051602061046483398151915282860960408051808201909152978852602088015250949b9a5050505050505050505050565b600082158061039757508183145b806103a0575081155b156103f2576040805162461bcd60e51b815260206004820152601960248201527f4572726f72206f6e20696e7075747320696e206d6f64696e7600000000000000604482015290519081900360640190fd5b8183600060015b821561043b5780868061040857fe5b878061041057fe5b8386888161041a57fe5b04098803840890925090508280858161042f57fe5b919550900692506103f9565b5095945050505050565b6040518060400160405280600290602082028038833950919291505056fe30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001a264697066735822122044ec224694dac494df2eadeb578ef057e8fda2d935a718a4837cd271c1f30ab164736f6c63430006010033"

// DeployBabyJubJub deploys a new Ethereum contract, binding an instance of BabyJubJub to it.
func DeployBabyJubJub(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BabyJubJub, error) {
	parsed, err := abi.JSON(strings.NewReader(BabyJubJubABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BabyJubJubBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BabyJubJub{BabyJubJubCaller: BabyJubJubCaller{contract: contract}, BabyJubJubTransactor: BabyJubJubTransactor{contract: contract}, BabyJubJubFilterer: BabyJubJubFilterer{contract: contract}}, nil
}

// BabyJubJub is an auto generated Go binding around an Ethereum contract.
type BabyJubJub struct {
	BabyJubJubCaller     // Read-only binding to the contract
	BabyJubJubTransactor // Write-only binding to the contract
	BabyJubJubFilterer   // Log filterer for contract events
}

// BabyJubJubCaller is an auto generated read-only Go binding around an Ethereum contract.
type BabyJubJubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BabyJubJubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BabyJubJubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BabyJubJubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BabyJubJubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BabyJubJubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BabyJubJubSession struct {
	Contract     *BabyJubJub       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BabyJubJubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BabyJubJubCallerSession struct {
	Contract *BabyJubJubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// BabyJubJubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BabyJubJubTransactorSession struct {
	Contract     *BabyJubJubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BabyJubJubRaw is an auto generated low-level Go binding around an Ethereum contract.
type BabyJubJubRaw struct {
	Contract *BabyJubJub // Generic contract binding to access the raw methods on
}

// BabyJubJubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BabyJubJubCallerRaw struct {
	Contract *BabyJubJubCaller // Generic read-only contract binding to access the raw methods on
}

// BabyJubJubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BabyJubJubTransactorRaw struct {
	Contract *BabyJubJubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBabyJubJub creates a new instance of BabyJubJub, bound to a specific deployed contract.
func NewBabyJubJub(address common.Address, backend bind.ContractBackend) (*BabyJubJub, error) {
	contract, err := bindBabyJubJub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BabyJubJub{BabyJubJubCaller: BabyJubJubCaller{contract: contract}, BabyJubJubTransactor: BabyJubJubTransactor{contract: contract}, BabyJubJubFilterer: BabyJubJubFilterer{contract: contract}}, nil
}

// NewBabyJubJubCaller creates a new read-only instance of BabyJubJub, bound to a specific deployed contract.
func NewBabyJubJubCaller(address common.Address, caller bind.ContractCaller) (*BabyJubJubCaller, error) {
	contract, err := bindBabyJubJub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BabyJubJubCaller{contract: contract}, nil
}

// NewBabyJubJubTransactor creates a new write-only instance of BabyJubJub, bound to a specific deployed contract.
func NewBabyJubJubTransactor(address common.Address, transactor bind.ContractTransactor) (*BabyJubJubTransactor, error) {
	contract, err := bindBabyJubJub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BabyJubJubTransactor{contract: contract}, nil
}

// NewBabyJubJubFilterer creates a new log filterer instance of BabyJubJub, bound to a specific deployed contract.
func NewBabyJubJubFilterer(address common.Address, filterer bind.ContractFilterer) (*BabyJubJubFilterer, error) {
	contract, err := bindBabyJubJub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BabyJubJubFilterer{contract: contract}, nil
}

// bindBabyJubJub binds a generic wrapper to an already deployed contract.
func bindBabyJubJub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BabyJubJubABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BabyJubJub *BabyJubJubRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BabyJubJub.Contract.BabyJubJubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BabyJubJub *BabyJubJubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BabyJubJub.Contract.BabyJubJubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BabyJubJub *BabyJubJubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BabyJubJub.Contract.BabyJubJubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BabyJubJub *BabyJubJubCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BabyJubJub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BabyJubJub *BabyJubJubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BabyJubJub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BabyJubJub *BabyJubJubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BabyJubJub.Contract.contract.Transact(opts, method, params...)
}

// Addition is a free data retrieval call binding the contract method 0x9c806946.
//
// Solidity: function addition(uint256[2] p, uint256[2] q) constant returns(uint256[2] r)
func (_BabyJubJub *BabyJubJubCaller) Addition(opts *bind.CallOpts, p [2]*big.Int, q [2]*big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _BabyJubJub.contract.Call(opts, out, "addition", p, q)
	return *ret0, err
}

// Addition is a free data retrieval call binding the contract method 0x9c806946.
//
// Solidity: function addition(uint256[2] p, uint256[2] q) constant returns(uint256[2] r)
func (_BabyJubJub *BabyJubJubSession) Addition(p [2]*big.Int, q [2]*big.Int) ([2]*big.Int, error) {
	return _BabyJubJub.Contract.Addition(&_BabyJubJub.CallOpts, p, q)
}

// Addition is a free data retrieval call binding the contract method 0x9c806946.
//
// Solidity: function addition(uint256[2] p, uint256[2] q) constant returns(uint256[2] r)
func (_BabyJubJub *BabyJubJubCallerSession) Addition(p [2]*big.Int, q [2]*big.Int) ([2]*big.Int, error) {
	return _BabyJubJub.Contract.Addition(&_BabyJubJub.CallOpts, p, q)
}

// Scalarmul is a free data retrieval call binding the contract method 0x4c8632cd.
//
// Solidity: function scalarmul(uint256 n, uint256[2] p) constant returns(uint256[2] r)
func (_BabyJubJub *BabyJubJubCaller) Scalarmul(opts *bind.CallOpts, n *big.Int, p [2]*big.Int) ([2]*big.Int, error) {
	var (
		ret0 = new([2]*big.Int)
	)
	out := ret0
	err := _BabyJubJub.contract.Call(opts, out, "scalarmul", n, p)
	return *ret0, err
}

// Scalarmul is a free data retrieval call binding the contract method 0x4c8632cd.
//
// Solidity: function scalarmul(uint256 n, uint256[2] p) constant returns(uint256[2] r)
func (_BabyJubJub *BabyJubJubSession) Scalarmul(n *big.Int, p [2]*big.Int) ([2]*big.Int, error) {
	return _BabyJubJub.Contract.Scalarmul(&_BabyJubJub.CallOpts, n, p)
}

// Scalarmul is a free data retrieval call binding the contract method 0x4c8632cd.
//
// Solidity: function scalarmul(uint256 n, uint256[2] p) constant returns(uint256[2] r)
func (_BabyJubJub *BabyJubJubCallerSession) Scalarmul(n *big.Int, p [2]*big.Int) ([2]*big.Int, error) {
	return _BabyJubJub.Contract.Scalarmul(&_BabyJubJub.CallOpts, n, p)
}

// EddsaBabyJubJubABI is the input ABI used to generate the binding from.
const EddsaBabyJubJubABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poseidonContractAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"pk\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"m\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"r\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"Verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// EddsaBabyJubJubFuncSigs maps the 4-byte function signature to its string representation.
var EddsaBabyJubJubFuncSigs = map[string]string{
	"da047321": "Verify(uint256[2],uint256,uint256[2],uint256)",
}

// EddsaBabyJubJubBin is the compiled bytecode used for deploying new contracts.
var EddsaBabyJubJubBin = "0x608060405234801561001057600080fd5b506040516105e23803806105e28339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b031990921691909117905561057d806100656000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063da04732114610030575b600080fd5b6100b1600480360360c081101561004657600080fd5b604080518082018252918301929181830191839060029083908390808284376000920191909152505060408051808201825292958435959094909360608201935091602090910190600290839083908082843760009201919091525091945050903591506100c59050565b604080519115158252519081900360200190f35b60408051600680825260e08201909252600091606091906020820160c080388339505085518251929350918391506000906100fc57fe5b602090810291909101015283600160200201518160018151811061011c57fe5b6020908102919091010152855181518290600290811061013857fe5b602090810291909101015285600160200201518160038151811061015857fe5b602002602001018181525050848160048151811061017257fe5b60200260200101818152505060008160058151811061018d57fe5b6020908102919091018101919091526000805460405163311083ed60e21b81526004810184815285516024830152855193946001600160a01b039093169363c4420fb4938793839260449091019185810191028083838b5b838110156101fd5781810151838201526020016101e5565b505050509050019250505060206040518083038186803b15801561022057600080fd5b505afa158015610234573d6000803e3d6000fd5b505050506040513d602081101561024a57600080fd5b50519050610256610529565b7f0bb77a6ad63e739b4eacb2e09d6277c12ab8d8010534e0b62893f3f6bb95705181527f25797203f7a0b24925572e1cd16bf9edfce0051fb9e133774b3c257a872d7d8b60208201526102a7610529565b73__$c062ba69f5356030b3a7f236ebd3c63787$__634c8632cd87846040518363ffffffff1660e01b81526004018083815260200182600260200280838360005b838110156103005781810151838201526020016102e8565b5050505090500192505050604080518083038186803b15801561032257600080fd5b505af4158015610336573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250604081101561035b57600080fd5b5090506008830261036a610529565b73__$c062ba69f5356030b3a7f236ebd3c63787$__634c8632cd838d6040518363ffffffff1660e01b81526004018083815260200182600260200280838360005b838110156103c35781810151838201526020016103ab565b5050505090500192505050604080518083038186803b1580156103e557600080fd5b505af41580156103f9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250604081101561041e57600080fd5b5060408051634e4034a360e11b815291925073__$c062ba69f5356030b3a7f236ebd3c63787$__91639c806946918c9185916004909101908190849080838360005b83811015610478578181015183820152602001610460565b5050505090500182600260200280838360005b838110156104a357818101518382015260200161048b565b5050505090500192505050604080518083038186803b1580156104c557600080fd5b505af41580156104d9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525060408110156104fe57600080fd5b508051845191925014801561051a575060208082015190840151145b9b9a5050505050505050505050565b6040518060400160405280600290602082028038833950919291505056fea2646970667358221220154ddb13d3acaec0728ec191f334a477b806d4aff915c39f8e7720cd3e5c74c864736f6c63430006010033"

// DeployEddsaBabyJubJub deploys a new Ethereum contract, binding an instance of EddsaBabyJubJub to it.
func DeployEddsaBabyJubJub(auth *bind.TransactOpts, backend bind.ContractBackend, _poseidonContractAddr common.Address) (common.Address, *types.Transaction, *EddsaBabyJubJub, error) {
	parsed, err := abi.JSON(strings.NewReader(EddsaBabyJubJubABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	babyJubJubAddr, _, _, _ := DeployBabyJubJub(auth, backend)
	EddsaBabyJubJubBin = strings.Replace(EddsaBabyJubJubBin, "__$c062ba69f5356030b3a7f236ebd3c63787$__", babyJubJubAddr.String()[2:], -1)

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EddsaBabyJubJubBin), backend, _poseidonContractAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EddsaBabyJubJub{EddsaBabyJubJubCaller: EddsaBabyJubJubCaller{contract: contract}, EddsaBabyJubJubTransactor: EddsaBabyJubJubTransactor{contract: contract}, EddsaBabyJubJubFilterer: EddsaBabyJubJubFilterer{contract: contract}}, nil
}

// EddsaBabyJubJub is an auto generated Go binding around an Ethereum contract.
type EddsaBabyJubJub struct {
	EddsaBabyJubJubCaller     // Read-only binding to the contract
	EddsaBabyJubJubTransactor // Write-only binding to the contract
	EddsaBabyJubJubFilterer   // Log filterer for contract events
}

// EddsaBabyJubJubCaller is an auto generated read-only Go binding around an Ethereum contract.
type EddsaBabyJubJubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EddsaBabyJubJubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EddsaBabyJubJubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EddsaBabyJubJubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EddsaBabyJubJubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EddsaBabyJubJubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EddsaBabyJubJubSession struct {
	Contract     *EddsaBabyJubJub  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EddsaBabyJubJubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EddsaBabyJubJubCallerSession struct {
	Contract *EddsaBabyJubJubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// EddsaBabyJubJubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EddsaBabyJubJubTransactorSession struct {
	Contract     *EddsaBabyJubJubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// EddsaBabyJubJubRaw is an auto generated low-level Go binding around an Ethereum contract.
type EddsaBabyJubJubRaw struct {
	Contract *EddsaBabyJubJub // Generic contract binding to access the raw methods on
}

// EddsaBabyJubJubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EddsaBabyJubJubCallerRaw struct {
	Contract *EddsaBabyJubJubCaller // Generic read-only contract binding to access the raw methods on
}

// EddsaBabyJubJubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EddsaBabyJubJubTransactorRaw struct {
	Contract *EddsaBabyJubJubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEddsaBabyJubJub creates a new instance of EddsaBabyJubJub, bound to a specific deployed contract.
func NewEddsaBabyJubJub(address common.Address, backend bind.ContractBackend) (*EddsaBabyJubJub, error) {
	contract, err := bindEddsaBabyJubJub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EddsaBabyJubJub{EddsaBabyJubJubCaller: EddsaBabyJubJubCaller{contract: contract}, EddsaBabyJubJubTransactor: EddsaBabyJubJubTransactor{contract: contract}, EddsaBabyJubJubFilterer: EddsaBabyJubJubFilterer{contract: contract}}, nil
}

// NewEddsaBabyJubJubCaller creates a new read-only instance of EddsaBabyJubJub, bound to a specific deployed contract.
func NewEddsaBabyJubJubCaller(address common.Address, caller bind.ContractCaller) (*EddsaBabyJubJubCaller, error) {
	contract, err := bindEddsaBabyJubJub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EddsaBabyJubJubCaller{contract: contract}, nil
}

// NewEddsaBabyJubJubTransactor creates a new write-only instance of EddsaBabyJubJub, bound to a specific deployed contract.
func NewEddsaBabyJubJubTransactor(address common.Address, transactor bind.ContractTransactor) (*EddsaBabyJubJubTransactor, error) {
	contract, err := bindEddsaBabyJubJub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EddsaBabyJubJubTransactor{contract: contract}, nil
}

// NewEddsaBabyJubJubFilterer creates a new log filterer instance of EddsaBabyJubJub, bound to a specific deployed contract.
func NewEddsaBabyJubJubFilterer(address common.Address, filterer bind.ContractFilterer) (*EddsaBabyJubJubFilterer, error) {
	contract, err := bindEddsaBabyJubJub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EddsaBabyJubJubFilterer{contract: contract}, nil
}

// bindEddsaBabyJubJub binds a generic wrapper to an already deployed contract.
func bindEddsaBabyJubJub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EddsaBabyJubJubABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EddsaBabyJubJub *EddsaBabyJubJubRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EddsaBabyJubJub.Contract.EddsaBabyJubJubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EddsaBabyJubJub *EddsaBabyJubJubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EddsaBabyJubJub.Contract.EddsaBabyJubJubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EddsaBabyJubJub *EddsaBabyJubJubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EddsaBabyJubJub.Contract.EddsaBabyJubJubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EddsaBabyJubJub *EddsaBabyJubJubCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EddsaBabyJubJub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EddsaBabyJubJub *EddsaBabyJubJubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EddsaBabyJubJub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EddsaBabyJubJub *EddsaBabyJubJubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EddsaBabyJubJub.Contract.contract.Transact(opts, method, params...)
}

// Verify is a free data retrieval call binding the contract method 0xda047321.
//
// Solidity: function Verify(uint256[2] pk, uint256 m, uint256[2] r, uint256 s) constant returns(bool)
func (_EddsaBabyJubJub *EddsaBabyJubJubCaller) Verify(opts *bind.CallOpts, pk [2]*big.Int, m *big.Int, r [2]*big.Int, s *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EddsaBabyJubJub.contract.Call(opts, out, "Verify", pk, m, r, s)
	return *ret0, err
}

// Verify is a free data retrieval call binding the contract method 0xda047321.
//
// Solidity: function Verify(uint256[2] pk, uint256 m, uint256[2] r, uint256 s) constant returns(bool)
func (_EddsaBabyJubJub *EddsaBabyJubJubSession) Verify(pk [2]*big.Int, m *big.Int, r [2]*big.Int, s *big.Int) (bool, error) {
	return _EddsaBabyJubJub.Contract.Verify(&_EddsaBabyJubJub.CallOpts, pk, m, r, s)
}

// Verify is a free data retrieval call binding the contract method 0xda047321.
//
// Solidity: function Verify(uint256[2] pk, uint256 m, uint256[2] r, uint256 s) constant returns(bool)
func (_EddsaBabyJubJub *EddsaBabyJubJubCallerSession) Verify(pk [2]*big.Int, m *big.Int, r [2]*big.Int, s *big.Int) (bool, error) {
	return _EddsaBabyJubJub.Contract.Verify(&_EddsaBabyJubJub.CallOpts, pk, m, r, s)
}

// PoseidonABI is the input ABI used to generate the binding from.
const PoseidonABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poseidonContractAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"inp\",\"type\":\"uint256[]\"}],\"name\":\"Hash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// PoseidonFuncSigs maps the 4-byte function signature to its string representation.
var PoseidonFuncSigs = map[string]string{
	"77d4ef3d": "Hash(uint256[])",
}

// PoseidonBin is the compiled bytecode used for deploying new contracts.
var PoseidonBin = "0x608060405234801561001057600080fd5b506040516102363803806102368339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b03199092169190911790556101d1806100656000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806377d4ef3d14610030575b600080fd5b6100d36004803603602081101561004657600080fd5b81019060208101813564010000000081111561006157600080fd5b82018360208201111561007357600080fd5b8035906020019184602083028401116401000000008311171561009557600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295506100e5945050505050565b60408051918252519081900360200190f35b6000805460405163311083ed60e21b81526020600482018181528551602484015285516001600160a01b039094169363c4420fb4938793839260449092019181860191028083838b5b8381101561014657818101518382015260200161012e565b505050509050019250505060206040518083038186803b15801561016957600080fd5b505afa15801561017d573d6000803e3d6000fd5b505050506040513d602081101561019357600080fd5b50519291505056fea2646970667358221220ac878c144f4b4e8ec153f9711145285a5d178411038ce62a445881fc2736b19064736f6c63430006010033"

// DeployPoseidon deploys a new Ethereum contract, binding an instance of Poseidon to it.
func DeployPoseidon(auth *bind.TransactOpts, backend bind.ContractBackend, _poseidonContractAddr common.Address) (common.Address, *types.Transaction, *Poseidon, error) {
	parsed, err := abi.JSON(strings.NewReader(PoseidonABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PoseidonBin), backend, _poseidonContractAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Poseidon{PoseidonCaller: PoseidonCaller{contract: contract}, PoseidonTransactor: PoseidonTransactor{contract: contract}, PoseidonFilterer: PoseidonFilterer{contract: contract}}, nil
}

// Poseidon is an auto generated Go binding around an Ethereum contract.
type Poseidon struct {
	PoseidonCaller     // Read-only binding to the contract
	PoseidonTransactor // Write-only binding to the contract
	PoseidonFilterer   // Log filterer for contract events
}

// PoseidonCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoseidonCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoseidonTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoseidonTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoseidonFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoseidonFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoseidonSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoseidonSession struct {
	Contract     *Poseidon         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoseidonCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoseidonCallerSession struct {
	Contract *PoseidonCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// PoseidonTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoseidonTransactorSession struct {
	Contract     *PoseidonTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// PoseidonRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoseidonRaw struct {
	Contract *Poseidon // Generic contract binding to access the raw methods on
}

// PoseidonCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoseidonCallerRaw struct {
	Contract *PoseidonCaller // Generic read-only contract binding to access the raw methods on
}

// PoseidonTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoseidonTransactorRaw struct {
	Contract *PoseidonTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoseidon creates a new instance of Poseidon, bound to a specific deployed contract.
func NewPoseidon(address common.Address, backend bind.ContractBackend) (*Poseidon, error) {
	contract, err := bindPoseidon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Poseidon{PoseidonCaller: PoseidonCaller{contract: contract}, PoseidonTransactor: PoseidonTransactor{contract: contract}, PoseidonFilterer: PoseidonFilterer{contract: contract}}, nil
}

// NewPoseidonCaller creates a new read-only instance of Poseidon, bound to a specific deployed contract.
func NewPoseidonCaller(address common.Address, caller bind.ContractCaller) (*PoseidonCaller, error) {
	contract, err := bindPoseidon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoseidonCaller{contract: contract}, nil
}

// NewPoseidonTransactor creates a new write-only instance of Poseidon, bound to a specific deployed contract.
func NewPoseidonTransactor(address common.Address, transactor bind.ContractTransactor) (*PoseidonTransactor, error) {
	contract, err := bindPoseidon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoseidonTransactor{contract: contract}, nil
}

// NewPoseidonFilterer creates a new log filterer instance of Poseidon, bound to a specific deployed contract.
func NewPoseidonFilterer(address common.Address, filterer bind.ContractFilterer) (*PoseidonFilterer, error) {
	contract, err := bindPoseidon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoseidonFilterer{contract: contract}, nil
}

// bindPoseidon binds a generic wrapper to an already deployed contract.
func bindPoseidon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PoseidonABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Poseidon *PoseidonRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Poseidon.Contract.PoseidonCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Poseidon *PoseidonRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Poseidon.Contract.PoseidonTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Poseidon *PoseidonRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Poseidon.Contract.PoseidonTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Poseidon *PoseidonCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Poseidon.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Poseidon *PoseidonTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Poseidon.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Poseidon *PoseidonTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Poseidon.Contract.contract.Transact(opts, method, params...)
}

// Hash is a free data retrieval call binding the contract method 0x77d4ef3d.
//
// Solidity: function Hash(uint256[] inp) constant returns(uint256)
func (_Poseidon *PoseidonCaller) Hash(opts *bind.CallOpts, inp []*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Poseidon.contract.Call(opts, out, "Hash", inp)
	return *ret0, err
}

// Hash is a free data retrieval call binding the contract method 0x77d4ef3d.
//
// Solidity: function Hash(uint256[] inp) constant returns(uint256)
func (_Poseidon *PoseidonSession) Hash(inp []*big.Int) (*big.Int, error) {
	return _Poseidon.Contract.Hash(&_Poseidon.CallOpts, inp)
}

// Hash is a free data retrieval call binding the contract method 0x77d4ef3d.
//
// Solidity: function Hash(uint256[] inp) constant returns(uint256)
func (_Poseidon *PoseidonCallerSession) Hash(inp []*big.Int) (*big.Int, error) {
	return _Poseidon.Contract.Hash(&_Poseidon.CallOpts, inp)
}

// PoseidonUnitABI is the input ABI used to generate the binding from.
const PoseidonUnitABI = "[{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"name\":\"poseidon\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// PoseidonUnitFuncSigs maps the 4-byte function signature to its string representation.
var PoseidonUnitFuncSigs = map[string]string{
	"c4420fb4": "poseidon(uint256[])",
}

// PoseidonUnitBin is the compiled bytecode used for deploying new contracts.
var PoseidonUnitBin = "0x608060405234801561001057600080fd5b50610118806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063c4420fb414602d575b600080fd5b60ca60048036036020811015604157600080fd5b810190602081018135640100000000811115605b57600080fd5b820183602082011115606c57600080fd5b80359060200191846020830284011164010000000083111715608d57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092955060dc945050505050565b60408051918252519081900360200190f35b5060009056fea2646970667358221220f5623e7bb2f730602eb13e5fe8999eb06d69f02039dd4bfa18e620cfb760b01564736f6c63430006010033"

// DeployPoseidonUnit deploys a new Ethereum contract, binding an instance of PoseidonUnit to it.
func DeployPoseidonUnit(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PoseidonUnit, error) {
	parsed, err := abi.JSON(strings.NewReader(PoseidonUnitABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PoseidonUnitBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PoseidonUnit{PoseidonUnitCaller: PoseidonUnitCaller{contract: contract}, PoseidonUnitTransactor: PoseidonUnitTransactor{contract: contract}, PoseidonUnitFilterer: PoseidonUnitFilterer{contract: contract}}, nil
}

// PoseidonUnit is an auto generated Go binding around an Ethereum contract.
type PoseidonUnit struct {
	PoseidonUnitCaller     // Read-only binding to the contract
	PoseidonUnitTransactor // Write-only binding to the contract
	PoseidonUnitFilterer   // Log filterer for contract events
}

// PoseidonUnitCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoseidonUnitCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoseidonUnitTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoseidonUnitTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoseidonUnitFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoseidonUnitFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoseidonUnitSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoseidonUnitSession struct {
	Contract     *PoseidonUnit     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoseidonUnitCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoseidonUnitCallerSession struct {
	Contract *PoseidonUnitCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PoseidonUnitTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoseidonUnitTransactorSession struct {
	Contract     *PoseidonUnitTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PoseidonUnitRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoseidonUnitRaw struct {
	Contract *PoseidonUnit // Generic contract binding to access the raw methods on
}

// PoseidonUnitCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoseidonUnitCallerRaw struct {
	Contract *PoseidonUnitCaller // Generic read-only contract binding to access the raw methods on
}

// PoseidonUnitTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoseidonUnitTransactorRaw struct {
	Contract *PoseidonUnitTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoseidonUnit creates a new instance of PoseidonUnit, bound to a specific deployed contract.
func NewPoseidonUnit(address common.Address, backend bind.ContractBackend) (*PoseidonUnit, error) {
	contract, err := bindPoseidonUnit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoseidonUnit{PoseidonUnitCaller: PoseidonUnitCaller{contract: contract}, PoseidonUnitTransactor: PoseidonUnitTransactor{contract: contract}, PoseidonUnitFilterer: PoseidonUnitFilterer{contract: contract}}, nil
}

// NewPoseidonUnitCaller creates a new read-only instance of PoseidonUnit, bound to a specific deployed contract.
func NewPoseidonUnitCaller(address common.Address, caller bind.ContractCaller) (*PoseidonUnitCaller, error) {
	contract, err := bindPoseidonUnit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoseidonUnitCaller{contract: contract}, nil
}

// NewPoseidonUnitTransactor creates a new write-only instance of PoseidonUnit, bound to a specific deployed contract.
func NewPoseidonUnitTransactor(address common.Address, transactor bind.ContractTransactor) (*PoseidonUnitTransactor, error) {
	contract, err := bindPoseidonUnit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoseidonUnitTransactor{contract: contract}, nil
}

// NewPoseidonUnitFilterer creates a new log filterer instance of PoseidonUnit, bound to a specific deployed contract.
func NewPoseidonUnitFilterer(address common.Address, filterer bind.ContractFilterer) (*PoseidonUnitFilterer, error) {
	contract, err := bindPoseidonUnit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoseidonUnitFilterer{contract: contract}, nil
}

// bindPoseidonUnit binds a generic wrapper to an already deployed contract.
func bindPoseidonUnit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PoseidonUnitABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoseidonUnit *PoseidonUnitRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PoseidonUnit.Contract.PoseidonUnitCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoseidonUnit *PoseidonUnitRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoseidonUnit.Contract.PoseidonUnitTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoseidonUnit *PoseidonUnitRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoseidonUnit.Contract.PoseidonUnitTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoseidonUnit *PoseidonUnitCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PoseidonUnit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoseidonUnit *PoseidonUnitTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoseidonUnit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoseidonUnit *PoseidonUnitTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoseidonUnit.Contract.contract.Transact(opts, method, params...)
}

// Poseidon is a free data retrieval call binding the contract method 0xc4420fb4.
//
// Solidity: function poseidon(uint256[] ) constant returns(uint256)
func (_PoseidonUnit *PoseidonUnitCaller) Poseidon(opts *bind.CallOpts, arg0 []*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PoseidonUnit.contract.Call(opts, out, "poseidon", arg0)
	return *ret0, err
}

// Poseidon is a free data retrieval call binding the contract method 0xc4420fb4.
//
// Solidity: function poseidon(uint256[] ) constant returns(uint256)
func (_PoseidonUnit *PoseidonUnitSession) Poseidon(arg0 []*big.Int) (*big.Int, error) {
	return _PoseidonUnit.Contract.Poseidon(&_PoseidonUnit.CallOpts, arg0)
}

// Poseidon is a free data retrieval call binding the contract method 0xc4420fb4.
//
// Solidity: function poseidon(uint256[] ) constant returns(uint256)
func (_PoseidonUnit *PoseidonUnitCallerSession) Poseidon(arg0 []*big.Int) (*big.Int, error) {
	return _PoseidonUnit.Contract.Poseidon(&_PoseidonUnit.CallOpts, arg0)
}

// StateABI is the input ABI used to generate the binding from.
const StateABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poseidonContractAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_eddsaBabyJubJubAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"state\",\"type\":\"bytes32\"}],\"name\":\"StateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint64\",\"name\":\"blockN\",\"type\":\"uint64\"}],\"name\":\"getStateDataByBlock\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"}],\"name\":\"getStateDataById\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getStateDataByTime\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"genesisState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint256[2]\",\"name\":\"kOp\",\"type\":\"uint256[2]\"},{\"internalType\":\"bytes\",\"name\":\"itp\",\"type\":\"bytes\"},{\"internalType\":\"uint256[2]\",\"name\":\"sigR8\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"sigS\",\"type\":\"uint256\"}],\"name\":\"initState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newState\",\"type\":\"bytes32\"},{\"internalType\":\"bytes31\",\"name\":\"id\",\"type\":\"bytes31\"},{\"internalType\":\"uint256[2]\",\"name\":\"kOp\",\"type\":\"uint256[2]\"},{\"internalType\":\"bytes\",\"name\":\"itp\",\"type\":\"bytes\"},{\"internalType\":\"uint256[2]\",\"name\":\"sigR8\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"sigS\",\"type\":\"uint256\"}],\"name\":\"setState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// StateFuncSigs maps the 4-byte function signature to its string representation.
var StateFuncSigs = map[string]string{
	"c1056e53": "getState(bytes31)",
	"c68631e1": "getStateDataByBlock(bytes31,uint64)",
	"4cabaefa": "getStateDataById(bytes31)",
	"5710773a": "getStateDataByTime(bytes31,uint64)",
	"9b0682df": "initState(bytes32,bytes32,bytes31,uint256[2],bytes,uint256[2],uint256)",
	"5f1fb4e0": "setState(bytes32,bytes31,uint256[2],bytes,uint256[2],uint256)",
}

// StateBin is the compiled bytecode used for deploying new contracts.
var StateBin = "0x608060405234801561001057600080fd5b50604051610e7d380380610e7d8339818101604052604081101561003357600080fd5b508051602090910151600180546001600160a01b039384166001600160a01b03199182161790915560008054939092169216919091179055610e038061007a6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80634cabaefa146100675780635710773a146100b25780635f1fb4e0146100e25780639b0682df146101f2578063c1056e5314610305578063c68631e114610338575b600080fd5b6100886004803603602081101561007d57600080fd5b503560ff1916610368565b604080516001600160401b0394851681529290931660208301528183015290519081900360600190f35b610088600480360360408110156100c857600080fd5b50803560ff191690602001356001600160401b0316610418565b6101f060048036036101008110156100f957600080fd5b60408051808201825283359360ff1960208201351693810192909160808301918084019060029083908390808284376000920191909152509194939260208101925035905064010000000081111561015057600080fd5b82018360208201111561016257600080fd5b8035906020019184600183028401116401000000008311171561018457600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080518082018252939695948181019493509150600290839083908082843760009201919091525091945050903591506108109050565b005b6101f0600480360361012081101561020957600080fd5b60408051808201825283359360208101359360ff198483013516939082019260a0830191606084019060029083908390808284376000920191909152509194939260208101925035905064010000000081111561026557600080fd5b82018360208201111561027757600080fd5b8035906020019184600183028401116401000000008311171561029957600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080518082018252939695948181019493509150600290839083908082843760009201919091525091945050903591506108f99050565b6103266004803603602081101561031b57600080fd5b503560ff1916610926565b60408051918252519081900360200190f35b6100886004803603604081101561034e57600080fd5b50803560ff191690602001356001600160401b0316610984565b60ff19811660009081526002602052604081205481908190610394575050600354600091508190610411565b61039c610d8c565b60ff1985166000908152600260205260409020805460001981019081106103bf57fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b909204169484018590526001909101549290910182905295509093509150505b9193909250565b600080600042846001600160401b03161061046f576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff198516600090815260026020526040902054610497575050600354600091508190610809565b60ff1985166000908152600260205260408120805460001981019081106104ba57fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116915085168110156105a75760ff19861660009081526002602052604090208054600019810190811061050b57fe5b600091825260208083206002928302015460ff198a168452919052604090912080546001600160401b0390921691600019810190811061054757fe5b600091825260208083206002928302015460ff198b16845291905260409091208054600160401b9092046001600160401b031691600019810190811061058957fe5b90600052602060002090600202016001015493509350935050610809565b60ff198616600090815260026020526040812054600019015b8082116107f9576000600282840160ff198b16600090815260026020526040902080549290910492506001600160401b038a1691839081106105fe57fe5b6000918252602090912060029091020154600160401b90046001600160401b031614156106db5760ff198916600090815260026020526040902080548290811061064457fe5b600091825260208083206002928302015460ff198d168452919052604090912080546001600160401b03909216918390811061067c57fe5b600091825260208083206002928302015460ff198e16845291905260409091208054600160401b9092046001600160401b031691849081106106ba57fe5b90600052602060002090600202016001015496509650965050505050610809565b60ff19891660009081526002602052604090208054829081106106fa57fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908916118015610772575060ff198916600090815260026020526040902080546001830190811061074a57fe5b60009182526020909120600290910201546001600160401b03600160401b9091048116908916105b156107965760ff198916600090815260026020526040902080548290811061064457fe5b60ff19891660009081526002602052604090208054829081106107b557fe5b60009182526020909120600290910201546001600160401b03600160401b909104811690891611156107ec578060010192506107f3565b6001810391505b506105c0565b5050600354600094508493509150505b9250925092565b60ff19851660009081526002602052604090205461082d57600080fd5b610835610d8c565b60ff19861660009081526002602052604090208054600019810190811061085857fe5b600091825260209182902060408051606081018252600290930290910180546001600160401b03808216808652600160401b9092041694840194909452600101549082015291504314156108dd5760405162461bcd60e51b8152600401808060200182810382526021815260200180610dad6021913960400191505060405180910390fd5b6108f08782604001518888888888610c0f565b50505050505050565b60ff1985166000908152600260205260409020541561091757600080fd5b6108f087878787878787610c0f565b60ff198116600090815260026020526040812054610947575060035461097f565b60ff19821660009081526002602052604090208054600019810190811061096a57fe5b90600052602060002090600202016001015490505b919050565b600080600043846001600160401b0316106109db576040805162461bcd60e51b8152602060048201526012602482015271195c9c939bd19d5d1d5c99505b1b1bddd95960721b604482015290519081900360640190fd5b60ff198516600090815260026020526040902054610a03575050600354600091508190610809565b60ff198516600090815260026020526040812080546000198101908110610a2657fe5b60009182526020909120600290910201546001600160401b0390811691508516811015610a705760ff19861660009081526002602052604090208054600019810190811061050b57fe5b60ff198616600090815260026020526040812054600019015b8082116107f9576000600282840160ff198b16600090815260026020526040902080549290910492506001600160401b038a169183908110610ac757fe5b60009182526020909120600290910201546001600160401b03161415610b065760ff198916600090815260026020526040902080548290811061064457fe5b60ff1989166000908152600260205260409020805482908110610b2557fe5b60009182526020909120600290910201546001600160401b03908116908916118015610b8f575060ff1989166000908152600260205260409020805460018301908110610b6e57fe5b60009182526020909120600290910201546001600160401b03908116908916105b15610bb35760ff198916600090815260026020526040902080548290811061064457fe5b60ff1989166000908152600260205260409020805482908110610bd257fe5b60009182526020909120600290910201546001600160401b039081169089161115610c0257806001019250610c09565b6001810391505b50610a89565b610c1a868885610d83565b1515600114610c2857600080fd5b600260008660ff191660ff191681526020019081526020016000206040518060600160405280436001600160401b03168152602001426001600160401b0316815260200189815250908060018154018082558091505060019003906000526020600020906002020160009091909190915060008201518160000160006101000a8154816001600160401b0302191690836001600160401b0316021790555060208201518160000160086101000a8154816001600160401b0302191690836001600160401b031602179055506040820151816001015550507fbe8b490a2c0932bd7e7463b73a740738a6a4f6a99ac6e50aa73ddc788673c88d8543428a604051808560ff191660ff19168152602001846001600160401b03166001600160401b03168152602001836001600160401b03166001600160401b0316815260200182815260200194505050505060405180910390a150505050505050565b60019392505050565b60408051606081018252600080825260208201819052918101919091529056fe6e6f206d756c7469706c652073657420696e207468652073616d6520626c6f636ba264697066735822122019b293b7f759af56d4c24222d078648514ad939a5079b0b7c7d94d94e41420c864736f6c63430006010033"

// DeployState deploys a new Ethereum contract, binding an instance of State to it.
func DeployState(auth *bind.TransactOpts, backend bind.ContractBackend, _poseidonContractAddr common.Address, _eddsaBabyJubJubAddr common.Address) (common.Address, *types.Transaction, *State, error) {
	parsed, err := abi.JSON(strings.NewReader(StateABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StateBin), backend, _poseidonContractAddr, _eddsaBabyJubJubAddr)
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

// InitState is a paid mutator transaction binding the contract method 0x9b0682df.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, uint256[2] kOp, bytes itp, uint256[2] sigR8, uint256 sigS) returns()
func (_State *StateTransactor) InitState(opts *bind.TransactOpts, newState [32]byte, genesisState [32]byte, id [31]byte, kOp [2]*big.Int, itp []byte, sigR8 [2]*big.Int, sigS *big.Int) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "initState", newState, genesisState, id, kOp, itp, sigR8, sigS)
}

// InitState is a paid mutator transaction binding the contract method 0x9b0682df.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, uint256[2] kOp, bytes itp, uint256[2] sigR8, uint256 sigS) returns()
func (_State *StateSession) InitState(newState [32]byte, genesisState [32]byte, id [31]byte, kOp [2]*big.Int, itp []byte, sigR8 [2]*big.Int, sigS *big.Int) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, kOp, itp, sigR8, sigS)
}

// InitState is a paid mutator transaction binding the contract method 0x9b0682df.
//
// Solidity: function initState(bytes32 newState, bytes32 genesisState, bytes31 id, uint256[2] kOp, bytes itp, uint256[2] sigR8, uint256 sigS) returns()
func (_State *StateTransactorSession) InitState(newState [32]byte, genesisState [32]byte, id [31]byte, kOp [2]*big.Int, itp []byte, sigR8 [2]*big.Int, sigS *big.Int) (*types.Transaction, error) {
	return _State.Contract.InitState(&_State.TransactOpts, newState, genesisState, id, kOp, itp, sigR8, sigS)
}

// SetState is a paid mutator transaction binding the contract method 0x5f1fb4e0.
//
// Solidity: function setState(bytes32 newState, bytes31 id, uint256[2] kOp, bytes itp, uint256[2] sigR8, uint256 sigS) returns()
func (_State *StateTransactor) SetState(opts *bind.TransactOpts, newState [32]byte, id [31]byte, kOp [2]*big.Int, itp []byte, sigR8 [2]*big.Int, sigS *big.Int) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "setState", newState, id, kOp, itp, sigR8, sigS)
}

// SetState is a paid mutator transaction binding the contract method 0x5f1fb4e0.
//
// Solidity: function setState(bytes32 newState, bytes31 id, uint256[2] kOp, bytes itp, uint256[2] sigR8, uint256 sigS) returns()
func (_State *StateSession) SetState(newState [32]byte, id [31]byte, kOp [2]*big.Int, itp []byte, sigR8 [2]*big.Int, sigS *big.Int) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, kOp, itp, sigR8, sigS)
}

// SetState is a paid mutator transaction binding the contract method 0x5f1fb4e0.
//
// Solidity: function setState(bytes32 newState, bytes31 id, uint256[2] kOp, bytes itp, uint256[2] sigR8, uint256 sigS) returns()
func (_State *StateTransactorSession) SetState(newState [32]byte, id [31]byte, kOp [2]*big.Int, itp []byte, sigR8 [2]*big.Int, sigS *big.Int) (*types.Transaction, error) {
	return _State.Contract.SetState(&_State.TransactOpts, newState, id, kOp, itp, sigR8, sigS)
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
