package discord

import (
	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// TODO: Probably dont panic so much

type Module interface {
	// Returns the name of the module
	Name() string

	// Enables the module, handles any setup and registration
	// of commands, writes config to file.
	Enable() error

	// Disables the module, handles any cleanup and deregistration
	// of commands, writes config to file.
	Disable() error

	// Loads the module, handles any setup and registration of
	// commands, *reads* config from file. To only be called once
	Load() error
}

func New(cfg *config.ConfigManager, log *zerolog.Logger) *discordgo.Session {
	s, err := discordgo.New("Bot " + cfg.GetAppConfig().DiscordBotToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create a new Discord session")
	}

	// Settings
	s.Identify.Intents = discordgo.IntentsAll

	log.Info().Msg("Created a new Discord session")
	return s
}
