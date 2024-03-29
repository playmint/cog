package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/model"
)

func (r *actionTransactionResolver) Nonce(ctx context.Context, obj *model.ActionTransaction) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *routerResolver) Sessions(ctx context.Context, obj *model.Router, owner *string) ([]*model.Session, error) {
	return r.Indexer.GetSessions(
		common.HexToAddress(obj.ID),
		owner,
	), nil
}

func (r *routerResolver) Session(ctx context.Context, obj *model.Router, id string) (*model.Session, error) {
	return r.Indexer.GetSession(
		common.HexToAddress(obj.ID),
		id,
	), nil
}

func (r *routerResolver) Transactions(ctx context.Context, obj *model.Router, owner *string, status []model.ActionTransactionStatus) ([]*model.ActionTransaction, error) {
	return r.Sequencer.GetTransactions(
		common.HexToAddress(obj.ID),
		owner, status,
	)
}

func (r *routerResolver) Transaction(ctx context.Context, obj *model.Router, id string) (*model.ActionTransaction, error) {
	return r.Sequencer.GetTransaction(
		common.HexToAddress(obj.ID),
		id,
	)
}

// ActionTransaction returns generated.ActionTransactionResolver implementation.
func (r *Resolver) ActionTransaction() generated.ActionTransactionResolver {
	return &actionTransactionResolver{r}
}

// Router returns generated.RouterResolver implementation.
func (r *Resolver) Router() generated.RouterResolver { return &routerResolver{r} }

type actionTransactionResolver struct{ *Resolver }
type routerResolver struct{ *Resolver }
