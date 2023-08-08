// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";

import { State } from "cog/IState.sol";

import { Game } from "src/Game.sol";
import { Actions, Direction } from "src/actions/Actions.sol";
import { Schema, Node, BiomeKind, ResourceKind } from "src/schema/Schema.sol";

using Schema for State;

contract CornSeekersTest is Test {

    Game internal game;
    State internal state;

    // accounts
    address aliceAccount;

    function setUp() public {
        // setup game
        game = new Game();

        // fetch the State to play with
        state = game.getState();

        // setup users
        uint256 alicePrivateKey = 0xA11CE;
        aliceAccount = vm.addr(alicePrivateKey);

        // reset map before all tests
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.RESET_MAP, ())
        );
    }

    function testHarvesting() public {
        // moving a seeker onto a CORN tile harvests the corn
        // this converts the tile to a GRASS tile and increases
        // the HAS_RESOURCE balance on the seeker

        // dispatch as alice
        vm.startPrank(aliceAccount);

        // spawn a seeker bottom left corner of map
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.SPAWN_SEEKER, (
                1,   // seeker id (sid)
                0,   // x
                0,   // y
                100  // strength attr
            ))
        );

        assertEq(
            state.getLocation(Node.Seeker(1)),
            Node.Tile(0,0),
            "expected seeker to start at tile 0,0"
        );

        assertEq(
            state.getResourceBalance(Node.Seeker(1), ResourceKind.CORN),
            0,
            "expected seeker's CORN resource balance to start at zero"
        );

        // hack in CORN at tile (0,1) to bypass scouting
        state.setBiome( Node.Tile(0,1), BiomeKind.CORN );

        // move the seeker NORTH to tile (0,1)
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.MOVE_SEEKER, (
                1,                   // seeker id (sid)
                Direction.NORTH      // direction to move
            ))
        );

        // comfirm seeker is now at tile (1,1)
        assertEq(
            state.getLocation(Node.Seeker(1)),
            Node.Tile(0,1),
            "expected seeker 1 at tile 1,1"
        );

        // confirm our corn balance is now 1
        assertEq(
            state.getResourceBalance(Node.Seeker(1), ResourceKind.CORN),
            1,
            "expected seeker 1 to have 1unit of corn"
        );

        // stop being alice
        vm.stopPrank();
    }

    function testScouting() public {
        // scouting is in two parts because it requires
        // some randomness.
        //
        // the first part occurs during a MOVE_SEEKER or SPAWN_SEEKER
        // action which requests a SEED for any surrounding tiles
        //
        // a seed request is an edge from a SEED node pointing to
        // a TILE node.
        //
        // later a REVEAL_SEED request will submit the required
        // entopy and perform any followup processing
        //

        // dispatch as alice
        vm.startPrank(aliceAccount);

        // spawn a seeker bottom left corner of map
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.SPAWN_SEEKER, (
                1,   // seeker id (sid)
                0,   // x
                0,   // y
                100  // strength attr
            ))
        );

        uint64 str = state.getStrength(Node.Seeker(1));
        assertEq(
            str,
            100,
            "expect seeker 1 to have a strength set"
        );

        assertEq(
            state.getLocation(Node.Seeker(1)),
            Node.Tile(0,0),
            "expect seeker to have a location"
        );

        bytes24[] memory pendingTiles = state.getEntropyCommitments(uint32(block.number));
        assertEq(
            pendingTiles.length,
            1,
            "expected 1 pending adjacent tile caused by spawning seeker at 0,0"
        );

        BiomeKind pendingContent = state.getBiome(pendingTiles[0]);
        assertEq(
            uint(pendingContent),
            uint(BiomeKind.UNDISCOVERED),
            "the one pending tile should be UNDISCOVERED"
        );

        // wait until the blockhash is revealed
        vm.roll(block.number + 1);

        // once we know the blockhash of the requested
        // seed, we can submit REVEAL_SEED action
        // to resolve it
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.REVEAL_SEED, (
                uint32(block.number - 1),
                uint32(uint(blockhash(block.number-1)))
            ))
        );

        BiomeKind discoveredContent = state.getBiome(pendingTiles[0]);
        assertGt(
            uint(discoveredContent),
            uint(BiomeKind.UNDISCOVERED),
            "the pending tile sholud now be discovered"
        );

        // move the seeker NORTH
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.MOVE_SEEKER, (
                1,                   // seeker id (sid)
                Direction.NORTH      // direction to move
            ))
        );

        assertEq(
            state.getLocation(Node.Seeker(1)),
            Node.Tile(0,1),
            "expected seeker 1 at location 0,1"
        );

        pendingTiles = state.getEntropyCommitments(uint32(block.number));
        assertEq(
            pendingTiles.length,
            1,
            "expected there to be 3 pending tiles after moving to tile 0,1"
        );

        // attempting to move EAST into an UNDISCOVERED tile
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.MOVE_SEEKER, (
                1,                   // seeker id (sid)
                Direction.EAST       // direction to move
            ))
        );
        assertEq(
            state.getLocation(Node.Seeker(1)),
            Node.Tile(0,1),
            "expected seeker to not have moved onto UNDISCOVERED tile"
        );

        // roll time forward
        vm.roll(block.number + 1);

        // submit the reveal action
        game.getDispatcher().dispatch(
            abi.encodeCall(Actions.REVEAL_SEED, (
                uint32(block.number - 1),
                uint32(uint(blockhash(block.number-1)))
            ))
        );

        for (uint i=0; i<pendingTiles.length; i++) {
            pendingContent = state.getBiome(pendingTiles[i]);
            assertGt(
                uint(pendingContent),
                uint(BiomeKind.UNDISCOVERED),
                string(abi.encodePacked("expected pending tile ",i+48, " to be discovered"))
            );
        }


        // stop being alice
        vm.stopPrank();
    }

}
