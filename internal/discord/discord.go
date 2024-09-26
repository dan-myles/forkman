package discord

import (
	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/discord/moderation"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Discord struct {
	session           *discordgo.Session
	moderationModules map[string]*moderation.ModerationModule
}

func New(cfg *config.SentinelConfig, log *zerolog.Logger, db *gorm.DB) *Discord {
	s, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		panic(err)
	}

	// Settings
	s.Identify.Intents = discordgo.IntentsAll // What do we need permission for?
	s.SyncEvents = false                      // Launch goroutines for handlers
	s.StateEnabled = true

	// Init modules
	m1 := make(map[string]*moderation.ModerationModule)

	// Global handlers
	s.AddHandler(onGuildCreateGuildUpdate(db, log, cfg, m1))

	// Open the session
	log.Info().Msg("Opening discord session")
	err = s.Open()
	if err != nil {
		panic(err)
	}

	return &Discord{
		session:           s,
		moderationModules: m1,
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
