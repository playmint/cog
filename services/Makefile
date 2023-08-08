SHELL := /bin/bash
CWD := $(shell cd -P -- '$(shell dirname -- "$0")' && pwd -P)
SRC := $(shell find . -name '*.go')

COG_CONTRACTS_DIR=../contracts
COG_CONTRACTS_SOL=$(wildcard $(COG_CONTRACTS_DIR)/src/*.sol)
COG_CONTRACTS_GO=$(patsubst $(COG_CONTRACTS_DIR)/src/%,pkg/contracts/%,$(COG_CONTRACTS_SOL:.sol=.go))

COG_EXAMPLES_CONTRACTS_DIR=../examples

.PHONY: all
all: bin/ds-node

bin/ds-node: bin/wait-for $(SRC) pkg/contracts/state/State.go pkg/contracts/router/SessionRouter.go pkg/contracts/game/BaseGame.go pkg/api/resolver/resolver.go
	mkdir -p bin
	go build -o $@ ./cmd/ds-node

bin/wait-for: $(SRC)
	mkdir -p bin
	go build -o $@ ./cmd/wait-for

pkg/api/resolver/resolver.go: schema/*.graphqls gqlgen.yml
	go generate ./...
	touch $@

pkg/contracts/state/State.go: $(COG_CONTRACTS_DIR)/src/IState.sol
	mkdir -p pkg/contracts/state
	(cd $(COG_CONTRACTS_DIR) && forge inspect State abi) | abigen -abi - -pkg state -type State --out $@

pkg/contracts/router/SessionRouter.go: $(COG_CONTRACTS_DIR)/src/BaseRouter.sol
	mkdir -p pkg/contracts/router
	(cd $(COG_CONTRACTS_DIR) && forge inspect BaseRouter abi) | abigen -abi - -pkg router -type SessionRouter --out $@

pkg/contracts/game/BaseGame.go: $(COG_CONTRACTS_DIR)/src/IGame.sol
	mkdir -p pkg/contracts/game
	(cd $(COG_CONTRACTS_DIR) && forge inspect BaseGame abi) | abigen -abi - -pkg game -type BaseGame --out $@

test/integration/fixtures/cornseekers/Actions.go: $(COG_EXAMPLES_CONTRACTS_DIR)/cornseekers/contracts/src/actions/Actions.sol
	mkdir -p test/integration/fixtures/cornseekers
	(cd $(COG_EXAMPLES_CONTRACTS_DIR)/cornseekers/contracts && forge inspect Actions abi) | abigen -abi - -pkg cornseekers -type Actions --out $@

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: dev-local
dev-local: bin/ds-node
	PORT=8181 \
	CHAIN_ID=1337 \
    SEQUENCER_PRIVATE_KEY=095a37ef5b5d87db7fe50551725cb64804c8c554868d3d729c0dd17f0e664c87 \
	SEQUENCER_PROVIDER_URL_HTTP=http://localhost:8545 \
	SEQUENCER_PROVIDER_URL_WS=ws://localhost:8545 \
	INDEXER_WATCH_PENDING=false \
	INDEXER_PROVIDER_URL_HTTP=http://localhost:8545 \
	INDEXER_PROVIDER_URL_WS=ws://localhost:8545 \
		./bin/ds-node -debug

.PHONY: test
test: pkg/api/resolver/resolver.go
	go test -v -count=1 ./...
