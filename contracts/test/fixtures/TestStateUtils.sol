// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {State} from "src/IState.sol";

// some wrappers to treat State as a single value that you can set or get

interface Rel {
    function HasValue() external;
}

interface Kind {
    function TheValue() external;
}

library StateTestUtils {
    function setUint(State s, uint64 value) internal {
        bytes24 valueNode = bytes24(abi.encodePacked(Kind.TheValue.selector, uint96(0), value));
        return s.set(Rel.HasValue.selector, 0x0, bytes24(0), valueNode, 0);
    }

    function setAddress(State s, address value) internal {
        bytes24 valueNode = bytes24(abi.encodePacked(Kind.TheValue.selector, value));
        return s.set(Rel.HasValue.selector, 0x0, bytes24(0), valueNode, 0);
    }

    function setBytes(State s, bytes20 value) internal {
        bytes24 valueNode = bytes24(abi.encodePacked(Kind.TheValue.selector, value));
        return s.set(Rel.HasValue.selector, 0x0, bytes24(0), valueNode, 0);
    }

    function getUint(State s) internal view returns (uint64) {
        (bytes24 valueNode,) = s.get(Rel.HasValue.selector, 0x0, 0x0);
        return uint64(uint192(valueNode));
    }

    function getAddress(State s) internal view returns (address) {
        (bytes24 valueNode,) = s.get(Rel.HasValue.selector, 0x0, 0x0);
        return address(uint160(uint192(valueNode)));
    }

    function getBytes(State s) internal view returns (bytes20) {
        (bytes24 valueNode,) = s.get(Rel.HasValue.selector, 0x0, 0x0);
        return bytes20(uint160(uint192(valueNode)));
    }
}
