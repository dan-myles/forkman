package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avvo-na/forkman/api/router"
	"github.com/avvo-na/forkman/common/logger"
	"github.com/avvo-na/forkman/config"
	"github.com/avvo-na/forkman/discord"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Main deps for application
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New()
	log := logger.New(cfg.GoEnv)

	// Create a new Discord bot
	disco := discord.New(cfg, log)
	disco.Setup()
	err := disco.Open()
	defer disco.Close()
	if err != nil {
		panic(err)
	}

	// Init new http server :D
	r := router.New(log, valid, disco)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: r,
	}

	// Wait for sigterm (Ctrl+C)
	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msgf("Shutting down server %v", s.Addr)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			1000*time.Millisecond, // TODO: Make this configurable
		)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server shutdown failure!")
		}

		// TODO: once we have a db, we should also close it here

		close(closed)
	}()

	log.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	log.Info().Msgf("Server shutdown successfully")
}
