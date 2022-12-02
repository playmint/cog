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

struct Attribute {
    string name;
    AttributeKind kind;
    bytes32 value;
}

struct AttributeTypeDef {
    AttributeKind kind;
}

type NodeID is uint224;

// EdgeData encodes a target NodeID the edge points to
// and a 32bit weight value. The weight value may have
// meaning to the EdgeType or it may be ignored.
struct EdgeData {
    NodeID nodeID;
    uint32 weight;
}

// NodeData is a 256bit memory slot that can be encoded and decoded
// by asking the NodeType
type NodeData is uint256;

// NodeTypeID is an address to a contract that implements NodeType
type NodeTypeID is address;

// EdgeTypeID is an address to a contract that implments EdgeType
type EdgeTypeID is address;

// NodeTypeDef is an introspection type
// struct NodeTypeDef {
//     string name; // a friendly human readable name, not used internaly
//     address id; // the NodeTypeID
//     AttributeTypeDef[] attrs; // info for decoding NodeData
// }

struct NodeMetadata {
    string typeName;
    address typeID;
    Attribute[] attrs;
}

library NodeIDUtils {
    function decodeID(NodeID id) internal pure returns (NodeType kind, uint32 s1, uint32 s2) {
        kind = NodeType(address(uint160((NodeID.unwrap(id) >> 64))));
        s1 = uint32((NodeID.unwrap(id) >> 32));
        s2 = uint32((NodeID.unwrap(id) >> 0));
        return (kind, s1, s2);
    }
    function getType(NodeID id) internal pure returns (NodeType kind) {
        return NodeType(address(uint160((NodeID.unwrap(id) >> 64))));
    }
}


interface NodeType {
    function getAttributes(NodeID, NodeData) external view returns (Attribute[] memory);
}

library NodeTypeUtils {
    function ID(NodeType kind, uint32 s1, uint32 s2) internal pure returns (NodeID) {
        return NodeID.wrap(
            uint224(s2)
            | (uint224(s1) << 32)
            | (uint224(uint160(address(kind))) << 64)
        );
    }
    function ID(NodeType kind, uint32 s2) internal pure returns (NodeID) {
        return NodeID.wrap(
            uint224(s2)
            | (uint224(0) << 32)
            | (uint224(uint160(address(kind))) << 64)
        );
    }
}


interface EdgeType {
    function getAttributes(NodeID srcNodeID, uint idx) external view returns (Attribute[] memory attrs);
}

library EdgeTypeUtils {
    function ID(EdgeType kind) internal pure returns (EdgeTypeID) {
        return EdgeTypeID.wrap(address(kind));
    }
}

interface State {

    event NodeSet(
        NodeID nodeID,
        NodeData nodeData
    );

    event EdgeSet(
        EdgeTypeID kind,
        NodeID srcNodeID,
        NodeID dstNodeID,
        uint idx,
        uint32 weight
    );

    function getNode(NodeID nodeID) external view returns (NodeData);
    function setNode(NodeID id, NodeData data) external returns (State);
    function setEdge(EdgeTypeID t, NodeID srcNodeID, uint idx, EdgeData memory data) external returns (State);
    function setEdge(EdgeTypeID t, NodeID srcNodeID, EdgeData memory data) external returns (State);
    function appendEdge(EdgeTypeID t, NodeID srcNodeID, EdgeData memory data) external returns (State);
    function getEdge(EdgeTypeID t, NodeID srcNodeID, uint idx) external view returns (EdgeData memory);
    function getEdge(EdgeTypeID t, NodeID srcNodeID) external view returns (EdgeData memory edge);
    function getEdges(EdgeTypeID t, NodeID srcNodeID) external view returns (EdgeData[] memory);
    function getNodeAttributes(NodeID, NodeData) external view returns (Attribute[] memory);
    function getEdgeAttributes(EdgeType, NodeID, uint) external view returns (Attribute[] memory);
}
