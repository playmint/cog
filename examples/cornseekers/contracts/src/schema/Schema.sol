// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {
    State,
    CompoundKeyEncoder,
    CompoundKeyDecoder
} from "cog/IState.sol";

interface Rel {
    function Owner() external;
    function Location() external;
    function Balance() external;
    function Biome() external;
    function Strength() external;
    function ProvidesEntropyTo() external;
}

interface Kind {
    function Seed() external;
    function Tile() external;
    function Resource() external;
    function Seeker() external;
    function Player() external;
}

enum ResourceKind {
    UNKNOWN,
    CORN
}

enum BiomeKind {
    UNDISCOVERED,
    BLOCKER,
    GRASS,
    CORN
}

library Node {
    function Seeker(uint64 id) internal pure returns (bytes24) {
        return CompoundKeyEncoder.UINT64(Kind.Seeker.selector, id);
    }
    function Tile(uint32 x, uint32 y) internal pure returns (bytes24) {
        return CompoundKeyEncoder.UINT32_ARRAY(Kind.Tile.selector, [x, y]);
    }
    function Resource(ResourceKind rk) internal pure returns (bytes24) {
        return CompoundKeyEncoder.UINT64(Kind.Resource.selector, uint64(rk));
    }
    function Seed(uint32 blk) internal pure returns (bytes24) {
        return CompoundKeyEncoder.UINT64(Kind.Seed.selector, blk);
    }
    function Player(address addr) internal pure returns (bytes24) {
        return CompoundKeyEncoder.ADDRESS(Kind.Player.selector, addr);
    }
}

using Schema for State;

library Schema {

    function setLocation(State state, bytes24 node, bytes24 locationNode) internal {
        return state.set(Rel.Location.selector, 0x0, node, locationNode, 0);
    }

    function getLocation(State state, bytes24 node) internal view returns (bytes24) {
        (bytes24 tile,) = state.get(Rel.Location.selector, 0x0, node);
        return tile;
    }

    function getLocationCoords(State state, bytes24 node) internal view returns (uint32 x, uint32 y) {
        bytes24 tile = getLocation(state, node);
        uint32[2] memory keys = CompoundKeyDecoder.UINT32_ARRAY(tile);
        return (keys[0], keys[1]);
    }

    function setBiome(State state, bytes24 node, BiomeKind biome) internal {
        return state.set(Rel.Biome.selector, 0x0, node, 0x0, uint64(biome));
    }

    function getBiome(State state, bytes24 node) internal view returns (BiomeKind) {
        (,uint160 biome) = state.get(Rel.Biome.selector, 0x0, node);
        return BiomeKind(uint8(biome));
    }

    function setResourceBalance(State state, bytes24 node, ResourceKind rk, uint64 balance) internal {
        return state.set(Rel.Balance.selector, uint8(rk), node, Node.Resource(rk), balance);
    }

    function getResourceBalance(State state, bytes24 node, ResourceKind rk) internal view returns (uint32) {
        (,uint160 balance) = state.get(Rel.Balance.selector, uint8(rk), node);
        return uint32(balance);
    }

    function setOwner(State state, bytes24 node, bytes24 ownerNode) internal {
        return state.set(Rel.Owner.selector, 0x0, node, ownerNode, 0);
    }

    function getOwner(State state, bytes24 node) internal view returns (bytes24, uint160) {
        return state.get(Rel.Owner.selector, 0x0, node);
    }

    function getOwnerAddress(State state, bytes24 ownerNode) internal view returns (address) {
        while (bytes4(ownerNode) != Kind.Player.selector) {
            (ownerNode,) = state.getOwner(ownerNode);
        }
        return address(uint160(uint192(ownerNode)));
    }

    function setStrength(State state, bytes24 node, uint64 v) internal {
        return state.set(Rel.Strength.selector, 0x0, node, 0x0, v);
    }

    function getStrength(State state, bytes24 node) internal view returns (uint64) {
        (,uint64 str) = state.get(Rel.Strength.selector, 0x0, node);
        return uint64(str);
    }

    function setEntropyCommitment(State state, uint32 blk, bytes24 node) internal {
        // we will treat the key as an idx and iterate to find a free slot
        // this is not a very effceient solution, but is a direct port from
        // how it worked before with appendEdge
        for (uint8 key=0; key<256; key++) {
            (bytes24 dstNodeID,) = state.get(Rel.ProvidesEntropyTo.selector, key, Node.Seed(blk));
            if (dstNodeID == bytes24(0)) {
                return state.set(Rel.ProvidesEntropyTo.selector, key, Node.Seed(blk), node, 0);
            }
        }
        revert("too many edges");
    }

    function getEntropyCommitments(State state, uint32 blk) internal view returns (bytes24[] memory) {
        // we will treat the key as an idx and iterate to find a free slot
        // this is not a very effceient solution, but is a direct port from
        // how it worked before with appendEdge
        bytes24[100] memory foundNodes;
        uint8 i;
        for (i=0; i<255; i++) {
            (foundNodes[i],) = state.get(Rel.ProvidesEntropyTo.selector, i, Node.Seed(blk));
            if (foundNodes[i] == bytes24(0)) {
                break;
            }
        }
        bytes24[] memory nodes = new bytes24[](i);
        for (uint8 j=0; j<i; j++) {
            nodes[j] = foundNodes[j];
        }
        return nodes;
    }

}
