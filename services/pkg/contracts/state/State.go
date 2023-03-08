// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package state

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

// StateMetaData contains all meta data concerning the State contract.
var StateMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes24\",\"name\":\"id\",\"type\":\"bytes24\"},{\"indexed\":false,\"internalType\":\"enumAnnotationKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"ref\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"AnnotationSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"relID\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"relKey\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes24\",\"name\":\"srcNodeID\",\"type\":\"bytes24\"}],\"name\":\"EdgeRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"relID\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"relKey\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes24\",\"name\":\"srcNodeID\",\"type\":\"bytes24\"},{\"indexed\":false,\"internalType\":\"bytes24\",\"name\":\"dstNodeID\",\"type\":\"bytes24\"},{\"indexed\":false,\"internalType\":\"uint160\",\"name\":\"weight\",\"type\":\"uint160\"}],\"name\":\"EdgeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"enumWeightKind\",\"name\":\"kind\",\"type\":\"uint8\"}],\"name\":\"EdgeTypeRegister\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"enumCompoundKeyKind\",\"name\":\"keyKind\",\"type\":\"uint8\"}],\"name\":\"NodeTypeRegister\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes24\",\"name\":\"nodeID\",\"type\":\"bytes24\"},{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"annotationData\",\"type\":\"string\"}],\"name\":\"annotate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"authorizeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"relID\",\"type\":\"bytes4\"},{\"internalType\":\"uint8\",\"name\":\"relKey\",\"type\":\"uint8\"},{\"internalType\":\"bytes24\",\"name\":\"srcNodeID\",\"type\":\"bytes24\"}],\"name\":\"get\",\"outputs\":[{\"internalType\":\"bytes24\",\"name\":\"dstNodeId\",\"type\":\"bytes24\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"relID\",\"type\":\"bytes4\"},{\"internalType\":\"string\",\"name\":\"relName\",\"type\":\"string\"},{\"internalType\":\"enumWeightKind\",\"name\":\"weightKind\",\"type\":\"uint8\"}],\"name\":\"registerEdgeType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"kindID\",\"type\":\"bytes4\"},{\"internalType\":\"string\",\"name\":\"kindName\",\"type\":\"string\"},{\"internalType\":\"enumCompoundKeyKind\",\"name\":\"keyKind\",\"type\":\"uint8\"}],\"name\":\"registerNodeType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"relID\",\"type\":\"bytes4\"},{\"internalType\":\"uint8\",\"name\":\"relKey\",\"type\":\"uint8\"},{\"internalType\":\"bytes24\",\"name\":\"srcNodeID\",\"type\":\"bytes24\"}],\"name\":\"remove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"relID\",\"type\":\"bytes4\"},{\"internalType\":\"uint8\",\"name\":\"relKey\",\"type\":\"uint8\"},{\"internalType\":\"bytes24\",\"name\":\"srcNodeID\",\"type\":\"bytes24\"},{\"internalType\":\"bytes24\",\"name\":\"dstNodeID\",\"type\":\"bytes24\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StateABI is the input ABI used to generate the binding from.
// Deprecated: Use StateMetaData.ABI instead.
var StateABI = StateMetaData.ABI

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
	parsed, err := StateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_State *StateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_State *StateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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

