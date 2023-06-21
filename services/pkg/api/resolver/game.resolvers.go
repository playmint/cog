package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/model"
)

func (r *gameResolver) State(ctx context.Context, obj *model.Game, block *int, simulated *bool) (*model.State, error) {
	if obj == nil {
		return nil, fmt.Errorf("nil game")
	}
	return obj.State(block, simulated), nil
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

type gameResolver struct{ *Resolver }
