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
	s, err := discordgo.New("Bot " + config.AppCfg.DiscordBotToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Discord session! Check your token.")
	}

	return &Sentinel{
		Session: s,
	}
}

func (s *Sentinel) Start() {
}

func (s *Sentinel) Stop() {
}
