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
	GetGraph(stateContractAddr common.Address, block int, simulated bool) *model.Graph
	GetSession(routerAddr common.Address, sessionID string) *model.Session
	GetSessions(routerAddr common.Address, owner *string) []*model.Session
	AddPendingOpSet(opset cog.OpSet)
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

	var contractAddrs []common.Address
	empty := common.Address{}
	if config.IndexerStateAddress != empty {
		contractAddrs = append(contractAddrs, config.IndexerStateAddress)
	}
	if config.IndexerGameAddress != empty {
		contractAddrs = append(contractAddrs, config.IndexerGameAddress)
	}
	if config.IndexerRouterAddress != empty {
		contractAddrs = append(contractAddrs, config.IndexerRouterAddress)
	}

	idxr.events, err = eventwatcher.New(eventwatcher.Config{
		HTTPClient:    idxr.httpClient,
		Websocket:     idxr.wsClient,
		LogRange:      config.IndexerMaxLogRange,
		Notifications: notifications,
		Addresses:     contractAddrs,
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

func (idxr *MemoryIndexer) Ready() chan struct{} {
	return idxr.events.Ready()
}

func (idxr *MemoryIndexer) GetGame(id string) *model.Game {
	return idxr.gameStore.GetGame(id)
}

func (idxr *MemoryIndexer) GetGames() []*model.Game {
	return idxr.gameStore.GetGames()
}

func (idxr *MemoryIndexer) AddPendingOpSet(opset cog.OpSet) {
	idxr.stateStore.AddPendingOpSet(opset)
}

func (idxr *MemoryIndexer) GetGraph(stateContractAddr common.Address, block int, simulated bool) *model.Graph {
	if simulated {
		return idxr.stateStore.GetPendingGraph()
	}
	return idxr.stateStore.GetGraph()
}
func (idxr *MemoryIndexer) GetSession(routerAddr common.Address, sessionID string) *model.Session {
	return idxr.sessionStore.GetSession(routerAddr, sessionID)
}

func (idxr *MemoryIndexer) GetSessions(routerAddr common.Address, owner *string) []*model.Session {
	return idxr.sessionStore.GetSessions(routerAddr, owner)
}
