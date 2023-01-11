// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {
    State,
    Rule,
    Context
} from "../../src/Dispatcher.sol";
import { TestActions } from "./TestActions.sol";

import {StateTestUtils} from "./TestStateUtils.sol";
using StateTestUtils for State;

contract LastWriteWinsRule is Rule {
    bytes20 b;
    constructor(bytes20 bytesToSet) {
        b = bytesToSet;
    }
    function reduce(State s, bytes calldata /*action*/, Context calldata /*ctx*/) public returns (State) {
        s.setBytes(b);
        return s;
    }
}

contract SetBytesRule is Rule {
    function reduce(State s, bytes calldata action, Context calldata /*ctx*/) public returns (State) {
        if (bytes4(action) == TestActions.SET_BYTES.selector) {
            (bytes memory b) = abi.decode(action[4:], (bytes));
            s.setBytes(bytes20(b));
        }
        return s;
    }
}

contract LogSenderRule is Rule {
    function reduce(State s, bytes calldata action, Context calldata ctx) public returns (State) {
        if (bytes4(action) == TestActions.SET_SENDER.selector) {
            s.setAddress(ctx.sender);
        }
        return s;
    }
}

