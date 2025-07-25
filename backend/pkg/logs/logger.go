package logs

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/errors"
)

func NewLogger(conf *config.ConfigLogs) (*zerolog.Logger, error) {
	file := os.Stderr
	if conf.File != "" {
		if f, err := os.OpenFile(conf.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			file = f
		} else {
			return nil, &errors.LogsFileError{
				FilePath: conf.File,
				Err:      err,
			}
		}
	}

	level := zerolog.InfoLevel
	switch strings.ToLower(conf.Level) {
	case "trace":
		level = zerolog.TraceLevel
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	default:
		return nil, &errors.LogsLevelError{
			Level: conf.Level,
		}
	}
	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(file).With().Timestamp().Logger()
	if conf.Pretty {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
		})
	}

	return &logger, nil
}
