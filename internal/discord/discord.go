package discord

import (
	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/discord/utility"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Discord struct {
	session *discordgo.Session
	util    *utility.UtilityModule
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

	// Load modules
	log.Info().Msg("Loading utility module")
	utility := utility.New(s, log, cfg, db)
	utility.Sync()

	// Open the session
	log.Info().Msg("Opening discord session")
	err = s.Open()
	if err != nil {
		panic(err)
	}

	return &Discord{
		session: s,
		util:    utility,
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

func (d *Discord) GetUtility() *utility.UtilityModule {
	return d.util
}
