package resolver

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/indexer"
	"github.com/playmint/ds-node/pkg/sequencer"
)

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	Indexer       indexer.Indexer
	Sequencer     sequencer.Sequencer
	Subscriptions *model.Subscriptions
}
