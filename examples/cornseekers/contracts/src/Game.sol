// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import { BaseGame } from "cog/IGame.sol";
import { BaseDispatcher } from "cog/IDispatcher.sol";
import { BaseRouter } from "cog/BaseRouter.sol";
import { BaseState, CompoundKeyKind, WeightKind } from "cog/BaseState.sol";

import { HarvestRule } from "src/rules/HarvestRule.sol";
import { ScoutingRule } from "src/rules/ScoutingRule.sol";
import { ResetRule } from "src/rules/ResetRule.sol";
import { MovementRule } from "src/rules/MovementRule.sol";
import { SpawnSeekerRule } from "src/rules/SpawnSeekerRule.sol";

import { BaseRouter } from "cog/BaseRouter.sol";
import { Actions } from "src/actions/Actions.sol";
import { Rel, Kind } from "src/schema/Schema.sol";

// -----------------------------------------------
// a Game sets up the State, Dispatcher and Router
//
// it sets up the rules our game uses and exposes
// the Game interface for discovery by cog-services
//
// we are using BasicGame to handle the boilerplate
// so all we need to do here is call registerRule()
// -----------------------------------------------

contract Game is BaseGame {

    constructor() BaseGame("CORNSEEKERS", "http://example.com") {
        // create a state
        BaseState state = new BaseState();

        // register the kind ids we are using
        state.registerNodeType(Kind.Seed.selector, "Seed", CompoundKeyKind.UINT160);
        state.registerNodeType(Kind.Tile.selector, "Tile", CompoundKeyKind.UINT32_ARRAY);
        state.registerNodeType(Kind.Resource.selector, "Resource", CompoundKeyKind.UINT160);
        state.registerNodeType(Kind.Seeker.selector, "Seeker", CompoundKeyKind.UINT160);
        state.registerNodeType(Kind.Player.selector, "Player", CompoundKeyKind.ADDRESS);

        // register the relationship ids we are using
        state.registerEdgeType(Rel.Owner.selector, "Owner", WeightKind.UINT64);
        state.registerEdgeType(Rel.Location.selector, "Location", WeightKind.UINT64);
        state.registerEdgeType(Rel.Balance.selector, "Balance", WeightKind.UINT64);
        state.registerEdgeType(Rel.Biome.selector, "Biome", WeightKind.UINT64);
        state.registerEdgeType(Rel.Strength.selector, "Strength", WeightKind.UINT64);
        state.registerEdgeType(Rel.ProvidesEntropyTo.selector, "ProvidesEntropyTo", WeightKind.UINT64);

        // create a session router
        BaseRouter router = new BaseRouter();

        // configure our dispatcher with state, rules and trust the router
        BaseDispatcher dispatcher = new BaseDispatcher();
        dispatcher.registerState(state);
        dispatcher.registerRule(new ResetRule());
        dispatcher.registerRule(new SpawnSeekerRule());
        dispatcher.registerRule(new MovementRule());
        dispatcher.registerRule(new ScoutingRule());
        dispatcher.registerRule(new HarvestRule());
        dispatcher.registerRouter(router);

        // update the game with this config
        _registerState(state);
        _registerRouter(router);
        _registerDispatcher(dispatcher);

        // playing...
        // TODO: REMOVE THESE - I'm just playing with the services
        // dispatcher.dispatch(
        //     abi.encodeCall(Actions.RESET_MAP, ())
        // );
        // BaseRouter(address(router)).authorizeAddr(dispatcher, 0, 0xffffffff, address(0x1));
    }

}
