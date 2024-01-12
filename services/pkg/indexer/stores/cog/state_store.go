package cog

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/contracts/state"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type OpSet struct {
	Expires int64
	Sig     string
	Ops     []interface{}
}

type StateStore struct {
	graph         *model.Graph
	pendingGraph  *model.Graph
	abi           *abi.ABI
	log           zerolog.Logger
	notifications chan interface{}
	pendingOpSets []OpSet
	sync.RWMutex
}

func NewStateStore(ctx context.Context, watcher *eventwatcher.Watcher, notifications chan interface{}) (*StateStore, error) {
	cabi, err := abi.JSON(strings.NewReader(state.StateABI))
	if err != nil {
		panic(err)
	}
	store := &StateStore{
		abi:           &cabi,
		log:           log.With().Str("service", "indexer").Str("component", "statestore").Str("name", "latest").Logger(),
		notifications: notifications,
	}
	store.watch(ctx, watcher)
	return store, nil
}

func (rs *StateStore) watch(ctx context.Context, watcher *eventwatcher.Watcher) {
	// watch all events from all contracts that match the GameDeployed topic
	query := [][]interface{}{{
		rs.abi.Events["EdgeRemove"].ID,
		rs.abi.Events["EdgeSet"].ID,
		rs.abi.Events["NodeTypeRegister"].ID,
		rs.abi.Events["EdgeTypeRegister"].ID,
		rs.abi.Events["AnnotationSet"].ID,
		rs.abi.Events["DataSet"].ID,
		rs.abi.Events["SeenOpSet"].ID,
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
	g := rs.graph

	if g == nil {
		g = model.NewGraph(0)
	}

	seenOps := map[string]bool{}
	for _, rawEvent := range block.Logs {
		if rawEvent.Removed {
			// FIXME: ignoring reorg
			rs.log.Warn().Msgf("unhandled reorg %v", rawEvent)
			continue
		}
		eventABI, err := rs.abi.EventByID(rawEvent.Topics[0])
		if err != nil {
			rs.log.Debug().Msgf("unhandleable event topic: %v", err)
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
			g, err = rs.setAnnotation(g, &evt)
			if err != nil {
				rs.log.Error().Err(err).Msgf("failed process %T event", evt)
			}
		case "DataSet":
			var evt state.StateDataSet
			if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
				rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
				continue
			}
			evt.Raw = rawEvent
			g, err = rs.setData(g, &evt)
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
		case "SeenOpSet":
			var evt state.StateSeenOpSet
			if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
				rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
				continue
			}
			seenOps[hexutil.Encode(evt.Sig)] = true
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
	}

	// update
	rs.pendingOpSets = rs.removePendingOpSets(rs.pendingOpSets, seenOps, block.ToBlock)
	rs.graph = g
	rs.pendingGraph = rs.rebuildPendingGraph()
	rs.Unlock()

	// notify

	sigs := []string{}
	for sig := range seenOps {
		sigs = append(sigs, sig)
	}
	rs.Notify(int(block.ToBlock), sigs, false)
}

func (rs *StateStore) Notify(blockNumber int, sigs []string, simulated bool) {
	rs.notifications <- &model.BlockEvent{
		ID:        fmt.Sprintf("block-%d", blockNumber),
		Block:     blockNumber,
		Sigs:      sigs,
		Simulated: simulated,
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

func (rs *StateStore) setData(g *model.Graph, evt *state.StateDataSet) (*model.Graph, error) {
	nodeID := hexutil.Encode(evt.Id[:])
	value := hexutil.Encode(evt.Data[:])

	// update graph
	g = g.SetData(
		nodeID,
		evt.Label,
		value,
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

func (rs *StateStore) GetGraph() *model.Graph {
	rs.Lock()
	defer rs.Unlock()
	return rs.graph
}

func (rs *StateStore) AddPendingOpSet(estimatedBlockNumber int, opset OpSet) {
	// default expiry to ~30 blocks in future this means we will stop waiting
	// for the pending sig to arrive if we don't hear anything within about 1m
	if opset.Expires == 0 {
		opset.Expires = int64(estimatedBlockNumber + 30)
	}

	rs.Lock()
	rs.pendingOpSets = append(rs.pendingOpSets, opset)
	rs.pendingGraph = rs.rebuildPendingGraph()
	rs.Unlock()

	rs.Notify(estimatedBlockNumber, []string{opset.Sig}, true)
}

func (rs *StateStore) RemovePendingOpSets(seenOps map[string]bool) {
	rs.Lock()
	defer rs.Unlock()
	rs.pendingOpSets = rs.removePendingOpSets(rs.pendingOpSets, seenOps, -1)
}

func (rs *StateStore) removePendingOpSets(existingOpSets []OpSet, seenOps map[string]bool, currentBlock int64) []OpSet {
	newPendingOpSets := []OpSet{}
	for _, opset := range existingOpSets {
		if seenOps[opset.Sig] {
			continue
		}
		if currentBlock > 0 && opset.Expires > 0 && currentBlock > opset.Expires {
			continue
		}
		newPendingOpSets = append(newPendingOpSets, opset)
	}
	return newPendingOpSets
}

func (rs *StateStore) GetPendingGraph() *model.Graph {
	return rs.pendingGraph
}

func (rs *StateStore) rebuildPendingGraph() *model.Graph {
	g := rs.graph
	if g == nil {
		return nil
	}
	for _, opset := range rs.pendingOpSets {
		for _, op := range opset.Ops {
			var err error
			var g2 *model.Graph
			switch evt := op.(type) {
			case *state.StateEdgeSet:
				g2, err = rs.setEdge(g, evt)
			case *state.StateEdgeRemove:
				g2, err = rs.removeEdge(g, evt)
			case *state.StateAnnotationSet:
				g2, err = rs.setAnnotation(g, evt)
			case *state.StateDataSet:
				g2, err = rs.setData(g, evt)
			default:
				err = fmt.Errorf("unexpected evt: %v", evt)
			}
			if err != nil {
				fmt.Printf("ERROR: failed to apply pending op: %v\n", err)
			} else {
				g = g2
			}
		}
	}
	return g
}
