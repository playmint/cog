package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/model"
)

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

// State returns generated.StateResolver implementation.
func (r *Resolver) State() generated.StateResolver { return &stateResolver{r} }

type stateResolver struct{ *Resolver }
