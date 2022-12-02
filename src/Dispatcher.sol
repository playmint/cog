// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { State } from "./State.sol";

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

struct Action {
    address owner; // who sent this action
    address id;    // address of ActionType
    bytes args;  // encoded payload, ActionType can decode
    uint256 block;  // when did this action arrive
}

// ActionArgDef describes the parameters for an action.
// The kind references an abi.encode type that is expected
// to be used to encode the args.
struct ActionArgDef {
    string name;
    bool required;
    ActionArgKind kind;
}

// ActionTypeDef describes the parameters this action expects
// and the kinds that it expects them to be abi.encoded to
// action names must be unique within the scope of the Dispatcher
struct ActionTypeDef {
    string name;
    address id;
    ActionArgDef arg0; // TODO: support dynamic number of args
    ActionArgDef arg1; //       instead of this fixed list of 4
    ActionArgDef arg2;
    ActionArgDef arg3;
}

interface ActionType {
    function getTypeDef() external view returns (ActionTypeDef memory);
}

interface Dispatcher {
    event ActionRegistered(
        address id,
        string name
    );
    event ActionDispatched(
        address indexed accountID,
        address indexed actionID,
        bytes32 actionNonce
    );

    function dispatch(Action memory action) external;
    function getActionTypeDefs() external returns (ActionTypeDef[] memory);
    function getActionID(string memory) external returns (address);
}

interface Rule {
    // reduce checks if the given action is relevent to this rule
    // and applies any state changes.
    // reduce functions should be idempotent and side-effect free
    // ideally they would be pure functions, but it is impractical
    // for solidity/gas/storage reasons to implement state in this kind of way
    // so instead we ask that you think of it as pure even if it's not
    function reduce(State state, Action memory action) external returns (State);

    // getActionTypeDefs returns metadata about the actions this rule uses
    // it can be used for introspection by clients to discover available actions
    function getActionTypeDefs() external view returns (ActionTypeDef[] memory);

    // getNodeTypeDefs returns metadata about the node types this rule uses
    // it can be used for introspection of nodetype and edgetype ids
    // function getNodeTypeDefs() external view returns (NodeTypeDef[] memory);
}

error ActionNameConflict();

// BaseDispatcher implements some basic structure around registering ActionTypes
// and Rules and executing those rules in the defined order against a given State
// implementation.
//
// To use it, inherit from BaseDispatcher and then override `dispatch()` to add
// any application specific validation/authorization for who can dispatch Actions
//
// TODO:
// * need way to remove, reorder or clear the rulesets
abstract contract BaseDispatcher is Dispatcher {
    mapping(string => address) private actionAddrs;
    ActionTypeDef[] private actionDefs;
    Rule[] private rules;
    State private gameState;

    constructor(State s) {
        gameState = s;
    }

    function registerRule(Rule rule) internal {
        rules.push() = rule;
        // register all the actions this rule uses
        ActionTypeDef[] memory defs = rule.getActionTypeDefs();
        for (uint i=0; i<defs.length; i++) {
            registerAction(defs[i]);
        }
    }

    function registerAction(ActionTypeDef memory def) internal {
        address id = actionAddrs[def.name];
        if (id == def.id) { // already set
            return;
        } else if (id != address(0)) { // two actions same name
            revert ActionNameConflict();
        }
        actionAddrs[def.name] = def.id;
        actionDefs.push() = def;
        emit ActionRegistered(def.id, def.name);
    }

    function getActionTypeDefs() public view returns (ActionTypeDef[] memory defs) {
        return actionDefs;
    }

    function getActionID(string memory name) public view returns (address) {
        return actionAddrs[name];
    }

    function reduce(State state, Action memory action) internal returns (State) {
        for (uint i=0; i<rules.length; i++) {
            state = rules[i].reduce(state, action);
        }
        return state;
    }

    // You there, implement this:
    // function dispatch(Action memory action) public {
    //     _dispatch(action);
    // }

    function _dispatch(Action memory action) internal {
        gameState = reduce(gameState, action);
        emit ActionDispatched(
            address(action.owner),
            action.id,
            "<nonce>" // TODO: unique ids, nonces, and replay protection
        );
    }

}
