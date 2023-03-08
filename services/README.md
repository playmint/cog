# cog-services

## Overview

Auxillary services supporting the games built with [cog](https://github.com/playmint/cog).

* API - GraphQL endpoint for displaying Indexer data
* Sequencer - Queues and submits player signed actions to chain
* Indexer - Fast cache of game state from chain
* Tests - Integration test running against [cog-examples](https://github.com/playmint/cog-examples)

```mermaid
flowchart TD
	API["GraphQL API"]
	Seq["Sequencer"]
	Idx["Indexer"]
	Unity["Unity Game Client"]
	Rule1["Rule"]
	Rule2["Rule"]
	Rule3["Rule"]
	Dispatcher
	Router
	subgraph cog-services
		API
		Seq -- status --> Idx
	end
	subgraph on-chain-cog-game
		Dispatcher
		State
		Router -- verifyActionSignature -->Router
		Router -- action -->Dispatcher
		Dispatcher -- action -->Rule1
		Dispatcher -- action -->Rule2
		Dispatcher -- action -->Rule3
		Rule1 -- modifyState -->State
		Rule2 -- modifyState -->State
		Rule3 -- modifyState -->State
	end
	subgraph game-client
		Unity
	end
	Unity -- query/action --> API
	API -- state --> Unity
	API -- action --> Seq
	API -- query --> Idx
	Idx -- state --> API
	Seq -- action --> Router
	State -- event -->Idx

```

## Quickstart

Use Docker Compose to provision a local instance of cog-services.

```
docker-compose up --build
```

To run the integration tests against the local deployment:

```
docker compose --profile=test up --build --exit-code-from cog-tests
```


