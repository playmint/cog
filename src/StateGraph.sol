// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {
    State,
    EdgeData
} from "./State.sol";

error StateUnauthorizedSender();

contract StateGraph is State {

    mapping(bytes12 => mapping(bytes4 => mapping(bytes8 => EdgeData))) edges;
    mapping(address => bool) allowlist;

    constructor() { }

    function set(bytes4 relID, bytes8 relKey, bytes12 srcNodeID, bytes12 dstNodeID, uint160 weight) external {
        // TODO: uncomment this
        // if (!allowlist[msg.sender]) {
        //     revert StateUnauthorizedSender();
        // }
        edges[srcNodeID][relID][relKey] = EdgeData(dstNodeID, weight);
        emit State.EdgeSet(
            relID,
            relKey,
            srcNodeID,
            dstNodeID,
            weight
        );
    }

    function get(bytes4 relID, bytes8 relKey, bytes12 srcNodeID) external view returns (bytes12 dstNodeID, uint160 weight) {
        EdgeData storage e = edges[srcNodeID][relID][relKey];
        return (e.dstNodeID, e.weight);
    }

    function registerNodeType(bytes4 kindID, string memory kindName) external {
        emit State.NodeTypeRegister(
            kindID,
            kindName
        );
    }

    function registerEdgeType(bytes4 relID, string memory relName) external {
        emit State.EdgeTypeRegister(
            relID,
            relName
        );
    }

    function authorizeContract(address addr) external {
        allowlist[addr] = true;
    }
}

