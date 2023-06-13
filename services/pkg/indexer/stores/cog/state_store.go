package cog

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/contracts/state"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type StateStore struct {
	graphs        map[common.Address]*model.Graph
	live          map[common.Address]*model.Graph
	abi           *abi.ABI
	notifications chan interface{}
	log           zerolog.Logger
	name          string
	sync.RWMutex
}

func NewStateStore(ctx context.Context, watcher *eventwatcher.Watcher, notifications chan interface{}) (*StateStore, error) {
	cabi, err := abi.JSON(strings.NewReader(state.StateABI))
	if err != nil {
		panic(err)
	}
	store := &StateStore{
		graphs:        map[common.Address]*model.Graph{},
		live:          map[common.Address]*model.Graph{},
		abi:           &cabi,
		notifications: notifications,
		name:          "og",
		log:           log.With().Str("service", "indexer").Str("component", "statestore").Logger(),
	}
	store.watch(ctx, watcher)
	return store, nil
}

func (rs *StateStore) Fork(ctx context.Context, watcher *eventwatcher.Watcher) *StateStore {
	newStore := &StateStore{
		graphs:        map[common.Address]*model.Graph{},
		live:          map[common.Address]*model.Graph{},
		abi:           rs.abi,
		notifications: rs.notifications,
		log:           rs.log,
		name:          "forked",
	}
	for addr, g := range rs.graphs {
		newStore.graphs[addr] = g
	}
	for addr, g := range rs.live {
		newStore.live[addr] = g
	}
	newStore.watch(ctx, watcher)
	return newStore
}

func (rs *StateStore) emitStateEvent(stateAddr common.Address, stateEvent model.Event) {
	// rs.notifications <- &model.StateEvent{
	// 	Event:   stateEvent,
	// 	StateID: stateAddr.Hex(),
	// }
}

func (rs *StateStore) watch(ctx context.Context, watcher *eventwatcher.Watcher) {
	// watch all events from all contracts that match the GameDeployed topic
	query := [][]interface{}{{
		rs.abi.Events["EdgeRemove"].ID,
		rs.abi.Events["EdgeSet"].ID,
		rs.abi.Events["NodeTypeRegister"].ID,
		rs.abi.Events["EdgeTypeRegister"].ID,
		rs.abi.Events["AnnotationSet"].ID,
	}}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		panic(err)
	}
	events := watcher.SubscribeTopic(topics[0])
	go rs.watchLoop(ctx, events)
}

func (rs *StateStore) watchLoop(ctx context.Context, blocks chan *eventwatcher.LogBatch) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-blocks:
			for _, rawEvent := range block.Logs {
				rs.log.Info().Str("indexer", rs.name).Msgf("got event %v", rawEvent)
				eventABI, err := rs.abi.EventByID(rawEvent.Topics[0])
				if err != nil {
					rs.log.Warn().Err(err).Msg("unhandleable event topic")
					continue
				}
				rs.log.Debug().Msgf("recv %v", eventABI.RawName)
				switch eventABI.RawName {
				case "AnnotationSet":
					var evt state.StateAnnotationSet
					if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
						rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
						continue
					}
					evt.Raw = rawEvent
					err := rs.setAnnotation(&evt)
					if err != nil {
						rs.log.Error().Err(err).Msgf("failed process %T event", evt)
					}
				case "EdgeSet":
					var evt state.StateEdgeSet
					if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
						rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
						continue
					}
					evt.Raw = rawEvent
					err := rs.setEdge(&evt)
					if err != nil {
						rs.log.Error().Err(err).Msgf("failed process %T event", evt)
					}
				case "EdgeRemove":
					var evt state.StateEdgeRemove
					if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
						rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
						continue
					}
					evt.Raw = rawEvent
					err := rs.removeEdge(&evt)
					if err != nil {
						rs.log.Error().Err(err).Msgf("failed process %T event", evt)
					}
				case "NodeTypeRegister":
					var evt state.StateNodeTypeRegister
					if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
						rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
						continue
					}
					evt.Raw = rawEvent
					err := rs.setNodeType(&evt)
					if err != nil {
						rs.log.Error().Err(err).Msgf("failed process %T event", evt)
					}
				case "EdgeTypeRegister":
					var evt state.StateEdgeTypeRegister
					if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
						rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
						continue
					}
					evt.Raw = rawEvent
					err := rs.setEdgeType(&evt)
					if err != nil {
						rs.log.Error().Err(err).Msgf("failed process %T event", evt)
					}
				default:
					rs.log.Warn().Msgf("ignoring unhandled event type %v", eventABI)
				}

			}
			for addr, g := range rs.graphs {
				rs.live[addr] = g
				rs.notifications <- &model.StateEvent{
					Event:   &model.SetAnnotationEvent{}, // just using this as a proxy to mean "an update arrived" for now
					StateID: addr.Hex(),
				}
			}
		}
	}
}

