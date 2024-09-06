package main

import (
	"os"
	"os/signal"

	"github.com/avvo-na/devil-guard/common/logger"
	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/avvo-na/devil-guard/internal/sentinel"
	"github.com/avvo-na/devil-guard/internal/sentinel/utility"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Load our configuration & logger etc.
	v := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New(v)
	logger := logger.New(cfg)
	session := sentinel.NewSession(cfg, logger)

	// Setup utility module
	utils := utility.New(session, logger, cfg)
	err := utils.Load()
	if err != nil {
		panic(err)
	}

	// Open a connection to Discord
	err = session.Open()
	if err != nil {
		panic(err)
	}
	logger.Info().Msg("Bot is now running, press CTRL+C to exit")

	// Wait for a signal to stop the bot
	defer func() {
		logger.Info().Msg("Stopping bot...")
		session.Close()
		logger.Info().Msg("Bot has stopped!")
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
