// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

// WeightKind is used to hint to the indexer what kind of value
// you are storing in an edge's weight field.
enum WeightKind {
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

// Builtin node types are well-known nodes that may be useful for referencing.
library Builtin {
    // The Null node type for example points to the zero value, and may be used
    // if you want a "dangling" edge that does not actually "point" anywhere
    // particular. Most of the time this kind of dangling edge is probably an indication
    // that the model is not quite right, but sometimes when you are treating
    // edges like properties, components or labels pointing at the null node is
    // reasonable.
    function Null() internal pure returns (bytes12) {
        return bytes12(0);
    }
}

interface State {

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

    function set(bytes4 relID, uint8 relKey, bytes12 srcNodeID, bytes12 dstNodeID, uint160 weight) external;
    function get(bytes4 relID, uint8 relKey, bytes12 srcNodeID) external view returns (bytes12 dstNodeId, uint160 weight);

    function registerNodeType(bytes4 kindID, string memory kindName) external;
    function registerEdgeType(bytes4 relID, string memory relName, WeightKind weightKind) external;
    function authorizeContract(address addr) external;
}
