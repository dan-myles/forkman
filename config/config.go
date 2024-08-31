package config

import (
	"encoding/json"
	"os"

	"github.com/avvo-na/devil-guard/validator"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	DiscordAppID        string `json:"discord_app_id" validate:"required"`
	DiscordClientID     string `json:"discord_client_id" validate:"required"`
	DiscordClientSecret string `json:"discord_client_secret" validate:"required"`
	DiscordBotToken     string `json:"discord_bot_token" validate:"required"`
	DiscordDevGuildID   string `json:"discord_dev_guild_id" validate:"required"`
	DiscordOwnerID      string `json:"discord_owner_id" validate:"required"`
	LogLevel            string `json:"log_level" validate:"required"`
	Environment         string `json:"environment" validate:"required"`
}

// Should be "disabled" or "enabled"
type PluginConfig struct {
	Utility string `json:"utility" validate:"required,oneof=disabled enabled"`
}

var (
	PluginCfg        PluginConfig
	PluginDefaultCfg PluginConfig = PluginConfig{
		Utility: "disabled",
	}
	AppCfg        AppConfig
	AppDefaultCfg AppConfig = AppConfig{
		DiscordAppID:        "",
		DiscordClientID:     "",
		DiscordClientSecret: "",
		DiscordBotToken:     "",
		DiscordDevGuildID:   "",
		DiscordOwnerID:      "",
		LogLevel:            "info",
		Environment:         "dev",
	}
)

func InitConfig() {
	// Both of these functions will just panic if they fail
	// This is because the bot cannot run without these files!
	loadAppCfg()
	loadPluginCfg()
}

func loadAppCfg() {
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
		err = encoder.Encode(AppDefaultCfg)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to write default config")
		}

		// this will exit the program!
		log.Fatal().Msg("Config file generated, please setup config.json and restart the bot")
	}

	// Read config.json
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&AppCfg)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read config.json")
	}

	// Validate config.json
	err = validator.Validate.Struct(AppCfg)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to validate config.json")
	}

	// Success!
	log.Info().Interface("config", AppCfg).Msg("Config loaded")
}

func loadPluginCfg() {
	file, err := os.Open("plugins.json")
	if err != nil {
		log.Warn().Msg("Plugin config file not found, generating...")

		// Create config.json
		file, err := os.Create("plugins.json")
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create plugins.json")
		}

		// Write default config
		encoder := json.NewEncoder(file)
		err = encoder.Encode(PluginDefaultCfg)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to write default config")
		}

		// this will exit the program!
		log.Fatal().Msg("Config file generated, please setup plugins.json and restart the bot")
	}

	// Read config.json
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&PluginCfg)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read plugins.json")
	}

	// Validate config.json
	err = validator.Validate.Struct(PluginCfg)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to validate plugins.json")
	}

	// Success!
	log.Info().Interface("config", PluginCfg).Msg("Plugin config loaded")
}
