// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {
    State,
    WeightKind,
    CompoundKeyKind
} from "./State.sol";


error StateUnauthorizedSender();

contract StateGraph is State {

    struct EdgeData {
        bytes24 dstNodeID;
        uint64 weight;
    }

    mapping(bytes24 => mapping(bytes4 => mapping(uint8 => EdgeData))) edges;
    mapping(address => bool) allowlist;

    constructor() {
        // register the zero value under the kind name NULL
        _registerNodeType(0, "NULL", CompoundKeyKind.NONE);
    }

    function set(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint64 weight) external {
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

    function get(bytes4 relID, uint8 relKey, bytes24 srcNodeID) external view returns (bytes24 dstNodeID, uint64 weight) {
        EdgeData storage e = edges[srcNodeID][relID][relKey];
        return (e.dstNodeID, e.weight);
    }

    // TODO: allowlist only
    function registerNodeType(bytes4 kindID, string memory kindName, CompoundKeyKind keyKind) external {
        _registerNodeType(kindID, kindName, keyKind);
    }

    function _registerNodeType(bytes4 kindID, string memory kindName, CompoundKeyKind keyKind) internal {
        emit State.NodeTypeRegister(
            kindID,
            kindName,
            keyKind
        );
    }

    // TODO: allowlist only
    function registerEdgeType(bytes4 relID, string memory relName, WeightKind weightKind) external {
        _registerEdgeType(relID, relName, weightKind);
    }

    function _registerEdgeType(bytes4 relID, string memory relName, WeightKind weightKind) internal {
        emit State.EdgeTypeRegister(
            relID,
            relName,
            weightKind
        );
    }

    // TODO: owner only
    function authorizeContract(address addr) external {
        allowlist[addr] = true;
    }
}

