package pendingtx

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	// "github.com/ethereum/go-ethereum/core/types"

	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/playmint/ds-node/pkg/client/types"
)

type Config struct {
	Watch     bool
	Websocket *alchemy.Client
	Addresses []common.Address
}

type Watcher struct {
	sink          chan types.Transaction
	subscriptions map[common.Address][]chan types.Transaction
	config        Config
}

func New(cfg Config) (*Watcher, error) {
	return &Watcher{
		sink:          make(chan types.Transaction, 64),
		subscriptions: map[common.Address][]chan types.Transaction{},
		config:        cfg,
	}, nil
}

func (rs *Watcher) Start(ctx context.Context) {
	// watching the pending pool is not supported by
	// all providers and nodes (ie local dev) so give an
	// option to turn it off
	if rs.config.Watch {
		for _, addr := range rs.config.Addresses {
			go rs.subscribePendingWithReconnect(ctx, addr)
		}
	}
	go rs.emitter(ctx)
}

// subscribePendingWithReconnect calls subscribePending in a loop to reconnect if dropped
func (rs *Watcher) subscribePendingWithReconnect(ctx context.Context, addr common.Address) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		rs.subscribePending(ctx, addr)
	}
}

// subscribePending opens websocket to the provider to watch pending pool
func (rs *Watcher) subscribePending(ctx context.Context, addr common.Address) {
	fmt.Println("started watching pending tx for", addr)
	defer fmt.Println("stopped watching pending tx for", addr)

	watchQuery := ethereum.FilterQuery{
		Addresses: []common.Address{addr},
	}
	sub, err := rs.config.Websocket.SubscribeFilterPending(ctx, watchQuery, rs.sink)
	if err != nil {
		log.Println("pending-watch-fail-backoff:", err)
		time.Sleep(time.Second * 2)
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Println("watcher-err:", err)
			return
		case <-ctx.Done():
			return
		}
	}
}

// emitter forwards tx to anyone who subscribes
func (rs *Watcher) emitter(ctx context.Context) {
	for {
		select {
		case tx := <-rs.sink:
			fmt.Println("pushing out tx", tx)
			subscriptions, ok := rs.subscriptions[tx.To]
			if !ok {
				continue
			}
			for _, sub := range subscriptions {
				sub <- tx
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rs *Watcher) Subscribe(address common.Address) chan types.Transaction {
	ch := make(chan types.Transaction, 64)
	rs.subscriptions[address] = append(rs.subscriptions[address], ch)
	return ch
}

func (rs *Watcher) Send(tx types.Transaction) {
	rs.sink <- tx
}
