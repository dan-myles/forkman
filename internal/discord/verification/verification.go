package verification

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/resend/resend-go/v2"
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
	email          *resend.Client
	repo           *Repository
	log            *zerolog.Logger
	unlisten       *func()
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
	email *resend.Client,
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
		email:          email,
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

		insert := &database.Module{
			GuildSnowflake: m.guildSnowflake,
			Name:           name,
			Description:    description,
			Enabled:        false,
			Config:         cfgJson,
		}

		if mod, err = m.repo.CreateModule(insert); err != nil {
			return fmt.Errorf("unable to create verification module: %w", err)
		}
	}

	fn := m.session.AddHandler(m.listen)
	m.unlisten = &fn

	m.log.Info().Msgf("module %s loaded", mod.Name)
	return nil
}

func (m *Verification) Disable() error {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err != nil {
		return err
	}

	if !mod.Enabled {
		return ErrModuleAlreadyDisabled
	}
	mod.Enabled = false

	_, err = m.repo.UpdateModule(mod)
	if err != nil {
		return err
	}

	if m.unlisten != nil {
		(*m.unlisten)()
		m.unlisten = nil
	}

	m.log.Info().Msg("module disabled")
	return nil
}

func (m *Verification) Enable() error {
	mod, err := m.repo.ReadModule(m.guildSnowflake)
	if err != nil {
		return err
	}

	if mod.Enabled {
		return ErrModuleAlreadyEnabled
	}
	mod.Enabled = true

	_, err = m.repo.UpdateModule(mod)
	if err != nil {
		return err
	}

	fn := m.session.AddHandler(m.listen)
	m.unlisten = &fn

	m.log.Info().Msg("module enabled")
	return nil
}
