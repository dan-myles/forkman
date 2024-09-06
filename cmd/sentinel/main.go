package main

import (
	"os"
	"os/signal"

	"github.com/avvo-na/devil-guard/common/logger"
	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/discord"
	"github.com/avvo-na/devil-guard/discord/utility"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Load our configuration & logger etc.
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New(valid)
	log := logger.New(cfg)
	sesh := discord.New(cfg, log)

	// Setup utility module
	utils := utility.New(sesh, log, cfg)
	err := utils.Load()
	if err != nil {
		panic(err)
	}

	// Open a connection to Discord
	err = sesh.Open()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Bot is now running, press CTRL+C to exit")

	// Wait for a signal to stop the bot
	defer func() {
		log.Info().Msg("Stopping bot...")
		sesh.Close()
		log.Info().Msg("Bot has stopped!")
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
