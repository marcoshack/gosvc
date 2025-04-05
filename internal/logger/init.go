package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/marcoshack/gosvc/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(ctx context.Context, config config.ServiceConfig) (context.Context, zerolog.Logger, error) {
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

	var logWritter io.Writer
	logWritter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}

	logFileName := config.GetLogFileName()

	if logFileName != "" {
		//#nosec: G304 (CWE-22): Potential file inclusion via variable
		//#nosec: G302 (CWE-276): Expect file permissions to be 0600 or less
		logfile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err).Str("filename", logFileName).Msg("failed to open log file")
		}

		logWritter = io.MultiWriter(logWritter, logfile)
	}

	logger := zerolog.New(logWritter).Level(level).With().Timestamp().Logger()

	ctx = logger.WithContext(ctx)
	return ctx, logger, nil
}
