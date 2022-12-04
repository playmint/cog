// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Dispatcher, BaseDispatcher, Router, Rule} from "./Dispatcher.sol";
import {SessionRouter} from "./SessionRouter.sol";
import {State,StateGraph} from "./StateGraph.sol";

struct GameMetadata {
    string name;
}
// Game links together the State and Dispatcher that together
//
// TODO: this interface is too specific to playmint's setup
//       which requires a SessionRouter, but not all games really
//       require session routing so we should make this optional.
interface Game {
    event GameDeployed(
        address gameAddr,
        address dispatcherAddr,
        address stateAddr,
        address routerAddr
    );
    function getMetadata() external returns (GameMetadata memory);
    function getDispatcher() external returns (Dispatcher);
    function getRouter() external returns (SessionRouter);
    function getState() external returns (State);
}

// BasicGame implements the Game interface in a very naive way
// that is useful for getting started quickly, but is intended more
// as a reference than as a building block.
//
// It handles the creation of State, Dispatcher and Router and expects
// the inheriter to call registerRule(x).
//
// You will likely want to create your own custom Game contract that keeps the
// seperation of deployment of State, Dispatchers, Routes to enable better
// upgradability.
abstract contract BasicGame is Game {

    string internal name;
    SessionRouter internal router;
    BaseDispatcher internal dispatcher;
    State internal state;

    constructor(
        string memory newName
    ) {
        name = newName;
        state = new StateGraph();
        router = new SessionRouter();
        dispatcher = new BaseDispatcher(state);

        dispatcher.registerRouter(address(router));

        emit GameDeployed(
            address(this),
            address(dispatcher),
            address(state),
            address(router)
        );
    }

    function getMetadata() external view returns (GameMetadata memory) {
        return GameMetadata({
            name: name
        });
    }

    function getDispatcher() external view returns (Dispatcher) {
        return dispatcher;
    }

    function getRouter() external view returns (SessionRouter) {
        return router;
    }

    function getState() external view returns (State) {
        return state;
    }

}
