// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {State} from "./State.sol";
import {Context, Router, Dispatcher} from "./Dispatcher.sol";

import {LibString} from "../src/utils/LibString.sol";

using {LibString.toString} for uint256;

bytes constant PREFIX_MESSAGE = "\x19Ethereum Signed Message:\n";
bytes constant AUTHEN_MESSAGE = "You are signing in with session: ";
bytes constant REVOKE_MESSAGE = "You are signing out of session: ";

error SessionExpiryTooLong();
error SessionUnauthorized();
error SessionExpired();

uint32 constant MAX_TTL = 40000;

contract SessionRouter is Router {
    event SessionCreate(address session, address owner, uint32 exp, uint32 scopes);

    event SessionDestroy(address session);

    // TODO: needs gasgolfing
    struct Session {
        Dispatcher dispatcher;
        address owner;
        uint32 exp;
        uint32 scopes;
    }

    mapping(address => Session) public sessions;

    // authorizeKey delegates permissions to key to act as msg.sender when talking to dispatcher
    function authorizeAddr(Dispatcher dispatcher, uint32 ttl, uint32 scopes, address addr) public {
        _authorizeAddr(dispatcher, ttl, scopes, addr, msg.sender);
    }

    // authorizeKey delegates permissions to key to act as the signer of v/r/s when talking to dispatcher
    function authorizeAddr(Dispatcher dispatcher, uint32 ttl, uint32 scopes, address addr, bytes calldata sig) public {
        address owner = ecrecover(
            keccak256(abi.encodePacked(PREFIX_MESSAGE, (AUTHEN_MESSAGE.length + 20).toString(), AUTHEN_MESSAGE, addr)),
            uint8(bytes1(sig[64:65])),
            bytes32(sig[0:32]),
            bytes32(sig[32:64])
        );
        if (owner == address(0)) {
            revert SessionUnauthorized();
        }
        _authorizeAddr(dispatcher, ttl, scopes, addr, owner);
    }

    function _authorizeAddr(Dispatcher dispatcher, uint32 ttl, uint32 scopes, address sessionAddr, address ownerAddr)
        internal
    {
        uint32 exp = expires(ttl);
        sessions[sessionAddr] = Session({dispatcher: dispatcher, exp: exp, scopes: scopes, owner: ownerAddr});
        emit SessionCreate(sessionAddr, ownerAddr, exp, scopes);
    }

    // revokeKey expires the session key, requires msg.sender to be owner of the session
    function revokeAddr(address addr) public {
        Session storage session = sessions[addr];
        if (session.owner != msg.sender) {
            revert SessionUnauthorized();
        }
        delete sessions[addr];
        emit SessionDestroy(addr);
    }

    // revokeKey expires the session key, requires signer of v/r/s to be session owner
    function revokeAddr(address addr, bytes calldata sig) public {
        address owner = ecrecover(
            keccak256(abi.encodePacked(PREFIX_MESSAGE, (REVOKE_MESSAGE.length + 20).toString(), REVOKE_MESSAGE, addr)),
            uint8(bytes1(sig[64:65])),
            bytes32(sig[0:32]),
            bytes32(sig[32:64])
        );
        Session storage session = sessions[addr];
        if (session.owner != owner) {
            revert SessionUnauthorized();
        }
        delete sessions[addr];
    }

    // dispatch expects actionSig to be either:
    // - a valid sig (v/r/s) of the action data, in which case we treat the SIGNER as the session key
    // - an empty sig, in which case we treat the SENDER as the session key
    // session
    // if the key has not expired, the target dispatcher is called with the generated context
    //
    // +-----------------------------------------------------------------------------------------+
    // | [!] CRITICAL TODO: there is currently no replay protection for session signed actions! |
    // +-----------------------------------------------------------------------------------------+
    //
    function _dispatch(bytes[] calldata actions, bytes calldata sig) private {
        Session storage session;
        if (sig.length == 0) {
            // no signature provided, so we treat the sender as the session key
            // this is useful for authorizing external contract addresses to act
            // on behalf of the player
            session = sessions[msg.sender];
        } else {
            // ecrecover sender from sig as key to lookup session info
            // this is the path for when a player is using a temporary
            // short lived session key in their client to sign actions
            address signer = ecrecover(
                keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", keccak256(abi.encode(actions)))),
                uint8(bytes1(sig[64:65])),
                bytes32(sig[0:32]),
                bytes32(sig[32:64])
            );
            session = sessions[signer];
        }
        if (session.owner == address(0)) {
            revert SessionUnauthorized();
        }
        if (block.number > session.exp) {
            revert SessionExpired();
        }
        // TODO: replay protection
        Context memory ctx = Context({sender: session.owner, scopes: session.scopes, clock: uint32(block.number)});
        // forward to the dispatcher registered with the session
        session.dispatcher.dispatch(actions, ctx);
    }

    // dispatch (batched)
    function dispatch(bytes[][] calldata actions, bytes[] calldata sig) public {
        for (uint256 i = 0; i < actions.length; i++) {
            _dispatch(actions[i], sig[i]);
        }
    }

    // expires converts a ttl to a future block number
    // reverts if requested ttl "too long"
    function expires(uint32 ttl) internal view returns (uint32) {
        if (ttl > MAX_TTL) {
            // TODO: make this configurable
            revert SessionExpiryTooLong();
        }
        return uint32(block.number + ttl);
    }

    // annotations are blobs of data stored in the transaction calldata
    // we take a hash of any annotations and pass the hash to the dispatcher
    // the hash can be used as a reference to data that we can guarentee has been
    // made available to off-chain clients
    function hashAnnotations(bytes[] calldata annotations) private pure returns (bytes32[] memory) {
        bytes32[] memory hashes = new bytes32[](annotations.length);
        for (uint256 i = 0; i < annotations.length; i++) {
            hashes[i] = keccak256(annotations[i]);
        }
        return hashes;
    }
}
