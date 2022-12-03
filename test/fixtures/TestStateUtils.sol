// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {State, NodeTypeUtils, NodeType, NodeData} from "../../src/State.sol";

// some wrappers to treat State as a single value that you can set or get
library StateTestUtils {
    function set(State s, uint256 value) internal returns (State) {
        return s.setNode(
            NodeTypeUtils.ID(NodeType(address(0)), 0, 1),
            NodeData.wrap(uint256(value))
        );
    }
    function set(State s, address value) internal returns (State) {
        return s.setNode(
            NodeTypeUtils.ID(NodeType(address(0)), 0, 1),
            NodeData.wrap(uint256(uint160(value)))
        );
    }
    function set(State s, bytes32 value) internal returns (State) {
        return s.setNode(
            NodeTypeUtils.ID(NodeType(address(0)), 0, 1),
            NodeData.wrap(uint256(value))
        );
    }
    function getUint(State s) internal view returns (uint256) {
        return uint256(uint160(NodeData.unwrap(s.getNode(
            NodeTypeUtils.ID(NodeType(address(0)), 0, 1)
        ))));
    }
    function getAddress(State s) internal view returns (address) {
        return address(uint160(NodeData.unwrap(s.getNode(
            NodeTypeUtils.ID(NodeType(address(0)), 0, 1)
        ))));
    }
    function getBytes32(State s) internal view returns (bytes32) {
        return bytes32(NodeData.unwrap(s.getNode(
            NodeTypeUtils.ID(NodeType(address(0)), 0, 1)
        )));
    }
}
