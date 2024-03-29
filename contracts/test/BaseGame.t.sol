// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";

import "../src/BaseRouter.sol";
import "../src/BaseState.sol";
import "../src/BaseDispatcher.sol";
import "../src/BaseGame.sol";

import "./fixtures/TestActions.sol";
import "./fixtures/TestRules.sol";
import "./fixtures/TestStateUtils.sol";

using StateTestUtils for State;

contract ExampleGame is BaseGame {
    constructor(State s, Dispatcher d, Router r) BaseGame("ExampleGame", "http://localhost:3000/") {
        _registerState(s);
        _registerRouter(r);
        _registerDispatcher(d);
    }
}

contract BaseGameTest is Test {
    event GameDeployed(address dispatcherAddr, address stateAddr, address routerAddr);

    Game game;
    BaseState state;

    uint256 ownerKey = 0xA11CE;
    address ownerAddr = vm.addr(ownerKey);

    uint256 sessionKey = 0x5e55;
    address sessionAddr = vm.addr(sessionKey);

    uint256 relayKey = 0x11111;
    address relayAddr = vm.addr(relayKey);

    function setUp() public {
        state = new BaseState();
        BaseRouter r = new BaseRouter();
        BaseDispatcher d = new BaseDispatcher();
        d.registerRouter(r);
        d.registerState(state);
        d.registerRule(new LogSenderRule());
        d.registerRule(new SetBytesRule());
        d.registerRule(new AnnotateNode());

        vm.expectEmit(true, true, true, true);
        emit GameDeployed(address(d), address(state), address(r));

        game = new ExampleGame(state, d, r);
    }

    // Ensure that we can setup sessions, dispatch signed actions and
    // have the registered rules executed to modify the state.
    function testRoutedActions() public {
        // setup a sessionkey with the router
        vm.startPrank(ownerAddr);
        game.getRouter().authorizeAddr(game.getDispatcher(), MAX_TTL, SCOPE_FULL_ACCESS, sessionAddr);
        vm.stopPrank();

        // pick a random nonce
        uint256 nonce = 4;

        // encode an action bundle
        bytes[] memory actions = new bytes[](2);
        actions[0] = abi.encodeCall(TestActions.SET_BYTES, ("MAGIC_BYTES"));
        actions[1] = abi.encodeCall(TestActions.ANNOTATE_NODE, ("A_POTENTIALLY_REALLY_LONG_UTF8_STRING"));

        // sign the action bundle with the sessionKey
        vm.startPrank(sessionAddr);
        (uint8 v, bytes32 r, bytes32 s) = sign(actions, nonce, sessionKey);
        bytes memory sig = abi.encodePacked(r, s, v);
        vm.stopPrank();

        // add the action bundle to a routing batch
        bytes[][] memory batchedActions = new bytes[][](1);
        bytes[] memory batchedSigs = new bytes[](1);
        batchedActions[0] = actions;
        batchedSigs[0] = sig;

        // dispatch the batch via a relayer
        vm.startPrank(relayAddr);
        game.getRouter().dispatch(batchedActions[0], batchedSigs[0], nonce);
        vm.stopPrank();

        // check that the state was modified as a reult of running
        // through the rules
        assertEq(game.getState().getBytes(), "MAGIC_BYTES");
        assertEq(state.getAnnotationRef(0x0, "name"), keccak256(bytes("A_POTENTIALLY_REALLY_LONG_UTF8_STRING")));
    }

    function testMetadata() public {
        GameMetadata memory metadata = game.getMetadata();
        assertEq(metadata.name, "ExampleGame");
        assertEq(metadata.url, "http://localhost:3000/");
    }

    function sign(bytes[] memory actions, uint256 nonce, uint256 privateKey)
        private
        pure
        returns (uint8 v, bytes32 r, bytes32 s)
    {
        bytes32 digest =
            keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", keccak256(abi.encode(actions, nonce))));
        return vm.sign(privateKey, digest);
    }
}
