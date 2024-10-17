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
	"github.com/avvo-na/forkman/internal/database"
	"github.com/avvo-na/forkman/internal/discord"
	"github.com/avvo-na/forkman/internal/server"
	"github.com/go-playground/validator/v10"
	"github.com/resend/resend-go/v2"
)

func main() {
	// Startup w/ API & Discord
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New()
	log := logger.New(cfg.GoEnv)
	db := database.New(log)
	email := resend.NewClient(cfg.ResendAPIKey)
	discord := discord.New(cfg, log, db, email)
	server := server.New(cfg, log, valid, discord, db)
	log.Info().Msg("Initialization complete, starting server!")

	// Cleanup on Interrupt/SIGTERM
	// We need to catch both incase we're running on Windows
	shutdown := make(chan struct{})
	go func() {
		ch := make(chan os.Signal, 2)
		signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
		<-ch

		log.Info().Msg("Attempting graceful shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Error().Err(err).Msg("HTTP Server shutdown failure")
		}
		log.Info().Msg("HTTP Server shutdown successfully")

		sqlDB, err := db.DB()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get database connection")
		}

		err = sqlDB.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
		log.Info().Msg("Database connection closed successfully")

		err = discord.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close discord session")
		}
		log.Info().Msg("Discord session closed successfully")

		log.Info().Msg("Application shutdown completed!")
		close(shutdown)
	}()

	// Listen & Serve
	log.Info().Msgf("Server starting on :%d", cfg.ServerPort)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	// Wait for shutdown
	<-shutdown
}
