// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Dispatcher} from "./Dispatcher.sol";
import {State} from "./State.sol";

// Game links together the State and Dispatcher that together
// form the "game" and gives it a name.
// This enables discovery of the game's entrypoint and (via the Dispatcher)
// the Actions supported.
interface Game {
    function getName() external returns (string memory);
    function getDispatcher() external returns (Dispatcher);
    function getState() external returns (State);
}
