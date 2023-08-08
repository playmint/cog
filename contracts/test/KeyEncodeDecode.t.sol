// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";

import {State, CompoundKeyEncoder, CompoundKeyDecoder} from "../src/IState.sol";

interface Kind {
    function Item() external;
}

contract KeyEncodeDecodeTest is Test {
    function testStringKey() public {
        // Test strings under 20 chars
        string memory testKey = "0123456789";
        bytes24 encodedKey = CompoundKeyEncoder.STRING(Kind.Item.selector, testKey);
        string memory decodedKey = CompoundKeyDecoder.STRING(encodedKey);

        assertEq(decodedKey, testKey, "expected decoded key to equal original key");

        // Test strings that are 20 chars.
        testKey = "01234567890123456789";
        encodedKey = CompoundKeyEncoder.STRING(Kind.Item.selector, testKey);
        decodedKey = CompoundKeyDecoder.STRING(encodedKey);

        assertEq(decodedKey, testKey, "expected decoded key to equal original key");

        // Test strings that are over 20 chars (Expected to truncate the key)
        testKey = "012345678901234567890123456789";
        encodedKey = CompoundKeyEncoder.STRING(Kind.Item.selector, testKey);
        decodedKey = CompoundKeyDecoder.STRING(encodedKey);

        string memory testKeyTruncated = "01234567890123456789";
        assertEq(decodedKey, testKeyTruncated, "expected decoded key to equal truncated 20 char key");
    }
}
