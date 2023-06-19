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
	wg            map[uint64]*sync.WaitGroup
	graphs        map[uint64]map[common.Address]*model.Graph
	latest        uint64
	abi           *abi.ABI
	notifications chan interface{}
	log           zerolog.Logger
	sync.RWMutex
}

func NewStateStore(ctx context.Context, watcher *eventwatcher.Watcher, notifications chan interface{}) (*StateStore, error) {
	cabi, err := abi.JSON(strings.NewReader(state.StateABI))
	if err != nil {
		panic(err)
	}
	store := &StateStore{
		graphs:        map[uint64]map[common.Address]*model.Graph{},
		wg:            map[uint64]*sync.WaitGroup{},
		abi:           &cabi,
		notifications: notifications,
		log:           log.With().Str("service", "indexer").Str("component", "statestore").Str("name", "latest").Logger(),
	}
	store.watch(ctx, watcher)
	return store, nil
}

func (rs *StateStore) Fork(ctx context.Context, watcher *eventwatcher.Watcher, blockNumber uint64) *StateStore {
	rs.Lock()
	defer rs.Unlock()

	// wait til store contains blockNumber if we are behind
	if rs.latest < blockNumber {
		wg, ok := rs.wg[blockNumber]
		if !ok {
			wg = &sync.WaitGroup{}
			wg.Add(1)
			rs.wg[blockNumber] = wg
		}
		rs.log.Warn().Uint64("block", blockNumber).Msgf("fork-wait")
		wg.Wait()
	}

	newStore := &StateStore{
		graphs:        map[uint64]map[common.Address]*model.Graph{},
		wg:            map[uint64]*sync.WaitGroup{},
		latest:        blockNumber,
		abi:           rs.abi,
		notifications: rs.notifications,
		log:           rs.log.With().Str("name", fmt.Sprintf("fork-%d", blockNumber)).Logger(),
	}
	newStore.graphs[blockNumber] = map[common.Address]*model.Graph{}
	for addr, g := range rs.graphs[blockNumber] {
		newStore.graphs[blockNumber][addr] = g
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

func (rs *StateStore) processBlock(ctx context.Context, block *eventwatcher.LogBatch) {
	rs.Lock()
	defer rs.Unlock()
	graphs, ok := rs.graphs[rs.latest]
	if !ok {
		graphs = map[common.Address]*model.Graph{}
	}
	for _, rawEvent := range block.Logs {
		if rawEvent.Removed {
			// FIXME: ignoring reorg
			rs.log.Warn().Msgf("unhandled reorg", rawEvent)
			continue
		}
		eventABI, err := rs.abi.EventByID(rawEvent.Topics[0])
		if err != nil {
			rs.log.Debug().Msgf("unhandleable event topic: %v", err)
			continue
		}
		g, ok := rs.graphs[rs.latest][rawEvent.Address]
		if !ok {
			g = model.NewGraph(rawEvent.BlockNumber)
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
			g, err = rs.setAnnotation(g, &evt)
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
			g, err = rs.setEdge(g, &evt)
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
			g, err = rs.removeEdge(g, &evt)
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
			g, err = rs.setNodeType(g, &evt)
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
			g, err = rs.setEdgeType(g, &evt)
			if err != nil {
				rs.log.Error().Err(err).Msgf("failed process %T event", evt)
			}
		default:
			rs.log.Warn().Msgf("ignoring unhandled event type %v", eventABI)
		}
		graphs[rawEvent.Address] = g
	}
	rs.graphs[uint64(block.ToBlock)] = map[common.Address]*model.Graph{}
	for addr, g := range graphs {
		rs.graphs[uint64(block.ToBlock)][addr] = g
	}
	rs.latest = uint64(block.ToBlock)
	rs.log.Warn().Msgf("latest now %d", rs.latest)

	wg, ok := rs.wg[uint64(block.ToBlock)]
	if ok {
		wg.Done()
	}
	// for addr := range graphs {
	// 	rs.notifications <- &model.StateEvent{
	// 		Event:   &model.SetAnnotationEvent{}, // FIXME: just using this as a proxy to mean "an update arrived" for now
	// 		StateID: addr.Hex(),
	// 	}
	// }
}

func (rs *StateStore) NotifyAll() {
	rs.RLock()
	defer rs.RUnlock()
	for addr := range rs.graphs[rs.latest] {
		rs.notifications <- &model.StateEvent{
			Event:   &model.SetAnnotationEvent{}, // FIXME: just using this as a proxy to mean "an update arrived" for now
			StateID: addr.Hex(),
		}
	}
}

func (rs *StateStore) watchLoop(ctx context.Context, blocks chan *eventwatcher.LogBatch) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-blocks:
			rs.processBlock(ctx, block)
		}
	}
}

func (rs *StateStore) setEdgeType(g *model.Graph, evt *state.StateEdgeTypeRegister) (*model.Graph, error) {
	g = g.SetRelData(evt)
	return g, nil
}

func (rs *StateStore) setNodeType(g *model.Graph, evt *state.StateNodeTypeRegister) (*model.Graph, error) {
	g = g.SetKindData(evt)
	return g, nil
}

func (rs *StateStore) setAnnotation(g *model.Graph, evt *state.StateAnnotationSet) (*model.Graph, error) {
	nodeID := hexutil.Encode(evt.Id[:])
	ref := hexutil.Encode(evt.Ref[:])

	// update graph
	g = g.SetAnnotationData(
		nodeID,
		evt.Label,
		ref,
		evt.Data,
		evt.Raw.BlockNumber,
	)

	return g, nil
}

func (rs *StateStore) setEdge(g *model.Graph, evt *state.StateEdgeSet) (*model.Graph, error) {
	relID := hexutil.Encode(evt.RelID[:])
	relKey := evt.RelKey
	srcNodeID := hexutil.Encode(evt.SrcNodeID[:])
	dstNodeID := hexutil.Encode(evt.DstNodeID[:])
	weight := evt.Weight

	// update graph
	g = g.SetEdge(
		relID,
		relKey,
		srcNodeID,
		dstNodeID,
		weight,
		evt.Raw.BlockNumber,
	)

	// commit the graph
	rs.emitStateEvent(evt.Raw.Address, &model.SetEdgeEvent{
		ID:   fmt.Sprintf("%d-%d", evt.Raw.BlockNumber, evt.Raw.Index),
		From: srcNodeID,
		To:   dstNodeID,
		Rel:  relID,
		Key:  int(relKey),
	})

	return g, nil
}

func (rs *StateStore) removeEdge(g *model.Graph, evt *state.StateEdgeRemove) (*model.Graph, error) {
	relID := hexutil.Encode(evt.RelID[:])
	relKey := evt.RelKey
	srcNodeID := hexutil.Encode(evt.SrcNodeID[:])

	// update graph
	g = g.RemoveEdge(
		relID,
		relKey,
		srcNodeID,
		evt.Raw.BlockNumber,
	)

	return g, nil
}

func (rs *StateStore) GetGraph(stateContractAddr common.Address) *model.Graph {
	rs.Lock()
	defer rs.Unlock()
	return rs.graphs[rs.latest][stateContractAddr]
}
