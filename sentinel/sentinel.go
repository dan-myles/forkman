package sentinel

import (
	"github.com/avvo-na/devil-guard/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Sentinel struct {
	Client *discordgo.Session
}

func New() *Sentinel {
	dg, err := discordgo.New("Bot " + utils.ConfigData.DiscordBotToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Discord session")
	}

	return &Sentinel{
		Client: dg,
	}
}

func (s *Sentinel) Start() {
	s.Client.Identify.Intents = discordgo.IntentsGuildMessages

	err := s.Client.Open()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open Discord session")
	}
}

func (s *Sentinel) Stop() {
	s.Client.Close()
}
