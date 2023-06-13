package eventwatcher

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/playmint/ds-node/pkg/client"
	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventBatch struct {
	FromBlock int64
	ToBlock   int64
}

type Config struct {
	Websocket   *alchemy.Client
	HTTPClient  *alchemy.Client
	Concurrency int
	EpochBlock  int64
	LogRange    int
}

type Watcher struct {
	sink      chan types.Log
	ready     chan struct{}
	contracts map[common.Address][]chan types.Log
	topics    map[common.Hash][]chan types.Log
	topic0    []common.Hash // event types to watch
	config    Config
	log       zerolog.Logger
}

func New(cfg Config) (*Watcher, error) {
	if cfg.LogRange < 1 {
		return nil, fmt.Errorf("invalid log range config")
	}
	if cfg.Concurrency < 1 {
		return nil, fmt.Errorf("invalid concurrency config")
	}
	return &Watcher{
		sink:      make(chan types.Log, 1024),
		contracts: map[common.Address][]chan types.Log{},
		topics:    map[common.Hash][]chan types.Log{},
		ready:     make(chan struct{}),
		config:    cfg,
		log:       log.With().Str("service", "indexer").Str("component", "eventwatcher").Logger(),
	}, nil
}

func (rs *Watcher) Start(ctx context.Context) {
	addrs := []common.Address{}
	for addr := range rs.contracts {
		addrs = append(addrs, addr)
	}

	if len(addrs) > 0 && len(rs.topic0) > 0 {
		panic("watching both topics and contracts via the same watcher not supported yet as we need to be able to keep the events ordered")
	}

	// watch all logs for contract addrs
	if len(addrs) > 0 {
		rs.log.Info().Msgf("watcher-start addrs:%v", addrs)
		contractQuery := ethereum.FilterQuery{Addresses: addrs}
		go rs.watch(ctx, contractQuery)
		go rs.fetchEvents(ctx, contractQuery)
	}

	// watch all event type accross all contracts
	if len(rs.topic0) > 0 {
		rs.log.Info().Msgf("watcher-start topic0:%v", rs.topic0)
		topicQuery := ethereum.FilterQuery{Topics: [][]common.Hash{rs.topic0}}
		go rs.watch(ctx, topicQuery)
		go rs.fetchEvents(ctx, topicQuery)
	}
	rs.log.Info().Msgf("watcher-started addrs:%v topic0:%v", addrs, rs.topic0)
}

func (rs *Watcher) fetchEvents(ctx context.Context, query ethereum.FilterQuery) {
	nowBlockRes, err := rs.config.HTTPClient.BlockNumber(ctx)
	if err != nil {
		rs.log.Fatal().Err(err).Msg("") // FIXME: error system
	}
	nowBlock := int64(nowBlockRes)
	fromBlock := rs.config.EpochBlock
	batchSize := int64(rs.config.LogRange)

	var wg sync.WaitGroup

	batchChan := make(chan EventBatch)
	for i := 0; i < rs.config.Concurrency; i++ {
		go rs.fetchEventsWorker(ctx, batchChan, query, &wg)
	}

	for {
		if fromBlock > nowBlock {
			break // done
		}
		toBlock := min(fromBlock+batchSize-1, nowBlock)
		if fromBlock >= toBlock {
			break
		}
		wg.Add(1)
		batchChan <- EventBatch{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
		}
		fromBlock = fromBlock + batchSize
	}

	wg.Wait()
	close(batchChan)
	close(rs.ready)
}

func (rs *Watcher) fetchEventsWorker(ctx context.Context, batches chan EventBatch, query ethereum.FilterQuery, wg *sync.WaitGroup) {
	for batch := range batches {
		if err := rs.getBatch(ctx, batch, query); err != nil {
			if client.IsRetryable(err) {
				rs.log.Warn().
					Err(err).
					Int64("from", batch.FromBlock).
					Int64("to", batch.ToBlock).
					Msg("get-batch-rate-limited")
				go func(batch EventBatch) {
					// requeue the failed batch with some jitter
					// TODO: exponetial backoff would be better here to avoid thunder
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
					rs.log.Info().
						Int64("from", batch.FromBlock).
						Int64("to", batch.ToBlock).
						Msg("requeue-batch")
					batches <- batch
				}(batch)
				continue
			}
			rs.log.Fatal().
				Err(err).
				Int64("from", batch.FromBlock).
				Int64("to", batch.ToBlock).
				Msg("get-batch-fail")
		}
		wg.Done()
	}
}

func (rs *Watcher) getBatch(ctx context.Context, batch EventBatch, query ethereum.FilterQuery) error {
	rs.log.Info().
		Int64("from", batch.FromBlock).
		Int64("to", batch.ToBlock).
		Msg("searching")
	query.FromBlock = big.NewInt(batch.FromBlock)
	query.ToBlock = big.NewInt(batch.ToBlock)
	logs, err := rs.config.HTTPClient.FilterLogs(ctx, query)
	if err != nil {
		return err
	}
	for _, log := range logs {
		rs.sink <- log
	}
	return nil
}

func (rs *Watcher) watch(ctx context.Context, watchQuery ethereum.FilterQuery) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		sub, err := rs.config.Websocket.SubscribeFilterLogs(ctx, watchQuery, rs.sink)
		if err != nil {
			rs.log.Error().
				Err(err).
				Msg("subscribe-fail")
			time.Sleep(time.Second * 2)
			continue
		}
		rs.watcher(ctx, sub)
	}
}

func (rs *Watcher) watcher(ctx context.Context, sub event.Subscription) {
	rs.log.Info().Msg("started")
	defer rs.log.Info().Msg("stopped")
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			rs.log.Error().
				Err(err).
				Msg("watcher-err")
			return
		case evt := <-rs.sink:
			contractSubscribers, ok := rs.contracts[evt.Address]
			if ok {
				for _, sub := range contractSubscribers {
					sub <- evt
				}
			}
			if len(evt.Topics) > 0 {
				topicSubscribers, ok := rs.topics[evt.Topics[0]]
				if ok {
					for _, sub := range topicSubscribers {
						sub <- evt
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rs *Watcher) Ready() chan struct{} {
	return rs.ready
}

func (rs *Watcher) Subscribe(address common.Address) chan types.Log {
	ch := make(chan types.Log, 1024)
	rs.contracts[address] = append(rs.contracts[address], ch)
	return ch
}

func (rs *Watcher) SubscribeTopic(eventTypes []common.Hash) chan types.Log {
	ch := make(chan types.Log, 1024)
	for _, topic := range eventTypes {
		rs.topic0 = append(rs.topic0, topic)
		rs.topics[topic] = append(rs.topics[topic], ch)
	}
	return ch
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