// Get is a free data retrieval call binding the contract method 0x0bf24542.
//
// Solidity: function get(bytes4 relID, uint8 relKey, bytes24 srcNodeID) view returns(bytes24 dstNodeId, uint64 weight)
func (_State *StateCaller) Get(opts *bind.CallOpts, relID [4]byte, relKey uint8, srcNodeID [24]byte) (struct {
	DstNodeId [24]byte
	Weight    uint64
}, error) {
	var out []interface{}
	err := _State.contract.Call(opts, &out, "get", relID, relKey, srcNodeID)

	outstruct := new(struct {
		DstNodeId [24]byte
		Weight    uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DstNodeId = *abi.ConvertType(out[0], new([24]byte)).(*[24]byte)
	outstruct.Weight = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// Get is a free data retrieval call binding the contract method 0x0bf24542.
//
// Solidity: function get(bytes4 relID, uint8 relKey, bytes24 srcNodeID) view returns(bytes24 dstNodeId, uint64 weight)
func (_State *StateSession) Get(relID [4]byte, relKey uint8, srcNodeID [24]byte) (struct {
	DstNodeId [24]byte
	Weight    uint64
}, error) {
	return _State.Contract.Get(&_State.CallOpts, relID, relKey, srcNodeID)
}

// Get is a free data retrieval call binding the contract method 0x0bf24542.
//
// Solidity: function get(bytes4 relID, uint8 relKey, bytes24 srcNodeID) view returns(bytes24 dstNodeId, uint64 weight)
func (_State *StateCallerSession) Get(relID [4]byte, relKey uint8, srcNodeID [24]byte) (struct {
	DstNodeId [24]byte
	Weight    uint64
}, error) {
	return _State.Contract.Get(&_State.CallOpts, relID, relKey, srcNodeID)
}

// Annotate is a paid mutator transaction binding the contract method 0xff271d48.
//
// Solidity: function annotate(bytes24 nodeID, string label, string annotationData) returns()
func (_State *StateTransactor) Annotate(opts *bind.TransactOpts, nodeID [24]byte, label string, annotationData string) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "annotate", nodeID, label, annotationData)
}

// Annotate is a paid mutator transaction binding the contract method 0xff271d48.
//
// Solidity: function annotate(bytes24 nodeID, string label, string annotationData) returns()
func (_State *StateSession) Annotate(nodeID [24]byte, label string, annotationData string) (*types.Transaction, error) {
	return _State.Contract.Annotate(&_State.TransactOpts, nodeID, label, annotationData)
}

// Annotate is a paid mutator transaction binding the contract method 0xff271d48.
//
// Solidity: function annotate(bytes24 nodeID, string label, string annotationData) returns()
func (_State *StateTransactorSession) Annotate(nodeID [24]byte, label string, annotationData string) (*types.Transaction, error) {
	return _State.Contract.Annotate(&_State.TransactOpts, nodeID, label, annotationData)
}

// AuthorizeContract is a paid mutator transaction binding the contract method 0x67561d93.
//
// Solidity: function authorizeContract(address addr) returns()
func (_State *StateTransactor) AuthorizeContract(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "authorizeContract", addr)
}

// AuthorizeContract is a paid mutator transaction binding the contract method 0x67561d93.
//
// Solidity: function authorizeContract(address addr) returns()
func (_State *StateSession) AuthorizeContract(addr common.Address) (*types.Transaction, error) {
	return _State.Contract.AuthorizeContract(&_State.TransactOpts, addr)
}

// AuthorizeContract is a paid mutator transaction binding the contract method 0x67561d93.
//
// Solidity: function authorizeContract(address addr) returns()
func (_State *StateTransactorSession) AuthorizeContract(addr common.Address) (*types.Transaction, error) {
	return _State.Contract.AuthorizeContract(&_State.TransactOpts, addr)
}

// RegisterEdgeType is a paid mutator transaction binding the contract method 0x27d9e1aa.
//
// Solidity: function registerEdgeType(bytes4 relID, string relName, uint8 weightKind) returns()
func (_State *StateTransactor) RegisterEdgeType(opts *bind.TransactOpts, relID [4]byte, relName string, weightKind uint8) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "registerEdgeType", relID, relName, weightKind)
}

// RegisterEdgeType is a paid mutator transaction binding the contract method 0x27d9e1aa.
//
// Solidity: function registerEdgeType(bytes4 relID, string relName, uint8 weightKind) returns()
func (_State *StateSession) RegisterEdgeType(relID [4]byte, relName string, weightKind uint8) (*types.Transaction, error) {
	return _State.Contract.RegisterEdgeType(&_State.TransactOpts, relID, relName, weightKind)
}

