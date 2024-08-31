package sentinel

import (
	"fmt"

	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
)

type Sentinel struct {
	Session *discordgo.Session
}

func New() (*Sentinel, error) {
	s, err := discordgo.New("Bot " + config.AppCfg.DiscordBotToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Discord session %w", err)
	}

	return &Sentinel{
		Session: s,
	}, nil
}

func (s *Sentinel) Start() error {
	err := s.Session.Open()
	if err != nil {
		return fmt.Errorf("Failed to open Discord session %w", err)
	}

	return nil
}

func (s *Sentinel) Stop() error {
	err := s.Session.Close()
	if err != nil {
		return fmt.Errorf("Failed to close Discord session %w", err)
	}

	return nil
}
