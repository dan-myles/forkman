package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/avvo-na/forkman/config"
	"github.com/rs/zerolog"
)

func New(c *config.ConfigManager) *zerolog.Logger {
	cfg := c.GetAppConfig()

	switch cfg.Environment {
	case "dev":
		return dev(cfg)
	case "prod":
		return prod(cfg)
	default:
		panic("Unknown environment, please check your configuration file")
	}
}

func dev(c config.AppConfig) *zerolog.Logger {
	// Default to console output
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	output.NoColor = false

	// Format the "Level" field
	output.FormatLevel = func(i interface{}) string {
		switch i {
		case "trace":
			return colorize(fmt.Sprintf("| %-6s|", i), traceColor)
		case "debug":
			return colorize(fmt.Sprintf("| %-6s|", i), debugColor)
		case "info":
			return colorize(fmt.Sprintf("| %-6s|", i), infoColor)
		case "warn":
			return colorize(fmt.Sprintf("| %-6s|", i), warnColor)
		case "error":
			return colorize(fmt.Sprintf("| %-6s|", i), errorColor)
		case "fatal":
			return colorize(fmt.Sprintf("| %-6s|", i), fatalColor)
		case "panic":
			return colorize(fmt.Sprintf("| %-6s|", i), panicColor)
		default:
			return colorize(fmt.Sprintf("| %-6s|", i), 0)
		}
	}

	// Format the "Message" field
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	// Format the "Caller" field (file and line number)
	output.FormatCaller = func(i interface{}) string {
		// chop off the path
		s := fmt.Sprintf("%s", i)
		s = strings.Split(s, "/")[len(strings.Split(s, "/"))-1]

		// Evenly space
		// but not too much
		if len(s) < 15 {
			s = fmt.Sprintf("%-15s", s)
		}

		s = colorize(s, sourceColor)
		separator := colorize(sourceSeparator, sourceSeparatorColor)

		return fmt.Sprintf("%s%s", s, separator)
	}

	// Format the "Timestamp" field
	output.FormatFieldName = func(i interface{}) string {
		// put on new line
		s := fmt.Sprintf("\n%+59s=", i)

		// colorize
		s = colorize(s, green)

		return s
	}

	// Format the "FieldKey" field
	output.FormatFieldValue = func(i interface{}) string {
		return colorize(fmt.Sprintf("%s", i), grayscaleStart+10)
	}

	// Format the "Timestamp" field
	output.FormatTimestamp = func(i interface{}) string {
		return colorize(fmt.Sprintf("%s", i), timeColor)
	}

	// Initialize the logger
	logger := zerolog.New(output).With().Timestamp().Caller().Logger()

	// Add the level
	switch c.LogLevel {
	case "trace":
		logger = logger.Level(zerolog.TraceLevel)
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	case "fatal":
		logger = logger.Level(zerolog.FatalLevel)
	case "panic":
		logger = logger.Level(zerolog.PanicLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}

	return &logger
}

// TODO: impl production logger, will just be json output
func prod(c config.AppConfig) *zerolog.Logger {
	return nil
}
