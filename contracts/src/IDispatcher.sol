// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

// const SCOPE_READ_SENSITIVE = 0x1;
// uint32 SCOPE_WRITE_SENSITIVE = 0x2;
uint32 constant SCOPE_FULL_ACCESS = 0xffff;

// Context is metadata that is submitted with the action
// it contains:
// - the sender (address of player who initiated action)
// - scopes (a uint32 intended to be used as role/auth bitwise flags)
struct Context {
    address sender; // action sender
    uint32 scopes; // authorized scopes
    uint32 clock; // block at time of action commit
}

// Dispatchers accept Actions and execute Rules to modify State
interface Dispatcher {
    event ActionRegistered(address id, string name);
    event ActionDispatched(address indexed sender, bytes32 actionNonce);

    // dispatch(action, session) applies action with the supplied session
    // to the rules.
    // - action is the abi encoded action + arguments
    // - session is contains data about the action sender
    // session data should be considered untrusted and implementations MUST
    // verify the session data or the sender before executing Rules.
    function dispatch(bytes calldata action, Context calldata ctx) external;
    function dispatch(bytes[] calldata actions, Context calldata ctx) external;

    // same as dispatch above, but ctx is built from msg.sender
    function dispatch(bytes calldata action) external;
    function dispatch(bytes[] calldata actions) external;
}
