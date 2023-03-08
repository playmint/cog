// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package game

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
	_ = abi.ConvertType
)

// GameMetadata is an auto generated low-level Go binding around an user-defined struct.
type GameMetadata struct {
	Name string
	Url  string
}

// BaseGameMetaData contains all meta data concerning the BaseGame contract.
var BaseGameMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dispatcherAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stateAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"routerAddr\",\"type\":\"address\"}],\"name\":\"GameDeployed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getDispatcher\",\"outputs\":[{\"internalType\":\"contractDispatcher\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMetadata\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"internalType\":\"structGameMetadata\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"contractRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"contractState\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BaseGameABI is the input ABI used to generate the binding from.
// Deprecated: Use BaseGameMetaData.ABI instead.
var BaseGameABI = BaseGameMetaData.ABI

// BaseGame is an auto generated Go binding around an Ethereum contract.
type BaseGame struct {
	BaseGameCaller     // Read-only binding to the contract
	BaseGameTransactor // Write-only binding to the contract
	BaseGameFilterer   // Log filterer for contract events
}

// BaseGameCaller is an auto generated read-only Go binding around an Ethereum contract.
type BaseGameCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseGameTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BaseGameTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseGameFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BaseGameFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseGameSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BaseGameSession struct {
	Contract     *BaseGame         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BaseGameCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BaseGameCallerSession struct {
	Contract *BaseGameCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BaseGameTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BaseGameTransactorSession struct {
	Contract     *BaseGameTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BaseGameRaw is an auto generated low-level Go binding around an Ethereum contract.
type BaseGameRaw struct {
	Contract *BaseGame // Generic contract binding to access the raw methods on
}

// BaseGameCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BaseGameCallerRaw struct {
	Contract *BaseGameCaller // Generic read-only contract binding to access the raw methods on
}

// BaseGameTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BaseGameTransactorRaw struct {
	Contract *BaseGameTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBaseGame creates a new instance of BaseGame, bound to a specific deployed contract.
func NewBaseGame(address common.Address, backend bind.ContractBackend) (*BaseGame, error) {
	contract, err := bindBaseGame(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BaseGame{BaseGameCaller: BaseGameCaller{contract: contract}, BaseGameTransactor: BaseGameTransactor{contract: contract}, BaseGameFilterer: BaseGameFilterer{contract: contract}}, nil
}

// NewBaseGameCaller creates a new read-only instance of BaseGame, bound to a specific deployed contract.
func NewBaseGameCaller(address common.Address, caller bind.ContractCaller) (*BaseGameCaller, error) {
	contract, err := bindBaseGame(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BaseGameCaller{contract: contract}, nil
}

// NewBaseGameTransactor creates a new write-only instance of BaseGame, bound to a specific deployed contract.
func NewBaseGameTransactor(address common.Address, transactor bind.ContractTransactor) (*BaseGameTransactor, error) {
	contract, err := bindBaseGame(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BaseGameTransactor{contract: contract}, nil
}

// NewBaseGameFilterer creates a new log filterer instance of BaseGame, bound to a specific deployed contract.
func NewBaseGameFilterer(address common.Address, filterer bind.ContractFilterer) (*BaseGameFilterer, error) {
	contract, err := bindBaseGame(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BaseGameFilterer{contract: contract}, nil
}

// bindBaseGame binds a generic wrapper to an already deployed contract.
func bindBaseGame(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BaseGameMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseGame *BaseGameRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseGame.Contract.BaseGameCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseGame *BaseGameRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseGame.Contract.BaseGameTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseGame *BaseGameRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseGame.Contract.BaseGameTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseGame *BaseGameCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseGame.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseGame *BaseGameTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseGame.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseGame *BaseGameTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseGame.Contract.contract.Transact(opts, method, params...)
}

// GetDispatcher is a free data retrieval call binding the contract method 0xebb3d589.
//
// Solidity: function getDispatcher() view returns(address)
func (_BaseGame *BaseGameCaller) GetDispatcher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BaseGame.contract.Call(opts, &out, "getDispatcher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDispatcher is a free data retrieval call binding the contract method 0xebb3d589.
//
// Solidity: function getDispatcher() view returns(address)
func (_BaseGame *BaseGameSession) GetDispatcher() (common.Address, error) {
	return _BaseGame.Contract.GetDispatcher(&_BaseGame.CallOpts)
}

// GetDispatcher is a free data retrieval call binding the contract method 0xebb3d589.
//
// Solidity: function getDispatcher() view returns(address)
func (_BaseGame *BaseGameCallerSession) GetDispatcher() (common.Address, error) {
	return _BaseGame.Contract.GetDispatcher(&_BaseGame.CallOpts)
}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns((string,string))
func (_BaseGame *BaseGameCaller) GetMetadata(opts *bind.CallOpts) (GameMetadata, error) {
	var out []interface{}
	err := _BaseGame.contract.Call(opts, &out, "getMetadata")

	if err != nil {
		return *new(GameMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(GameMetadata)).(*GameMetadata)

	return out0, err

}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns((string,string))
func (_BaseGame *BaseGameSession) GetMetadata() (GameMetadata, error) {
	return _BaseGame.Contract.GetMetadata(&_BaseGame.CallOpts)
}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns((string,string))
func (_BaseGame *BaseGameCallerSession) GetMetadata() (GameMetadata, error) {
	return _BaseGame.Contract.GetMetadata(&_BaseGame.CallOpts)
}

// GetRouter is a free data retrieval call binding the contract method 0xb0f479a1.
//
// Solidity: function getRouter() view returns(address)
func (_BaseGame *BaseGameCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BaseGame.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRouter is a free data retrieval call binding the contract method 0xb0f479a1.
//
// Solidity: function getRouter() view returns(address)
func (_BaseGame *BaseGameSession) GetRouter() (common.Address, error) {
	return _BaseGame.Contract.GetRouter(&_BaseGame.CallOpts)
}

// GetRouter is a free data retrieval call binding the contract method 0xb0f479a1.
//
// Solidity: function getRouter() view returns(address)
func (_BaseGame *BaseGameCallerSession) GetRouter() (common.Address, error) {
	return _BaseGame.Contract.GetRouter(&_BaseGame.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address)
func (_BaseGame *BaseGameCaller) GetState(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BaseGame.contract.Call(opts, &out, "getState")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address)
func (_BaseGame *BaseGameSession) GetState() (common.Address, error) {
	return _BaseGame.Contract.GetState(&_BaseGame.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address)
func (_BaseGame *BaseGameCallerSession) GetState() (common.Address, error) {
	return _BaseGame.Contract.GetState(&_BaseGame.CallOpts)
}

// BaseGameGameDeployedIterator is returned from FilterGameDeployed and is used to iterate over the raw logs and unpacked data for GameDeployed events raised by the BaseGame contract.
type BaseGameGameDeployedIterator struct {
	Event *BaseGameGameDeployed // Event containing the contract specifics and raw log

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
func (it *BaseGameGameDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseGameGameDeployed)
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
		it.Event = new(BaseGameGameDeployed)
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
func (it *BaseGameGameDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseGameGameDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseGameGameDeployed represents a GameDeployed event raised by the BaseGame contract.
type BaseGameGameDeployed struct {
	DispatcherAddr common.Address
	StateAddr      common.Address
	RouterAddr     common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterGameDeployed is a free log retrieval operation binding the contract event 0xf983dbfb3d7ef46fae38b82f406ef6895a8239899ab04e0d0cbd6d78a84c6a73.
//
// Solidity: event GameDeployed(address dispatcherAddr, address stateAddr, address routerAddr)
func (_BaseGame *BaseGameFilterer) FilterGameDeployed(opts *bind.FilterOpts) (*BaseGameGameDeployedIterator, error) {

	logs, sub, err := _BaseGame.contract.FilterLogs(opts, "GameDeployed")
	if err != nil {
		return nil, err
	}
	return &BaseGameGameDeployedIterator{contract: _BaseGame.contract, event: "GameDeployed", logs: logs, sub: sub}, nil
}

// WatchGameDeployed is a free log subscription operation binding the contract event 0xf983dbfb3d7ef46fae38b82f406ef6895a8239899ab04e0d0cbd6d78a84c6a73.
//
// Solidity: event GameDeployed(address dispatcherAddr, address stateAddr, address routerAddr)
func (_BaseGame *BaseGameFilterer) WatchGameDeployed(opts *bind.WatchOpts, sink chan<- *BaseGameGameDeployed) (event.Subscription, error) {

	logs, sub, err := _BaseGame.contract.WatchLogs(opts, "GameDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseGameGameDeployed)
				if err := _BaseGame.contract.UnpackLog(event, "GameDeployed", log); err != nil {
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

// ParseGameDeployed is a log parse operation binding the contract event 0xf983dbfb3d7ef46fae38b82f406ef6895a8239899ab04e0d0cbd6d78a84c6a73.
//
// Solidity: event GameDeployed(address dispatcherAddr, address stateAddr, address routerAddr)
func (_BaseGame *BaseGameFilterer) ParseGameDeployed(log types.Log) (*BaseGameGameDeployed, error) {
	event := new(BaseGameGameDeployed)
	if err := _BaseGame.contract.UnpackLog(event, "GameDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
