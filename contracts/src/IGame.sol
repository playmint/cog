// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Dispatcher} from "./IDispatcher.sol";
import {Router} from "./IRouter.sol";
import {State} from "./IState.sol";

struct GameMetadata {
    string name;
    string url;
}

// A Game advertises the it's Dispatcher (entrypoint), State (data),
// Router (session manager), and GameMetadata (name) to indexers and
// serves as a jumping off point to interacting with the game contracts
interface Game {
    event GameDeployed(address dispatcherAddr, address stateAddr, address routerAddr);

    function getMetadata() external returns (GameMetadata memory);
    function getDispatcher() external returns (Dispatcher);
    function getRouter() external returns (Router);
    function getState() external returns (State);
}
