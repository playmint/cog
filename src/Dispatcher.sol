// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {State} from "./State.sol";

enum ActionArgKind {
    BOOL,
    INT8,
    INT16,
    INT32,
    INT64,
    INT128,
    INT256,
    INT,
    UINT8,
    UINT16,
    UINT32,
    UINT64,
    UINT128,
    UINT256,
    BYTES,
    STRING,
    ADDRESS,
    NODEID,
    ENUM
}

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
    function dispatch(bytes[] calldata action, Context calldata ctx) external;

    // same as dispatch above, but ctx is built from msg.sender
    function dispatch(bytes calldata actions) external;
    function dispatch(bytes[] calldata actions) external;
}

// Routers accept "signed" Actions and forwards them to Dispatcher.dispatch
// They might be a seperate contract or an extension of the Dispatcher
interface Router {
    function dispatch(bytes[][] calldata actions, bytes[] calldata sig) external;

    function authorizeAddr(Dispatcher dispatcher, uint32 ttl, uint32 scopes, address addr) external;

    function authorizeAddr(Dispatcher dispatcher, uint32 ttl, uint32 scopes, address addr, bytes calldata sig)
        external;

    function revokeAddr(address addr) external;

    function revokeAddr(address addr, bytes calldata sig) external;
}

interface Rule {
    // reduce checks if the given action is relevent to this rule
    // and applies any state changes.
    // reduce functions should be idempotent and side-effect free
    // ideally they would be pure functions, but it is impractical
    // for solidity/gas/storage reasons to implement state in this kind of way
    // so instead we ask that you think of it as pure even if it's not
    function reduce(State state, bytes calldata action, Context calldata ctx) external returns (State);
}

error DispatchUntrustedSender();

// BaseDispatcher implements some basic structure around registering ActionTypes
// and Rules and executing those rules in the defined order against a given State
// implementation.
//
// To use it, inherit from BaseDispatcher and then override `dispatch()` to add
// any application specific validation/authorization for who can dispatch Actions
//
// TODO:
// * need way to remove, reorder or clear the rulesets
contract BaseDispatcher is Dispatcher {
    mapping(Router => bool) private trustedRouters;
    mapping(string => address) private actionAddrs;
    Rule[] private rules;
    State private state;

    constructor() {
        // allow calling ourself
        registerRouter(Router(address(this)));
    }

    // TODO: this should be an owneronly func
    function registerRule(Rule rule) public {
        _registerRule(rule);
    }

    function _registerRule(Rule rule) internal {
        rules.push() = rule;
    }

    // registerRouter(r) will implicitly trust the Context data submitted
    // by r to dispatch(action, ctx) calls.
    //
    // this is useful if there is an external contract managing authN, authZ
    // or when using a "session key" pattern like the SessionRouter.
    //
    // TODO: this should be an owneronly func
    function registerRouter(Router r) public {
        _registerRouter(r);
    }

    function _registerRouter(Router r) internal {
        trustedRouters[r] = true;
    }

    // TODO: this should be an owneronly func
    function registerState(State s) public {
        _registerState(s);
    }

    function _registerState(State s) internal {
        state = s;
    }

    function isRegisteredRouter(address r) internal view returns (bool) {
        return trustedRouters[Router(r)];
    }

    function dispatch(bytes calldata action, Context calldata ctx) public {
        // check ctx can be trusted
        // we trust ctx built from ourself see the dispatch(action) function above that builds a full-access session for the msg.sender
        // we trust ctx built from any registered routers
        if (!isRegisteredRouter(msg.sender)) {
            revert DispatchUntrustedSender();
        }
        for (uint256 i = 0; i < rules.length; i++) {
            rules[i].reduce(state, action, ctx);
        }
        emit ActionDispatched(
            address(ctx.sender),
            "<nonce>" // TODO: unique ids, nonces, and replay protection
        );
    }

    // dispatch from router trusted context
    function dispatch(bytes[] calldata actions, Context calldata ctx) public {
        for (uint256 i = 0; i < actions.length; i++) {
            dispatch(actions[i], ctx);
        }
    }

    function dispatch(bytes calldata action) public {
        Context memory ctx = Context({sender: msg.sender, scopes: SCOPE_FULL_ACCESS, clock: uint32(block.number)});
        this.dispatch(action, ctx);
    }

    function dispatch(bytes[] calldata actions) public {
        for (uint256 i = 0; i < actions.length; i++) {
            dispatch(actions[i]);
        }
    }
}
