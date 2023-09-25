package eventwatcher

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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
	Websocket  *alchemy.Client
	HTTPClient *alchemy.Client
	EpochBlock int64
	LogRange   int
	Simulated  bool
	Addresses  []common.Address
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
	return &Watcher{
		sink:        make(chan *LogBatch, 1024),
		subscribers: []chan *LogBatch{},
		ready:       make(chan struct{}),
		config:      cfg,
		log:         log.With().Str("service", "indexer").Str("component", "eventwatcher").Bool("simulated", cfg.Simulated).Int64("epoch", cfg.EpochBlock).Logger(),
	}, nil
}

func (rs *Watcher) Stop() {
	if rs.stop != nil {
		rs.stop()
	}
}

func (rs *Watcher) Start(ctx context.Context) {
	rs.log.Info().Msg("watcher-start")
	topicQuery := ethereum.FilterQuery{Topics: [][]common.Hash{rs.topic0}, Addresses: rs.config.Addresses}
	ctx, rs.stop = context.WithCancel(ctx)
	go rs.watch(ctx, topicQuery)
	go rs.publisher(ctx)
}

func (rs *Watcher) fetchEvents(ctx context.Context, query ethereum.FilterQuery) {
	nowBlockRes, err := rs.config.HTTPClient.BlockNumber(ctx)
	if err != nil {
		rs.log.Error().Err(err).Msg("fetch-events-get-block")
		return
	}
	nowBlock := int64(nowBlockRes)
	fromBlock := rs.config.EpochBlock
	batchSize := int64(rs.config.LogRange)

	for {
		if fromBlock > nowBlock {
			break // done
		}
		toBlock := min(fromBlock+batchSize-1, nowBlock)
		if fromBlock >= toBlock {
			break
		}
		batch := EventBatch{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
		}
		rs.getBatchWithRetry(ctx, batch, query)
		fromBlock = fromBlock + batchSize
	}

	close(rs.ready)
}

func (rs *Watcher) getBatchWithRetry(ctx context.Context, batch EventBatch, query ethereum.FilterQuery) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := rs.getBatch(ctx, batch, query); err != nil {
				rs.log.Error().
					Err(err).
					Int64("from", batch.FromBlock).
					Int64("to", batch.ToBlock).
					Msg("get-batch-with-retry")
				time.Sleep(time.Second)
				continue
			}
			return
		}
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
	rs.fetchEvents(ctx, query)
	<-rs.ready
	for {
		select {
		case <-ctx.Done():
			rs.log.Info().Msg("done")
			return
		default:
			if err := rs.subscribeHead(ctx, query); err != nil {
				rs.log.Error().
					Err(err).
					Msg("subscribe-fail")
				time.Sleep(time.Second * 2)
			}
		}
	}
}

func (rs *Watcher) subscribeHead(ctx context.Context, query ethereum.FilterQuery) error {
	rs.log.Info().Msg("subscribed")
	defer rs.log.Info().Msg("unsubscribed")
	blocks := make(chan *types.Header, 32)
	sub, err := rs.config.Websocket.SubscribeNewHead(context.Background(), blocks)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	return rs.watcher(ctx, sub, blocks, query)
}

func (rs *Watcher) watcher(ctx context.Context, sub event.Subscription, blocks chan *types.Header, query ethereum.FilterQuery) error {
	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("suberr: %v", err)
		case block := <-blocks:
			rs.getBatchWithRetry(ctx, EventBatch{
				FromBlock: block.Number.Int64(),
				ToBlock:   block.Number.Int64(),
			}, query)
		case <-ctx.Done():
			return nil
		}
	}
}

func (rs *Watcher) publisher(ctx context.Context) {
	for {
		select {
		case logs := <-rs.sink:
			for _, sub := range rs.subscribers {
				sub <- logs
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
	ch := make(chan *LogBatch, 0)
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
