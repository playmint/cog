


## Concepts

COG Games are structured around Actions (thing a player commits to doing) and
Rules (the state modifications that occur at some point as a result of that
Action).

The pattern has it's roots in the [Command Pattern](https://en.wikipedia.org/wiki/Command_pattern), [CQRS](https://en.wikipedia.org/wiki/Command%E2%80%93query_separation) and [event-sourced](https://www.confluent.io/blog/event-sourcing-cqrs-stream-processing-apache-kafka-whats-connection/)
systems, and shares architectual simularities to state management frameworks
like [Flux or Redux](https://redux.js.org/tutorials/essentials/part-1-overview-concepts).

### Actions

Actions are the user _intent_.

Conceptually they are the log of user intent, practically they a just an immutable id and definition of arguments (like a function signature).

We define them in solidity as a contract that implements the `ActionType` interface:

```solidity
contract SpawnCharacter is ActionType {

    function getTypeDef() external view returns (ActionTypeDef memory def) {
        def.name = "SPAWN_CHAPACTER";
        def.id = address(this);
        def.arg0.name = "characterID";
        def.arg0.required = true;
        def.arg0.kind = ActionArgKind.NODEID;
    }

}
```

### Rules

Rules apply `Actions` to modify `State`.

The logic of what state modifications to make for any given `Action` is through a pure-like `reduce` function. Many Rules may operate on many Actions.

We define then in solidity as a contract that implements the `Rule` interface:

```solidity
contract SpawnCharacterRule is Rule {

    function reduce(State state, Action memory action) public returns (State) {
        if (action.id == address(SPAWN_CHAPACTER)) {
            // do some logic here
            state = state.set(....)
        }
        return state;
    }

}

```

### Dispatchers

Dispatchers are the entrypoint to games. They define the set of `Rules` that will be applied in order, the State to use, and the `dispatch` function entrypoint to apply incoming `Actions`.

We define them in solidity as a contract that implements the `Dispatcher` interface. A `BaseDispatcher` abstract contract can be used to avoid some boilerplate:

```solidity
contract MyGame is BaseDispatcher {

    constructor(State s, Rule[] memory rs) BaseDispatcher(s) {
        for (uint i=0; i<rs.length; i++) {
            registerRule(rs[i]);
        }
    }

    function dispatch(Action memory action) public {
        // do any custom action validation/authorization here
        // ...

        // call _dispatch to send action through the registered rules
        _dispatch(action);
    }

}
```

### State

To ensure that there is a common way for Rules to modify the game state, and ensure we have a consistent way to index that state we perform all state modifications through a common `State` type. Think of it like a really simple database or key/value store.

State is a modelled as sparse directed graph of Nodes, Attributes and Edges. `State` is an interface:

```solidity
interface State {
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
```

`StateGraph` is an implementation of that interface that aims to minimise the amount of storage writes by limiting the amount of attribute storage per NodeType to a single uint256 slot.

`Nodes` in the State are basically just an immutable ID. But we define them in solidity as something that implements the `NodeType` interface which describes how to decode the node's attribute data:

```solidity
contract Character is NodeType {
    function getAttributes(NodeID id, NodeData data) public pure returns (Attribute[] memory attrs) {
        (uint8 str) = getAttributeValues(id, data);
        attrs = new Attribute[](3);
        attrs[0].name = "kind";
        attrs[0].kind = AttributeKind.STRING;
        attrs[0].value = bytes32("CHAPACTER");
        attrs[1].name = "strength";
        attrs[1].kind = AttributeKind.UINT8;
        attrs[1].value = bytes32(uint(str));
    }
}
```

Edges are weighted directional relationships between Nodes.

We define them in solidity as a contract that implements the `EdgeType` interface:

```solidity
contract HasLocation is EdgeType {
    function getAttributes(NodeID /*id*/, uint /*idx*/) public pure returns (Attribute[] memory attrs) {
        attrs = new Attribute[](1);
        attrs[0].name = "kind";
        attrs[0].kind = AttributeKind.STRING;
        attrs[0].value = bytes32("HAS_LOCATION");
    }
}
```

