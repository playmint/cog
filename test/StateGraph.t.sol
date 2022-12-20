// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {State, WeightKind} from "../src/State.sol";
import {StateGraph} from "../src/StateGraph.sol";

interface Rel {
    function Friend() external;
}

interface Kind {
    function Person() external;
}

contract StateGraphTest is Test {

    event EdgeTypeRegister(
        bytes4 id,
        string name,
        WeightKind kind
    );
    event NodeTypeRegister(
        bytes4 id,
        string name
    );
    event EdgeSet(
        bytes4 relID,
        uint8 relKey,
        bytes12 srcNodeID,
        bytes12 dstNodeID,
        uint160 weight
    );


    StateGraph internal state;

    function setUp() public {
        state = new StateGraph();
    }

    function testSetEdge() public {
        bytes12 srcPersonID = bytes12(abi.encodePacked(Kind.Person.selector, uint64(1)));
        bytes12 dstPersonID = bytes12(abi.encodePacked(Kind.Person.selector, uint64(2)));

        bytes4 relID = Rel.Friend.selector;
        uint8 relKey = 100;
        uint160 weight = 1;

        vm.expectEmit(true, true, true, true, address(state));
        emit EdgeSet(
            relID,
            relKey,
            srcPersonID,
            dstPersonID,
            weight
        );

        state.set(
            relID,
            relKey,
            srcPersonID,
            dstPersonID,
            weight
        );

        (bytes12 gotPersonID, uint160 gotWeight) = state.get(
            relID,
            relKey,
            srcPersonID
        );


        assertEq(
            gotPersonID,
            dstPersonID
        );
        assertEq(
            gotWeight,
            weight
        );
    }

    function testRegisterEdgeType() public {
        bytes4 relID = bytes4(uint32(1));
        string memory relName = "TESTING_EDGE_NAME";
        WeightKind weightKind = WeightKind.UINT64;
        vm.expectEmit(true, true, true, true, address(state));
        emit EdgeTypeRegister(
            relID,
            relName,
            weightKind
        );
        state.registerEdgeType(relID, relName, weightKind);
    }

    function testRegisterNodeType() public {
        bytes4 relID = bytes4(uint32(2));
        string memory relName = "TESTING_NODE_NAME";
        vm.expectEmit(true, true, true, true, address(state));
        emit NodeTypeRegister(
            relID,
            relName
        );
        state.registerNodeType(relID, relName);
    }

}
