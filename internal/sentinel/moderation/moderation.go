package moderation

import (
	"fmt"

	"github.com/avvo-na/devil-guard/common/log"
	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "kick",
		Description: "kicks a member from the guild",
		Type:        discordgo.ChatApplicationCommand,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "to kick",
				Required:    true,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"kick": kick,
}

func kick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Kicking user... (NOT IMPL)",
		},
	})
}

type ModerationModule struct{}

func New() *ModerationModule {
	return &ModerationModule{}
}

func (u *ModerationModule) Name() string {
	return "moderation"
}

func (u *ModerationModule) Load(s *discordgo.Session) error {
	// Grab the config
	config := config.GetConfig()
	config.RWMutex.RLock()
	defer config.RWMutex.RUnlock()

	appID := config.AppCfg.DiscordAppID
	guildID := config.AppCfg.DiscordDevGuildID

	// If the module is disabled, skip registration
	if !*config.ModuleCfg.Moderation.Enabled {
		log.Info().
			Str("module", u.Name()).
			Msg("Module is disabled, skipping registration")
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
		Msg("Registering moderation command handlers...")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, ok := commandHandlers[i.ApplicationCommandData().Name]
		if !ok {
			return
		}

		handler(s, i)
	})

	log.Debug().Msg("Moderation module registration complete")
	return nil
}

func (u *ModerationModule) Enable(s *discordgo.Session) error {
	// Grab the config
	config := config.GetConfig()
	config.RWMutex.Lock()
	defer config.RWMutex.Unlock()

	// Grab the app ID and guild ID
	appID := config.AppCfg.DiscordAppID
	guildID := config.AppCfg.DiscordDevGuildID

	// Write new config
	enable := true
	config.ModuleCfg.Moderation.Enabled = &enable

	err := config.WriteConfig()
	if err != nil {
		return fmt.Errorf("Failed to write config: %w", err)
	}
	log.Debug().
		Interface("config", config).
		Msg("Updated moderation module config")

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
		Msg("Registering moderation command handlers...")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, ok := commandHandlers[i.ApplicationCommandData().Name]
		if !ok {
			return
		}

		handler(s, i)
	})

	log.Debug().Msg("Moderation module registration complete")
	return nil
}

func (u *ModerationModule) Disable(s *discordgo.Session) error {
	// Grab the config
	config := config.GetConfig()
	config.RWMutex.Lock()
	defer config.RWMutex.Unlock()

	// Grab the app ID and guild ID
	appID := config.AppCfg.DiscordAppID
	guildID := config.AppCfg.DiscordDevGuildID

	// Disable the module
	disable := false
	config.ModuleCfg.Moderation.Enabled = &disable

	// Write new config
	err := config.WriteConfig()
	if err != nil {
		return fmt.Errorf("Failed to write config: %w", err)
	}
	log.Debug().Msg("Updated moderation module config")

	// Get all registered commands
	registeredCommands, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		return fmt.Errorf("Failed to get registered commands for moderation module: %w", err)
	}

	// Filter out moderation commands
	log.Debug().Msg("Filtering moderation commands for deletion...")
	cmds := make([]*discordgo.ApplicationCommand, 0)
	for _, globalCmd := range registeredCommands {
		for _, modCmd := range commands {
			if globalCmd.Name == modCmd.Name {
				cmds = append(cmds, globalCmd)
			}
		}
	}

	// Delete moderation commands
	log.Debug().Msgf("Deleting %d moderation commands", len(cmds))
	for _, command := range cmds {
		err := s.ApplicationCommandDelete(appID, guildID, command.ID)
		if err != nil {
			return err
		}
	}

	log.Debug().Msg("Moderation module deregistration complete")
	return nil
}
