// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package router

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Zispatch is a free data retrieval call binding the contract method 0x5a1087dc.
//
// Solidity: function zispatch(bytes[] actions, bytes sig) view returns((uint8,bytes4,uint8,bytes24,bytes24,uint160,string,string)[])
func (_SessionRouter *SessionRouterCaller) Zispatch(opts *bind.CallOpts, actions [][]byte, sig []byte) ([]Op, error) {
	var out []interface{}
	err := _SessionRouter.contract.Call(opts, &out, "dispatch", actions, sig)

	if err != nil {
		return *new([]Op), err
	}

	out0 := *abi.ConvertType(out[0], new([]Op)).(*[]Op)

	return out0, err

}

// Zispatch is a free data retrieval call binding the contract method 0x5a1087dc.
//
// Solidity: function zispatch(bytes[] actions, bytes sig) view returns((uint8,bytes4,uint8,bytes24,bytes24,uint160,string,string)[])
func (_SessionRouter *SessionRouterSession) Zispatch(actions [][]byte, sig []byte) ([]Op, error) {
	return _SessionRouter.Contract.Zispatch(&_SessionRouter.CallOpts, actions, sig)
}

// Zispatch is a free data retrieval call binding the contract method 0x5a1087dc.
//
// Solidity: function zispatch(bytes[] actions, bytes sig) view returns((uint8,bytes4,uint8,bytes24,bytes24,uint160,string,string)[])
func (_SessionRouter *SessionRouterCallerSession) Zispatch(actions [][]byte, sig []byte) ([]Op, error) {
	return _SessionRouter.Contract.Zispatch(&_SessionRouter.CallOpts, actions, sig)
}
