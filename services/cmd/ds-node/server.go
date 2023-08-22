package main

import (
	"context"
	"flag"
	"os"

	"github.com/playmint/ds-node/pkg/api"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/config"
	"github.com/playmint/ds-node/pkg/indexer"
	"github.com/playmint/ds-node/pkg/mgmt"
	"github.com/playmint/ds-node/pkg/sequencer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Main(ctx context.Context) error {

	// init logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// start the management server
	go func() {
		if err := mgmt.ListenAndServe(":9090"); err != nil {
			log.Fatal().Err(err).Str("service", "mgmt").Msg("exited")
		}
	}()

	// configure subscriptions consumer
	subscriptions, notifications := model.NewSubscriptions()
	go subscriptions.Listen(ctx)

	// start an indexer
	idxr, err := indexer.NewMemoryIndexer(ctx, notifications, config.IndexerProviderHTTP, config.IndexerProviderWS)
	if err != nil {
		return err
	}

	// start a sequencer
	seqr, err := sequencer.NewMemorySequencer(
		ctx,
		config.SequencerPrivateKey,
		notifications,
		config.SequencerProviderHTTP,
		idxr,
	)
	if err != nil {
		return err
	}

	// wait for ready
	<-idxr.Ready()
	log.Info().Str("service", "indexer").Msg("ready")
	<-seqr.Ready()
	log.Info().Str("service", "sequencer").Msg("ready")

	// start graphql api server
	api := api.Server{
		Indexer:   idxr,
		Sequencer: seqr,
	}
	if err := api.Start(ctx, subscriptions); err != nil {
		log.Fatal().Err(err).Str("service", "api").Msg("exited")
	}

	return nil
}

func main() {
	if err := Main(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
