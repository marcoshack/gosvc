package gosvc

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/marcoshack/gosvc/internal/config"
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

	config, err := loadConfig[ConfigType](ctx, awsConfig, configFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}

	ctx, logger, err := initLogger(ctx, config)
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

func loadConfig[ConfigType config.ServiceConfig](ctx context.Context, awsConfig aws.Config, configFilename *string) (ConfigType, error) {
	var c ConfigType
	var err error

	// try to load configuration from file
	if *configFilename != "" {
		c, err = config.LoadFromFile[ConfigType](&config.LoadFromFileInput{Filename: *configFilename})
		if err != nil {
			return c, errors.Wrap(err, "failed to load configuration from file")
		}
	}

	// try to load configuration from AppConfig
	stageName := os.Getenv("STAGE")
	appName := os.Getenv("APP_NAME")
	configProfile := os.Getenv("CONFIG_PROFILE")
	if !c.IsValid() && stageName != "" && appName != "" && configProfile != "" {
		c, err = config.LoadFromAppConfig[ConfigType](ctx, &config.LoadFromAppConfigInput{
			AWSConfig:                awsConfig,
			ApplicationName:          appName,
			ConfigurationProfileName: configProfile,
			EnvironmentName:          stageName,
		})
		if err != nil {
			return c, errors.Wrap(err, "failed to load configuration from AppConfig")
		}
	}

	if !c.IsValid() {
		return c, fmt.Errorf("no configuration file provided, nor STAGE, APP_NAME and CONFIG_PROFILE environment variables set to load from AppConfig")
	}

	return c, nil
}

func initLogger(ctx context.Context, config config.ServiceConfig) (context.Context, zerolog.Logger, error) {
	level := zerolog.InfoLevel

	if config != nil && config.GetLogLevel() != "" {
		var err error
		level, err = zerolog.ParseLevel(config.GetLogLevel())
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse log level, using default INFO")
			level = zerolog.InfoLevel
		}
	}
	zerolog.SetGlobalLevel(level)

	var logWriter io.Writer
	logWriter = os.Stdout

	logFileName := config.GetLogFileName()

	if logFileName != "" {
		//#nosec: G304 (CWE-22): Potential file inclusion via variable
		//#nosec: G302 (CWE-276): Expect file permissions to be 0600 or less
		logfile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err).Str("filename", logFileName).Msg("failed to open log file")
		}

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}
		logWriter = io.MultiWriter(consoleWriter, logfile)
	}

	logger := zerolog.New(logWriter).Level(level).With().Timestamp().Logger()

	ctx = logger.WithContext(ctx)
	return ctx, logger, nil
}
