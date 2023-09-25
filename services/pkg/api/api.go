package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/hooks"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/api/resolver"
	"github.com/playmint/ds-node/pkg/config"
	"github.com/playmint/ds-node/pkg/indexer"
	"github.com/playmint/ds-node/pkg/sequencer"
	"github.com/rs/cors"
)

type Server struct {
	Indexer   indexer.Indexer
	Sequencer sequencer.Sequencer
}

func (api *Server) Start(ctx context.Context, subscriptions *model.Subscriptions) error {

	// configure resolver
	resolver := &resolver.Resolver{
		Indexer:       api.Indexer,
		Sequencer:     api.Sequencer,
		Subscriptions: subscriptions,
	}

	// start server

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // FIXME: this is lazy and potentially bad, allow setting from env/config
		AllowCredentials: true,
		MaxAge:           86400,
		Debug:            false,
	}).Handler)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 15 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
	// srv.SetQueryCache(lru.New(1000))
	srv.Use(hooks.Prometheus{})
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	addr := fmt.Sprintf(":%d", config.APIPort)
	log.Info().Str("service", "api").Msg("ready")
	return http.ListenAndServe(addr, router)
}
