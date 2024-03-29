package config

import "github.com/ethereum/go-ethereum/common"

var IndexerProviderHTTP = getRequiredEnvString("INDEXER_PROVIDER_URL_HTTP")
var IndexerProviderWS = getRequiredEnvString("INDEXER_PROVIDER_URL_WS")
var IndexerMaxConcurrency = getOptionalEnvInt("INDEXER_MAX_CONCURRENCY", 200)
var IndexerMaxLogRange = getOptionalEnvInt("INDEXER_MAX_LOG_RANGE", 1000)
var IndexerWatchPending = getOptionalEnvBool("INDEXER_WATCH_PENDING", "true")
var IndexerGameAddress = getOptionalEnvAddress("INDEXER_GAME_ADDRESS", common.Address{})
var IndexerStateAddress = getOptionalEnvAddress("INDEXER_STATE_ADDRESS", common.Address{})
var IndexerRouterAddress = getOptionalEnvAddress("INDEXER_ROUTER_ADDRESS", common.Address{})

var SequencerProviderHTTP = getRequiredEnvString("SEQUENCER_PROVIDER_URL_HTTP")
var SequencerProviderWS = getRequiredEnvString("SEQUENCER_PROVIDER_URL_WS")
var SequencerPrivateKey = getRequiredEnvKey("SEQUENCER_PRIVATE_KEY")
var SequencerMaxConcurrency = getOptionalEnvInt("SEQUENCER_MAX_CONCURRENCY", 200)
var SequencerMinBatchDelayMilliseconds = getOptionalEnvInt("SEQUENCER_MIN_BATCH_DELAY_MS", 100)
var SequencerMineEmpty = getOptionalEnvBool("SEQUENCER_MINE_EMPTY", "true")
var SequencerPendingSim = getOptionalEnvBool("SEQUENCER_PENDING_SIM", "false")

var APIPort = getOptionalEnvInt("API_PORT", 8080)
