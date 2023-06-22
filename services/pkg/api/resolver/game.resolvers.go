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

func (r *gameResolver) Subscribers(ctx context.Context, obj *model.Game) (int, error) {
	if r.Subscriptions == nil && r.Subscriptions.Events == nil {
		return 0, nil
	}
	subs, ok := r.Subscriptions.Events[obj.StateAddress.Hex()]
	if !ok {
		return 0, nil
	}
	return len(subs), nil
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

type gameResolver struct{ *Resolver }
