// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {
    Attribute,
    State,
    NodeData,
    EdgeData,
    EdgeTypeID,
    EdgeType,
    NodeIDUtils,
    NodeTypeID,
    NodeID,
    NodeType
} from "./State.sol";

contract StateGraph is State {

    mapping(NodeID => NodeData) nodes;
    mapping(NodeID => mapping(EdgeTypeID => EdgeData[])) edges;

    using NodeIDUtils for NodeID;

    constructor() { }

    function getNode(NodeID nodeID) public view returns (NodeData) {
        return nodes[nodeID];
    }

    function getNodeAttributes(NodeID id, NodeData data) public view returns (Attribute[] memory attrs) {
        NodeType kind = id.getType();
        return kind.getAttributes(id, data);
    }

    function setNode(NodeID id, NodeData data) public returns (State) {
        // TODO: require only registered node type can call this
        nodes[id] = data;
        emit State.NodeSet(
            id,
            data
        );
        return this;
    }

    function setEdge(EdgeTypeID t, NodeID srcNodeID, uint idx, EdgeData memory data) public returns (State) {
        if (edges[srcNodeID][t].length == idx) {
            edges[srcNodeID][t].push() = data;
        } else {
            edges[srcNodeID][t][idx] = data;
        }
        emit State.EdgeSet(
            t,
            srcNodeID,
            data.nodeID,
            idx,
            data.weight
        );
        return this;
    }

    // setEdge without index. Use when you know there is only ever a single edge of the given type.
    function setEdge(EdgeTypeID t, NodeID srcNodeID, EdgeData memory data) public returns (State) {
        return setEdge(t, srcNodeID, uint(0), data);
    }

    // appendEdge performs a setEdge at the end of the list.
    function appendEdge(EdgeTypeID t, NodeID srcNodeID, EdgeData memory data) public returns (State) {
        return setEdge(t, srcNodeID, uint(edges[srcNodeID][t].length), data);
    }

    function getEdge(EdgeTypeID t, NodeID srcNodeID, uint idx) public view returns (EdgeData memory) {
        return edges[srcNodeID][t][idx];
    }

    function getEdge(EdgeTypeID t, NodeID srcNodeID) public view returns (EdgeData memory edge) {
        if (edges[srcNodeID][t].length == 0) {
            return edge;
        }
        return edges[srcNodeID][t][0];
    }

    function getEdges(EdgeTypeID t, NodeID srcNodeID) public view returns (EdgeData[] memory) {
        return edges[srcNodeID][t];
    }

    function getEdgeAttributes(EdgeType t, NodeID srcNodeID, uint idx) public view returns (Attribute[] memory attrs) {
        return t.getAttributes(srcNodeID, idx);
    }

}

