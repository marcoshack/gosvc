package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/marcoshack/gosvc"
	"github.com/marcoshack/gosvc/internal/config"
)

func main() {
	bs, err := gosvc.NewBootstrap[config.DefaultServiceConfig](context.Background(), gosvc.BootstrapInput{
		ServiceName: "gosvcsample",
		Args:        os.Args,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bootstrap")
	}

	ctx := bs.Ctx

	log.Ctx(ctx).Info().Msg("starting")
	log.Ctx(ctx).Info().Interface("config", bs.Config).Msg("configuration")
	log.Ctx(ctx).Info().Msg("stopped")

	os.Exit(0)
}
