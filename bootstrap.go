package gosvc

import (
	"context"
	"flag"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/marcoshack/gosvc/internal/config"
	"github.com/marcoshack/gosvc/internal/logger"
)

type Bootstrap[ConfigType config.ServiceConfig] struct {
	Name           string
	Ctx            context.Context
	Config         ConfigType
	ConfigFileName string
	Logger         zerolog.Logger
}

type BootstrapInput struct {
	ServiceName string
	Args        []string
}

func NewBootstrap[ConfigType config.ServiceConfig](ctx context.Context, input BootstrapInput) (*Bootstrap[ConfigType], error) {
	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed to create AWS config")
	}

	// TODO : add support to pass additional CLI options from input
	fs := flag.NewFlagSet(input.ServiceName, flag.ExitOnError)
	configFilename := fs.String("c", "", "configuration filepath")
	err = fs.Parse(input.Args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse flags")
	}

	config, err := config.LoadConfig[ConfigType](ctx, awsConfig, configFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}

	ctx, logger, err := logger.InitLogger(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize logger")
	}

	return &Bootstrap[ConfigType]{
		Name:   input.ServiceName,
		Ctx:    ctx,
		Config: config,
		Logger: logger.With().Str("service", input.ServiceName).Logger(),
	}, nil
}
