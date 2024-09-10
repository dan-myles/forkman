package main

import (
	"os"
	"os/signal"

	"github.com/avvo-na/forkman/common/logger"
	"github.com/avvo-na/forkman/config"
	"github.com/avvo-na/forkman/discord"
	"github.com/go-playground/validator/v10"
)

// NOTE: Eventually when (if and) when we have more services
// would probably want to add context to all logs that these
// are from the "sentinel" service. But... we only have one 😁
func main() {
	// Load our configuration & logger etc.
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New(valid)
	log := logger.New(cfg)

	// Create a new Discord bot
	disco := discord.New(cfg, log)
	disco.Setup()
	err := disco.Open()
	if err != nil {
		panic(err)
	}

	// Wait for a signal to stop the app
	log.Info().Msg("Sentinel started 🔥, press CTRL+C to shutdown.")
	defer func() {
		log.Info().Msg("Stopping...")
		disco.Close()
		log.Info().Msg("Sentinel stopped :D")
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
