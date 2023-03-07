// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

// WeightKind is used to hint to the indexer what kind of value
// you are storing in an edge's weight field.
enum WeightKind {
    UINT64,
    INT64,
    BYTES,
    STRING
}

// CompoundKeyKind is a hint to how to decode the last bytes of a node ID.
// This hint allows the indexer to split out the key parts of the id
enum CompoundKeyKind {
    NONE, // key is not expected to be anything other than 0
    UINT160, // key is a single uint64
    UINT8_ARRAY, // key is 20 uint8s
    INT8_ARRAY, // key is 20 int8s
    UINT16_ARRAY, // key is 10 uint16s
    INT16_ARRAY, // key is 10 int16s
    UINT32_ARRAY, // key is 5 uint32s
    INT32_ARRAY, // key is 5 int32s
    UINT64_ARRAY, // key is 2 uint64s
    INT64_ARRAY, // key is 2 int64s
    ADDRESS, // key is 20 byte address
    BYTES, // key is an 20 byte blob of data
    STRING // key is an 20 byte string
}

enum AnnotationKind {
    CALLDATA
}

library CompoundKeyEncoder {
    function UINT64(bytes4 kindID, uint64 key) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, uint96(0), key));
    }

    function BYTES(bytes4 kindID, bytes20 key) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, key));
    }

    function UINT8_ARRAY(bytes4 kindID, uint8[8] memory keys) internal pure returns (bytes24) {
        return bytes24(
            abi.encodePacked(kindID, uint96(0), keys[0], keys[1], keys[2], keys[3], keys[4], keys[5], keys[6], keys[7])
        );
    }

    function UINT16_ARRAY(bytes4 kindID, uint16[4] memory keys) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, uint96(0), keys[0], keys[1], keys[2], keys[3]));
    }

    function INT16_ARRAY(bytes4 kindID, int16[4] memory keys) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, uint96(0), keys[0], keys[1], keys[2], keys[3]));
    }

    function UINT32_ARRAY(bytes4 kindID, uint32[2] memory keys) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, uint96(0), keys[0], keys[1]));
    }

    function INT32_ARRAY(bytes4 kindID, int32[2] memory keys) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, uint96(0), keys[0], keys[1]));
    }

    function ADDRESS(bytes4 kindID, address addr) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, uint160(addr)));
    }

    function STRING(bytes4 kindID, string memory id) internal pure returns (bytes24) {
        return bytes24(abi.encodePacked(kindID, id));
    }
}

library CompoundKeyDecoder {
    function UINT64(bytes24 id) internal pure returns (uint64) {
        return uint64(uint192(id));
    }

    function BYTES8(bytes24 id) internal pure returns (bytes8) {
        return bytes8(uint64(uint192(id)));
    }

    function UINT8_ARRAY(bytes24 id) internal pure returns (uint8[8] memory keys) {
        keys[0] = uint8(uint192(id) >> 56);
        keys[1] = uint8(uint192(id) >> 48);
        keys[2] = uint8(uint192(id) >> 40);
        keys[3] = uint8(uint192(id) >> 32);
        keys[4] = uint8(uint192(id) >> 24);
        keys[5] = uint8(uint192(id) >> 16);
        keys[6] = uint8(uint192(id) >> 8);
        keys[7] = uint8(uint192(id));
    }

    function UINT16_ARRAY(bytes24 id) internal pure returns (uint16[4] memory keys) {
        keys[0] = uint16(uint192(id) >> 48);
        keys[1] = uint16(uint192(id) >> 32);
        keys[2] = uint16(uint192(id) >> 16);
        keys[3] = uint16(uint192(id));
    }

    function INT16_ARRAY(bytes24 id) internal pure returns (int16[4] memory keys) {
        keys[0] = int16(int192(uint192(id) >> 48));
        keys[1] = int16(int192(uint192(id) >> 32));
        keys[2] = int16(int192(uint192(id) >> 16));
        keys[3] = int16(int192(uint192(id)));
    }

    function UINT32_ARRAY(bytes24 id) internal pure returns (uint32[2] memory keys) {
        keys[0] = uint32(uint192(id) >> 32);
        keys[1] = uint32(uint192(id));
    }

    function INT32_ARRAY(bytes24 id) internal pure returns (int32[2] memory keys) {
        keys[0] = int32(int192(uint192(id) >> 32));
        keys[1] = int32(int192(uint192(id)));
    }

    function ADDRESS(bytes24 id) internal pure returns (address) {
        return address(uint160(uint192(id)));
    }

    function STRING(bytes24 id) internal pure returns (string memory) {
        // Find string length. Keys are fixed at 20 bytes so treat first 0 as null terminator
        uint8 len;
        while (len < 20 && id[4 + len] != 0) {
            len++;
        }

        // Copy string bytes
        bytes memory stringBytes = new bytes(len);
        for (uint8 i = 0; i < len; i++) {
            stringBytes[i] = id[4 + i];
        }

        return string(stringBytes);
    }
}

interface State {
    event EdgeTypeRegister(bytes4 id, string name, WeightKind kind);
    event NodeTypeRegister(bytes4 id, string name, CompoundKeyKind keyKind);
    event EdgeSet(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint160 weight);
    event EdgeRemove(bytes4 relID, uint8 relKey, bytes24 srcNodeID);
    event AnnotationSet(bytes24 id, AnnotationKind kind, string label, bytes32 ref);

    function set(bytes4 relID, uint8 relKey, bytes24 srcNodeID, bytes24 dstNodeID, uint64 weight) external;
    function remove(bytes4 relID, uint8 relKey, bytes24 srcNodeID) external;
    function get(bytes4 relID, uint8 relKey, bytes24 srcNodeID)
        external
        view
        returns (bytes24 dstNodeId, uint64 weight);

    function registerNodeType(bytes4 kindID, string memory kindName, CompoundKeyKind keyKind) external;
    function registerEdgeType(bytes4 relID, string memory relName, WeightKind weightKind) external;
    function authorizeContract(address addr) external;

    // an annotation is an on-chain tag that points to a (potentially large)
    // string stored in transaction calldata. The content of annotations are
    // not available on-chain, only a content addressable reference to it.
    // indexers/clients are expected to watch for AnnotationSet event and store
    // the annotationData for later lookup
    function setAnnotation(bytes24 nodeID, string memory label, string memory annotationData) external;
    function getAnnotationRef(bytes24 nodeID, string memory label) external returns (bytes32 annotationRef);
}
