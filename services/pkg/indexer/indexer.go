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
	BlockNumber(context.Context) (uint64, error)
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
}

func NewMemoryIndexer(ctx context.Context, notifications chan interface{}) (*MemoryIndexer, error) {
	var err error

	idxr := &MemoryIndexer{}

	idxr.notifications = notifications

	idxr.httpClient, err = alchemy.Dial(
		config.IndexerProviderHTTP,
		config.IndexerMaxConcurrency,
		nil,
	)
	if err != nil {
		return nil, err
	}

	idxr.wsClient, err = alchemy.Dial(
		config.IndexerProviderWS,
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
		idxr.httpClient,
		idxr.events,
		idxr.notifications,
	)
	if err != nil {
		return nil, err
	}

	// start listening for SessionCreate events
	idxr.sessionStore, err = cog.NewSessionStore(
		ctx,
		idxr.httpClient,
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

func (idxr *MemoryIndexer) BlockNumber(ctx context.Context) (uint64, error) {
	return idxr.httpClient.BlockNumber(ctx)
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
	return idxr.stateStore.GetGraph(stateContractAddr)
}
func (idxr *MemoryIndexer) GetSession(routerAddr common.Address, sessionID string) *model.Session {
	return idxr.sessionStore.GetSession(routerAddr, sessionID)
}

func (idxr *MemoryIndexer) GetSessions(routerAddr common.Address, owner *string) []*model.Session {
	return idxr.sessionStore.GetSessions(routerAddr, owner)
}
