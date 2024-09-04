package utility

import "github.com/bwmarrin/discordgo"

type UtilityModule struct{}

func (u *UtilityModule) Name() string {
	return "utility"
}

func (u *UtilityModule) Enable(s *discordgo.Session) error {
	return nil
}

func (u *UtilityModule) Disable(s *discordgo.Session) error {
	return nil
}
