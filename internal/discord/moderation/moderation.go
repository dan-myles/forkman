package moderation

import (
	"encoding/json"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ModerationState struct {
	ModerationCommands []ModerationCommand `json:"moderation_commands"`
}

type ModerationCommand struct {
	Enabled bool                         `json:"enabled"`
	Command discordgo.ApplicationCommand `json:"command"`
}

// Default Configuration, DB State takes precedence
var (
	name         = "Moderation"
	description  = "Fork hammer to-go, please!"
	defaultState = ModerationState{
		ModerationCommands: []ModerationCommand{
			{
				Enabled: false,
				Command: discordgo.ApplicationCommand{
					Name:        "ban",
					Description: "Ban a user",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "User to ban",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    true,
						},
					},
				},
			},
		},
	}
)

type ModerationModule struct {
	guildSnowflake string
	session        *discordgo.Session
	db             *gorm.DB
	log            *zerolog.Logger
	cfg            *config.SentinelConfig
	unhandlers     map[string]*func()
}

func New(
	gID string,
	s *discordgo.Session,
	l *zerolog.Logger,
	db *gorm.DB,
	cfg *config.SentinelConfig,
) *ModerationModule {
	subLogger := l.With().Str("module", "moderation").Logger()

	return &ModerationModule{
		guildSnowflake: gID,
		session:        s,
		db:             db,
		log:            &subLogger,
		cfg:            cfg,
		unhandlers:     make(map[string]*func()),
	}
}

func (m *ModerationModule) Sync() error {
	// Context for logger
	log := m.log.With().
		Str("module", "moderation").
		Str("guild_snowflake", m.guildSnowflake).
		Logger()

	// we need to grab mod from the database
	mod := database.Module{}
	result := m.db.First(&mod, "name = ? AND guild_snowflake = ?", name, m.guildSnowflake)
	if result.RowsAffected == 0 {
		log.Debug().Msg("Module state not found, creating")
		jsonState, _ := json.Marshal(defaultState)
		mod = database.Module{
			Name:           name,
			Description:    description,
			Enabled:        true,
			GuildSnowflake: m.guildSnowflake,
			State:          jsonState,
		}
		m.db.Create(&mod)
		return nil
	}

	// Unmarshal the state
	state := ModerationState{}
	err := json.Unmarshal(mod.State, &state)
	if err != nil {
		return err
	}

	// Grab all command names
	var commandNames []string
	for _, command := range state.ModerationCommands {
		commandNames = append(commandNames, command.Command.Name)
	}

	// Check if we are enabled
	if !mod.Enabled {
		log.Debug().Msg("Module is disabled, removing commands and handlers")

		// disable all commands
		remoteCommands, err := m.session.ApplicationCommands(m.cfg.DiscordAppID, m.guildSnowflake)
		if err != nil {
			return err
		}
		for _, remoteCommand := range remoteCommands {
			for _, commandName := range commandNames {
				if remoteCommand.Name == commandName {
					m.session.ApplicationCommandDelete(
						m.cfg.DiscordAppID,
						m.guildSnowflake,
						remoteCommand.ID,
					)
					log.Debug().Str("command", remoteCommand.Name).Msg("Command disabled")
				}
			}
		}

		// remove all handlers
		for name, unhandle := range m.unhandlers {
			if unhandle != nil {
				log.Debug().Str("command", name).Msg("Removing handler")
				(*unhandle)()
			}
		}

		return nil
	}

	// Sync commands
	log.Info().Msg("Syncing commands")
	remoteCommands, _ := m.session.ApplicationCommands(m.cfg.DiscordAppID, m.guildSnowflake)
	for _, command := range state.ModerationCommands {
		// If the command is enabled, create it (or overwrite it)
		if command.Enabled {
			m.session.ApplicationCommandCreate(
				m.cfg.DiscordAppID,
				m.guildSnowflake,
				&command.Command,
			)
			log.Debug().Str("command", command.Command.Name).Msg("Command enabled")
		}

		// If the command is disabled, delete it
		if !command.Enabled {
			for _, remoteCommand := range remoteCommands {
				if remoteCommand.Name == command.Command.Name {
					m.session.ApplicationCommandDelete(
						m.cfg.DiscordAppID,
						m.guildSnowflake,
						remoteCommand.ID,
					)
					log.Debug().Str("command", remoteCommand.Name).Msg("Command disabled")
				}
			}
		}
	}

	return nil
}
