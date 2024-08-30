package main

import (
	"os"

	"github.com/avvo-na/devil-guard/sentinel"
	"github.com/avvo-na/devil-guard/utils"
	"github.com/eiannone/keyboard"
	"github.com/rs/zerolog/log"
)

func main() {
	// Init config and logger
	utils.InitConfig()
	utils.InitLogger()

	// start the bot on a goroutine
	// so we can run other things in the main thread
	// without blocking
	var s *sentinel.Sentinel
	s = sentinel.New()
	s.Start()

	// Wait here until q is pressed
	log.Info().Msg("Bot is now running. Press q to exit.")
	key, _, err := keyboard.GetSingleKey()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read key")
	}

	if key == rune('q') {
		log.Info().Msg("Exiting...")
		s.Stop()
		os.Exit(1)
	}
}
