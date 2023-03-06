// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {BaseDispatcher, Dispatcher, Router, DispatchUntrustedSender, Rule, Context} from "../src/Dispatcher.sol";
import {State} from "../src/State.sol";
import {StateGraph} from "../src/StateGraph.sol";

import "./fixtures/TestActions.sol";
import "./fixtures/TestRules.sol";
import "./fixtures/TestStateUtils.sol";

using StateTestUtils for State;

contract ExampleDispatcher is Dispatcher, BaseDispatcher {
    constructor(State s) BaseDispatcher() {
        _registerState(s);
    }
}

contract BaseDispatcherTest is Test {
    State s;
    BaseDispatcher d;

    function setUp() public {
        s = new StateGraph(); // TODO: replace with a mock
        d = new ExampleDispatcher(s);
    }

    function testReigsterRuleOrder() public {
        bytes memory action = abi.encodeCall(TestActions.NOOP, ());

        d.registerRule(new LastWriteWinsRule("x"));
        d.registerRule(new LastWriteWinsRule("y"));
        d.dispatch(action);

        assertEq(s.getBytes(), "y");
    }

    function testActionArgs() public {
        bytes memory action = abi.encodeCall(TestActions.SET_BYTES, ("MAGIC_BYTES"));

        d.registerRule(new SetBytesRule());
        d.dispatch(action);

        assertEq(s.getBytes(), "MAGIC_BYTES");
    }

    function testDispatchAsSender() public {
        address sender = vm.addr(0xcafe);
        bytes memory action = abi.encodeCall(TestActions.SET_SENDER, ());

        d.registerRule(new LogSenderRule());
        vm.prank(sender);
        d.dispatch(action);

        assertEq(s.getAddress(), sender);
    }

    function testDispatchWithContext() public {
        address router = vm.addr(0x88888);
        address sender = vm.addr(0x11111);

        d.registerRouter(Router(router));

        Context memory ctx = newContext(sender);
        bytes memory action = abi.encodeCall(TestActions.SET_SENDER, ());

        d.registerRule(new LogSenderRule());
        vm.prank(router);
        d.dispatch(action, ctx);

        assertEq(s.getAddress(), sender);
    }

    function testRevertUntrustedRouter() public {
        address router = vm.addr(0x88888);
        address sender = vm.addr(0x11111);

        Context memory ctx = newContext(sender);
        bytes memory action = abi.encodeCall(TestActions.SET_SENDER, ());

        d.registerRule(new LogSenderRule());
        vm.expectRevert(DispatchUntrustedSender.selector);
        vm.prank(router);
        d.dispatch(action, ctx);

        assertEq(s.getAddress(), address(0));
    }

    function newContext(address sender) private view returns (Context memory) {
        return Context({
            sender: sender,
            scopes: 0,
            clock: uint32(block.number),
            annotations: new bytes32[](0)
        });
    }
}
