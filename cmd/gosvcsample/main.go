package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/marcoshack/gosvc/bootstrap"
	"github.com/marcoshack/gosvc/internal/config"
)

func main() {
	defaultConfig := config.SampleServiceConfig{
		Host: "localhost",
		Port: 3000,
	}

	bs, err := bootstrap.New[config.SampleServiceConfig](context.Background(), bootstrap.Input{
		ServiceName:   "gosvcsample",
		AWSRegion:     "us-east-1",
		Args:          os.Args,
		DefaultConfig: defaultConfig,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bootstrap")
	}

	ctx, cancel := context.WithCancel(bs.Ctx)
	defer cancel()

	log.Ctx(ctx).Info().Msg("service started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Ctx(ctx).Info().Msg("stopping")
	cancel()

	// simulate some work
	time.Sleep(2 * time.Second)

	log.Ctx(ctx).Info().Msg("stopped")
}
