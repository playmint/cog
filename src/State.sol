// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

// WeightKind is used to hint to the indexer what kind of value
// you are storing in an edge's weight field.
enum WeightKind {
    UINT160,
    INT160,
    ADDRESS,
    BYTES,
    STRING
}

// CompoundKeyKind is a hint to how to decode the last 8bytes of a node ID.
// This hint allows the indexer to split out the key parts of the id in
// useful way.
enum CompoundKeyKind {
    NONE,         // key is not expected to be anything other than 0
    UINT64,       // key is a single uint64
    BYTES8,       // key is an 8 byte blob
    UINT8_ARRAY,  // key is 8 uint8s
    UINT16_ARRAY, // key is 4 uint16s
    UINT32_ARRAY, // key is 2 uint32s
    BYTES4_ARRAY, // key is 2 4byte blobs
    STRING        // key is an 8byte string
}

library CompoundKey {
    function UINT64(bytes4 kindID, uint64 key) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, key));
    }
    function BYTES8(bytes4 kindID, bytes8 key) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, key));
    }
    function UINT8_ARRAY(bytes4 kindID, uint8[8] memory keys) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, keys));
    }
    function UINT16_ARRAY(bytes4 kindID, uint16[4] memory keys) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, keys));
    }
    function UINT32_ARRAY(bytes4 kindID, uint32[2] memory keys) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, keys));
    }
    function BYTES4_ARRAY(bytes4 kindID, bytes4[2] memory keys) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, keys));
    }
    function STRING(bytes4 kindID, string memory key) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, key));
    }
    // The NULL node type for example points to the zero value, and may be used
    // if you want a "dangling" edge that does not actually "point" anywhere
    // particular. Most of the time this kind of dangling edge is probably an indication
    // that the model is not quite right, but sometimes when you are treating
    // edges like properties, components or labels pointing at the null node is
    // reasonable.
    function NULL() internal pure returns (bytes12) {
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
        string name,
        CompoundKeyKind keyKind
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

    function registerNodeType(bytes4 kindID, string memory kindName, CompoundKeyKind keyKind) external;
    function registerEdgeType(bytes4 relID, string memory relName, WeightKind weightKind) external;
    function authorizeContract(address addr) external;
}