func (rs *StateStore) setEdgeType(evt *state.StateEdgeTypeRegister) error {
	g, ok := rs.graphs[evt.Raw.Address]
	if !ok {
		g = model.NewGraph(evt.Raw.BlockNumber)
	}
	g = g.SetRelData(evt)
	// commit
	rs.graphs[evt.Raw.Address] = g
	return nil
}

func (rs *StateStore) setNodeType(evt *state.StateNodeTypeRegister) error {
	g, ok := rs.graphs[evt.Raw.Address]
	if !ok {
		g = model.NewGraph(evt.Raw.BlockNumber)
	}
	g = g.SetKindData(evt)
	// commit
	rs.graphs[evt.Raw.Address] = g
	return nil
}

func (rs *StateStore) setAnnotation(evt *state.StateAnnotationSet) error {
	rs.Lock()
	defer rs.Unlock()

	if evt.Raw.Removed {
		// hmmm blockchain reorg occured so what do we do?
		// just ignore for now, but this probably needs more thought
		return nil
	}

	nodeID := hexutil.Encode(evt.Id[:])
	ref := hexutil.Encode(evt.Ref[:])

	// update graph
	g, ok := rs.graphs[evt.Raw.Address]
	if !ok {
		g = model.NewGraph(evt.Raw.BlockNumber)
	}
	g = g.SetAnnotationData(
		nodeID,
		evt.Label,
		ref,
		evt.Data,
		evt.Raw.BlockNumber,
	)

	// commit the graph
	rs.graphs[evt.Raw.Address] = g
	rs.emitStateEvent(evt.Raw.Address, &model.SetAnnotationEvent{
		ID:   fmt.Sprintf("%d-%d", evt.Raw.BlockNumber, evt.Raw.Index),
		From: nodeID,
		Name: evt.Label,
	})

	return nil
}

func (rs *StateStore) setEdge(evt *state.StateEdgeSet) error {
	rs.Lock()
	defer rs.Unlock()

	if evt.Raw.Removed {
		// hmmm blockchain reorg occured so what do we do?
		// just ignore for now, but this probably needs more thought
		return nil
	}

	relID := hexutil.Encode(evt.RelID[:])
	relKey := evt.RelKey
	srcNodeID := hexutil.Encode(evt.SrcNodeID[:])
	dstNodeID := hexutil.Encode(evt.DstNodeID[:])
	weight := evt.Weight

	// update graph
	g, ok := rs.graphs[evt.Raw.Address]
	if !ok {
		g = model.NewGraph(evt.Raw.BlockNumber)
	}
	g = g.SetEdge(
		relID,
		relKey,
		srcNodeID,
		dstNodeID,
		weight,
		evt.Raw.BlockNumber,
	)

	// commit the graph
	rs.graphs[evt.Raw.Address] = g
	rs.emitStateEvent(evt.Raw.Address, &model.SetEdgeEvent{
		ID:   fmt.Sprintf("%d-%d", evt.Raw.BlockNumber, evt.Raw.Index),
		From: srcNodeID,
		To:   dstNodeID,
		Rel:  relID,
		Key:  int(relKey),
	})

	return nil
}

func (rs *StateStore) removeEdge(evt *state.StateEdgeRemove) error {
	rs.Lock()
	defer rs.Unlock()

	if evt.Raw.Removed {
		// hmmm blockchain reorg occured so what do we do?
		// just ignore for now, but this probably needs more thought
		return nil
	}

	relID := hexutil.Encode(evt.RelID[:])
	relKey := evt.RelKey
	srcNodeID := hexutil.Encode(evt.SrcNodeID[:])

	// update graph
	g, ok := rs.graphs[evt.Raw.Address]
	if !ok {
		g = model.NewGraph(evt.Raw.BlockNumber)
	}
	g = g.RemoveEdge(
		relID,
		relKey,
		srcNodeID,
		evt.Raw.BlockNumber,
	)

	// commit the graph
	rs.graphs[evt.Raw.Address] = g
	rs.emitStateEvent(evt.Raw.Address, &model.RemoveEdgeEvent{
		ID:   fmt.Sprintf("%d-%d", evt.Raw.BlockNumber, evt.Raw.Index),
		From: srcNodeID,
		Rel:  relID,
		Key:  int(relKey),
	})

	return nil
}

func (rs *StateStore) GetGraph(stateContractAddr common.Address) *model.Graph {
	rs.RLock()
	defer rs.RUnlock()
	return rs.live[stateContractAddr]
}
