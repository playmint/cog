// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {
    State,
    EdgeData,
    Attr
} from "../src/State.sol";
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
        string name
    );
    event NodeTypeRegister(
        bytes4 id,
        string name
    );
    event EdgeSet(
        bytes4 relID,
        bytes8 relKey,
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
        bytes8 relKey = bytes8("best");
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

}
