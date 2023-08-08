// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { State } from "cog/IState.sol";
import { Context, Rule } from "cog/IDispatcher.sol";
import { Schema, Node, BiomeKind, ResourceKind } from "src/schema/Schema.sol";
import { Actions, Direction } from "src/actions/Actions.sol";

using Schema for State;

contract HarvestRule is Rule {

    function reduce(State state, bytes calldata action, Context calldata /*ctx*/) public returns (State) {
        // harvesting is triggered when you move to tile with CORN on it
        // standing on a CORN tile converts the tile to a GRASS tile
        // and increases the seeker's CORN balance in their STORAGE
        if (bytes4(action) == Actions.MOVE_SEEKER.selector) {
            (uint32 sid,) = abi.decode(action[4:], (uint32, Direction));
            bytes24 seeker = Node.Seeker(sid);

            (uint32 x, uint32 y) = state.getLocationCoords(seeker);
            bytes24 targetTile = Node.Tile(x,y);
            BiomeKind biome = state.getBiome(targetTile);
            if (biome == BiomeKind.CORN) {
                // convert tile to grass
                state.setBiome(targetTile, BiomeKind.GRASS);
                // get seeker's current balance of corn
                uint32 balance = state.getResourceBalance(Node.Seeker(sid), ResourceKind.CORN);
                // increase the balance
                balance++;
                // store new balance
                state.setResourceBalance(Node.Seeker(sid), ResourceKind.CORN, balance);
            }

        }
        return state;
    }

}

