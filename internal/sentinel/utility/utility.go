package utility

import (
	"fmt"

	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping the bot",
	},
	{
		Name:        "role",
		Description: "Manage roles",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "all",
				Description: "Give to all users",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": ping,
	"role": role,
}

type UtilityModule struct{}

func New() *UtilityModule {
	return &UtilityModule{}
}

func (u *UtilityModule) Name() string {
	return "utility"
}

// NOTE: For enabling and disabling modules, we do not need to worry about
// the mutex that holds the configuration. This is because we only ever call
// these methods from the moduleManager, which will already have a lock.
func (u *UtilityModule) Enable(s *discordgo.Session) error {
	// Grab the config
	config := config.GetConfig()

	// Grab the app ID and guild ID
	appID := config.AppCfg.DiscordAppID
	guildID := config.AppCfg.DiscordDevGuildID

	// Write new config
	enable := true
	config.ModuleCfg.Utility.Enabled = &enable

	err := config.WriteConfig()
	if err != nil {
		return fmt.Errorf("Failed to write config: %w", err)
	}
	log.Debug().Msg("Updated utility module config")

	// Register all commands
	log.Debug().Msg("Registering utility commands...")
	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(appID, guildID, command)
		if err != nil {
			return err
		}
	}

	log.Debug().Msg("Registering utility command handlers...")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, ok := commandHandlers[i.ApplicationCommandData().Name]
		if !ok {
			return
		}

		handler(s, i)
	})

	log.Info().Msg("Utility module enabled!")
	return nil
}

func (u *UtilityModule) Disable(s *discordgo.Session) error {
	config := config.GetConfig()

	// Grab the app ID and guild ID
	appID := config.AppCfg.DiscordAppID
	guildID := config.AppCfg.DiscordDevGuildID

	// Write new config
	disable := false
	config.ModuleCfg.Utility.Enabled = &disable
	err := config.WriteConfig()
	if err != nil {
		return fmt.Errorf("Failed to write config: %w", err)
	}
	log.Debug().Msg("Updated utility module config")

	// Get all registered commands
	registeredCommands, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		return fmt.Errorf("Failed to get registered commands for utility module: %w", err)
	}

	// Filter out utility commands
	log.Debug().Msg("Filtering utility commands for deletion...")
	utilCommands := make([]*discordgo.ApplicationCommand, 0)
	for _, command := range registeredCommands {
		for _, utilityCommand := range commands {
			if command.Name == utilityCommand.Name {
				utilCommands = append(utilCommands, command)
			}
		}
	}

	// Delete utility commands
	log.Debug().Msgf("Deleting %d utility commands", len(utilCommands))
	for _, command := range utilCommands {
		err := s.ApplicationCommandDelete(appID, guildID, command.ID)
		if err != nil {
			return err
		}
	}

	log.Info().Msg("Utility module disabled!")
	return nil
}
