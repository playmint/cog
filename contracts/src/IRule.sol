// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "./IState.sol";
import "./IDispatcher.sol";

interface Rule {
    // reduce checks if the given action is relevent to this rule
    // and applies any state changes.
    // reduce functions should be idempotent and side-effect free
    // ideally they would be pure functions, but it is impractical
    // for solidity/gas/storage reasons to implement state in this kind of way
    // so instead we ask that you think of it as pure even if it's not
    function reduce(State state, bytes calldata action, Context calldata ctx) external returns (State);
}
