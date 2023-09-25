// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "./IState.sol";
import "./IDispatcher.sol";
import "./IRouter.sol";
import "./IRule.sol";
import {Op, BaseState} from "./BaseState.sol";

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
    BaseState private state;

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
    // or when using a "session key" pattern like the BaseRouter.
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
        state = BaseState(address(s));
    }

    function isRegisteredRouter(address r) internal view returns (bool) {
        return trustedRouters[Router(r)];
    }

    function dispatch(bytes calldata action, Context calldata ctx) public returns (Op[] memory) {
        uint256 fromHead = state.getHead();
        // check ctx can be trusted
        // we trust ctx built from ourself see the dispatch(action) function above that builds a full-access session for the msg.sender
        // we trust ctx built from any registered routers
        if (!isRegisteredRouter(msg.sender)) {
            revert("DispatchUntrustedSender");
        }
        for (uint256 i = 0; i < rules.length; i++) {
            rules[i].reduce(state, action, ctx);
        }
        emit ActionDispatched(
            address(ctx.sender),
            "<nonce>" // TODO: unique ids, nonces, and replay protection
        );
        uint256 toHead = state.getHead();
        return state.getOps(fromHead, toHead);
    }

    // dispatch from router trusted context
    function dispatch(bytes[] calldata actions, Context calldata ctx) public returns (Op[] memory) {
        uint256 fromHead = state.getHead();
        for (uint256 i = 0; i < actions.length; i++) {
            dispatch(actions[i], ctx);
        }
        uint256 toHead = state.getHead();
        return state.getOps(fromHead, toHead);
    }

    function dispatch(bytes calldata action) public returns (Op[] memory) {
        Context memory ctx = Context({sender: msg.sender, scopes: SCOPE_FULL_ACCESS, clock: uint32(block.number)});
        return this.dispatch(action, ctx);
    }

    function dispatch(bytes[] calldata actions) public returns (Op[] memory) {
        uint256 fromHead = state.getHead();
        for (uint256 i = 0; i < actions.length; i++) {
            dispatch(actions[i]);
        }
        uint256 toHead = state.getHead();
        return state.getOps(fromHead, toHead);
    }
}
