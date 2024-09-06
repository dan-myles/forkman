package logger

import "fmt"

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
	traceColor = grayscaleStart + 15
	debugColor = grayscaleStart + 14
	infoColor  = pastelMint
	warnColor  = yellow
	errorColor = red
	fatalColor = red
	panicColor = red

	timeColor            = grayscaleStart + 5
	sourceColor          = grayscaleStart + 10
	sourceSeparator      = ">"
	sourceSeparatorColor = pastelSkyBlue
)

func colorize(s string, color int) string {
	if color < 0 || color > 255 {
		return s
	}
	return fmt.Sprintf("\033[38;5;%dm%s\033[0m", color, s)
}
