package config

import (
	"encoding/json"
	"os"

	"github.com/avvo-na/devil-guard/validator"
	"github.com/rs/zerolog/log"
)

type Config struct {
	OwnerID             string `json:"owner_id" validate:"required"`
	DiscordClientID     string `json:"discord_client_id" validate:"required"`
	DiscordClientSecret string `json:"discord_client_secret" validate:"required"`
	DiscordBotToken     string `json:"discord_bot_token" validate:"required"`
	LogLevel            string `json:"log_level" validate:"required"`
	Environment         string `json:"environment" validate:"required"`
}

var (
	ConfigData    Config
	DefaultConfig Config = Config{
		OwnerID:             "",
		DiscordClientID:     "",
		DiscordClientSecret: "",
		DiscordBotToken:     "",
		LogLevel:            "info",
		Environment:         "dev",
	}
)

func InitConfig() {
	// Check current directory for config.json
	// Generate one if not found
	file, err := os.Open("config.json")
	if err != nil {
		log.Warn().Msg("Config file not found, generating...")

		// Create config.json
		file, err := os.Create("config.json")
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create config.json")
		}

		// Write default config
		encoder := json.NewEncoder(file)
		err = encoder.Encode(DefaultConfig)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to write default config")
		}

		// this will exit the program!
		log.Fatal().Msg("Config file generated, please setup config.json and restart the bot")
	}

	// Read config.json
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ConfigData)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read config.json")
	}

	// Validate config.json
	err = validator.Validate.Struct(ConfigData)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to validate config.json")
	}

	// Success!
	log.Info().Interface("config", ConfigData).Msg("Config loaded")
}
