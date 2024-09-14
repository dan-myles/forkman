package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/common/logger"
	"github.com/avvo-na/forkman/internal/discord"
	"github.com/avvo-na/forkman/internal/server"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Main deps for application
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New()
	log := logger.New(cfg.GoEnv)

	// Create a new Discord bot
	discord := discord.New(cfg, log)
	discord.Setup()
	err := discord.Open()
	defer discord.Close()
	if err != nil {
		panic(err)
	}

	// Init new http server :D
	server := server.New(cfg, log, valid, discord)

	// Wait for sigterm (Ctrl+C)
	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msgf("Shutting down server %v", server.Addr)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			1000*time.Millisecond, // TODO: Make this configurable
		)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server shutdown failure!")
		}

		// TODO: once we have a db, we should also close it here

		close(closed)
	}()

	log.Info().Msgf("Starting server %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panic().Err(err).Msg("Server startup failure")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}
