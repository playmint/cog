// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { State } from "cog/IState.sol";
import { Context, Rule } from "cog/IDispatcher.sol";
import { Actions } from "src/actions/Actions.sol";
import { Schema, Node, ResourceKind } from "src/schema/Schema.sol";

using Schema for State;

contract SpawnSeekerRule is Rule {

    function reduce(State state, bytes calldata action, Context calldata ctx) public returns (State) {
        if (bytes4(action) == Actions.SPAWN_SEEKER.selector) {

            // decode action
            (uint32 sid, uint8 x, uint8 y, uint8 str) = abi.decode(action[4:], (uint32, uint8, uint8, uint8));

            // build full seeker node id
            bytes24 seeker = Node.Seeker(sid);

            // set a strength attr for no reason
            state.setStrength( seeker, str );

            // set the seeker's owner relationship to the action sender
            state.setOwner( seeker, Node.Player(ctx.sender) );

            // set location by pointing a location relationship at the tile
            state.setLocation( seeker, Node.Tile(x,y) );

            // start with 0 corn
            state.setResourceBalance(Node.Seeker(sid), ResourceKind.CORN, 0);
        }
        return state;
    }

}

