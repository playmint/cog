package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/sequencer"
)

func (r *mutationResolver) Signup(ctx context.Context, gameID string, authorization string) (bool, error) {
	// if err := r.Sequencer.Signup(ctx, common.HexToAddress(account)); err != nil {
	// 	return false, err
	// }
	// this is currently a noop - just signin
	return true, nil
}

func (r *mutationResolver) Signin(ctx context.Context, gameID string, session string, ttl int, scope string, authorization string) (bool, error) {
	game := r.Indexer.GetGame(gameID)
	if game == nil {
		return false, fmt.Errorf("no game found with id %v", game)
	}
	scopesBig, err := hexutil.DecodeBig(scope)
	if err != nil {
		return false, err
	}
	scopes := uint32(scopesBig.Uint64())
	if err := r.Sequencer.Signin(ctx,
		game.RouterAddress,
		game.DispatcherAddress,
		common.HexToAddress(session),
		uint32(ttl),
		scopes,
		authorization,
	); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) Signout(ctx context.Context, gameID string, session string, authorization string) (bool, error) {
	game := r.Indexer.GetGame(gameID)
	if game == nil {
		return false, fmt.Errorf("no game found with id %v", game)
	}
	if err := r.Sequencer.Signout(ctx, game.RouterAddress, common.HexToAddress(session), authorization); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) Dispatch(ctx context.Context, gameID string, actions []string, authorization string) (*model.ActionTransaction, error) {
	game := r.Indexer.GetGame(gameID)
	if game == nil {
		return nil, fmt.Errorf("no game found with id %v", game)
	}
	// extract the signer of the action
	signer, err := sequencer.ValidateActions(actions, authorization)
	if err != nil {
		return nil, err
	} else if signer == nil {
		return nil, fmt.Errorf("invalid action: failed to extract signer")
	}
	// check that the signer has a session
	session := r.Indexer.GetSession(game.RouterAddress, (*signer).Hex())
	if session == nil {
		return nil, fmt.Errorf("invalid action: no session for signer or invalid signature")
	}
	// TODO: check the session is not expired, needs the current block number
	// push it to the pending batch
	tx, err := r.Sequencer.Enqueue(
		ctx,
		game.RouterAddress,
		common.HexToAddress(session.Owner),
		actions,
		authorization,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
