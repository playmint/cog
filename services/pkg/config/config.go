package config

var IndexerProviderHTTP = getRequiredEnvString("INDEXER_PROVIDER_URL_HTTP")
var IndexerProviderWS = getRequiredEnvString("INDEXER_PROVIDER_URL_WS")
var IndexerMaxConcurrency = getOptionalEnvInt("INDEXER_MAX_CONCURRENCY", 200)
var IndexerMaxLogRange = getOptionalEnvInt("INDEXER_MAX_LOG_RANGE", 1000)
var IndexerWatchPending = getOptionalEnvBool("INDEXER_WATCH_PENDING", "true")

var SequencerProviderHTTP = getRequiredEnvString("SEQUENCER_PROVIDER_URL_HTTP")
var SequencerProviderWS = getRequiredEnvString("SEQUENCER_PROVIDER_URL_WS")
var SequencerPrivateKey = getRequiredEnvKey("SEQUENCER_PRIVATE_KEY")
var SequencerMaxConcurrency = getOptionalEnvInt("SEQUENCER_MAX_CONCURRENCY", 200)
var SequencerMinBatchDelaySeconds = getOptionalEnvInt("SEQUENCER_MIN_BATCH_DELAY_SECONDS", 1)

var APIPort = getOptionalEnvInt("API_PORT", 8080)