// RegisterEdgeType is a paid mutator transaction binding the contract method 0x27d9e1aa.
//
// Solidity: function registerEdgeType(bytes4 relID, string relName, uint8 weightKind) returns()
func (_State *StateTransactorSession) RegisterEdgeType(relID [4]byte, relName string, weightKind uint8) (*types.Transaction, error) {
	return _State.Contract.RegisterEdgeType(&_State.TransactOpts, relID, relName, weightKind)
}

// RegisterNodeType is a paid mutator transaction binding the contract method 0x72efd7ac.
//
// Solidity: function registerNodeType(bytes4 kindID, string kindName, uint8 keyKind) returns()
func (_State *StateTransactor) RegisterNodeType(opts *bind.TransactOpts, kindID [4]byte, kindName string, keyKind uint8) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "registerNodeType", kindID, kindName, keyKind)
}

// RegisterNodeType is a paid mutator transaction binding the contract method 0x72efd7ac.
//
// Solidity: function registerNodeType(bytes4 kindID, string kindName, uint8 keyKind) returns()
func (_State *StateSession) RegisterNodeType(kindID [4]byte, kindName string, keyKind uint8) (*types.Transaction, error) {
	return _State.Contract.RegisterNodeType(&_State.TransactOpts, kindID, kindName, keyKind)
}

// RegisterNodeType is a paid mutator transaction binding the contract method 0x72efd7ac.
//
// Solidity: function registerNodeType(bytes4 kindID, string kindName, uint8 keyKind) returns()
func (_State *StateTransactorSession) RegisterNodeType(kindID [4]byte, kindName string, keyKind uint8) (*types.Transaction, error) {
	return _State.Contract.RegisterNodeType(&_State.TransactOpts, kindID, kindName, keyKind)
}

// Remove is a paid mutator transaction binding the contract method 0x8c7a9e38.
//
// Solidity: function remove(bytes4 relID, uint8 relKey, bytes24 srcNodeID) returns()
func (_State *StateTransactor) Remove(opts *bind.TransactOpts, relID [4]byte, relKey uint8, srcNodeID [24]byte) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "remove", relID, relKey, srcNodeID)
}

// Remove is a paid mutator transaction binding the contract method 0x8c7a9e38.
//
// Solidity: function remove(bytes4 relID, uint8 relKey, bytes24 srcNodeID) returns()
func (_State *StateSession) Remove(relID [4]byte, relKey uint8, srcNodeID [24]byte) (*types.Transaction, error) {
	return _State.Contract.Remove(&_State.TransactOpts, relID, relKey, srcNodeID)
}

// Remove is a paid mutator transaction binding the contract method 0x8c7a9e38.
//
// Solidity: function remove(bytes4 relID, uint8 relKey, bytes24 srcNodeID) returns()
func (_State *StateTransactorSession) Remove(relID [4]byte, relKey uint8, srcNodeID [24]byte) (*types.Transaction, error) {
	return _State.Contract.Remove(&_State.TransactOpts, relID, relKey, srcNodeID)
}

// Set is a paid mutator transaction binding the contract method 0xf4602114.
//
// Solidity: function set(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint64 weight) returns()
func (_State *StateTransactor) Set(opts *bind.TransactOpts, relID [4]byte, relKey uint8, srcNodeID [24]byte, dstNodeID [24]byte, weight uint64) (*types.Transaction, error) {
	return _State.contract.Transact(opts, "set", relID, relKey, srcNodeID, dstNodeID, weight)
}

// Set is a paid mutator transaction binding the contract method 0xf4602114.
//
// Solidity: function set(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint64 weight) returns()
func (_State *StateSession) Set(relID [4]byte, relKey uint8, srcNodeID [24]byte, dstNodeID [24]byte, weight uint64) (*types.Transaction, error) {
	return _State.Contract.Set(&_State.TransactOpts, relID, relKey, srcNodeID, dstNodeID, weight)
}

// Set is a paid mutator transaction binding the contract method 0xf4602114.
//
// Solidity: function set(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint64 weight) returns()
func (_State *StateTransactorSession) Set(relID [4]byte, relKey uint8, srcNodeID [24]byte, dstNodeID [24]byte, weight uint64) (*types.Transaction, error) {
	return _State.Contract.Set(&_State.TransactOpts, relID, relKey, srcNodeID, dstNodeID, weight)
}

