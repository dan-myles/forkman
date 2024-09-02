package main

import (
	"os"

	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/logger"
	"github.com/avvo-na/devil-guard/sentinel"
	"github.com/avvo-na/devil-guard/validator"
	"github.com/eiannone/keyboard"
	"github.com/rs/zerolog/log"
)

func init() {
	validator.InitValidator()
	config.InitConfig()
	logger.InitLogger()
	sentinel.InitSentinel()
}

func main() {
	sentinel.Start()

	// Wait here until q is pressed
	log.Info().Msg("Bot is now running!")
	log.Info().Msg("Press 'q' to exit")

	for {
		key, _, err := keyboard.GetSingleKey()
		defer keyboard.Close()

		if err != nil {
			log.Panic().Err(err).Msg("Failed to read key")
		}

		if key == rune('q') {
			log.Info().Msg("Exiting...")
			sentinel.Stop()
			os.Exit(0)
		}
	}
}
