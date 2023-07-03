//go:build js && wasm
// +build js,wasm

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	wasmhttp "github.com/nlepage/go-wasm-http-server"
	"github.com/playmint/ds-node/pkg/api"
	"github.com/playmint/ds-node/pkg/api/generated"
	"github.com/playmint/ds-node/pkg/api/resolver"
	"github.com/playmint/ds-node/pkg/indexer"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var counter int32

func Main(ctx context.Context) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// fake event watcher
	events, err := eventwatcher.New(eventwatcher.Config{
		LogRange:             1000,
		Concurrency:          1, // config.IndexerMaxConcurrency, - NodeSet/EdgeSet cannot arrive out of order yet
		NotificationsEnabled: false,
	})
	if err != nil {
		return err
	}

	// start an indexer
	idxr, err := indexer.NewMemoryIndexer2(ctx, events)
	if err != nil {
		return err
	}

	// start graphql api server
	api := api.Server{
		Indexer: idxr,
	}

	// configure resolver
	resolver := &resolver.Resolver{
		Indexer: api.Indexer,
	}

	// start server

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	// router.Use(cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"}, // FIXME: this is lazy and potentially bad, allow setting from env/config
	// 	AllowCredentials: true,
	// 	Debug:            false,
	// }).Handler)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	// srv.AddTransport(&transport.Websocket{
	// 	KeepAlivePingInterval: 15 * time.Second,
	// 	Upgrader: websocket.Upgrader{
	// 		CheckOrigin: func(r *http.Request) bool {
	// 			return true
	// 		},
	// 		ReadBufferSize:  1024,
	// 		WriteBufferSize: 1024,
	// 	},
	// })
	// srv.SetQueryCache(lru.New(1000))
	// srv.Use(hooks.Prometheus{})
	// srv.Use(extension.Introspection{})
	// srv.Use(extension.AutomaticPersistedQuery{
	// 	Cache: lru.New(100),
	// })

	router.Handle("/sync", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sr indexer.SyncRequest
		if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
			writeError(w, err)
			return
		}
		if err := idxr.Sync(&sr); err != nil {
			writeError(w, err)
			return
		}
		response := map[string]bool{}
		response["sync"] = true
		writeJSON(w, response)
	}))

	router.Handle("/event", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var batch eventwatcher.LogBatch
		if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
			writeError(w, err)
			return
		}
		if err := idxr.Push(&batch); err != nil {
			writeError(w, err)
		}
		response := map[string]bool{}
		response["sync"] = true
		writeJSON(w, response)
	}))

	router.Handle("/query", srv)

	wasmhttp.Serve(router)
	select {}
}

func writeJSON(w http.ResponseWriter, obj interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		fmt.Printf("failed to write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = fmt.Sprintf("%v", err)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("writeError failed to write error: %v", err)
		return
	}
	w.Write(jsonResp)
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
