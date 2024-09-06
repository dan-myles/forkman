package utility

import (
	"fmt"

	"github.com/avvo-na/devil-guard/common/log"
	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "role",
		Description: "role management commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "all",
				Description: "gives roles to all members",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "role",
						Description: "role to add",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
			{
				Name:        "remove",
				Description: "Removes role from specified user",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "role",
						Description: "role to add",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
					{
						Name:        "user",
						Description: "user to remove role from",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
				},
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"role": role,
}

type UtilityModule struct{}

func New() *UtilityModule {
	return &UtilityModule{}
}

func (u *UtilityModule) Name() string {
	return "utility"
}

func (u *UtilityModule) Load(s *discordgo.Session) error {
	// Grab the config
	config := config.GetConfig()
	config.RWMutex.RLock()
	defer config.RWMutex.RUnlock()

	appID := config.AppCfg.DiscordAppID
	guildID := config.AppCfg.DiscordDevGuildID

	// If the module is disabled, skip registration
	if !*config.ModuleCfg.Utility.Enabled {
		log.Debug().Msg("Utility module is disabled, skipping...")
		return nil
	}

	// Register all commands
	for _, command := range commands {
		log.Debug().
			Str("appID", appID).
			Str("guildID", guildID).
			Str("command", command.Name).
			Msg("Registering command")

		_, err := s.ApplicationCommandCreate(appID, guildID, command)
		if err != nil {
			return err
		}
	}

	log.Debug().
		Interface("commands", commands).
		Msg("Registering utility command handlers...")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, ok := commandHandlers[i.ApplicationCommandData().Name]
		if !ok {
			return
		}

		handler(s, i)
	})

	log.Debug().Msg("Utility module registration complete")
	return nil
}

func (u *UtilityModule) Enable(s *discordgo.Session) error {
	// Grab the config
	config := config.GetConfig()
	config.RWMutex.Lock()
	defer config.RWMutex.Unlock()

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
	log.Debug().
		Interface("config", config).
		Msg("Updated utility module config")

	// Register all commands
	for _, command := range commands {
		log.Debug().
			Str("appID", appID).
			Str("guildID", guildID).
			Str("command", command.Name).
			Msg("Registering command")

		_, err := s.ApplicationCommandCreate(appID, guildID, command)
		if err != nil {
			return err
		}
	}

	log.Debug().
		Interface("commands", commands).
		Msg("Registering utility command handlers...")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, ok := commandHandlers[i.ApplicationCommandData().Name]
		if !ok {
			return
		}

		handler(s, i)
	})

	log.Debug().Msg("Utility module registration complete")
	return nil
}

func (u *UtilityModule) Disable(s *discordgo.Session) error {
	// Grab the config
	config := config.GetConfig()
	config.RWMutex.Lock()
	defer config.RWMutex.Unlock()

	// Grab the app ID and guild ID
	appID := config.AppCfg.DiscordAppID
	guildID := config.AppCfg.DiscordDevGuildID

	// Disable the module
	disable := false
	config.ModuleCfg.Utility.Enabled = &disable

	// Write new config
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

	log.Debug().Msg("Utility module deregistration complete")
	return nil
}
