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
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/client"
	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EventBatch struct {
	FromBlock int64
	ToBlock   int64
}

type LogBatch struct {
	EventBatch
	Logs []types.Log
}

type Config struct {
	Websocket            *alchemy.Client
	HTTPClient           *alchemy.Client
	Concurrency          int
	EpochBlock           int64
	LogRange             int
	Simulated            bool
	Notifications        chan interface{}
	NotificationsEnabled bool
}

type Watcher struct {
	sink        chan *LogBatch
	ready       chan struct{}
	subscribers []chan *LogBatch
	topic0      []common.Hash
	config      Config
	log         zerolog.Logger
	stop        func()
}

func New(cfg Config) (*Watcher, error) {
	if cfg.LogRange < 1 {
		return nil, fmt.Errorf("invalid log range config")
	}
	if cfg.Concurrency < 1 {
		return nil, fmt.Errorf("invalid concurrency config")
	}
	return &Watcher{
		sink:        make(chan *LogBatch, 1024),
		subscribers: []chan *LogBatch{},
		ready:       make(chan struct{}),
		config:      cfg,
		log:         log.With().Str("service", "indexer").Str("component", "eventwatcher").Bool("simulated", cfg.Simulated).Logger(),
	}, nil
}

func (rs *Watcher) Stop() {
	if rs.stop != nil {
		rs.stop()
	}
}

func (rs *Watcher) Start(ctx context.Context) {
	rs.log.Info().Msg("watcher-start")
	topicQuery := ethereum.FilterQuery{Topics: [][]common.Hash{rs.topic0}}
	ctx, rs.stop = context.WithCancel(ctx)
	go rs.watch(ctx, topicQuery)
	go rs.fetchEvents(ctx, topicQuery)
	go rs.publisher(ctx)
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
	query.FromBlock = big.NewInt(batch.FromBlock)
	query.ToBlock = big.NewInt(batch.ToBlock)
	logs, err := rs.config.HTTPClient.FilterLogs(ctx, query)
	if err != nil {
		return err
	}
	rs.sink <- &LogBatch{
		EventBatch: batch,
		Logs:       logs,
	}
	rs.log.Info().
		Int64("from", batch.FromBlock).
		Int64("to", batch.ToBlock).
		Int("logs", len(logs)).
		Msg("searching")
	return nil
}

func (rs *Watcher) watch(ctx context.Context, query ethereum.FilterQuery) {
	<-rs.ready
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		blocks := make(chan *types.Header, 32)
		sub, err := rs.config.Websocket.SubscribeNewHead(ctx, blocks)
		if err != nil {
			rs.log.Error().
				Err(err).
				Msg("subscribe-fail")
			time.Sleep(time.Second * 2)
			continue
		}
		rs.watcher(ctx, sub, blocks, query)
	}
}

func (rs *Watcher) watcher(ctx context.Context, sub event.Subscription, blocks chan *types.Header, query ethereum.FilterQuery) {
	rs.log.Info().Msg("started")
	defer rs.log.Info().Msg("stopped")
	defer sub.Unsubscribe()
	for {
		select {
		case err := <-sub.Err():
			rs.log.Error().
				Err(err).
				Msg("watcher-sub-err")
			return
		case block := <-blocks:
			retries := 0
		retry:
			err := rs.getBatch(ctx, EventBatch{
				FromBlock: block.Number.Int64() - 1,
				ToBlock:   block.Number.Int64(),
			}, query)
			if err != nil {
				if client.IsRetryable(err) {
					if retries < 5 {
						retries++
						time.Sleep(1 * time.Second)
						goto retry
					} else {
						panic(err)
					}
				} else {
					panic(err) // if we can't read a block we cannot sync so bang!
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rs *Watcher) SetNotificationsEnabled(enable bool) {
	rs.config.NotificationsEnabled = enable
}

func (rs *Watcher) Notify(blockNumber int64) {
	rs.config.Notifications <- &model.BlockEvent{
		ID:        fmt.Sprintf("block-%d", blockNumber),
		Block:     int(blockNumber),
		Simulated: rs.config.Simulated,
	}
}

func (rs *Watcher) publisher(ctx context.Context) {
	for {
		select {
		case logs := <-rs.sink:
			for _, sub := range rs.subscribers {
				sub <- logs
			}
			if rs.config.NotificationsEnabled {
				rs.Notify(logs.ToBlock)
			}
		case <-ctx.Done():
			return
		}

	}
}

func (rs *Watcher) Ready() chan struct{} {
	return rs.ready
}

func (rs *Watcher) SubscribeTopic(eventTypes []common.Hash) chan *LogBatch {
	ch := make(chan *LogBatch, 1024)
	rs.subscribers = append(rs.subscribers, ch)
	for _, topic := range eventTypes {
		rs.topic0 = append(rs.topic0, topic)
	}
	return ch
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
