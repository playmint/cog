// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package router

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

// SessionRouterMetaData contains all meta data concerning the SessionRouter contract.
var SessionRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"SessionExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SessionExpiryTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SessionUnauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"session\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"exp\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"scopes\",\"type\":\"uint32\"}],\"name\":\"SessionCreate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"session\",\"type\":\"address\"}],\"name\":\"SessionDestroy\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractDispatcher\",\"name\":\"dispatcher\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"ttl\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"scopes\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sessionAddr\",\"type\":\"address\"}],\"name\":\"authorizeAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractDispatcher\",\"name\":\"dispatcher\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"ttl\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"scopes\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sessionAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"authorizeAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[][]\",\"name\":\"actions\",\"type\":\"bytes[][]\"},{\"internalType\":\"bytes[]\",\"name\":\"sig\",\"type\":\"bytes[]\"}],\"name\":\"dispatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"revokeAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"revokeAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sessions\",\"outputs\":[{\"internalType\":\"contractDispatcher\",\"name\":\"dispatcher\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"exp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"scopes\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SessionRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use SessionRouterMetaData.ABI instead.
var SessionRouterABI = SessionRouterMetaData.ABI

// SessionRouter is an auto generated Go binding around an Ethereum contract.
type SessionRouter struct {
	SessionRouterCaller     // Read-only binding to the contract
	SessionRouterTransactor // Write-only binding to the contract
	SessionRouterFilterer   // Log filterer for contract events
}

// SessionRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type SessionRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SessionRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SessionRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SessionRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SessionRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SessionRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SessionRouterSession struct {
	Contract     *SessionRouter    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SessionRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SessionRouterCallerSession struct {
	Contract *SessionRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SessionRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SessionRouterTransactorSession struct {
	Contract     *SessionRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SessionRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type SessionRouterRaw struct {
	Contract *SessionRouter // Generic contract binding to access the raw methods on
}

// SessionRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SessionRouterCallerRaw struct {
	Contract *SessionRouterCaller // Generic read-only contract binding to access the raw methods on
}

// SessionRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SessionRouterTransactorRaw struct {
	Contract *SessionRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSessionRouter creates a new instance of SessionRouter, bound to a specific deployed contract.
func NewSessionRouter(address common.Address, backend bind.ContractBackend) (*SessionRouter, error) {
	contract, err := bindSessionRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SessionRouter{SessionRouterCaller: SessionRouterCaller{contract: contract}, SessionRouterTransactor: SessionRouterTransactor{contract: contract}, SessionRouterFilterer: SessionRouterFilterer{contract: contract}}, nil
}

// NewSessionRouterCaller creates a new read-only instance of SessionRouter, bound to a specific deployed contract.
func NewSessionRouterCaller(address common.Address, caller bind.ContractCaller) (*SessionRouterCaller, error) {
	contract, err := bindSessionRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SessionRouterCaller{contract: contract}, nil
}

// NewSessionRouterTransactor creates a new write-only instance of SessionRouter, bound to a specific deployed contract.
func NewSessionRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*SessionRouterTransactor, error) {
	contract, err := bindSessionRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SessionRouterTransactor{contract: contract}, nil
}

// NewSessionRouterFilterer creates a new log filterer instance of SessionRouter, bound to a specific deployed contract.
func NewSessionRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*SessionRouterFilterer, error) {
	contract, err := bindSessionRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SessionRouterFilterer{contract: contract}, nil
}

// bindSessionRouter binds a generic wrapper to an already deployed contract.
func bindSessionRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SessionRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SessionRouter *SessionRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SessionRouter.Contract.SessionRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SessionRouter *SessionRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SessionRouter.Contract.SessionRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SessionRouter *SessionRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SessionRouter.Contract.SessionRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SessionRouter *SessionRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SessionRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SessionRouter *SessionRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SessionRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SessionRouter *SessionRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SessionRouter.Contract.contract.Transact(opts, method, params...)
}

