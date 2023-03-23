package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/model"
)

func (r *subscriptionResolver) Events(ctx context.Context, gameID string) (<-chan model.Event, error) {
	game := r.Indexer.GetGame(gameID)
	if game == nil {
		return nil, fmt.Errorf("no game found with id %v", gameID)
	}
	return r.Subscriptions.SubscribeStateEvent(ctx, game.StateAddress.Hex()), nil
}

func (r *subscriptionResolver) Transaction(ctx context.Context, gameID string, owner *string) (<-chan *model.ActionTransaction, error) {
	game := r.Indexer.GetGame(gameID)
	if game == nil {
		return nil, fmt.Errorf("no game found with id %v", gameID)
	}
	return r.Subscriptions.SubscribeTransaction(ctx, game.RouterAddress.Hex(), owner), nil
}

func (r *subscriptionResolver) Session(ctx context.Context, gameID string, owner *string) (<-chan *model.Session, error) {
	game := r.Indexer.GetGame(gameID)
	if game == nil {
		return nil, fmt.Errorf("no game found with id %v", gameID)
	}
	return r.Subscriptions.SubscribeSession(ctx, game.RouterAddress.Hex(), owner), nil
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
