package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	black = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
	darkGray = 90
)

func Init() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	output.NoColor = false
	output.FormatLevel = func(i interface{}) string {
		switch i {
		case "debug":
			return colorize(fmt.Sprintf("| %-6s|", i), darkGray)
		case "info":
			return colorize(fmt.Sprintf("| %-6s|", i), blue)
		case "warn":
			return colorize(fmt.Sprintf("| %-6s|", i), yellow)
		case "error":
			return colorize(fmt.Sprintf("| %-6s|", i), red)
		}

		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		if i == "debug" {
			return colorize(fmt.Sprintf("%s", i), darkGray)
		}

		return fmt.Sprintf("%s", i)
	}
	output.FormatCaller = func(i interface{}) string {
		if i == "debug" {
			return colorize(fmt.Sprintf("%s >", i), darkGray)
		}

		return colorize(fmt.Sprintf("%s", i), darkGray)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\n%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	output.FormatTimestamp = func(i interface{}) string {
		return colorize(fmt.Sprintf("%s", i), darkGray)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(output)
	log.Logger = log.With().Caller().Timestamp().Logger()
}

func colorize(s string, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
