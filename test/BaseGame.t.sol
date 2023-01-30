// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {
    BaseDispatcher,
    Dispatcher,
    Router,
    DispatchUntrustedSender,
    Rule,
    Context,
    SCOPE_FULL_ACCESS
} from "../src/Dispatcher.sol";
import {State} from "../src/State.sol";
import {StateGraph} from "../src/StateGraph.sol";
import {Game, BaseGame, GameMetadata} from "../src/Game.sol";
import {SessionRouter, MAX_TTL} from "../src/SessionRouter.sol";

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

    event GameDeployed(
        address dispatcherAddr,
        address stateAddr,
        address routerAddr
    );

    Game game;

    uint256 ownerKey = 0xA11CE;
    address ownerAddr = vm.addr(ownerKey);

    uint256 sessionKey = 0x5e55;
    address sessionAddr = vm.addr(sessionKey);

    uint256 relayKey = 0x11111;
    address relayAddr = vm.addr(relayKey);

    function setUp() public {
        State s = new StateGraph();
        SessionRouter r = new SessionRouter();
        BaseDispatcher d = new BaseDispatcher();
        d.registerRouter(r);
        d.registerState(s);
        d.registerRule(new LogSenderRule());
        d.registerRule(new SetBytesRule());

        vm.expectEmit(true, true, true, true);
        emit GameDeployed(
            address(d),
            address(s),
            address(r)
        );

        game = new ExampleGame(s, d, r);
    }

    // Ensure that we can setup sessions, dispatch signed actions and
    // have the registered rules executed to modify the state.
    function testRoutedActions() public {

        // setup a sessionkey with the router
        vm.startPrank(ownerAddr);
        game.getRouter().authorizeAddr(
            game.getDispatcher(),
            MAX_TTL,
            SCOPE_FULL_ACCESS,
            sessionAddr
        );
        vm.stopPrank();

        // sign an action with the sessionKey
        vm.startPrank(sessionAddr);
        bytes memory action = abi.encodeCall(TestActions.SET_BYTES, ("MAGIC_BYTES"));
        (uint8 v, bytes32 r, bytes32 s) = sign(action, sessionKey);
        bytes memory sig = abi.encodePacked(r,s,v);
        vm.stopPrank();

        // dispatch the signed action via a relayer
        vm.startPrank(relayAddr);
        game.getRouter().dispatch(action, sig);
        vm.stopPrank();

        // check that the state was modified as a reult of running
        // through the rules
        assertEq(
            game.getState().getBytes(),
            "MAGIC_BYTES"
        );

    }

    function testMetadata() public {
        GameMetadata memory metadata = game.getMetadata();
        assertEq(metadata.name, "ExampleGame");
        assertEq(metadata.url, "http://localhost:3000/");
    }

    function sign(bytes memory action, uint256 privateKey) private pure returns (uint8 v, bytes32 r, bytes32 s) {
        bytes32 digest = keccak256(abi.encodePacked(
            "\x19Ethereum Signed Message:\n32",
            keccak256(action)
        ));
        return vm.sign(privateKey, digest);
    }

}
