// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {
    Attribute,
    AttributeKind,
    State,
    NodeData,
    EdgeType,
    EdgeData,
    EdgeTypeID,
    EdgeTypeUtils,
    NodeType,
    NodeIDUtils,
    NodeTypeID,
    NodeTypeUtils,
    NodeID
} from "../src/State.sol";
import {StateGraph} from "../src/StateGraph.sol";

using NodeTypeUtils for NodeType;
using EdgeTypeUtils for EdgeType;
using NodeIDUtils for NodeID;

contract Thing is NodeType {
    function getAttributes(NodeID, NodeData) public pure returns (Attribute[] memory attrs) {
        attrs = new Attribute[](3);
        attrs[0].name = "kind";
        attrs[0].kind = AttributeKind.STRING;
        attrs[0].value = bytes32("THING");
    }
}

contract HasOne is EdgeType {
    function getAttributes(NodeID /*id*/, uint /*idx*/) public pure returns (Attribute[] memory attrs) {
        attrs = new Attribute[](1);
        attrs[0].name = "kind";
        attrs[0].kind = AttributeKind.STRING;
        attrs[0].value = bytes32("HAS_ONE");
    }
}

contract StateGraphTest is Test {

    event NodeSet(
        NodeID nodeID,
        NodeData nodeData
    );

    event EdgeSet(
        EdgeTypeID kind,
        NodeID srcNodeID,
        NodeID dstNodeID,
        uint idx,
        uint32 weight
    );


    StateGraph internal g;

    NodeType thing;
    EdgeType hasOne;

    function setUp() public {
        g = new StateGraph();
        thing = new Thing();
        hasOne = new HasOne();
    }

    function testSetNode() public {
        NodeID nodeA = thing.ID(1);

        uint256 inData = 1;

        vm.expectEmit(true, true, true, true, address(g));
        emit NodeSet(
            nodeA,
            NodeData.wrap(inData)
        );

        g.setNode(nodeA, NodeData.wrap(inData));
        uint256 outData = NodeData.unwrap(g.getNode(nodeA));

        assertEq(inData, outData);
    }

    function testSetEdge() public {
        NodeID nodeA = thing.ID(1);
        NodeID nodeB = thing.ID(2);

        vm.expectEmit(true, true, true, true, address(g));
        emit EdgeSet(
            hasOne.ID(),
            nodeA,
            nodeB,
            0,            // idx
            uint32(65000) // weight
        );

        g.setEdge(hasOne.ID(), nodeA, EdgeData({
            nodeID: nodeB,
            weight: 65000
        }));

        EdgeData memory outEdge = g.getEdge(hasOne.ID(), nodeA);

        assertEq(
            NodeID.unwrap(outEdge.nodeID),
            NodeID.unwrap(nodeB)
        );
        assertEq(
            outEdge.weight,
            65000
        );
    }

    function testAppendEdges() public {
        NodeID nodeA = thing.ID(100);

        for (uint32 i=0; i<3; i++) {
            NodeID nodeB = thing.ID(i);

            vm.expectEmit(true, true, true, true, address(g));
            emit EdgeSet(
                hasOne.ID(),
                nodeA,
                nodeB,
                i,        // idx
                uint32(i) // weight
            );

            g.appendEdge(hasOne.ID(), nodeA, EdgeData({
                nodeID: nodeB,
                weight: uint32(i)
            }));
        }

        EdgeData[] memory outEdges = g.getEdges(hasOne.ID(), nodeA);

        for (uint32 i=0; i<3; i++) {
            NodeID nodeB = thing.ID(i);
            assertEq(
                NodeID.unwrap(outEdges[i].nodeID),
                NodeID.unwrap(nodeB)
            );
        }
    }

}
