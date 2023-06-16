package cog

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/benbjohnson/immutable"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/playmint/ds-node/pkg/contracts/game"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const LATEST = "latest"

type GameStore struct {
	games         *immutable.Map[string, *model.Game]
	latest        *model.Game
	latestByName  map[string]*model.Game
	abi           *abi.ABI
	events        *eventwatcher.Watcher
	client        *alchemy.Client
	notifications chan interface{}
	log           zerolog.Logger
	sync.RWMutex
}

func NewGameStore(ctx context.Context, client *alchemy.Client, watcher *eventwatcher.Watcher, notifications chan interface{}) (*GameStore, error) {
	cabi, err := abi.JSON(strings.NewReader(game.BaseGameABI))
	if err != nil {
		return nil, err
	}
	store := &GameStore{
		client:        client,
		abi:           &cabi,
		events:        watcher,
		notifications: notifications,
		games:         immutable.NewMap[string, *model.Game](nil),
		log:           log.With().Str("service", "indexer").Str("component", "gamestore").Logger(),
		latestByName:  map[string]*model.Game{},
	}

	// watch all events from all contracts that match the GameDeployed topic
	query := [][]interface{}{{cabi.Events["GameDeployed"].ID}}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, err
	}
	queue := watcher.SubscribeTopic(topics[0])

	go store.watch(ctx, queue)
	return store, nil
}

func (rs *GameStore) Fork(ctx context.Context, watcher *eventwatcher.Watcher, client *alchemy.Client) *GameStore {
	return rs
}

func (rs *GameStore) emitGame(game *model.Game) {
	rs.notifications <- game
}

func (rs *GameStore) watch(ctx context.Context, blocks chan *eventwatcher.LogBatch) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-blocks:
			for _, rawEvent := range block.Logs {
				eventABI, err := rs.abi.EventByID(rawEvent.Topics[0])
				if err != nil {
					rs.log.Debug().Msgf("unhandleable event topic: %v", err)
					continue
				}
				switch eventABI.RawName {
				case "GameDeployed":
					var evt game.BaseGameGameDeployed
					if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
						rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
						continue
					}
					evt.Raw = rawEvent
					err := rs.setGame(&evt)
					if err != nil {
						rs.log.Error().Err(err).Msgf("failed process %T event", evt)
					}
				default:
					rs.log.Warn().Msgf("ignoring unhandled event type %v", eventABI)
				}
			}
		}
	}
}

// An action was registered, update the mapping for the game
func (rs *GameStore) setGame(evt *game.BaseGameGameDeployed) error {
	rs.Lock()
	defer rs.Unlock()

	if evt.Raw.Removed {
		// hmmm blockchain reorg occured so what do we do?
		// just ignore for now, but this probably needs more thought
		return nil
	}

	// fetch the metadata
	gameContract, err := game.NewBaseGame(evt.Raw.Address, rs.client)
	if err != nil {
		return err
	}
	meta, err := gameContract.GetMetadata(nil)
	if err != nil {
		return err
	}

	// create new game object
	game := &model.Game{
		ID:                evt.Raw.Address.Hex(),
		Name:              meta.Name,
		URL:               meta.Url,
		DispatcherAddress: evt.DispatcherAddr,
		RouterAddress:     evt.RouterAddr,
		StateAddress:      evt.StateAddr,
	}

	// add to the set
	rs.games = rs.games.Set(game.ID, game)

	// update the "LATEST" tag, handy in local development
	// TODO: probably disable this from config as it's a bit weird in prod
	rs.latest = game
	rs.latestByName[meta.Name] = game

	rs.emitGame(game)

	return nil
}

func (rs *GameStore) GetGame(id string) *model.Game {
	rs.RLock()
	defer rs.RUnlock()

	// [!] deprecated, use latestByName going forward
	if id == LATEST {
		return rs.latest
	}
	// fetch latest deployment by name
	// this is only for local-dev convinence
	// TODO: disable this when in production mode
	if game, ok := rs.latestByName[id]; ok {
		return game
	}

	// fetch by full id
	game, ok := rs.games.Get(id)
	if !ok {
		return nil
	}
	return game
}

func (rs *GameStore) GetGames() []*model.Game {
	rs.RLock()
	defer rs.RUnlock()

	games := []*model.Game{}
	itr := rs.games.Iterator()
	for !itr.Done() {
		_, game, ok := itr.Next()
		if !ok {
			continue
		}
		games = append(games, game)
	}
	return games
}

func unpackLog(cabi *abi.ABI, out interface{}, event string, log types.Log) error {
	if log.Topics[0] != cabi.Events[event].ID {
		return fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := cabi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range cabi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}
