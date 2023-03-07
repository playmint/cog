// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {State, WeightKind, CompoundKeyKind, AnnotationKind} from "../src/State.sol";
import {StateGraph} from "../src/StateGraph.sol";

interface Rel {
    function Friend() external;
}

interface Kind {
    function Person() external;
}

contract StateGraphTest is Test {
    event EdgeTypeRegister(bytes4 id, string name, WeightKind kind);
    event NodeTypeRegister(bytes4 id, string name, CompoundKeyKind keyKind);
    event EdgeSet(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint160 weight);
    event EdgeRemove(bytes4 relID, uint8 relKey, bytes24 srcNodeID);
    event AnnotationSet(bytes24 id, AnnotationKind kind, string label, bytes32 ref);

    StateGraph internal state;

    function setUp() public {
        state = new StateGraph();
    }

    function testSetEdge() public {
        bytes24 srcPersonID = bytes24(abi.encodePacked(Kind.Person.selector, uint64(1)));
        bytes24 dstPersonID = bytes24(abi.encodePacked(Kind.Person.selector, uint64(2)));

        bytes4 relID = Rel.Friend.selector;
        uint8 relKey = 100;
        uint64 weight = 1;

        vm.expectEmit(true, true, true, true, address(state));
        emit EdgeSet(relID, relKey, srcPersonID, dstPersonID, weight);

        state.set(relID, relKey, srcPersonID, dstPersonID, weight);

        (bytes24 gotPersonID, uint160 gotWeight) = state.get(relID, relKey, srcPersonID);

        assertEq(gotPersonID, dstPersonID);
        assertEq(gotWeight, weight);
    }

    function testRemoveEdge() public {
        bytes24 srcPersonID = bytes24(abi.encodePacked(Kind.Person.selector, uint64(1)));
        bytes24 dstPersonID = bytes24(abi.encodePacked(Kind.Person.selector, uint64(2)));
        bytes4 relID = Rel.Friend.selector;
        uint8 relKey = 100;
        uint64 weight = 1;

        state.set(relID, relKey, srcPersonID, dstPersonID, weight);

        (bytes24 gotPersonID, uint160 gotWeight) = state.get(relID, relKey, srcPersonID);
        assertEq(gotPersonID, dstPersonID);
        assertEq(gotWeight, weight);

        vm.expectEmit(true, true, true, true, address(state));
        emit EdgeRemove(relID, relKey, srcPersonID);

        state.remove(relID, relKey, srcPersonID);

        (bytes24 gotPersonIDAfterRemove, uint160 gotWeightAfterRemove) = state.get(relID, relKey, srcPersonID);
        assertEq(gotPersonIDAfterRemove, 0);
        assertEq(gotWeightAfterRemove, 0);
    }

    function testRegisterEdgeType() public {
        bytes4 relID = bytes4(uint32(1));
        string memory relName = "TESTING_EDGE_NAME";
        WeightKind weightKind = WeightKind.UINT64;
        vm.expectEmit(true, true, true, true, address(state));
        emit EdgeTypeRegister(relID, relName, weightKind);
        state.registerEdgeType(relID, relName, weightKind);
    }

    function testRegisterNodeType() public {
        bytes4 relID = bytes4(uint32(2));
        string memory relName = "TESTING_NODE_NAME";
        CompoundKeyKind keyKind = CompoundKeyKind.UINT160;
        vm.expectEmit(true, true, true, true, address(state));
        emit NodeTypeRegister(relID, relName, keyKind);
        state.registerNodeType(relID, relName, keyKind);
    }

    function testAnnotateNode() public {
        bytes24 nodeID = bytes24(abi.encodePacked(Kind.Person.selector, uint64(1)));
        string memory label = "ann";
        string memory data = "A_STRING_LONGER_THAN_32_BYTES_1234567890123456789012345678901234567890";
        vm.expectEmit(true, true, true, true, address(state));
        emit AnnotationSet(nodeID, AnnotationKind.CALLDATA, label, keccak256(bytes(data)));
        state.setAnnotation(nodeID, label, data);
    }
}
