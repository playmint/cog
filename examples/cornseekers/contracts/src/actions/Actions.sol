// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

// ----------------------------------
// define some constants/enums
// ----------------------------------

enum Direction {
    NORTH,
    NORTHEAST,
    EAST,
    SOUTHEAST,
    SOUTH,
    SOUTHWEST,
    WEST,
    NORTHWEST
}

// ----------------------------------
// define some actions
// ----------------------------------

interface Actions {
    function RESET_MAP() external;
    function REVEAL_SEED(uint32 blk, uint32 entropy) external;
    function SPAWN_SEEKER(uint32 sid, uint8 x, uint8 y, uint8 str) external;
    function MOVE_SEEKER(uint32 sid, Direction dir) external;
}
