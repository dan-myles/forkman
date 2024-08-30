package utils

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/rs/zerolog/log"
)

type Config struct {
	OwnerID             string `json:"owner_id"`
	DiscordClientID     string `json:"discord_client_id"`
	DiscordClientSecret string `json:"discord_client_secret"`
	DiscordBotToken     string `json:"discord_bot_token"`
	LogLevel            string `json:"log_level"`
	Environment         string `json:"environment"`
}

var (
	ConfigData    Config
	DefaultConfig Config = Config{
		OwnerID:             "0",
		DiscordClientID:     "0",
		DiscordClientSecret: "0",
		DiscordBotToken:     "0",
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
	err = decoder.Decode(&ConfigData)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read config.json")
	}

	// Check for missing fields
	v := reflect.ValueOf(ConfigData)
	var fields []string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == reflect.Zero(v.Field(i).Type()).Interface() {
			fields = append(fields, v.Type().Field(i).Name)
		}
	}

	if len(fields) > 0 {
		log.Panic().Interface("fields", fields).Msg("Missing fields in config.json")
	}

	// Success!
	log.Info().Interface("config", ConfigData).Msg("Config loaded")
}
