package discord

import (
	"errors"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/discord/moderation"
	"github.com/avvo-na/forkman/internal/discord/verification"
	"github.com/bwmarrin/discordgo"
	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Discord struct {
	session             *discordgo.Session
	moderationModules   map[string]*moderation.Moderation
	verificationModules map[string]*verification.Verification
}

var ErrModuleNotFound = errors.New("module not found")

func New(cfg *config.SentinelConfig, log *zerolog.Logger, db *gorm.DB, email *resend.Client) *Discord {
	s, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		panic(err)
	}

	// Settings
	s.Identify.Intents = discordgo.IntentsAll // What do we need permission for?
	s.SyncEvents = false                      // Launch goroutines for handlers
	s.StateEnabled = true

	// Module stores
	mm := make(map[string]*moderation.Moderation)
	vm := make(map[string]*verification.Verification)

	// Global handlers
	s.AddHandler(onGuildCreateGuildUpdate(db, log, cfg, email, mm, vm))
	s.AddHandler(onReadyNotify(log))

	// Open the session
	log.Info().Msg("Opening discord session")
	err = s.Open()
	if err != nil {
		panic(err)
	}

	return &Discord{
		session:             s,
		moderationModules:   mm,
		verificationModules: vm,
	}
}

func (d *Discord) Open() error {
	err := d.session.Open()
	if err != nil {
		return err
	}

	return nil
}

func (d *Discord) Close() error {
	err := d.session.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d *Discord) GetSession() *discordgo.Session {
	return d.session
}

func (d *Discord) GetModerationModule(guildSnowflake string) (*moderation.Moderation, error) {
	mod, ok := d.moderationModules[guildSnowflake]
	if !ok {
		return nil, ErrModuleNotFound
	}

	return mod, nil
}

func (d *Discord) GetVerificationModule(guildSnowflake string) (*verification.Verification, error) {
	mod, ok := d.verificationModules[guildSnowflake]
	if !ok {
		return nil, ErrModuleNotFound
	}

	return mod, nil
}
