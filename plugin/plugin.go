package plugin

import (
	"encoding/json"
	"os"

	"github.com/avvo-na/devil-guard/validator"
	"github.com/rs/zerolog/log"
)

type Plugin interface {
	// Name returns the name of the plugin
	Name() string
	// Register is used to register the plugin
	Register()
	// Run is the main function of the plugin
	Run() error
}

type PluginManager struct {
	Plugins []Plugin
}

// Should be "disabled" or "enabled"
type PluginConfig struct {
	Utility string `json:"utility" validate:"required,oneof=disabled enabled"`
}

var (
	ConfigData    PluginConfig
	DefaultConfig PluginConfig = PluginConfig{
		Utility: "",
	}
)

func New() *PluginManager {
	return &PluginManager{}
}

func (pm *PluginManager) Register(p Plugin) {
	pm.Plugins = append(pm.Plugins, p)
}

func (pm *PluginManager) RegisterAll() {
	for _, p := range pm.Plugins {
		p.Register()
	}
}

func (pm *PluginManager) Run() {
	for _, p := range pm.Plugins {
		p.Run()
	}
}

func InitConfig() {
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
		err = encoder.Encode(DefaultConfig)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to write default config")
		}

		// this will exit the program!
		log.Fatal().Msg("Config file generated, please setup plugins.json and restart the bot")
	}

	// Read config.json
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ConfigData)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read plugins.json")
	}

	// Validate config.json
	err = validator.Validate.Struct(ConfigData)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to validate plugins.json")
	}

	// Success!
	log.Info().Interface("config", ConfigData).Msg("Plugin config loaded")
}
