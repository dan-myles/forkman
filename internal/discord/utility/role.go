package utility

import (
	"github.com/avvo-na/forkman/internal/discord/templates"
	"github.com/bwmarrin/discordgo"
)

func (u *UtilityModule) role(s *discordgo.Session, i *discordgo.InteractionCreate) {
	templates.Message(s, i, "Role command")
}
