package config

import (
	"encoding/json"
	"os"
	"sync"

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

type UtilityConfig struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

type VerificationConfig struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

type ModuleConfig struct {
	Utility      *UtilityConfig      `json:"utility" validate:"required"`
	Verification *VerificationConfig `json:"verification" validate:"required"`
}

type Config struct {
	AppCfg    *AppConfig    `json:"app" validate:"required"`
	ModuleCfg *ModuleConfig `json:"modules" validate:"required"`
	RWMutex   *sync.RWMutex `json:"-"`
}

var (
	instance   *Config = &Config{}
	defaultCfg *Config = &Config{
		AppCfg: &AppConfig{
			DiscordAppID:        "",
			DiscordClientID:     "",
			DiscordClientSecret: "",
			DiscordBotToken:     "",
			DiscordDevGuildID:   "",
			DiscordOwnerID:      "",
			LogLevel:            "info",
			Environment:         "dev",
		},
		ModuleCfg: &ModuleConfig{
			Utility: &UtilityConfig{
				Enabled: new(bool),
			},
			Verification: &VerificationConfig{
				Enabled: new(bool),
			},
		},
	}
)

func Init() {
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
		err = encoder.Encode(defaultCfg)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to write default config")
		}

		// this will exit the program!
		log.Fatal().Msg("Config file generated, please setup config.json and restart the bot")
	}

	// Read config.json
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(instance)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read config.json")
	}

	// Set the RWMutex
	instance.RWMutex = &sync.RWMutex{}

	// Validate config.json
	err = validator.Validate.Struct(instance)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to validate config.json")
	}

	// Success!
	log.Info().Interface("config", instance).Msg("Config loaded")
}

func GetConfig() *Config {
	return instance
}

func (c *Config) WriteConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		return err
	}

	return nil
}