// Sessions is a free data retrieval call binding the contract method 0x431a1b97.
//
// Solidity: function sessions(address ) view returns(address dispatcher, address owner, uint32 exp, uint32 scopes)
func (_SessionRouter *SessionRouterCaller) Sessions(opts *bind.CallOpts, arg0 common.Address) (struct {
	Dispatcher common.Address
	Owner      common.Address
	Exp        uint32
	Scopes     uint32
}, error) {
	var out []interface{}
	err := _SessionRouter.contract.Call(opts, &out, "sessions", arg0)

	outstruct := new(struct {
		Dispatcher common.Address
		Owner      common.Address
		Exp        uint32
		Scopes     uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Dispatcher = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Owner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Exp = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.Scopes = *abi.ConvertType(out[3], new(uint32)).(*uint32)

	return *outstruct, err

}

// Sessions is a free data retrieval call binding the contract method 0x431a1b97.
//
// Solidity: function sessions(address ) view returns(address dispatcher, address owner, uint32 exp, uint32 scopes)
func (_SessionRouter *SessionRouterSession) Sessions(arg0 common.Address) (struct {
	Dispatcher common.Address
	Owner      common.Address
	Exp        uint32
	Scopes     uint32
}, error) {
	return _SessionRouter.Contract.Sessions(&_SessionRouter.CallOpts, arg0)
}

// Sessions is a free data retrieval call binding the contract method 0x431a1b97.
//
// Solidity: function sessions(address ) view returns(address dispatcher, address owner, uint32 exp, uint32 scopes)
func (_SessionRouter *SessionRouterCallerSession) Sessions(arg0 common.Address) (struct {
	Dispatcher common.Address
	Owner      common.Address
	Exp        uint32
	Scopes     uint32
}, error) {
	return _SessionRouter.Contract.Sessions(&_SessionRouter.CallOpts, arg0)
}

// AuthorizeAddr is a paid mutator transaction binding the contract method 0x02dd3ed4.
//
// Solidity: function authorizeAddr(address dispatcher, uint32 ttl, uint32 scopes, address sessionAddr) returns()
func (_SessionRouter *SessionRouterTransactor) AuthorizeAddr(opts *bind.TransactOpts, dispatcher common.Address, ttl uint32, scopes uint32, sessionAddr common.Address) (*types.Transaction, error) {
	return _SessionRouter.contract.Transact(opts, "authorizeAddr", dispatcher, ttl, scopes, sessionAddr)
}

// AuthorizeAddr is a paid mutator transaction binding the contract method 0x02dd3ed4.
//
// Solidity: function authorizeAddr(address dispatcher, uint32 ttl, uint32 scopes, address sessionAddr) returns()
func (_SessionRouter *SessionRouterSession) AuthorizeAddr(dispatcher common.Address, ttl uint32, scopes uint32, sessionAddr common.Address) (*types.Transaction, error) {
	return _SessionRouter.Contract.AuthorizeAddr(&_SessionRouter.TransactOpts, dispatcher, ttl, scopes, sessionAddr)
}

// AuthorizeAddr is a paid mutator transaction binding the contract method 0x02dd3ed4.
//
// Solidity: function authorizeAddr(address dispatcher, uint32 ttl, uint32 scopes, address sessionAddr) returns()
func (_SessionRouter *SessionRouterTransactorSession) AuthorizeAddr(dispatcher common.Address, ttl uint32, scopes uint32, sessionAddr common.Address) (*types.Transaction, error) {
	return _SessionRouter.Contract.AuthorizeAddr(&_SessionRouter.TransactOpts, dispatcher, ttl, scopes, sessionAddr)
}

// AuthorizeAddr0 is a paid mutator transaction binding the contract method 0x401870d9.
//
// Solidity: function authorizeAddr(address dispatcher, uint32 ttl, uint32 scopes, address sessionAddr, bytes sig) returns()
func (_SessionRouter *SessionRouterTransactor) AuthorizeAddr0(opts *bind.TransactOpts, dispatcher common.Address, ttl uint32, scopes uint32, sessionAddr common.Address, sig []byte) (*types.Transaction, error) {
	return _SessionRouter.contract.Transact(opts, "authorizeAddr0", dispatcher, ttl, scopes, sessionAddr, sig)
}

// AuthorizeAddr0 is a paid mutator transaction binding the contract method 0x401870d9.
//
// Solidity: function authorizeAddr(address dispatcher, uint32 ttl, uint32 scopes, address sessionAddr, bytes sig) returns()
func (_SessionRouter *SessionRouterSession) AuthorizeAddr0(dispatcher common.Address, ttl uint32, scopes uint32, sessionAddr common.Address, sig []byte) (*types.Transaction, error) {
	return _SessionRouter.Contract.AuthorizeAddr0(&_SessionRouter.TransactOpts, dispatcher, ttl, scopes, sessionAddr, sig)
}

// AuthorizeAddr0 is a paid mutator transaction binding the contract method 0x401870d9.
//
// Solidity: function authorizeAddr(address dispatcher, uint32 ttl, uint32 scopes, address sessionAddr, bytes sig) returns()
func (_SessionRouter *SessionRouterTransactorSession) AuthorizeAddr0(dispatcher common.Address, ttl uint32, scopes uint32, sessionAddr common.Address, sig []byte) (*types.Transaction, error) {
	return _SessionRouter.Contract.AuthorizeAddr0(&_SessionRouter.TransactOpts, dispatcher, ttl, scopes, sessionAddr, sig)
}

// Dispatch is a paid mutator transaction binding the contract method 0xd491289b.
//
// Solidity: function dispatch(bytes[][] actions, bytes[] sig) returns()
func (_SessionRouter *SessionRouterTransactor) Dispatch(opts *bind.TransactOpts, actions [][][]byte, sig [][]byte) (*types.Transaction, error) {
	return _SessionRouter.contract.Transact(opts, "dispatch", actions, sig)
}

// Dispatch is a paid mutator transaction binding the contract method 0xd491289b.
//
// Solidity: function dispatch(bytes[][] actions, bytes[] sig) returns()
func (_SessionRouter *SessionRouterSession) Dispatch(actions [][][]byte, sig [][]byte) (*types.Transaction, error) {
	return _SessionRouter.Contract.Dispatch(&_SessionRouter.TransactOpts, actions, sig)
}

// Dispatch is a paid mutator transaction binding the contract method 0xd491289b.
//
// Solidity: function dispatch(bytes[][] actions, bytes[] sig) returns()
func (_SessionRouter *SessionRouterTransactorSession) Dispatch(actions [][][]byte, sig [][]byte) (*types.Transaction, error) {
	return _SessionRouter.Contract.Dispatch(&_SessionRouter.TransactOpts, actions, sig)
}

// RevokeAddr is a paid mutator transaction binding the contract method 0x3f43ebca.
//
// Solidity: function revokeAddr(address addr, bytes sig) returns()
func (_SessionRouter *SessionRouterTransactor) RevokeAddr(opts *bind.TransactOpts, addr common.Address, sig []byte) (*types.Transaction, error) {
	return _SessionRouter.contract.Transact(opts, "revokeAddr", addr, sig)
}

// RevokeAddr is a paid mutator transaction binding the contract method 0x3f43ebca.
//
// Solidity: function revokeAddr(address addr, bytes sig) returns()
func (_SessionRouter *SessionRouterSession) RevokeAddr(addr common.Address, sig []byte) (*types.Transaction, error) {
	return _SessionRouter.Contract.RevokeAddr(&_SessionRouter.TransactOpts, addr, sig)
}

// RevokeAddr is a paid mutator transaction binding the contract method 0x3f43ebca.
//
// Solidity: function revokeAddr(address addr, bytes sig) returns()
func (_SessionRouter *SessionRouterTransactorSession) RevokeAddr(addr common.Address, sig []byte) (*types.Transaction, error) {
	return _SessionRouter.Contract.RevokeAddr(&_SessionRouter.TransactOpts, addr, sig)
}

// RevokeAddr0 is a paid mutator transaction binding the contract method 0xbb68f5e7.
//
// Solidity: function revokeAddr(address addr) returns()
func (_SessionRouter *SessionRouterTransactor) RevokeAddr0(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _SessionRouter.contract.Transact(opts, "revokeAddr0", addr)
}

// RevokeAddr0 is a paid mutator transaction binding the contract method 0xbb68f5e7.
//
// Solidity: function revokeAddr(address addr) returns()
func (_SessionRouter *SessionRouterSession) RevokeAddr0(addr common.Address) (*types.Transaction, error) {
	return _SessionRouter.Contract.RevokeAddr0(&_SessionRouter.TransactOpts, addr)
}

// RevokeAddr0 is a paid mutator transaction binding the contract method 0xbb68f5e7.
//
// Solidity: function revokeAddr(address addr) returns()
func (_SessionRouter *SessionRouterTransactorSession) RevokeAddr0(addr common.Address) (*types.Transaction, error) {
	return _SessionRouter.Contract.RevokeAddr0(&_SessionRouter.TransactOpts, addr)
}

// SessionRouterSessionCreateIterator is returned from FilterSessionCreate and is used to iterate over the raw logs and unpacked data for SessionCreate events raised by the SessionRouter contract.
type SessionRouterSessionCreateIterator struct {
	Event *SessionRouterSessionCreate // Event containing the contract specifics and raw log

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
func (it *SessionRouterSessionCreateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SessionRouterSessionCreate)
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
		it.Event = new(SessionRouterSessionCreate)
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
func (it *SessionRouterSessionCreateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SessionRouterSessionCreateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SessionRouterSessionCreate represents a SessionCreate event raised by the SessionRouter contract.
type SessionRouterSessionCreate struct {
	Session common.Address
	Owner   common.Address
	Exp     uint32
	Scopes  uint32
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSessionCreate is a free log retrieval operation binding the contract event 0x703437513e11491a538cde074889536a9cb295a07dc2d8de3ec316435737e4c5.
//
// Solidity: event SessionCreate(address session, address owner, uint32 exp, uint32 scopes)
func (_SessionRouter *SessionRouterFilterer) FilterSessionCreate(opts *bind.FilterOpts) (*SessionRouterSessionCreateIterator, error) {

	logs, sub, err := _SessionRouter.contract.FilterLogs(opts, "SessionCreate")
	if err != nil {
		return nil, err
	}
	return &SessionRouterSessionCreateIterator{contract: _SessionRouter.contract, event: "SessionCreate", logs: logs, sub: sub}, nil
}

// WatchSessionCreate is a free log subscription operation binding the contract event 0x703437513e11491a538cde074889536a9cb295a07dc2d8de3ec316435737e4c5.
//
// Solidity: event SessionCreate(address session, address owner, uint32 exp, uint32 scopes)
func (_SessionRouter *SessionRouterFilterer) WatchSessionCreate(opts *bind.WatchOpts, sink chan<- *SessionRouterSessionCreate) (event.Subscription, error) {

	logs, sub, err := _SessionRouter.contract.WatchLogs(opts, "SessionCreate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SessionRouterSessionCreate)
				if err := _SessionRouter.contract.UnpackLog(event, "SessionCreate", log); err != nil {
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

// ParseSessionCreate is a log parse operation binding the contract event 0x703437513e11491a538cde074889536a9cb295a07dc2d8de3ec316435737e4c5.
//
// Solidity: event SessionCreate(address session, address owner, uint32 exp, uint32 scopes)
func (_SessionRouter *SessionRouterFilterer) ParseSessionCreate(log types.Log) (*SessionRouterSessionCreate, error) {
	event := new(SessionRouterSessionCreate)
	if err := _SessionRouter.contract.UnpackLog(event, "SessionCreate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SessionRouterSessionDestroyIterator is returned from FilterSessionDestroy and is used to iterate over the raw logs and unpacked data for SessionDestroy events raised by the SessionRouter contract.
type SessionRouterSessionDestroyIterator struct {
	Event *SessionRouterSessionDestroy // Event containing the contract specifics and raw log

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
func (it *SessionRouterSessionDestroyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SessionRouterSessionDestroy)
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
		it.Event = new(SessionRouterSessionDestroy)
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
func (it *SessionRouterSessionDestroyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SessionRouterSessionDestroyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SessionRouterSessionDestroy represents a SessionDestroy event raised by the SessionRouter contract.
type SessionRouterSessionDestroy struct {
	Session common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSessionDestroy is a free log retrieval operation binding the contract event 0xd2dd19770a4e71b8a94a6f9865e6028f2b4d94ae88504113827e9f4979980a0f.
//
// Solidity: event SessionDestroy(address session)
func (_SessionRouter *SessionRouterFilterer) FilterSessionDestroy(opts *bind.FilterOpts) (*SessionRouterSessionDestroyIterator, error) {

	logs, sub, err := _SessionRouter.contract.FilterLogs(opts, "SessionDestroy")
	if err != nil {
		return nil, err
	}
	return &SessionRouterSessionDestroyIterator{contract: _SessionRouter.contract, event: "SessionDestroy", logs: logs, sub: sub}, nil
}

// WatchSessionDestroy is a free log subscription operation binding the contract event 0xd2dd19770a4e71b8a94a6f9865e6028f2b4d94ae88504113827e9f4979980a0f.
//
// Solidity: event SessionDestroy(address session)
func (_SessionRouter *SessionRouterFilterer) WatchSessionDestroy(opts *bind.WatchOpts, sink chan<- *SessionRouterSessionDestroy) (event.Subscription, error) {

	logs, sub, err := _SessionRouter.contract.WatchLogs(opts, "SessionDestroy")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SessionRouterSessionDestroy)
				if err := _SessionRouter.contract.UnpackLog(event, "SessionDestroy", log); err != nil {
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

// ParseSessionDestroy is a log parse operation binding the contract event 0xd2dd19770a4e71b8a94a6f9865e6028f2b4d94ae88504113827e9f4979980a0f.
//
// Solidity: event SessionDestroy(address session)
func (_SessionRouter *SessionRouterFilterer) ParseSessionDestroy(log types.Log) (*SessionRouterSessionDestroy, error) {
	event := new(SessionRouterSessionDestroy)
	if err := _SessionRouter.contract.UnpackLog(event, "SessionDestroy", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
