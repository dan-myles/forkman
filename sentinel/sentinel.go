package sentinel

import (
	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Sentinel struct {
	Session *discordgo.Session
}

func New() *Sentinel {
	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + config.AppCfg.DiscordBotToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Discord session")
	}

	// Declare intents
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	return &Sentinel{
		Session: dg,
	}
}

func (s *Sentinel) Start() {
	log.Info().Msg("Starting bot...")

	err := s.Session.Open()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open Discord session")
	}
}

func (s *Sentinel) Stop() {
	log.Info().Msg("Stopping bot...")
	s.Session.Close()
}

func (s *Sentinel) RegisterPlugins() {
	log.Info().Msg("Registering plugins...")
}
