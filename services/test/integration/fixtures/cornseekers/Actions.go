// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cornseekers

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ActionsMetaData contains all meta data concerning the Actions contract.
var ActionsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"sid\",\"type\":\"uint32\"},{\"internalType\":\"enumDirection\",\"name\":\"dir\",\"type\":\"uint8\"}],\"name\":\"MOVE_SEEKER\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RESET_MAP\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"NodeID\",\"name\":\"seedID\",\"type\":\"uint224\"},{\"internalType\":\"uint32\",\"name\":\"entropy\",\"type\":\"uint32\"}],\"name\":\"REVEAL_SEED\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"sid\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"x\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"y\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"str\",\"type\":\"uint8\"}],\"name\":\"SPAWN_SEEKER\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ActionsABI is the input ABI used to generate the binding from.
// Deprecated: Use ActionsMetaData.ABI instead.
var ActionsABI = ActionsMetaData.ABI

// Actions is an auto generated Go binding around an Ethereum contract.
type Actions struct {
	ActionsCaller     // Read-only binding to the contract
	ActionsTransactor // Write-only binding to the contract
	ActionsFilterer   // Log filterer for contract events
}

// ActionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ActionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ActionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ActionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ActionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ActionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ActionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ActionsSession struct {
	Contract     *Actions          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ActionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ActionsCallerSession struct {
	Contract *ActionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ActionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ActionsTransactorSession struct {
	Contract     *ActionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ActionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ActionsRaw struct {
	Contract *Actions // Generic contract binding to access the raw methods on
}

// ActionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ActionsCallerRaw struct {
	Contract *ActionsCaller // Generic read-only contract binding to access the raw methods on
}

// ActionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ActionsTransactorRaw struct {
	Contract *ActionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewActions creates a new instance of Actions, bound to a specific deployed contract.
func NewActions(address common.Address, backend bind.ContractBackend) (*Actions, error) {
	contract, err := bindActions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Actions{ActionsCaller: ActionsCaller{contract: contract}, ActionsTransactor: ActionsTransactor{contract: contract}, ActionsFilterer: ActionsFilterer{contract: contract}}, nil
}

// NewActionsCaller creates a new read-only instance of Actions, bound to a specific deployed contract.
func NewActionsCaller(address common.Address, caller bind.ContractCaller) (*ActionsCaller, error) {
	contract, err := bindActions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ActionsCaller{contract: contract}, nil
}

// NewActionsTransactor creates a new write-only instance of Actions, bound to a specific deployed contract.
func NewActionsTransactor(address common.Address, transactor bind.ContractTransactor) (*ActionsTransactor, error) {
	contract, err := bindActions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ActionsTransactor{contract: contract}, nil
}

// NewActionsFilterer creates a new log filterer instance of Actions, bound to a specific deployed contract.
func NewActionsFilterer(address common.Address, filterer bind.ContractFilterer) (*ActionsFilterer, error) {
	contract, err := bindActions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ActionsFilterer{contract: contract}, nil
}

// bindActions binds a generic wrapper to an already deployed contract.
func bindActions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ActionsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Actions *ActionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Actions.Contract.ActionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Actions *ActionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Actions.Contract.ActionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Actions *ActionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Actions.Contract.ActionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Actions *ActionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Actions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Actions *ActionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Actions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Actions *ActionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Actions.Contract.contract.Transact(opts, method, params...)
}

// MOVESEEKER is a paid mutator transaction binding the contract method 0x449a80c3.
//
// Solidity: function MOVE_SEEKER(uint32 sid, uint8 dir) returns()
func (_Actions *ActionsTransactor) MOVESEEKER(opts *bind.TransactOpts, sid uint32, dir uint8) (*types.Transaction, error) {
	return _Actions.contract.Transact(opts, "MOVE_SEEKER", sid, dir)
}

// MOVESEEKER is a paid mutator transaction binding the contract method 0x449a80c3.
//
// Solidity: function MOVE_SEEKER(uint32 sid, uint8 dir) returns()
func (_Actions *ActionsSession) MOVESEEKER(sid uint32, dir uint8) (*types.Transaction, error) {
	return _Actions.Contract.MOVESEEKER(&_Actions.TransactOpts, sid, dir)
}

