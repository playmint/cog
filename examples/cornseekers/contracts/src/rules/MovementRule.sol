// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { State } from "cog/IState.sol";
import { Context, Rule } from "cog/IDispatcher.sol";
import { Schema, Node, BiomeKind } from "src/schema/Schema.sol";
import { Actions, Direction } from "src/actions/Actions.sol";
import { Context, Rule } from "cog/IDispatcher.sol";

import { Actions, Direction } from "src/actions/Actions.sol";

using Schema for State;

contract MovementRule is Rule {

    function reduce(State state, bytes calldata action, Context calldata /*ctx*/) public returns (State) {


        // movement is one tile at a time
        // you can only move onto an discovered tile
        if (bytes4(action) == Actions.MOVE_SEEKER.selector) {
            // decode the action
            (uint32 sid, Direction dir) = abi.decode(action[4:], (uint32, Direction));

            // encode the full seeker node id
            bytes24 seeker = Node.Seeker(sid);

            // fetch the seeker's current location
            (uint32 x, uint32 y) = state.getLocationCoords(seeker);

            // find new location
            bytes24 targetTile = getTargetLocation(state, x, y, dir);
            state.setAnnotation(targetTile, "tag", "MY_ANNOTATION");

            // update seeker location
            state.setLocation(seeker, targetTile);
        }

        return state;
    }

    function getTargetLocation(State state, uint32 x, uint32 y, Direction dir) internal view returns (bytes24) {
        int xx = int(uint(x));
        int yy = int(uint(y));
        if (dir == Direction.NORTH) {
            yy++;
        } else if (dir == Direction.NORTHEAST) {
            xx++;
            yy++;
        } else if (dir == Direction.EAST) {
            xx++;
        } else if (dir == Direction.SOUTHEAST) {
            xx++;
            yy--;
        } else if (dir == Direction.SOUTH) {
            yy--;
        } else if (dir == Direction.SOUTHWEST) {
            xx--;
            yy--;
        } else if (dir == Direction.WEST) {
            xx--;
        } else if (dir == Direction.NORTHWEST) {
            xx--;
            yy--;
        }
        if (xx<0) {
            xx = 0;
        } else if (xx>31) {
            xx = 31;
        }
        if (yy<0) {
            yy = 0;
        } else if (yy>31) {
            yy = 31;
        }

        // check where we are moving to is legit
        bytes24 targetTile = Node.Tile(uint32(uint(xx)),uint32(uint(yy)));
        BiomeKind biome = state.getBiome(targetTile);
        if (biome == BiomeKind.UNDISCOVERED || biome == BiomeKind.BLOCKER) {
            // illegal move, just return original tile
            // revert(string(abi.encodePacked("tile x=",x+48, " y=", y+48, " is UNDISCOVERED")));
            return Node.Tile(x, y);
        }
        return targetTile;
    }

}

