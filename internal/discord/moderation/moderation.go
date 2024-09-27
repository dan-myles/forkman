package moderation

import (
	"encoding/json"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	ErrDiscord "github.com/avvo-na/forkman/internal/discord/common"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const (
	name        = "Moderation"
	description = "Fork hammer to-go, please!"
)

type ModerationModule struct {
	guildSnowflake string
	session        *discordgo.Session
	db             *gorm.DB
	log            *zerolog.Logger
	cfg            *config.SentinelConfig
	unhandler      *func()
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
		unhandler:      nil,
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
			Enabled:        false,
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
		commandNames = append(commandNames, command.CommandData.Name)
	}

	// Check if we are enabled
	if !mod.Enabled {
		log.Debug().Msg("Module is disabled, removing commands and handlers")

		// disable all commands
		remoteCommands, _ := m.session.ApplicationCommands(m.cfg.DiscordAppID, m.guildSnowflake)
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
		if m.unhandler != nil {
			log.Debug().Msg("Removing moderation handler")
			(*m.unhandler)()
			m.unhandler = nil
		}

		log.Info().Msg("Module synced & disabled")
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
				&command.CommandData,
			)
			log.Debug().Str("command", command.CommandData.Name).Msg("Command enabled")
		}

		// If the command is disabled, delete it
		if !command.Enabled {
			for _, remoteCommand := range remoteCommands {
				if remoteCommand.Name == command.CommandData.Name {
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

	// If our handler is not there, we need to add it :)
	if m.unhandler == nil {
		log.Debug().Msg("Adding handler")
		fn := m.session.AddHandler(m.handle)
		m.unhandler = &fn
	}

	log.Info().Msg("Module synced")
	return nil
}

// TODO: IMPLEMENT THESE
func (m *ModerationModule) DisableAndSync() error {
	log := m.log.With().
		Str("module", "moderation").
		Str("guild_snowflake", m.guildSnowflake).
		Logger()

	// we need to grab mod from the database
	mod := database.Module{}
	result := m.db.First(&mod, "name = ? AND guild_snowflake = ?", name, m.guildSnowflake)
	if result.RowsAffected == 0 {
		return ErrDiscord.ErrGuildNotFound
	}

	if !mod.Enabled {
		return nil
	}

	// Unmarshal the state
	mod.Enabled = false
	m.db.Save(&mod)

	// Sync
	log.Info().Msg("Disabling and syncing!")
	return m.Sync()
}

func (m *ModerationModule) EnableAndSync() error {
	log := m.log.With().
		Str("module", "moderation").
		Str("guild_snowflake", m.guildSnowflake).
		Logger()

	// we need to grab mod from the database
	mod := database.Module{}
	result := m.db.First(&mod, "name = ? AND guild_snowflake = ?", name, m.guildSnowflake)
	if result.RowsAffected == 0 {
		return ErrDiscord.ErrGuildNotFound
	}
	if mod.Enabled {
		return nil
	}

	// Update the state
	mod.Enabled = true
	m.db.Save(&mod)

	// Sync
	log.Info().Msg("Enabling and syncing!")
	return m.Sync()
}
