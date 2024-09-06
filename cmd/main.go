package main

import (
	"os"
	"os/signal"

	"github.com/avvo-na/devil-guard/common/log"
	"github.com/avvo-na/devil-guard/common/validator"
	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/avvo-na/devil-guard/internal/sentinel"
)

// This function runs before the main entry point
// No error handling here as if we fail, we can't continue
// anyway, it is a fatal error.
func init() {
	validator.Init()
	config.Init()
	log.Init()
}

func main() {
	// Start the bot
	s := sentinel.New()
	s.Start()

	// Wait here until q is pressed
	log.Info().Msg("Bot is now running!")
	log.Info().Msg("Press 'CTRL-C' to exit")

	// Wait for ctrl-c to exit
	defer func() {
		log.Info().Msg("Stopping bot...")
		s.Stop()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
