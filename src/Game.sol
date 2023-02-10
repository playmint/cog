// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Dispatcher, Router, Rule} from "./Dispatcher.sol";
import {SessionRouter} from "./SessionRouter.sol";
import {State, StateGraph} from "./StateGraph.sol";

struct GameMetadata {
    string name;
    string url;
}
// Game links together the State and Dispatcher that together
//
// TODO: this interface is too specific to playmint's setup
//       which requires a SessionRouter, but not all games really
//       require session routing so we should make this optional.

interface Game {
    event GameDeployed(address dispatcherAddr, address stateAddr, address routerAddr);

    function getMetadata() external returns (GameMetadata memory);
    function getDispatcher() external returns (Dispatcher);
    function getRouter() external returns (Router);
    function getState() external returns (State);
}

// BaseGame implements a basic shell for implementing Game
abstract contract BaseGame is Game {
    string internal name;
    string internal url;
    Router internal router;
    Dispatcher internal dispatcher;
    State internal state;

    constructor(string memory newName, string memory newURL) {
        name = newName;
        url = newURL;
    }

    function getMetadata() public view returns (GameMetadata memory) {
        return GameMetadata({name: name, url: url});
    }

    // TODO: should be OwnerOnly
    function _registerDispatcher(Dispatcher d) internal {
        dispatcher = d;
        emitUpdate();
    }

    function getDispatcher() external view returns (Dispatcher) {
        return dispatcher;
    }

    // TODO: should be OwnerOnly
    function _registerRouter(Router r) internal {
        router = r;
        emitUpdate();
    }

    function getRouter() external view returns (Router) {
        return router;
    }

    // TODO: should be OwnerOnly
    function _registerState(State s) internal {
        state = s;
        emitUpdate();
    }

    function getState() external view returns (State) {
        return state;
    }

    function emitUpdate() internal {
        if (address(dispatcher) == address(0)) {
            return;
        }
        if (address(state) == address(0)) {
            return;
        }
        if (address(router) == address(0)) {
            return;
        }
        emit GameDeployed(address(dispatcher), address(state), address(router));
    }
}