// MOVESEEKER is a paid mutator transaction binding the contract method 0x449a80c3.
//
// Solidity: function MOVE_SEEKER(uint32 sid, uint8 dir) returns()
func (_Actions *ActionsTransactorSession) MOVESEEKER(sid uint32, dir uint8) (*types.Transaction, error) {
	return _Actions.Contract.MOVESEEKER(&_Actions.TransactOpts, sid, dir)
}

// RESETMAP is a paid mutator transaction binding the contract method 0x16dfc800.
//
// Solidity: function RESET_MAP() returns()
func (_Actions *ActionsTransactor) RESETMAP(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Actions.contract.Transact(opts, "RESET_MAP")
}

// RESETMAP is a paid mutator transaction binding the contract method 0x16dfc800.
//
// Solidity: function RESET_MAP() returns()
func (_Actions *ActionsSession) RESETMAP() (*types.Transaction, error) {
	return _Actions.Contract.RESETMAP(&_Actions.TransactOpts)
}

// RESETMAP is a paid mutator transaction binding the contract method 0x16dfc800.
//
// Solidity: function RESET_MAP() returns()
func (_Actions *ActionsTransactorSession) RESETMAP() (*types.Transaction, error) {
	return _Actions.Contract.RESETMAP(&_Actions.TransactOpts)
}

// REVEALSEED is a paid mutator transaction binding the contract method 0x5453a01a.
//
// Solidity: function REVEAL_SEED(uint224 seedID, uint32 entropy) returns()
func (_Actions *ActionsTransactor) REVEALSEED(opts *bind.TransactOpts, seedID *big.Int, entropy uint32) (*types.Transaction, error) {
	return _Actions.contract.Transact(opts, "REVEAL_SEED", seedID, entropy)
}

// REVEALSEED is a paid mutator transaction binding the contract method 0x5453a01a.
//
// Solidity: function REVEAL_SEED(uint224 seedID, uint32 entropy) returns()
func (_Actions *ActionsSession) REVEALSEED(seedID *big.Int, entropy uint32) (*types.Transaction, error) {
	return _Actions.Contract.REVEALSEED(&_Actions.TransactOpts, seedID, entropy)
}

// REVEALSEED is a paid mutator transaction binding the contract method 0x5453a01a.
//
// Solidity: function REVEAL_SEED(uint224 seedID, uint32 entropy) returns()
func (_Actions *ActionsTransactorSession) REVEALSEED(seedID *big.Int, entropy uint32) (*types.Transaction, error) {
	return _Actions.Contract.REVEALSEED(&_Actions.TransactOpts, seedID, entropy)
}

// SPAWNSEEKER is a paid mutator transaction binding the contract method 0x496ac2d3.
//
// Solidity: function SPAWN_SEEKER(uint32 sid, uint8 x, uint8 y, uint8 str) returns()
func (_Actions *ActionsTransactor) SPAWNSEEKER(opts *bind.TransactOpts, sid uint32, x uint8, y uint8, str uint8) (*types.Transaction, error) {
	return _Actions.contract.Transact(opts, "SPAWN_SEEKER", sid, x, y, str)
}

// SPAWNSEEKER is a paid mutator transaction binding the contract method 0x496ac2d3.
//
// Solidity: function SPAWN_SEEKER(uint32 sid, uint8 x, uint8 y, uint8 str) returns()
func (_Actions *ActionsSession) SPAWNSEEKER(sid uint32, x uint8, y uint8, str uint8) (*types.Transaction, error) {
	return _Actions.Contract.SPAWNSEEKER(&_Actions.TransactOpts, sid, x, y, str)
}

// SPAWNSEEKER is a paid mutator transaction binding the contract method 0x496ac2d3.
//
// Solidity: function SPAWN_SEEKER(uint32 sid, uint8 x, uint8 y, uint8 str) returns()
func (_Actions *ActionsTransactorSession) SPAWNSEEKER(sid uint32, x uint8, y uint8, str uint8) (*types.Transaction, error) {
	return _Actions.Contract.SPAWNSEEKER(&_Actions.TransactOpts, sid, x, y, str)
}
