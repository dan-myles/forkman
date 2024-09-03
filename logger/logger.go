package logger

import (
	"os"

	"github.com/avvo-na/devil-guard/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	c := config.GetConfig()
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()

	switch c.AppCfg.LogLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Dev Only Settings
	env := c.AppCfg.Environment
	if env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("Development environment detected, using console output")
	}

	// General Settings
	log.Logger = log.With().Caller().Logger()

	log.Info().Msg("Initialized logger")
}
