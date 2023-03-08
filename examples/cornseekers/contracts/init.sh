#!/bin/bash

set -eu
set -o pipefail

RPC_URL="${RPC_URL:-"http://localhost:8545"}"
PRIVATE_KEY="${PRIVATE_KEY:-"0x6335c92c05660f35b36148bbfb2105a68dd40275ebf16eff9524d487fb5d57a8"}"

forge create \
	--offline \
	--rpc-url "${RPC_URL}" \
	--private-key "${PRIVATE_KEY}" \
	"./src/Game.sol:Game"

echo OK
