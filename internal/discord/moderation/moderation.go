package moderation

import (
	"encoding/json"
	"fmt"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

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
	cfg            *ModerationConfig
	unhandle       *func()
}

const (
	name        = "Moderation"
	description = "Fork-tilities to-go please!"
)

func New(
	gn string,
	gs string,
	id string,
	s *discordgo.Session,
	db *gorm.DB,
	log *zerolog.Logger,
) *Moderation {
	l := log.With().
		Str("module", name).
		Str("guild_snowflake", gs).
		Str("guild_name", gn).
		Logger()

	return &Moderation{
		guildName:      gn,
		guildSnowflake: gs,
		appId:          id,
		session:        s,
		repo:           NewRepository(db),
		log:            &l,
		cfg:            &ModerationConfig{},
		unhandle:       nil,
	}
}

func (m *Moderation) Load() error {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err == gorm.ErrRecordNotFound {
		// Default general config
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

		if _, err := m.repo.CreateModule(insert); err != nil {
			return fmt.Errorf("unable to create moderation module: %w", err)
		}
	}

	// If we are not enabled, don't do anything!
	if !mod.Enabled {
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
			log.Debug().Str("command", command.Name).Msg("Command disabled, skipping...")
			continue
		}

		_, err := m.session.ApplicationCommandCreate(m.appId, m.guildSnowflake, command)
		if err != nil {
			log.Error().Err(err).Str("command", command.Name).Msg("Error registering command!")
		}
	}

	// Register handler
	fn := m.session.AddHandler(m.handle)
	m.unhandle = &fn

	log.Info().Msg("Module loaded")
	return nil
}
