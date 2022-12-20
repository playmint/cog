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

library CompoundKeyEncoder {
    function UINT64(bytes4 kindID, uint64 key) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, key));
    }
    function BYTES8(bytes4 kindID, bytes8 key) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, key));
    }
    function UINT8_ARRAY(bytes4 kindID, uint8[8] memory keys) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, keys[0], keys[1], keys[2], keys[3], keys[4], keys[5], keys[6], keys[7]));
    }
    function UINT16_ARRAY(bytes4 kindID, uint16[4] memory keys) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, keys[0], keys[1], keys[2], keys[3]));
    }
    function UINT32_ARRAY(bytes4 kindID, uint32[2] memory keys) internal pure returns (bytes12) {
        return bytes12(abi.encodePacked(kindID, keys[0], keys[1]));
    }
}

library CompoundKeyDecoder {
    function UINT64(bytes12 id) internal pure returns (uint64) {
        return uint64(uint96(id));
    }
    function BYTES8(bytes12 id) internal pure returns (bytes8) {
        return bytes8(uint64(uint96(id)));
    }
    function UINT8_ARRAY(bytes12 id) internal pure returns (uint8[8] memory keys) {
        keys[0] = uint8(uint96(id) >> 56);
        keys[1] = uint8(uint96(id) >> 48);
        keys[2] = uint8(uint96(id) >> 40);
        keys[3] = uint8(uint96(id) >> 32);
        keys[4] = uint8(uint96(id) >> 24);
        keys[5] = uint8(uint96(id) >> 16);
        keys[6] = uint8(uint96(id) >> 8);
        keys[7] = uint8(uint96(id));
    }
    function UINT16_ARRAY(bytes12 id) internal pure returns (uint16[4] memory keys) {
        keys[0] = uint8(uint96(id) >> 48);
        keys[1] = uint8(uint96(id) >> 32);
        keys[2] = uint8(uint96(id) >> 16);
        keys[3] = uint8(uint96(id));
    }
    function UINT32_ARRAY(bytes12 id) internal pure returns (uint32[2] memory keys) {
        keys[0] = uint8(uint96(id) >> 32);
        keys[1] = uint8(uint96(id));
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
