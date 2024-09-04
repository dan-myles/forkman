package main

import (
	"os"

	"github.com/avvo-na/devil-guard/common/logger"
	"github.com/avvo-na/devil-guard/common/validator"
	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/avvo-na/devil-guard/internal/sentinel"
	"github.com/eiannone/keyboard"
	"github.com/rs/zerolog/log"
)

// This function runs before the main entry point
// No error handling here as if we fail, we can't continue
// anyway, it is a fatal error.
func init() {
	validator.Init()
	config.Init()
	logger.Init()
}

func main() {
	s := sentinel.New()
	s.Start()

	// Wait here until q is pressed
	log.Info().Msg("Bot is now running!")
	log.Info().Msg("Press 'q' to exit")

	for {
		key, _, err := keyboard.GetSingleKey()
		defer keyboard.Close()

		if err != nil {
			log.Panic().Err(err).Msg("Failed to read key")
		}

		if key == rune('q') || key == rune(keyboard.KeyCtrlC) {
			log.Info().Msg("Exiting...")
			s.Stop()
			os.Exit(0)
		}
	}
}
