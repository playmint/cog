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
    address id;    // address of ActionType
    bytes args;  // encoded payload, ActionType can decode
}

// const SCOPE_READ_SENSITIVE = 0x1;
// uint32 SCOPE_WRITE_SENSITIVE = 0x2;
uint32 constant SCOPE_FULL_ACCESS = 0xffff;

struct Context {
    address sender; // action sender
    uint32 scopes; // authorized scopes
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

// Dispatchers accept Actions and execute Rules to modify State
interface Dispatcher {
    event ActionRegistered(
        address id,
        string name
    );
    event ActionDispatched(
        address indexed sender,
        bytes32 actionNonce
    );

    // dispatch(action, session) applies action with the supplied session
    // to the rules.
    // - action is the abi encoded action + arguments
    // - session is contains data about the action sender
    // session data should be considered untrusted and implementations MUST
    // verify the session data or the sender before executing Rules.
    function dispatch(
        bytes calldata action,
        Context calldata ctx
    ) external;

    // same as dispatch above, but ctx is built from msg.sender
    function dispatch(
        bytes calldata action
    ) external;

    function getActionTypeDefs() external returns (ActionTypeDef[] memory);

    function getActionID(string memory) external returns (address);
}

// Routers accept "signed" Actions,
interface Router {
    function dispatch(
        bytes calldata action,
        uint8 v, bytes32 r, bytes32 s // sig
    ) external;
}

interface Rule {
    // reduce checks if the given action is relevent to this rule
    // and applies any state changes.
    // reduce functions should be idempotent and side-effect free
    // ideally they would be pure functions, but it is impractical
    // for solidity/gas/storage reasons to implement state in this kind of way
    // so instead we ask that you think of it as pure even if it's not
    function reduce(State state, bytes calldata action, Context calldata ctx) external returns (State);

    // getActionTypeDefs returns metadata about the actions this rule uses
    // it can be used for introspection by clients to discover available actions
    function getActionTypeDefs() external view returns (ActionTypeDef[] memory);

    // getNodeTypeDefs returns metadata about the node types this rule uses
    // it can be used for introspection of nodetype and edgetype ids
    // function getNodeTypeDefs() external view returns (NodeTypeDef[] memory);
}

error ActionNameConflict();
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
abstract contract BaseDispatcher is Dispatcher {
    mapping(address => bool) private trustedRouters;
    mapping(string => address) private actionAddrs;
    ActionTypeDef[] private actionDefs;
    Rule[] private rules;
    State private gameState;

    constructor(State s) {
        gameState = s;
        registerRouter(address(this));
    }

    // TODO: this should be an owneronly func
    function registerRule(Rule rule) public {
        rules.push() = rule;
        // register all the actions this rule uses
        ActionTypeDef[] memory defs = rule.getActionTypeDefs();
        for (uint i=0; i<defs.length; i++) {
            registerAction(defs[i]);
        }
    }

    // TODO: this should be an owneronly func
    function registerAction(ActionTypeDef memory def) public {
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

    // registerRouter(r) will implicitly trust the Context data submitted
    // by r to dispatch(action, ctx) calls.
    //
    // this is useful if there is an external contract managing authN, authZ
    // or when using a "session key" pattern like the SessionRouter.
    //
    // TODO: this should be an owneronly func
    function registerRouter(address r) public {
        trustedRouters[r] = true;
    }

    function isRegisteredRouter(address r) internal view returns (bool) {
        return trustedRouters[r];
    }

    // getActionTypeDefs returns type metadata about the actions this Dispatcher
    // can process.
    function getActionTypeDefs() public view returns (ActionTypeDef[] memory defs) {
        return actionDefs;
    }

    // getActionID - I want to get rid of this!
    function getActionID(string memory name) public view returns (address) {
        return actionAddrs[name];
    }

    // reduce applies the action+ctx to all the rules
    function reduce(State state, bytes calldata action, Context calldata ctx) internal returns (State) {
        for (uint i=0; i<rules.length; i++) {
            state = rules[i].reduce(state, action, ctx);
        }
        return state;
    }

    function dispatch(bytes calldata action, Context calldata ctx) public {
        // check sender is trusted
        // we trust sessions built from ourself see the dispatch(action) function above that builds a full-access session for the msg.sender
        // we trust sessions built from any registered routers
        if (!isRegisteredRouter(msg.sender)) {
            revert DispatchUntrustedSender();
        }
        gameState = reduce(gameState, action, ctx);
        emit ActionDispatched(
            address(ctx.sender),
            "<nonce>" // TODO: unique ids, nonces, and replay protection
        );
    }

    function dispatch(
        bytes calldata action
    ) public {
        Context memory ctx = Context({
            sender: msg.sender,
            scopes: SCOPE_FULL_ACCESS
        });
        this.dispatch(action, ctx);
    }

}


