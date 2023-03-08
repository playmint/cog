package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/model"
)

func (r *queryResolver) Game(ctx context.Context, id string) (*model.Game, error) {
	game := r.Indexer.GetGame(id)
	if game == nil {
		return nil, fmt.Errorf("no game found with id %v", id)
	}
	return game, nil
}

func (r *queryResolver) Games(ctx context.Context) ([]*model.Game, error) {
	return r.Indexer.GetGames(), nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
