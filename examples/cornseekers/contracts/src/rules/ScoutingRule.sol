// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { State } from "cog/IState.sol";
import { Context, Rule } from "cog/IDispatcher.sol";
import { Actions, Direction } from "src/actions/Actions.sol";
import { Schema, Node, BiomeKind, ResourceKind } from "src/schema/Schema.sol";

using Schema for State;

contract ScoutingRule is Rule {

    function reduce(State state, bytes calldata action, Context calldata ctx) public returns (State) {
        // scouting tiles is performed in two stages
        // stage1: we commit to a SEED during a MOVE_SEEKER or SPAWN_SEEKER action
        // stage2: occurs when a REVEAL_SEED action is processed
        if (bytes4(action) == Actions.SPAWN_SEEKER.selector) {

            (, uint8 x, uint8 y,) = abi.decode(action[4:], (uint32, uint8, uint8, uint8));
            state = commitAdjacent(state, ctx, int(uint(x)), int(uint(y)));

        } else if (bytes4(action) == Actions.MOVE_SEEKER.selector) {

            (uint32 sid,) = abi.decode(action[4:], (uint32, Direction));
            (uint32 x, uint32 y) = state.getLocationCoords( Node.Seeker(sid) );
            state = commitAdjacent(state, ctx, int(uint(x)), int(uint(y)));

        } else if (bytes4(action) == Actions.REVEAL_SEED.selector) {

            (uint32 blk, uint32 entropy) = abi.decode(action[4:], (uint32, uint32));
            state = revealTiles(state, blk, entropy);

        }
        return state;
    }

    function commitAdjacent(State state, Context calldata ctx, int x, int y) private returns (State) {
        int xx;
        int yy;
        for (uint8 i=0; i<8; i++) {
            if (i == 0) {
                xx = x-1;
                yy = y+1;
            } else if (i == 1) {
                xx = x;
                yy = y+1;
            } else if (i == 2) {
                xx = x+1;
                yy = y+1;
            } else if (i == 3) {
                xx = x+1;
                yy = y;
            } else if (i == 4) {
                xx = x-1;
                yy = y-1;
            } else if (i == 5) {
                xx = x;
                yy = y-1;
            } else if (i == 6) {
                xx = x+1;
                yy = y-1;
            } else if (i == 7) {
                xx = x-1;
                yy = y;
            }

            if (xx < 0 || yy < 0) {
                continue;
            }
            if (xx > 31 || yy > 31) {
                continue;
            }
            bytes24 tile = Node.Tile(uint32(uint(xx)),uint32(uint(yy)));
            BiomeKind biome = state.getBiome(tile);
            if (biome == BiomeKind.UNDISCOVERED) {
                // commit tile contents to a future SEED
                state.setEntropyCommitment(ctx.clock, tile);
            }
        }
        return state;
    }

    function revealTiles(State state, uint32 blk, uint32 entropy) private returns (State) {
        bytes24[] memory targetTiles = state.getEntropyCommitments(blk);
        for (uint i=0; i<targetTiles.length; i++) {
            (uint32 x, uint32 y) = state.getLocationCoords(targetTiles[i]);
            BiomeKind biome = state.getBiome(targetTiles[i]);
            if (biome != BiomeKind.UNDISCOVERED) {
                continue;
            }
            uint8 r = random(entropy, x, y, i);
            if (r > 220) {
                biome = BiomeKind.CORN;
            } else if (r > 50) {
                biome = BiomeKind.GRASS;
            } else {
                biome = BiomeKind.BLOCKER;
            }
            state.setBiome( targetTiles[i], biome );
        }
        return state;
    }

    function random(uint32 entropy, uint32 x, uint32 y, uint i) public pure returns(uint8){
        return uint8(uint( keccak256(abi.encodePacked(x, y, entropy, i)) ) % 255);
    }

}

