package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/avvo-na/devil-guard/validator"
	"github.com/rs/zerolog/log"
)

// This module is responsible for loading configuration files
// However, module config can be changed at runtime, and is NOT
// done here. This is just for initial loading.

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
	Utility      string `json:"utility" validate:"required,oneof=disabled enabled"`
	Verification string `json:"verification" validate:"required,oneof=disabled enabled"`
}

var (
	ModuleCfg        *ModuleConfig
	ModuleDefaultCfg *ModuleConfig = &ModuleConfig{
		Utility: "disabled",
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
	mutex = &sync.Mutex{}
)

func InitConfig() {
	// Both of these functions will just panic if they fail
	// This is because the bot cannot run without these files!
	loadAppCfg()
	loadPluginCfg()
}

// TODO: check if the module is already enabled/disabled
// also check if its an actual module
func WriteDisableModule(module string) error {
	return writeModule(module, "disabled")
}

func WriteEnableModule(module string) error {
	return writeModule(module, "enabled")
}

func writeModule(module string, value string) error {
	// Lock the mutex
	mutex.Lock()
	defer mutex.Unlock()

	// Find the module that matches the name
	r := reflect.ValueOf(ModuleCfg)
	f := reflect.Indirect(r).FieldByNameFunc(func(f string) bool {
		if strings.EqualFold(f, module) {
			return true
		}

		return false
	})
	if !f.IsValid() {
		return fmt.Errorf("Module not found: %s", module)
	}

	// WARN: This can panic and will need more error handling in the future
	f.Set(reflect.ValueOf(value))

	// Write the new config
	file, err := os.Create("modules.json")
	if err != nil {
		fmt.Errorf("Failed to open modules.json: %v", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(ModuleCfg)
	if err != nil {
		fmt.Errorf("Failed to write modules.json: %v", err)
	}

	log.Debug().Interface("config", ModuleCfg).Msg("Module config updated")
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
