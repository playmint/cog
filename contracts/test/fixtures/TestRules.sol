// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Context} from "../../src/IDispatcher.sol";
import {State} from "../../src/IState.sol";
import {Rule} from "../../src/IRule.sol";
import {TestActions} from "./TestActions.sol";

import {StateTestUtils} from "./TestStateUtils.sol";

using StateTestUtils for State;

contract LastWriteWinsRule is Rule {
    bytes20 b;

    constructor(bytes20 bytesToSet) {
        b = bytesToSet;
    }

    function reduce(State s, bytes calldata, /*action*/ Context calldata /*ctx*/ ) public returns (State) {
        s.setBytes(b);
        return s;
    }
}

contract SetBytesRule is Rule {
    function reduce(State s, bytes calldata action, Context calldata /*ctx*/ ) public returns (State) {
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

contract AnnotateNode is Rule {
    function reduce(State s, bytes calldata action, Context calldata /*ctx*/ ) public returns (State) {
        if (bytes4(action) == TestActions.ANNOTATE_NODE.selector) {
            (string memory data) = abi.decode(action[4:], (string));
            s.annotate(0x0, "name", data);
        }
        return s;
    }
}
