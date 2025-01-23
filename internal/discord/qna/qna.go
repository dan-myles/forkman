package qna

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type QNAConfig struct{}

type QNA struct {
	guildName       string
	guildSnowflake  string
	appId           string
	session         *discordgo.Session
	bedrock         *bedrockagentruntime.Client
	forumChannelId  string
	knowledgeBaseId string
	repo            *Repository
	log             *zerolog.Logger
}

const (
	name        = "QNA"
	description = "AI Q&A for your server!"
	modelId     = "amazon.nova-pro-v1:0"
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
	bedrock *bedrockagentruntime.Client,
	forumChannelId string,
	knowledgeBaseId string,
	db *gorm.DB,
	log *zerolog.Logger,
) *QNA {
	l := log.With().
		Str("module", name).
		Str("guild_name", guildName).
		Str("guild_snowflake", guildSnowflake).
		Logger()

	return &QNA{
		guildName:       guildName,
		guildSnowflake:  guildSnowflake,
		appId:           appId,
		session:         session,
		bedrock:         bedrock,
		forumChannelId:  forumChannelId,
		knowledgeBaseId: knowledgeBaseId,
		repo:            NewRepository(db),
		log:             &l,
	}
}

func (m *QNA) Load() error {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err == gorm.ErrRecordNotFound {
		m.log.Debug().Msg("module not found, creating...")

		// Default general config (empty)
		cfgJson, _ := json.Marshal(QNAConfig{})

		// Default command config (all enabled)
		cmdMap := make(map[string]bool)
		for _, command := range commands {
			cmdMap[command.Name] = true
		}
		cmdJson, _ := json.Marshal(cmdMap)

		// Default module config
		insert := &database.Module{
			GuildSnowflake: m.guildSnowflake,
			Name:           name,
			Description:    description,
			Enabled:        false,
			Config:         cfgJson,
			Commands:       cmdJson,
		}

		if mod, err = m.repo.CreateModule(insert); err != nil {
			return fmt.Errorf("unable to create qna module: %w", err)
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

	// If new commands are added, add them to the DB state
	for _, command := range commands {
		if _, ok := cmds[command.Name]; !ok {
			cmds[command.Name] = true
			mod.Commands, _ = json.Marshal(cmds)

			_, err = m.repo.UpdateModule(mod)
			if err != nil {
				return fmt.Errorf("unable to update module: %w", err)
			}

			m.log.Info().Msgf("Added new command %s to DB state", command.Name)
			continue
		}
	}

	// TODO: only register commands that are not registered on remote
	// so eventually check remote and do some sort of "merge"
	for _, command := range commands {
		if !cmds[command.Name] {
			m.log.Debug().Str("command", command.Name).Msg("command disabled, skipping...")
			continue
		}

		_, err := m.session.ApplicationCommandCreate(m.appId, m.guildSnowflake, command)
		if err != nil {
			m.log.Error().Err(err).Str("command", command.Name).Msg("error registering command!")
		}

		m.log.Debug().
			Str("command_name", command.Name).
			Str("command_id", command.ID).
			Msg("command enabled")
	}

	m.log.Debug().Msgf("module %s loaded", mod.Name)
	return nil
}

func (m *QNA) Disable() error {
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

	m.log.Info().Msg("module disabled")
	return nil
}

func (m *QNA) Enable() error {
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
			m.log.Debug().
				Str("command_name", command.Name).
				Msg("command disabled, skipping...")
			continue
		}

		ret, err := m.session.ApplicationCommandCreate(m.appId, m.guildSnowflake, command)
		if err != nil {
			m.log.Error().Err(err).Str("command", command.Name).Msg("error registering command")
		}

		m.log.Debug().
			Str("command_name", ret.Name).
			Str("command_id", ret.ID).
			Msg("command enabled")
	}

	return nil
}

func (m *QNA) Status() (bool, error) {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err != nil {
		return false, err
	}

	if !mod.Enabled {
		return false, nil
	}

	return true, nil
}

func (m *QNA) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	mod, err := m.repo.ReadModule(i.GuildID)
	if err != nil {
		return
	}

	if !mod.Enabled {
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		m.handleCommand(s, i)
	case discordgo.InteractionMessageComponent:
		m.handleComponent(s, i)
	}
}

func (m *QNA) OnMessageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.Bot {
		return
	}

	if msg == nil {
		return
	}

	mod, err := m.repo.ReadModule(msg.GuildID)
	if err != nil {
		return
	}

	if !mod.Enabled {
		return
	}

	m.handleQNARequest(s, msg)
}

func (m *QNA) handleCommand(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name

	switch name {
	default:
		m.log.Info().Msg("command not found")
	}
}

func (m *QNA) handleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cid := i.MessageComponentData().CustomID

	switch cid {
	case CIDAdditionalAssistanceBtn:
		m.handleCIDAdditionalAssistanceBtn(s, i)
	case CIDSatisfactoryAnswerBtn:
		m.handleCIDSatisfactoryAnswerBtn(s, i)
	default:
		m.log.Error().
			Str("custom_id", cid).
			Msg("unhandled interaction")
	}
}
