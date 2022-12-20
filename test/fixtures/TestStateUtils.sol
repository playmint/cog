// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { State } from "src/State.sol";

// some wrappers to treat State as a single value that you can set or get

interface Rel {
    function HasValue() external;
}
interface Kind {
    function TheValue() external;
}

library StateTestUtils {
    function valueNode() internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(Kind.TheValue.selector, uint64(1)));
    }
    function setUint(State s, uint160 value) internal {
        return s.set(Rel.HasValue.selector, 0x0, valueNode(), bytes12(0), uint160(value));
    }
    function setAddress(State s, address value) internal {
        return s.set(Rel.HasValue.selector, 0x0, valueNode(), bytes12(0), uint160(value));
    }
    function setBytes(State s, bytes20 value) internal {
        return s.set(Rel.HasValue.selector, 0x0, valueNode(), bytes12(0), uint160(value));
    }
    function getUint(State s) internal view returns (uint160) {
        (, uint160 value) = s.get(Rel.HasValue.selector, 0x0, valueNode());
        return value;
    }
    function getAddress(State s) internal view returns (address) {
        (, uint160 value) = s.get(Rel.HasValue.selector, 0x0, valueNode());
        return address(value);
    }
    function getBytes(State s) internal view returns (bytes20) {
        (, uint160 value) = s.get(Rel.HasValue.selector, 0x0, valueNode());
        return bytes20(value);
    }
}
