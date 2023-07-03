package indexer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/playmint/ds-node/pkg/indexer/stores/cog"
	"github.com/playmint/ds-node/pkg/indexer/stores/configstore"
)

type SyncRequest struct {
	Game  *model.Game
	Graph *model.SerializableGraph
}

type Indexer interface {
	Ready() chan struct{}

	// query funcs
	GetGame(id string) *model.Game
	GetGames() []*model.Game
	GetGraph(stateContractAddr common.Address, block int, simulated bool) *model.Graph
	GetSession(routerAddr common.Address, sessionID string) *model.Session
	GetSessions(routerAddr common.Address, owner *string) []*model.Session
	NewSim(ctx context.Context, blockNumber uint64, httpSimClient *alchemy.Client, wsSimClient *alchemy.Client) (*MemoryIndexer, error)
	GetSim() *MemoryIndexer
	SetSim(*MemoryIndexer)
	Notify(blockNumber int64, simulated bool)
	SetNotificationsEnabled(enabled bool, simulated bool)
	Sync(req *SyncRequest) error
	Push(batch *eventwatcher.LogBatch) error
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
	game          *model.Game // FIXME: this is a hack
	sim           *MemoryIndexer
}

func NewMemoryIndexer2(ctx context.Context, events *eventwatcher.Watcher) (*MemoryIndexer, error) {
	var err error

	idxr := &MemoryIndexer{}

	idxr.notifications = make(chan interface{})
	idxr.events = events

	// index cog games, dispatchers, state
	idxr.gameStore, err = cog.NewGameStore(
		ctx,
		idxr.httpClient,
		idxr.events,
	)
	if err != nil {
		return nil, err
	}

	// start listening for NodeSet and EdgeSet events
	idxr.stateStore, err = cog.NewStateStore(
		ctx,
		idxr.events,
	)
	if err != nil {
		return nil, err
	}

	// start listening for SessionCreate events
	idxr.sessionStore, err = cog.NewSessionStore(
		ctx,
		idxr.events,
	)
	if err != nil {
		return nil, err
	}

	// index config data
	idxr.configStore = configstore.New()

	// start event collection
	// idxr.events.Start(ctx)

	return idxr, nil
}

func NewMemoryIndexer(ctx context.Context, notifications chan interface{}, httpProviderURL string, wsProviderURL string) (*MemoryIndexer, error) {
	var err error

	idxr := &MemoryIndexer{}

	idxr.notifications = notifications

	idxr.httpClient, err = alchemy.Dial(
		httpProviderURL,
		1, //config.IndexerMaxConcurrency,
		nil,
	)
	if err != nil {
		return nil, err
	}

	idxr.wsClient, err = alchemy.Dial(
		wsProviderURL,
		1, //config.IndexerMaxConcurrency,
		nil,
	)
	if err != nil {
		return nil, err
	}

	idxr.events, err = eventwatcher.New(eventwatcher.Config{
		HTTPClient:    idxr.httpClient,
		Websocket:     idxr.wsClient,
		Concurrency:   1,    // config.IndexerMaxConcurrency, - NodeSet/EdgeSet cannot arrive out of order yet
		LogRange:      1000, // config.IndexerMaxLogRange,
		Notifications: notifications,
	})
	if err != nil {
		return nil, err
	}

	// index cog games, dispatchers, state
	idxr.gameStore, err = cog.NewGameStore(
		ctx,
		idxr.httpClient,
		idxr.events,
	)
	if err != nil {
		return nil, err
	}

	// start listening for NodeSet and EdgeSet events
	idxr.stateStore, err = cog.NewStateStore(
		ctx,
		idxr.events,
	)
	if err != nil {
		return nil, err
	}

	// start listening for SessionCreate events
	idxr.sessionStore, err = cog.NewSessionStore(
		ctx,
		idxr.events,
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

func (idxr *MemoryIndexer) Push(batch *eventwatcher.LogBatch) error {
	return idxr.stateStore.Push(batch)
}

func (idxr *MemoryIndexer) Sync(req *SyncRequest) error {
	idxr.game = req.Game
	g, err := model.LoadGraph(req.Graph)
	if err != nil {
		return err
	}
	idxr.stateStore.Sync(req.Game.StateAddress.Hex(), g)
	return nil
}

func (idxr *MemoryIndexer) NewSim(ctx context.Context, blockNumber uint64, httpSimClient *alchemy.Client, wsSimClient *alchemy.Client) (*MemoryIndexer, error) {
	if idxr.sim != nil {
		idxr.sim.events.Stop()
	}
	events, err := eventwatcher.New(eventwatcher.Config{
		HTTPClient:    httpSimClient,
		Websocket:     wsSimClient,
		Concurrency:   1,
		LogRange:      1000, //config.IndexerMaxLogRange,
		EpochBlock:    int64(blockNumber),
		Simulated:     true,
		Notifications: idxr.notifications,
	})
	if err != nil {
		return nil, err
	}
	// clone this indexer
	newIdxr := &MemoryIndexer{
		configStore:   idxr.configStore,
		gameStore:     idxr.gameStore,
		stateStore:    idxr.stateStore.Fork(ctx, events, blockNumber),
		sessionStore:  idxr.sessionStore,
		notifications: idxr.notifications,
		events:        events,
		httpClient:    httpSimClient,
		wsClient:      wsSimClient,
	}
	newIdxr.events.Start(ctx)
	<-newIdxr.events.Ready()
	return newIdxr, nil
}

func (idxr *MemoryIndexer) SetNotificationsEnabled(enabled bool, simulated bool) {
	if simulated {
		if idxr.sim != nil {
			idxr.sim.events.SetNotificationsEnabled(enabled)
		}
	} else {
		idxr.events.SetNotificationsEnabled(enabled)
	}
}
func (idxr *MemoryIndexer) Notify(blockNumber int64, simulated bool) {
	if simulated {
		if idxr.sim != nil {
			idxr.sim.events.NotifyFullSync(blockNumber)
		}
	} else {
		idxr.events.NotifyFullSync(blockNumber)
	}
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
	if idxr.game != nil {
		return idxr.game
	}
	return idxr.gameStore.GetGame(id)
}

func (idxr *MemoryIndexer) GetGames() []*model.Game {
	if idxr.game != nil {
		return []*model.Game{idxr.game}
	}
	return idxr.gameStore.GetGames()
}

func (idxr *MemoryIndexer) GetGraph(stateContractAddr common.Address, block int, simulated bool) *model.Graph {
	idx := idxr
	if simulated && idxr.sim != nil {
		idx = idxr.sim
	}
	return idx.stateStore.GetGraph(stateContractAddr, block)
}
func (idxr *MemoryIndexer) GetSession(routerAddr common.Address, sessionID string) *model.Session {
	return idxr.sessionStore.GetSession(routerAddr, sessionID)
}

func (idxr *MemoryIndexer) GetSessions(routerAddr common.Address, owner *string) []*model.Session {
	return idxr.sessionStore.GetSessions(routerAddr, owner)
}
