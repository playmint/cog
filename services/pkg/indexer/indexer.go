package indexer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/playmint/ds-node/pkg/config"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/playmint/ds-node/pkg/indexer/stores/cog"
	"github.com/playmint/ds-node/pkg/indexer/stores/configstore"
)

type Indexer interface {
	Ready() chan struct{}

	// query funcs
	GetGame(id string) *model.Game
	GetGames() []*model.Game
	GetGraph(stateContractAddr common.Address) *model.Graph
	GetSession(routerAddr common.Address, sessionID string) *model.Session
	GetSessions(routerAddr common.Address, owner *string) []*model.Session
	NewSim(ctx context.Context, blockNumber uint64, httpSimClient *alchemy.Client, wsSimClient *alchemy.Client) (*MemoryIndexer, error)
	GetSim() *MemoryIndexer
	SetSim(*MemoryIndexer)
}

var _ Indexer = &MemoryIndexer{}

type MemoryIndexer struct {
	configStore   *configstore.ConfigStore
	gameStore     *cog.GameStore
	stateStore    *cog.StateStore
	sessionStore  *cog.SessionStore
	notifications chan interface{}
	events        *eventwatcher.Watcher
	httpClient    *alchemy.Client
	wsClient      *alchemy.Client
	sim           *MemoryIndexer
}

func NewMemoryIndexer(ctx context.Context, notifications chan interface{}, httpProviderURL string, wsProviderURL string) (*MemoryIndexer, error) {
	var err error

	idxr := &MemoryIndexer{}

	idxr.notifications = notifications

	idxr.httpClient, err = alchemy.Dial(
		httpProviderURL,
		config.IndexerMaxConcurrency,
		nil,
	)
	if err != nil {
		return nil, err
	}

	idxr.wsClient, err = alchemy.Dial(
		wsProviderURL,
		config.IndexerMaxConcurrency,
		nil,
	)
	if err != nil {
		return nil, err
	}

	idxr.events, err = eventwatcher.New(eventwatcher.Config{
		HTTPClient:  idxr.httpClient,
		Websocket:   idxr.wsClient,
		Concurrency: 1, // config.IndexerMaxConcurrency, - NodeSet/EdgeSet cannot arrive out of order yet
		LogRange:    config.IndexerMaxLogRange,
	})
	if err != nil {
		return nil, err
	}

	// layerTwoPendingTx, err := pendingtx.New(pendingtx.Config{
	// 	Watch:     config.IndexerWatchPending,
	// 	Websocket: layerTwoWSClient,
	// 	Addresses: []common.Address{
	// 		common.HexToAddress(config.Contracts.Router),
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }
	// layerTwoPendingTx.Start(ctx)

	// index cog games, dispatchers, state
	idxr.gameStore, err = cog.NewGameStore(
		ctx,
		idxr.httpClient,
		idxr.events,
		idxr.notifications,
	)
	if err != nil {
		return nil, err
	}

	// start listening for NodeSet and EdgeSet events
	idxr.stateStore, err = cog.NewStateStore(
		ctx,
		idxr.events,
		idxr.notifications,
	)
	if err != nil {
		return nil, err
	}

	// start listening for SessionCreate events
	idxr.sessionStore, err = cog.NewSessionStore(
		ctx,
		idxr.events,
		idxr.notifications,
	)
	if err != nil {
		return nil, err
	}

	// index config data
	idxr.configStore = configstore.New()

	// start event collection
	idxr.events.Start(ctx)

	return idxr, nil
}

func (idxr *MemoryIndexer) NewSim(ctx context.Context, blockNumber uint64, httpSimClient *alchemy.Client, wsSimClient *alchemy.Client) (*MemoryIndexer, error) {
	if idxr.sim != nil {
		idxr.sim.events.Stop()
	}
	events, err := eventwatcher.New(eventwatcher.Config{
		HTTPClient:  httpSimClient,
		Websocket:   wsSimClient,
		Concurrency: 1,
		LogRange:    config.IndexerMaxLogRange,
		EpochBlock:  int64(blockNumber),
	})
	if err != nil {
		return nil, err
	}
	// clone this indexer
	newIdxr := &MemoryIndexer{
		configStore:   idxr.configStore,
		gameStore:     idxr.gameStore.Fork(ctx, events, httpSimClient), // FIXME: fork it proper
		stateStore:    idxr.stateStore.Fork(ctx, events),
		sessionStore:  idxr.sessionStore, // FIXME fork it
		notifications: idxr.notifications,
		events:        events,
		httpClient:    httpSimClient,
		wsClient:      wsSimClient,
	}
	newIdxr.events.Start(ctx)
	<-newIdxr.events.Ready()
	return newIdxr, nil
}

func (idxr *MemoryIndexer) SetSim(sim *MemoryIndexer) {
	idxr.sim = sim
}

func (idxr *MemoryIndexer) GetSim() *MemoryIndexer {
	return idxr.sim
}

func (idxr *MemoryIndexer) Ready() chan struct{} {
	return idxr.events.Ready()
}

func (idxr *MemoryIndexer) GetGame(id string) *model.Game {
	return idxr.gameStore.GetGame(id)
}

func (idxr *MemoryIndexer) GetGames() []*model.Game {
	return idxr.gameStore.GetGames()
}

func (idxr *MemoryIndexer) GetGraph(stateContractAddr common.Address) *model.Graph {
	if idxr.sim != nil {
		return idxr.sim.GetGraph(stateContractAddr)
	}
	return idxr.stateStore.GetGraph(stateContractAddr)
}
func (idxr *MemoryIndexer) GetSession(routerAddr common.Address, sessionID string) *model.Session {
	return idxr.sessionStore.GetSession(routerAddr, sessionID)
}

func (idxr *MemoryIndexer) GetSessions(routerAddr common.Address, owner *string) []*model.Session {
	return idxr.sessionStore.GetSessions(routerAddr, owner)
}
