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
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/go-playground/validator/v10"
)

func main() {
	// General deps
	// Adding a new change
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New()
	log := logger.New(cfg.GoEnv, cfg.LogLevel)
	db := database.New(log)

	// AWS
	acfg, err := awscfg.LoadDefaultConfig(context.TODO(),
		awscfg.WithRegion(cfg.AWS_REGION),
		awscfg.WithCredentialsProvider(aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				cfg.AWS_ACCESS_KEY_ID,
				cfg.AWS_SECRET_ACCESS_KEY,
				"",
			),
		)),
	)
	if err != nil {
		panic(err)
	}

	// Discord & http server
	discord := discord.New(cfg, log, db, acfg)
	server := server.New(cfg, log, valid, discord, db)

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

		close(shutdown)
	}()

	// Listen & Serve
	log.Info().Msgf("Server starting on :%d", cfg.ServerPort)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	// Wait for shutdown
	<-shutdown
	log.Info().Msg("Application shutdown completed!")
}
