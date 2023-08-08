// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { State } from "cog/IState.sol";
import { Context, Rule } from "cog/IDispatcher.sol";
import { Actions } from "src/actions/Actions.sol";
import { Schema, Node, BiomeKind } from "src/schema/Schema.sol";

using Schema for State;

contract ResetRule is Rule {

    function reduce(State state, bytes calldata action, Context calldata /*ctx*/) public returns (State) {
        if (bytes4(action) == Actions.RESET_MAP.selector) {
            // draw a grid of tiles encoding the x/y into the ID
            for (uint8 x=0; x<32; x++) {
                for (uint8 y=0; y<32; y++) {
                    state.setBiome( Node.Tile(x,y), getInitialBiome(x,y) );
                }
            }
        }
        return state;
    }

    function getInitialBiome(uint8 x, uint8 y) private pure returns (BiomeKind) {
        if (x == 0 || y == 0 || x == 31 || y == 31) { // grass around the edge
            return BiomeKind.GRASS;
        } else { // everything else unknown
            return BiomeKind.UNDISCOVERED;
        }
    }

}
