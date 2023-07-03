package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/model"
)

func (r *edgeResolver) Weight(ctx context.Context, obj *model.Edge) (int, error) {
	return obj.WeightInt(), nil
}

func (r *edgeResolver) Key(ctx context.Context, obj *model.Edge) (int, error) {
	return obj.KeyInt(), nil
}

func (r *stateResolver) Block(ctx context.Context, obj *model.State) (int, error) {
	graph := r.Indexer.GetGraph(common.HexToAddress(obj.ID), obj.Block, obj.Simulated)
	if graph == nil {
		graph = model.NewGraph(0)
	}
	return int(graph.BlockNumber()), nil // bad cast will stop working in ~100yrs
}

func (r *stateResolver) Nodes(ctx context.Context, obj *model.State, match *model.Match) ([]*model.Node, error) {
	graph := r.Indexer.GetGraph(common.HexToAddress(obj.ID), obj.Block, obj.Simulated)
	if graph == nil {
		graph = model.NewGraph(0)
	}
	return graph.GetNodes(match), nil
}

func (r *stateResolver) Node(ctx context.Context, obj *model.State, match *model.Match) (*model.Node, error) {
	graph := r.Indexer.GetGraph(common.HexToAddress(obj.ID), obj.Block, obj.Simulated)
	if graph == nil {
		graph = model.NewGraph(0)
	}
	return graph.GetNode(match), nil
}

func (r *stateResolver) JSON(ctx context.Context, obj *model.State) (string, error) {
	graph := r.Indexer.GetGraph(common.HexToAddress(obj.ID), obj.Block, obj.Simulated)
	if graph == nil {
		graph = model.NewGraph(0)
	}
	return graph.Dump()
}

// Edge returns generated.EdgeResolver implementation.
func (r *Resolver) Edge() generated.EdgeResolver { return &edgeResolver{r} }

// State returns generated.StateResolver implementation.
func (r *Resolver) State() generated.StateResolver { return &stateResolver{r} }

type edgeResolver struct{ *Resolver }
type stateResolver struct{ *Resolver }
