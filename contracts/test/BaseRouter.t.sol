// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {Dispatcher, Context} from "../src/IDispatcher.sol";
import {State} from "../src/IState.sol";
import {BaseState} from "../src/BaseState.sol";
import {BaseDispatcher} from "../src/BaseDispatcher.sol";
import {BaseRouter, PREFIX_MESSAGE, REVOKE_MESSAGE} from "../src/BaseRouter.sol";

import "./fixtures/TestActions.sol";
import "./fixtures/TestRules.sol";
import "./fixtures/TestStateUtils.sol";

using StateTestUtils for State;

import {LibString} from "../src/utils/LibString.sol";

using {LibString.toString} for uint256;
using LibString for address;

contract ExampleDispatcher is Dispatcher, BaseDispatcher {
    constructor(State s) BaseDispatcher() {
        _registerState(s);
        _registerRule(new LogSenderRule());
    }
}

contract BaseRouterTest is Test {
    State state;
    BaseDispatcher dispatcher;
    BaseRouter router;

    uint256 ownerKey = 0xA11CE;
    address ownerAddr = vm.addr(ownerKey);

    uint256 sessionKey = 0x5e55;
    address sessionAddr = vm.addr(sessionKey);

    uint256 relayKey = 0x11111;
    address relayAddr = vm.addr(relayKey);

    function setUp() public {
        state = new BaseState(); // TODO: replace with a mock
        dispatcher = new ExampleDispatcher(state);
        router = new BaseRouter();
        dispatcher.registerRouter(router);
    }

    function testSanityCheckAddrs() public {
        assertFalse(ownerAddr == address(0));
        assertFalse(sessionAddr == address(0));
        assertFalse(sessionAddr == ownerAddr);
        assertFalse(relayAddr == address(0));
        assertFalse(relayAddr == ownerAddr);
    }

    function testUnauthorizeSignerAsOwner() public {
        // should not be able to just sign actions with any old key
        vm.expectRevert("SessionUnauthorized");
        dispatchSigned(0x666);
        assertEq(state.getAddress(), address(0));
    }

    function testAuthorizeAddrWithSenderAsOwner() public {
        // authorize a session key by talking direct to
        // contract as sender
        vm.prank(ownerAddr);
        router.authorizeAddr(dispatcher, 0, 0, sessionAddr);

        // should now be able to use sessionKey to submit signed actions
        // that act as the original owner
        dispatchSigned(sessionKey);
        assertEq(state.getAddress(), ownerAddr);
    }

    function testAuthorizeAddrWithSignerAsOwner() public {
        // expected auth message
        bytes memory authMessage = abi.encodePacked(
            "Welcome!",
            "\n\nThis site is requesting permission to create a temporary session key.",
            "\n\nSigning this message will not incur any fees.",
            "\n\nValid: 5 blocks",
            "\n\nSession: ",
            sessionAddr.toHexString()
        );

        // owner signs the message authorizing the session
        (uint8 v, bytes32 r, bytes32 s) =
            vm.sign(ownerKey, keccak256(abi.encodePacked(PREFIX_MESSAGE, authMessage.length.toString(), authMessage)));
        bytes memory sig = abi.encodePacked(r, s, v);

        // relay submits the auth request on behalf of owner
        vm.prank(relayAddr);
        router.authorizeAddr(dispatcher, 5, 0, sessionAddr, sig);

        // should now be able to use sessionKey to act as owner
        dispatchSigned(sessionKey);
        assertEq(state.getAddress(), ownerAddr);
    }

    function testRevokeAddrWithSenderAsOwner() public {
        vm.prank(ownerAddr);
        router.authorizeAddr(dispatcher, 0, 0, sessionAddr);
        dispatchSigned(sessionKey);
        assertEq(state.getAddress(), ownerAddr);

        // sender trusted to destroy their session
        vm.prank(ownerAddr);
        router.revokeAddr(sessionAddr);

        // session signed actions should now fail...
        vm.expectRevert("SessionUnauthorized");
        dispatchSigned(sessionKey);
    }

    function testRevokeAddrWithSignerAsOwner() public {
        vm.prank(ownerAddr);
        router.authorizeAddr(dispatcher, 0, 0, sessionAddr);
        dispatchSigned(sessionKey);
        assertEq(state.getAddress(), ownerAddr);

        // owner signs the message destroying the session
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(
            ownerKey,
            keccak256(
                abi.encodePacked(PREFIX_MESSAGE, (REVOKE_MESSAGE.length + 20).toString(), REVOKE_MESSAGE, sessionAddr)
            )
        );
        bytes memory sig = abi.encodePacked(r, s, v);
        vm.prank(relayAddr);
        router.revokeAddr(sessionAddr, sig);

        // session signed actions should now fail...
        vm.expectRevert("SessionUnauthorized");
        dispatchSigned(sessionKey);
    }

    // dispatches a SET_SENDER action with msg.sender set to relayAddr
    // the LogSenderRule sets the state to the action's owner so we
    // can confirm what the action got processed as
    function dispatchSigned(uint256 privateKey) internal {
        uint256 nonce = 111; // random pick
        bytes[] memory actions = new bytes[](1);
        actions[0] = abi.encodeCall(TestActions.SET_SENDER, ());
        bytes32 digest =
            keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", keccak256(abi.encode(actions, nonce))));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privateKey, digest);
        bytes memory sig = abi.encodePacked(r, s, v);

        vm.prank(relayAddr);
        router.dispatch(actions, sig, nonce);
    }
}
