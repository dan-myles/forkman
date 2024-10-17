package moderation

import (
	"github.com/avvo-na/forkman/internal/discord/templates"
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "mute",
		Description: "a user for a certain duration",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "user to mute",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "length",
				Description: "length of the timeout (ie. '1d', '30m', '60s')",
				Required:    true,
			},
		},
	},
}

func (m *Moderation) handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	cmd := i.ApplicationCommandData().Name
	m.log.Info().
		Str("command_name", cmd).
		Str("user_id", i.Member.User.ID).
		Str("user_name", i.Member.User.Username).
		Msg("served discord interaction request")

	switch cmd {
	case "mute":
		m.mute(s, i)
	}
}

func (m *Moderation) mute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	templates.Message(s, i, "hi!")
}
