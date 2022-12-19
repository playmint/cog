// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

enum AttributeKind {
    BOOL,
    INT8,
    INT16,
    INT32,
    INT64,
    INT128,
    INT256,
    INT,
    UINT8,
    UINT16,
    UINT32,
    UINT64,
    UINT128,
    UINT256,
    BYTES,
    STRING,
    ADDRESS,
    BYTES4,
    BYTES12,
    BYTES20,
    BOOL_ARRAY,
    INT8_ARRAY,
    INT16_ARRAY,
    INT32_ARRAY,
    INT64_ARRAY,
    INT128_ARRAY,
    INT256_ARRAY,
    INT_ARRAY,
    UINT8_ARRAY,
    UINT16_ARRAY,
    UINT32_ARRAY,
    UINT64_ARRAY,
    UINT128_ARRAY,
    UINT256_ARRAY,
    BYTES_ARRAY,
    STRING_ARRAY
}

struct Edge {
    bytes4 rel;
    bytes8 key;
    bytes12 dst;
    uint160 val;
}

// builtin Attr node types. the Attr node types are special in that
// they have well-known ids mapping to simple scaler types that are
// understood by indexing services. Combined with edges these are used
// to assign property-style reltionships
library Attr {
    function Int() internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(bytes4(0xbbf0aac1), uint64(AttributeKind.INT32)));
    }
    function UInt() internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(bytes4(0xbbf0aac1), uint64(AttributeKind.UINT32)));
    }
    function Address() internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(bytes4(0xbbf0aac1), uint64(AttributeKind.ADDRESS)));
    }
    function Bytes() internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(bytes4(0xbbf0aac1), uint64(AttributeKind.BYTES)));
    }
}

struct EdgeData {
    bytes12 dstNodeID;
    uint160 weight;
}

interface State {

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

    function set(bytes4 relID, bytes8 relKey, bytes12 srcNodeID, bytes12 dstNodeID, uint160 weight) external;
    function get(bytes4 relID, bytes8 relKey, bytes12 srcNodeID) external view returns (bytes12 dstNodeId, uint160 weight);

    function registerNodeType(bytes4 kindID, string memory kindName) external;
    function registerEdgeType(bytes4 relID, string memory relName) external;
    function authorizeContract(address addr) external;
}
