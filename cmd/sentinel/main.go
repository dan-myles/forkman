package main

import (
	"os"
	"os/signal"

	"github.com/avvo-na/devil-guard/common/logger"
	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/discord"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Load our configuration & logger etc.
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New(valid)
	log := logger.New(cfg)

	// Create a new Discord bot
	disco := discord.New(cfg, log)
	disco.Setup()
	disco.Open()

	// Wait for a signal to stop the bot
	log.Info().Msg("Bot started 🔥, press CTRL+C to shutdown.")
	defer func() {
		log.Info().Msg("Stopping bot...")
		disco.Close()
		log.Info().Msg("Bot stopped :D")
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
