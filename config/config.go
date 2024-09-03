package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/avvo-na/devil-guard/validator"
	"github.com/rs/zerolog/log"
)

// TODO: WriteDisable/Enable needs to be updated to handle
// more complex config structs. This will be done by just
// offering a public "WriteConfig" function that assumes
// you've just unlocked the config mutex and have made changes
// to it. This is mainly for module config.

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
type ModuleConfig struct {
	Utility struct {
		Enabled bool `json:"enabled" validate:"required"`
	} `json:"utility"`
	Verification struct {
		Enabled bool `json:"enabled" validate:"required"`
	} `json:"verification"`
}

var (
	ModuleCfg        *ModuleConfig
	ModuleDefaultCfg *ModuleConfig = &ModuleConfig{
		Utility: struct {
			Enabled bool `json:"enabled" validate:"required"`
		}{
			Enabled: false,
		},
		Verification: struct {
			Enabled bool `json:"enabled" validate:"required"`
		}{
			Enabled: false,
		},
	}
	AppCfg        *AppConfig
	AppDefaultCfg *AppConfig = &AppConfig{
		DiscordAppID:        "",
		DiscordClientID:     "",
		DiscordClientSecret: "",
		DiscordBotToken:     "",
		DiscordDevGuildID:   "",
		DiscordOwnerID:      "",
		LogLevel:            "info",
		Environment:         "dev",
	}
	Mutex = &sync.Mutex{}
)

func InitConfig() {
	// Both of these functions will just panic if they fail
	// This is because the bot cannot run without these files!
	loadAppCfg()
	loadPluginCfg()
}

// WARN: This function assumes you have a lock on the config mutex
func WriteModuleConfig() error {
	// Write app config
	file, err := os.Create("modules.json")
	if err != nil {
		return fmt.Errorf("Failed to open modules.json: %w", err)
	}

	// Write our new config to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(ModuleCfg)
	if err != nil {
		return fmt.Errorf("Failed to write modules.json: %w", err)
	}

	return nil
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
		encoder.SetIndent("", "  ")
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
	file, err := os.Open("modules.json")
	if err != nil {
		log.Warn().Msg("Plugin config file not found, generating...")

		// Create config.json
		file, err := os.Create("modules.json")
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create modules.json")
		}

		// Write default config
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(ModuleDefaultCfg)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to write default config")
		}

		// this will exit the program!
		log.Fatal().Msg("Config file generated, please setup modules.json and restart the bot")
	}

	// Read config.json
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ModuleCfg)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read modules.json")
	}

	// Validate config.json
	err = validator.Validate.Struct(ModuleCfg)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to validate modules.json")
	}

	// Success!
	log.Info().Interface("config", ModuleCfg).Msg("Plugin config loaded")
}
