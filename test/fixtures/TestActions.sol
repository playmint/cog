// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

interface TestActions {
    function NOOP() external;
    function SET_SENDER() external;
    function SET_BYTES(bytes memory) external;
}
