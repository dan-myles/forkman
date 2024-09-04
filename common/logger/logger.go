package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Define ANSI 256 color codes as uncapitalized constants
const (
	// Standard Colors (0-15)
	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	white        = 37
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	lightWhite   = 97

	// Pastel Colors (specific palette)
	pastelPink      = 213
	pastelPeach     = 216
	pastelPurple    = 141
	pastelLavender  = 165
	pastelMint      = 121
	pastelSkyBlue   = 117
	pastelCoral     = 210
	pastelOrange    = 208
	pastelTurquoise = 79
	pastelLilac     = 171
	pastelYellow    = 227
	pastelBlueGray  = 135
	pastelGreenGray = 149
	pastelBeige     = 230

	// Color Cube (16-231)
	colorCubeStart = 16
	colorCubeEnd   = 231

	// Grayscale (232-255)
	grayscaleStart = 232
	grayscaleEnd   = 255
)

var (
	levelColors = map[zerolog.Level]int{
		zerolog.TraceLevel: grayscaleStart + 15,
		zerolog.DebugLevel: grayscaleStart + 14,
		zerolog.InfoLevel:  pastelMint,
		zerolog.WarnLevel:  yellow,
		zerolog.ErrorLevel: red,
		zerolog.FatalLevel: red,
		zerolog.PanicLevel: red,
	}

	timeColor            = grayscaleStart + 5
	sourceColor          = grayscaleStart + 10
	sourceSeparator      = ">"
	sourceSeparatorColor = pastelSkyBlue
)

func Init() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	output.NoColor = false
	output.FormatLevel = func(i interface{}) string {
		switch i {
		case "trace":
			return colorize(fmt.Sprintf("| %-6s|", i), levelColors[zerolog.TraceLevel])
		case "debug":
			return colorize(fmt.Sprintf("| %-6s|", i), levelColors[zerolog.DebugLevel])
		case "info":
			return colorize(fmt.Sprintf("| %-6s|", i), levelColors[zerolog.InfoLevel])
		case "warn":
			return colorize(fmt.Sprintf("| %-6s|", i), levelColors[zerolog.WarnLevel])
		case "error":
			return colorize(fmt.Sprintf("| %-6s|", i), levelColors[zerolog.ErrorLevel])
		case "fatal":
			return colorize(fmt.Sprintf("| %-6s|", i), levelColors[zerolog.FatalLevel])
		case "panic":
			return colorize(fmt.Sprintf("| %-6s|", i), levelColors[zerolog.PanicLevel])
		default:
			return colorize(fmt.Sprintf("| %-6s|", i), 0)
		}
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
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
	output.FormatFieldName = func(i interface{}) string {
		// put on new line
		s := fmt.Sprintf("\n%+59s=", i)

		// colorize
		s = colorize(s, green)

		return s
	}
	output.FormatFieldValue = func(i interface{}) string {
		return colorize(fmt.Sprintf("%s", i), grayscaleStart+10)
	}
	output.FormatTimestamp = func(i interface{}) string {
		return colorize(fmt.Sprintf("%s", i), timeColor)
	}

	switch config.GetConfig().AppCfg.LogLevel {
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

	log.Logger = log.Output(output)
	log.Logger = log.With().Caller().Timestamp().Logger()
}

// colorize formats a string with the given color code.
func colorize(s string, color int) string {
	// 38;5 is the escape sequence for 256-color mode
	if color < 0 || color > 255 {
		return s // Return the original string if the color code is out of range
	}
	return fmt.Sprintf("\033[38;5;%dm%s\033[0m", color, s)
}
