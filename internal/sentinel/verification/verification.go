package verification

import "github.com/bwmarrin/discordgo"

type VerificationModule struct{}

func New() *VerificationModule {
	return &VerificationModule{}
}

func (m *VerificationModule) Name() string {
	return "verification"
}

func (m *VerificationModule) Enable(s *discordgo.Session) error {
	return nil
}

func (m *VerificationModule) Disable(s *discordgo.Session) error {
	return nil
}