// StateAnnotationSetIterator is returned from FilterAnnotationSet and is used to iterate over the raw logs and unpacked data for AnnotationSet events raised by the State contract.
type StateAnnotationSetIterator struct {
	Event *StateAnnotationSet // Event containing the contract specifics and raw log

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
func (it *StateAnnotationSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateAnnotationSet)
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
		it.Event = new(StateAnnotationSet)
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
func (it *StateAnnotationSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateAnnotationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateAnnotationSet represents a AnnotationSet event raised by the State contract.
type StateAnnotationSet struct {
	Id    [24]byte
	Kind  uint8
	Label string
	Ref   [32]byte
	Data  string
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAnnotationSet is a free log retrieval operation binding the contract event 0x79838124986a3a7ac303823ad382114a566cc5168f73baf207b82beab566d4af.
//
// Solidity: event AnnotationSet(bytes24 id, uint8 kind, string label, bytes32 ref, string data)
func (_State *StateFilterer) FilterAnnotationSet(opts *bind.FilterOpts) (*StateAnnotationSetIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "AnnotationSet")
	if err != nil {
		return nil, err
	}
	return &StateAnnotationSetIterator{contract: _State.contract, event: "AnnotationSet", logs: logs, sub: sub}, nil
}

// WatchAnnotationSet is a free log subscription operation binding the contract event 0x79838124986a3a7ac303823ad382114a566cc5168f73baf207b82beab566d4af.
//
// Solidity: event AnnotationSet(bytes24 id, uint8 kind, string label, bytes32 ref, string data)
func (_State *StateFilterer) WatchAnnotationSet(opts *bind.WatchOpts, sink chan<- *StateAnnotationSet) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "AnnotationSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateAnnotationSet)
				if err := _State.contract.UnpackLog(event, "AnnotationSet", log); err != nil {
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

// ParseAnnotationSet is a log parse operation binding the contract event 0x79838124986a3a7ac303823ad382114a566cc5168f73baf207b82beab566d4af.
//
// Solidity: event AnnotationSet(bytes24 id, uint8 kind, string label, bytes32 ref, string data)
func (_State *StateFilterer) ParseAnnotationSet(log types.Log) (*StateAnnotationSet, error) {
	event := new(StateAnnotationSet)
	if err := _State.contract.UnpackLog(event, "AnnotationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateEdgeRemoveIterator is returned from FilterEdgeRemove and is used to iterate over the raw logs and unpacked data for EdgeRemove events raised by the State contract.
type StateEdgeRemoveIterator struct {
	Event *StateEdgeRemove // Event containing the contract specifics and raw log

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
func (it *StateEdgeRemoveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateEdgeRemove)
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
		it.Event = new(StateEdgeRemove)
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
func (it *StateEdgeRemoveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateEdgeRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateEdgeRemove represents a EdgeRemove event raised by the State contract.
type StateEdgeRemove struct {
	RelID     [4]byte
	RelKey    uint8
	SrcNodeID [24]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEdgeRemove is a free log retrieval operation binding the contract event 0x1e0b44403284c5f71ea7241d483d0e3f7c194af79a5bed65e25c2012dd229c22.
//
// Solidity: event EdgeRemove(bytes4 relID, uint8 relKey, bytes24 srcNodeID)
func (_State *StateFilterer) FilterEdgeRemove(opts *bind.FilterOpts) (*StateEdgeRemoveIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "EdgeRemove")
	if err != nil {
		return nil, err
	}
	return &StateEdgeRemoveIterator{contract: _State.contract, event: "EdgeRemove", logs: logs, sub: sub}, nil
}

// WatchEdgeRemove is a free log subscription operation binding the contract event 0x1e0b44403284c5f71ea7241d483d0e3f7c194af79a5bed65e25c2012dd229c22.
//
// Solidity: event EdgeRemove(bytes4 relID, uint8 relKey, bytes24 srcNodeID)
func (_State *StateFilterer) WatchEdgeRemove(opts *bind.WatchOpts, sink chan<- *StateEdgeRemove) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "EdgeRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateEdgeRemove)
				if err := _State.contract.UnpackLog(event, "EdgeRemove", log); err != nil {
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

// ParseEdgeRemove is a log parse operation binding the contract event 0x1e0b44403284c5f71ea7241d483d0e3f7c194af79a5bed65e25c2012dd229c22.
//
// Solidity: event EdgeRemove(bytes4 relID, uint8 relKey, bytes24 srcNodeID)
func (_State *StateFilterer) ParseEdgeRemove(log types.Log) (*StateEdgeRemove, error) {
	event := new(StateEdgeRemove)
	if err := _State.contract.UnpackLog(event, "EdgeRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateEdgeSetIterator is returned from FilterEdgeSet and is used to iterate over the raw logs and unpacked data for EdgeSet events raised by the State contract.
type StateEdgeSetIterator struct {
	Event *StateEdgeSet // Event containing the contract specifics and raw log

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
func (it *StateEdgeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateEdgeSet)
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
		it.Event = new(StateEdgeSet)
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
func (it *StateEdgeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateEdgeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateEdgeSet represents a EdgeSet event raised by the State contract.
type StateEdgeSet struct {
	RelID     [4]byte
	RelKey    uint8
	SrcNodeID [24]byte
	DstNodeID [24]byte
	Weight    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEdgeSet is a free log retrieval operation binding the contract event 0xa54dd4022502bfa0d2b33a2acc90999866ede17ec6ca83c8124bbfdc5eac8651.
//
// Solidity: event EdgeSet(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint160 weight)
func (_State *StateFilterer) FilterEdgeSet(opts *bind.FilterOpts) (*StateEdgeSetIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "EdgeSet")
	if err != nil {
		return nil, err
	}
	return &StateEdgeSetIterator{contract: _State.contract, event: "EdgeSet", logs: logs, sub: sub}, nil
}

// WatchEdgeSet is a free log subscription operation binding the contract event 0xa54dd4022502bfa0d2b33a2acc90999866ede17ec6ca83c8124bbfdc5eac8651.
//
// Solidity: event EdgeSet(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint160 weight)
func (_State *StateFilterer) WatchEdgeSet(opts *bind.WatchOpts, sink chan<- *StateEdgeSet) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "EdgeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateEdgeSet)
				if err := _State.contract.UnpackLog(event, "EdgeSet", log); err != nil {
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

// ParseEdgeSet is a log parse operation binding the contract event 0xa54dd4022502bfa0d2b33a2acc90999866ede17ec6ca83c8124bbfdc5eac8651.
//
// Solidity: event EdgeSet(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint160 weight)
func (_State *StateFilterer) ParseEdgeSet(log types.Log) (*StateEdgeSet, error) {
	event := new(StateEdgeSet)
	if err := _State.contract.UnpackLog(event, "EdgeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateEdgeTypeRegisterIterator is returned from FilterEdgeTypeRegister and is used to iterate over the raw logs and unpacked data for EdgeTypeRegister events raised by the State contract.
type StateEdgeTypeRegisterIterator struct {
	Event *StateEdgeTypeRegister // Event containing the contract specifics and raw log

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
func (it *StateEdgeTypeRegisterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateEdgeTypeRegister)
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
		it.Event = new(StateEdgeTypeRegister)
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
func (it *StateEdgeTypeRegisterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateEdgeTypeRegisterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateEdgeTypeRegister represents a EdgeTypeRegister event raised by the State contract.
type StateEdgeTypeRegister struct {
	Id   [4]byte
	Name string
	Kind uint8
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterEdgeTypeRegister is a free log retrieval operation binding the contract event 0xdffeac8fd30891d16d4fd6d375e13822f93ea2022cd53fff5807c05379f1426b.
//
// Solidity: event EdgeTypeRegister(bytes4 id, string name, uint8 kind)
func (_State *StateFilterer) FilterEdgeTypeRegister(opts *bind.FilterOpts) (*StateEdgeTypeRegisterIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "EdgeTypeRegister")
	if err != nil {
		return nil, err
	}
	return &StateEdgeTypeRegisterIterator{contract: _State.contract, event: "EdgeTypeRegister", logs: logs, sub: sub}, nil
}

// WatchEdgeTypeRegister is a free log subscription operation binding the contract event 0xdffeac8fd30891d16d4fd6d375e13822f93ea2022cd53fff5807c05379f1426b.
//
// Solidity: event EdgeTypeRegister(bytes4 id, string name, uint8 kind)
func (_State *StateFilterer) WatchEdgeTypeRegister(opts *bind.WatchOpts, sink chan<- *StateEdgeTypeRegister) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "EdgeTypeRegister")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateEdgeTypeRegister)
				if err := _State.contract.UnpackLog(event, "EdgeTypeRegister", log); err != nil {
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

// ParseEdgeTypeRegister is a log parse operation binding the contract event 0xdffeac8fd30891d16d4fd6d375e13822f93ea2022cd53fff5807c05379f1426b.
//
// Solidity: event EdgeTypeRegister(bytes4 id, string name, uint8 kind)
func (_State *StateFilterer) ParseEdgeTypeRegister(log types.Log) (*StateEdgeTypeRegister, error) {
	event := new(StateEdgeTypeRegister)
	if err := _State.contract.UnpackLog(event, "EdgeTypeRegister", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StateNodeTypeRegisterIterator is returned from FilterNodeTypeRegister and is used to iterate over the raw logs and unpacked data for NodeTypeRegister events raised by the State contract.
type StateNodeTypeRegisterIterator struct {
	Event *StateNodeTypeRegister // Event containing the contract specifics and raw log

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
func (it *StateNodeTypeRegisterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StateNodeTypeRegister)
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
		it.Event = new(StateNodeTypeRegister)
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
func (it *StateNodeTypeRegisterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StateNodeTypeRegisterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StateNodeTypeRegister represents a NodeTypeRegister event raised by the State contract.
type StateNodeTypeRegister struct {
	Id      [4]byte
	Name    string
	KeyKind uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNodeTypeRegister is a free log retrieval operation binding the contract event 0xbed80c5919aa762c9dc404f9c40660a62258ca74aa60222d9e107b4c07dd5977.
//
// Solidity: event NodeTypeRegister(bytes4 id, string name, uint8 keyKind)
func (_State *StateFilterer) FilterNodeTypeRegister(opts *bind.FilterOpts) (*StateNodeTypeRegisterIterator, error) {

	logs, sub, err := _State.contract.FilterLogs(opts, "NodeTypeRegister")
	if err != nil {
		return nil, err
	}
	return &StateNodeTypeRegisterIterator{contract: _State.contract, event: "NodeTypeRegister", logs: logs, sub: sub}, nil
}

// WatchNodeTypeRegister is a free log subscription operation binding the contract event 0xbed80c5919aa762c9dc404f9c40660a62258ca74aa60222d9e107b4c07dd5977.
//
// Solidity: event NodeTypeRegister(bytes4 id, string name, uint8 keyKind)
func (_State *StateFilterer) WatchNodeTypeRegister(opts *bind.WatchOpts, sink chan<- *StateNodeTypeRegister) (event.Subscription, error) {

	logs, sub, err := _State.contract.WatchLogs(opts, "NodeTypeRegister")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StateNodeTypeRegister)
				if err := _State.contract.UnpackLog(event, "NodeTypeRegister", log); err != nil {
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

// ParseNodeTypeRegister is a log parse operation binding the contract event 0xbed80c5919aa762c9dc404f9c40660a62258ca74aa60222d9e107b4c07dd5977.
//
// Solidity: event NodeTypeRegister(bytes4 id, string name, uint8 keyKind)
func (_State *StateFilterer) ParseNodeTypeRegister(log types.Log) (*StateNodeTypeRegister, error) {
	event := new(StateNodeTypeRegister)
	if err := _State.contract.UnpackLog(event, "NodeTypeRegister", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
