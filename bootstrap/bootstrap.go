package bootstrap

import (
	"context"
	"flag"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/marcoshack/gosvc/internal/config"
	"github.com/marcoshack/gosvc/internal/logger"
)

type Bootstrap[ConfigType config.ServiceConfig] struct {
	Name           string
	Ctx            context.Context
	Config         ConfigType
	ConfigFileName string
	Logger         zerolog.Logger
	AWSConfig      awsconfig.Config
}

type Input struct {
	ServiceName string
	AWSRegion   string
	Args        []string
}

func (i *Input) validate() error {
	if i.ServiceName == "" {
		return errors.New("service name is required")
	}
	if i.AWSRegion == "" {
		return errors.New("AWS region is required")
	}
	return nil
}

func New[ConfigType config.ServiceConfig](ctx context.Context, input Input) (*Bootstrap[ConfigType], error) {
	if err := input.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid input")
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(input.AWSRegion))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create AWS config")
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
		Name:      input.ServiceName,
		Ctx:       ctx,
		Config:    config,
		Logger:    logger.With().Str("service", input.ServiceName).Logger(),
		AWSConfig: awsConfig,
	}, nil
}
