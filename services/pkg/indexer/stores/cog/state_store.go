package cog

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/benbjohnson/immutable"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/contracts/state"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type OpSet struct {
	Sig string
	Ops []interface{}
}

type StateStore struct {
	wg            map[uint64]*sync.WaitGroup
	graphs        *immutable.Map[uint64, *immutable.Map[string, *model.Graph]]
	latest        uint64
	abi           *abi.ABI
	log           zerolog.Logger
	seenOpSet     map[string]bool
	pendingOpSets []OpSet
	sync.RWMutex
}

func NewStateStore(ctx context.Context, watcher *eventwatcher.Watcher) (*StateStore, error) {
	cabi, err := abi.JSON(strings.NewReader(state.StateABI))
	if err != nil {
		panic(err)
	}
	store := &StateStore{
		graphs:    immutable.NewMap[uint64, *immutable.Map[string, *model.Graph]](nil),
		wg:        map[uint64]*sync.WaitGroup{},
		abi:       &cabi,
		log:       log.With().Str("service", "indexer").Str("component", "statestore").Str("name", "latest").Logger(),
		seenOpSet: map[string]bool{},
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
	defer rs.Unlock()
	graphs, ok := rs.graphs.Get(rs.latest)
	if !ok {
		graphs = immutable.NewMap[string, *model.Graph](nil)
	}
	seenOps := []string{}
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
		g, ok := graphs.Get(rawEvent.Address.Hex())
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
		case "SeenOpSet":
			var evt state.StateSeenOpSet
			if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
				rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
				continue
			}
			seenOps = append(seenOps, hexutil.Encode(evt.Sig))
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
		graphs = graphs.Set(rawEvent.Address.Hex(), g)
	}
	allGraphs := rs.graphs.Set(uint64(block.ToBlock), graphs)

	// only keep last $maxKeep graphs
	// [!] this can be removed once we persist data to disk, but it isn't
	//     practical to keep history while using in-memory storage
	maxKeep := uint64(250)
	oldestToKeep := uint64(0)
	if uint64(block.ToBlock) > maxKeep {
		oldestToKeep = uint64(block.ToBlock) - maxKeep
	}
	itr := allGraphs.Iterator()
	for !itr.Done() {
		k, _, ok := itr.Next()
		if !ok {
			continue
		}
		if k == rs.latest {
			continue
		}
		if k == uint64(block.ToBlock) {
			continue
		}
		if k > oldestToKeep {
			continue
		}
		allGraphs = allGraphs.Delete(k)
	}

	// switch to latest
	rs.graphs = allGraphs
	rs.latest = uint64(block.ToBlock)

	// update pending
	for _, seen := range seenOps {
		fmt.Println("SEEEEEEEEEEEEEEEEEEEEN: ", seen)
		rs.seenOpSet[seen] = true
	}
	newPendingOpSets := []OpSet{}
	for _, opset := range rs.pendingOpSets {
		if rs.seenOpSet[opset.Sig] {
			continue
		}
		newPendingOpSets = append(newPendingOpSets, opset)
	}
	rs.pendingOpSets = newPendingOpSets

	for block, wg := range rs.wg {
		if rs.latest >= block {
			wg.Done()
			delete(rs.wg, block)
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

func (rs *StateStore) GetGraph(stateContractAddr common.Address, block int) *model.Graph {
	rs.RLock()
	defer rs.RUnlock()
	b := rs.latest
	if block != 0 {
		b = uint64(block)
	}

	graphs, ok := rs.graphs.Get(b)
	if !ok {
		return nil
	}
	g, ok := graphs.Get(stateContractAddr.Hex())
	if !ok {
		return nil
	}
	return g
}

func (rs *StateStore) AddPendingOpSet(stateContractAddr common.Address, opset OpSet) {
	rs.Lock()
	defer rs.Unlock()
	rs.pendingOpSets = append(rs.pendingOpSets, opset)
}

func (rs *StateStore) GetPendingGraph(stateContractAddr common.Address, block int) *model.Graph {
	g := rs.GetGraph(stateContractAddr, block)
	if g == nil {
		return nil
	}
	rs.RLock()
	defer rs.RUnlock()
	fmt.Println("GetPendingGraph pendingcount=", len(rs.pendingOpSets), "seencount=", len(rs.seenOpSet))
	for _, opset := range rs.pendingOpSets {
		fmt.Printf("sig=%s", opset.Sig)
		if rs.seenOpSet[opset.Sig] {
			// todo remove from pending
			continue
		}
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
			default:
				err = fmt.Errorf("unexpected evt: %v", evt)
			}
			if err != nil {
				fmt.Printf("ERROR: failed to apply pending op: %v\n", err)
			} else {
				// fmt.Println("....INCLUDED PENDING")
				g = g2
			}
		}
	}
	return g
}
