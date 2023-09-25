// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {State, WeightKind, CompoundKeyKind, AnnotationKind} from "./IState.sol";

error StateUnauthorizedSender();

enum OpKind {
    EdgeSet,
    EdgeRemove,
    AnnotationSet
}

struct Op {
    OpKind kind;
    bytes4 relID;
    uint8 relKey;
    bytes24 srcNodeID;
    bytes24 dstNodeID;
    uint160 weight;
    string annName;
    string annData;
}

contract BaseState is State {
    struct EdgeData {
        bytes24 dstNodeID;
        uint64 weight;
    }

    mapping(bytes24 => mapping(bytes4 => mapping(uint8 => EdgeData))) edges;
    mapping(bytes24 => mapping(bytes32 => bytes32)) annotations;
    mapping(address => bool) allowlist;

    Op[] ops;

    constructor() {
        // register the zero value under the kind name NULL
        _registerNodeType(0, "NULL", CompoundKeyKind.NONE);
    }

    function getHead() public view returns (uint256) {
        return ops.length;
    }

    function getOps(uint256 from, uint256 to) public view returns (Op[] memory res) {
        if (from == to) {
            res = new Op[](0);
            return res;
        }
        require(to > from, "TO must be after FROM");
        res = new Op[](to - from);
        for (uint256 i = 0; i < res.length; i++) {
            res[i] = ops[from + i];
        }
        return res;
    }

    function set(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint64 weight) external {
        // TODO: uncomment this
        // if (!allowlist[msg.sender]) {
        //     revert StateUnauthorizedSender();
        // }
        edges[srcNodeID][relID][relKey] = EdgeData(dstNodeID, weight);
        emit State.EdgeSet(relID, relKey, srcNodeID, dstNodeID, weight);

        Op storage op = ops.push();
        op.kind = OpKind.EdgeSet;
        op.relID = relID;
        op.relKey = relKey;
        op.srcNodeID = srcNodeID;
        op.dstNodeID = dstNodeID;
        op.weight = weight;
    }

    function remove(bytes4 relID, uint8 relKey, bytes24 srcNodeID) external {
        // TODO: uncomment this
        // if (!allowlist[msg.sender]) {
        //     revert StateUnauthorizedSender();
        // }
        delete edges[srcNodeID][relID][relKey];
        emit State.EdgeRemove(relID, relKey, srcNodeID);

        Op storage op = ops.push();
        op.kind = OpKind.EdgeRemove;
        op.relID = relID;
        op.relKey = relKey;
        op.srcNodeID = srcNodeID;
    }

    function get(bytes4 relID, uint8 relKey, bytes24 srcNodeID)
        external
        view
        returns (bytes24 dstNodeID, uint64 weight)
    {
        EdgeData storage e = edges[srcNodeID][relID][relKey];
        return (e.dstNodeID, e.weight);
    }

    function annotate(bytes24 nodeID, string memory label, string memory annotationData) external {
        bytes32 annotationRef = keccak256(bytes(annotationData));
        annotations[nodeID][keccak256(bytes(label))] = annotationRef;
        emit State.AnnotationSet(nodeID, AnnotationKind.CALLDATA, label, annotationRef, annotationData);

        Op storage op = ops.push();
        op.kind = OpKind.AnnotationSet;
        op.srcNodeID = nodeID;
        op.annName = label;
        op.annData = annotationData;
    }

    function getAnnotationRef(bytes24 nodeID, string memory annotationLabel) external view returns (bytes32) {
        return annotations[nodeID][keccak256(bytes(annotationLabel))];
    }

    // TODO: allowlist only
    function registerNodeType(bytes4 kindID, string memory kindName, CompoundKeyKind keyKind) external {
        _registerNodeType(kindID, kindName, keyKind);
    }

    function _registerNodeType(bytes4 kindID, string memory kindName, CompoundKeyKind keyKind) internal {
        emit State.NodeTypeRegister(kindID, kindName, keyKind);
    }

    // TODO: allowlist only
    function registerEdgeType(bytes4 relID, string memory relName, WeightKind weightKind) external {
        _registerEdgeType(relID, relName, weightKind);
    }

    function _registerEdgeType(bytes4 relID, string memory relName, WeightKind weightKind) internal {
        emit State.EdgeTypeRegister(relID, relName, weightKind);
    }

    // TODO: owner only
    function authorizeContract(address addr) external {
        allowlist[addr] = true;
    }
}
