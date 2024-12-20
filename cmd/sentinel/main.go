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
	valid := validator.New(validator.WithRequiredStructEnabled())
	cfg := config.New()
	log := logger.New(cfg.GoEnv)
	db := database.New(log)

	// AWS Config
	var awsConfig aws.Config
	var err error
	if cfg.AWSEnabled {
		awsConfig, err = awscfg.LoadDefaultConfig(context.TODO(),
			awscfg.WithRegion(cfg.AWSConfig.AWSRegion),
			awscfg.WithCredentialsProvider(aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					cfg.AWSConfig.AWSAccessKeyID,
					cfg.AWSConfig.AWSSecretAccessKey,
					"",
				),
			)),
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load AWS configuration")
			os.Exit(1)
		}
	}

	// Discord & HTTP Server
	discordClient := discord.New(cfg, log, db, awsConfig)
	httpServer := server.New(cfg, log, valid, discordClient, db)
	log.Info().Msg("Initialization complete, starting server!")

	// Cleanup on Interrupt/SIGTERM
	shutdown := make(chan struct{})
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
		<-ch

		log.Info().Msg("Attempting graceful shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown HTTP server
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to shut down HTTP server")
		} else {
			log.Info().Msg("HTTP server shut down successfully")
		}

		// Close database connection
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Error().Err(err).Msg("Failed to close database connection")
			} else {
				log.Info().Msg("Database connection closed successfully")
			}
		} else {
			log.Error().Err(err).Msg("Failed to retrieve database connection")
		}

		// Close Discord client
		if err := discordClient.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Discord session")
		} else {
			log.Info().Msg("Discord session closed successfully")
		}

		close(shutdown)
	}()

	// Listen & Serve
	log.Info().Msgf("Server starting on :%d", cfg.ServerConfig.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Failed to start HTTP server")
		os.Exit(1)
	}

	// Wait for shutdown
	<-shutdown
	log.Info().Msg("Application shutdown completed!")
}
