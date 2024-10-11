package moderation

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TODO:
// - add some sort of system to reject enable/disable calls
//   when they are in progress, prob some sort of mutex

type ModerationConfig struct {
	ImmuneRoles []string `json:"immune_roles"`
}

type Moderation struct {
	guildName      string
	guildSnowflake string
	appId          string
	session        *discordgo.Session
	repo           *Repository
	log            *zerolog.Logger
	unhandle       *func()
}

const (
	name        = "Moderation"
	description = "Fork-tilities to-go please!"
)

var (
	ErrModuleAlreadyDisabled = errors.New("module is already disabled")
	ErrModuleAlreadyEnabled  = errors.New("module is already enabled")
)

func New(
	guildName string,
	guildSnowflake string,
	appId string,
	session *discordgo.Session,
	db *gorm.DB,
	log *zerolog.Logger,
) *Moderation {
	l := log.With().
		Str("module", name).
		Str("guild_snowflake", guildSnowflake).
		Str("guild_name", guildName).
		Logger()

	return &Moderation{
		guildName:      guildName,
		guildSnowflake: guildSnowflake,
		appId:          appId,
		session:        session,
		repo:           NewRepository(db),
		log:            &l,
		unhandle:       nil,
	}
}

func (m *Moderation) Load() error {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err == gorm.ErrRecordNotFound {
		m.log.Debug().Msg("module not found, creating...")

		// Default general config (empty)
		cfgJson, _ := json.Marshal(ModerationConfig{ImmuneRoles: []string{""}})

		// Default command config (all enabled)
		cmdMap := make(map[string]bool)
		for _, command := range commands {
			cmdMap[command.Name] = true
		}
		cmdJson, _ := json.Marshal(cmdMap)

		insert := &database.Module{
			GuildSnowflake: m.guildSnowflake,
			Name:           name,
			Description:    description,
			Enabled:        false,
			Config:         cfgJson,
			Commands:       cmdJson,
		}

		if mod, err = m.repo.CreateModule(insert); err != nil {
			return fmt.Errorf("unable to create moderation module: %w", err)
		}
	}

	// If we are not enabled, don't do anything!
	if !mod.Enabled {
		m.log.Debug().Msg("module disabled, skipping...")
		return nil
	}

	// Grab command state
	var cmds map[string]bool
	err = json.Unmarshal([]byte(mod.Commands), &cmds)
	if err != nil {
		return fmt.Errorf("critical error unmarshalling cmd json: %w", err)
	}

	// Register commands (or overwrite!)
	for _, command := range commands {
		if !cmds[command.Name] {
			m.log.Debug().Str("command", command.Name).Msg("command disabled, skipping...")
			continue
		}

		m.log.Debug().
			Str("command_name", command.Name).
			Str("command_id", command.ID).
			Msg("command enabled")

		_, err := m.session.ApplicationCommandCreate(m.appId, m.guildSnowflake, command)
		if err != nil {
			m.log.Error().Err(err).Str("command", command.Name).Msg("error registering command!")
		}
	}

	// Register handler
	fn := m.session.AddHandler(m.handle)
	m.unhandle = &fn

	m.log.Info().Msgf("module %s loaded", mod.Name)
	return nil
}

func (m *Moderation) Disable() error {
	// Read DB state
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err != nil {
		return err
	}

	if !mod.Enabled {
		return ErrModuleAlreadyDisabled
	}
	mod.Enabled = false

	// Save DB state
	_, err = m.repo.UpdateModule(mod)
	if err != nil {
		return err
	}

	// Grab remote commands
	remoteCommands, err := m.session.ApplicationCommands(m.appId, m.guildSnowflake)
	if err != nil {
		return fmt.Errorf("unable to grab remote commands from guild: %w", err)
	}

	// Delete all commands
	for _, command := range remoteCommands {
		m.log.Debug().
			Str("command_name", command.Name).
			Str("command_id", command.ID).
			Msg("command disabled")

		err := m.session.ApplicationCommandDelete(m.appId, m.guildSnowflake, command.ID)
		if err != nil {
			m.log.Error().
				Err(err).
				Str("command_name", command.Name).
				Str("command_id", command.ID).
				Msg("error deleting application command")
		}
	}

	// Unregister handler
	if m.unhandle != nil {
		(*m.unhandle)()
		m.unhandle = nil
	}

	m.log.Info().Msg("module disabled")
	return nil
}

func (m *Moderation) Enable() error {
	// Read DB state
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err != nil {
		return err
	}

	if mod.Enabled {
		return ErrModuleAlreadyEnabled
	}
	mod.Enabled = true

	// Save DB state
	_, err = m.repo.UpdateModule(mod)
	if err != nil {
		return err
	}

	// Grab command state
	var cmds map[string]bool
	err = json.Unmarshal([]byte(mod.Commands), &cmds)
	if err != nil {
		return fmt.Errorf("critical error unmarshalling cmd json: %w", err)
	}

	// Register commands (or overwrite!)
	for _, command := range commands {
		if !cmds[command.Name] {
			m.log.Debug().Str("command", command.Name).Msg("command disabled, skipping...")
			continue
		}

		m.log.Debug().
			Str("command_name", command.Name).
			Str("command_id", command.ID).
			Msg("command enabled")

		_, err := m.session.ApplicationCommandCreate(m.appId, m.guildSnowflake, command)
		if err != nil {
			m.log.Error().Err(err).Str("command", command.Name).Msg("error registering command")
		}
	}

	// Register handler
	fn := m.session.AddHandler(m.handle)
	m.unhandle = &fn

	return nil
}
