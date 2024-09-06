package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
)

type AppConfig struct {
	DiscordAppID        string `json:"discord_app_id"        validate:"required"`
	DiscordClientID     string `json:"discord_client_id"     validate:"required"`
	DiscordClientSecret string `json:"discord_client_secret" validate:"required"`
	DiscordBotToken     string `json:"discord_bot_token"     validate:"required"`
	DiscordDevGuildID   string `json:"discord_dev_guild_id"  validate:"required"`
	DiscordOwnerID      string `json:"discord_owner_id"      validate:"required"`
	LogLevel            string `json:"log_level"             validate:"required"`
	Environment         string `json:"environment"           validate:"required"`
}

type UtilityConfig struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

type ModerationConfig struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

type ModuleConfig struct {
	Utility    UtilityConfig    `json:"utility"    validate:"required"`
	Moderation ModerationConfig `json:"moderation" validate:"required"`
}

type Config struct {
	AppCfg    AppConfig    `json:"app"     validate:"required"`
	ModuleCfg ModuleConfig `json:"modules" validate:"required"`
}

type ConfigManager struct {
	cfg *Config
	mtx *sync.RWMutex
	val *validator.Validate
}

// NOTE: Used to generate a default config file,
// only happens if the config file is not found.
var (
	defaultCfg = Config{
		AppCfg: AppConfig{
			DiscordAppID:        "",
			DiscordClientID:     "",
			DiscordClientSecret: "",
			DiscordBotToken:     "",
			DiscordDevGuildID:   "",
			DiscordOwnerID:      "",
			LogLevel:            "info",
			Environment:         "dev",
		},
		ModuleCfg: ModuleConfig{
			Utility: UtilityConfig{
				Enabled: new(bool),
			},
			Moderation: ModerationConfig{
				Enabled: new(bool),
			},
		},
	}
)

func New(v *validator.Validate) *ConfigManager {
	file, err := os.Open("config.json")
	if err != nil {
		// Generate a new config file
		file, err := os.Create("config.json")
		if err != nil {
			panic(err)
		}
		// Write the JSON to the file
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		enc.Encode(defaultCfg)

		// Close the file
		file.Close()
		panic(
			"Config file not found, one has been generated, please fill it out and restart the bot",
		)
	}

	// Load the config file
	cfg := Config{}
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	// Validate the config
	err = v.Struct(&cfg)
	if err != nil {
		panic(err)
	}

	return &ConfigManager{
		cfg: &cfg,
		mtx: &sync.RWMutex{},
		val: v,
	}
}

func (c *ConfigManager) GetAppConfig() AppConfig {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.cfg.AppCfg
}

func (c *ConfigManager) GetModuleConfig() ModuleConfig {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.cfg.ModuleCfg
}

func (c *ConfigManager) WriteModerationConfig(cfg ModerationConfig) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cfg.ModuleCfg.Moderation = cfg

	file, err := os.Create("config.json")
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(c.cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func (c *ConfigManager) WriteUtilityConfig(cfg UtilityConfig) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cfg.ModuleCfg.Utility = cfg

	file, err := os.Create("config.json")
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(c.cfg)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
