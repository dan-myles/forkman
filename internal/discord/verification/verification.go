package verification

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type VerificationConfig struct {
	Provider      string `json:"provider"`
	SenderAddress string `json:"sender_address"`
}

type Verification struct {
	guildName      string
	guildSnowflake string
	appId          string
	session        *discordgo.Session
	emailClient    *ses.Client
	repo           *Repository
	log            *zerolog.Logger
}

const (
	name        = "Verification"
	description = "Protect against raids!"
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
	email *ses.Client,
	log *zerolog.Logger,
) *Verification {
	l := log.With().
		Str("module", name).
		Str("guild_snowflake", guildSnowflake).
		Str("guild_name", guildName).
		Logger()

	return &Verification{
		guildName:      guildName,
		guildSnowflake: guildSnowflake,
		appId:          appId,
		session:        session,
		emailClient:    email,
		repo:           NewRepository(db),
		log:            &l,
	}
}

func (m *Verification) Load() error {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err == gorm.ErrRecordNotFound {
		m.log.Debug().Msg("module not found, creating...")

		// Default general config (empty)
		cfgJson, _ := json.Marshal(VerificationConfig{Provider: "", SenderAddress: ""})

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
			return fmt.Errorf("unable to create verification module: %w", err)
		}
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

			m.log.Info().Msgf("added new command %s to DB state", command.Name)
			continue
		}
	}

	// If we are not enabled, don't do anything!
	if !mod.Enabled {
		m.log.Debug().Msg("module disabled, skipping...")
		return nil
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

func (m *Verification) Disable() error {
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

func (m *Verification) Enable() error {
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

func (m *Verification) Status() (bool, error) {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err != nil {
		return false, err
	}

	if !mod.Enabled {
		return false, nil
	}

	return true, nil
}

func (m *Verification) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	mod, err := m.repo.ReadModule(i.GuildID)
	if err != nil {
		return
	}

	if !mod.Enabled {
		m.log.Info().Msg("interaction request interuppted, module is disabled")
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		m.handleCommand(s, i)
	case discordgo.InteractionMessageComponent:
		m.handleComponent(s, i)
	case discordgo.InteractionModalSubmit:
		m.handleModal(s, i)
	}
}

func (m *Verification) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name

	switch name {
	case "email":
		m.email(s, i)
	case "verify":
		m.verify(s, i)
	}
}

func (m *Verification) handleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cid := i.MessageComponentData().CustomID
	switch cid {
	case CIDVerifyEmailBtn:
		m.handleCIDVerifyEmailBtn(s, i)
	case CIDVerifyEmailCodeBtn:
		m.handleCIDVerifyEmailCodeBtn(s, i)
	default:
		m.log.Error().
			Str("custom_id", cid).
			Msg("unhandled interaction")
	}
}

func (m *Verification) handleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cid := i.ModalSubmitData().CustomID
	switch cid {
	case CIDVerifyEmailModal:
		m.handleCIDVerifyEmailModal(s, i)
	case CIDVerifyEmailCodeModal:
		m.handleCIDVerifyEmailCodeModal(s, i)
	default:
		m.log.Error().
			Str("custom_id", cid).
			Msg("unhandled interaction")
	}
}
