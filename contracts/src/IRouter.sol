// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Dispatcher} from "./IDispatcher.sol";
import {Op} from "./BaseState.sol";

// Routers accept "signed" bundles of Actions and forwards them to Dispatcher.dispatch
// They might be a seperate contract or an extension of the Dispatcher
// A "bundle" here means; one or more actions all signed by the same session key
interface Router {
    function dispatch(bytes[] calldata actionBundles, bytes calldata bundleSignatures, uint256 nonce)
        external
        returns (Op[] memory);

    function authorizeAddr(Dispatcher dispatcher, uint32 ttl, uint32 scopes, address addr) external;

    function authorizeAddr(Dispatcher dispatcher, uint32 ttl, uint32 scopes, address addr, bytes calldata sig)
        external;

    function revokeAddr(address addr) external;

    function revokeAddr(address addr, bytes calldata sig) external;
}
