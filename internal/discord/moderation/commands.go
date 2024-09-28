package moderation

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ban",
			Description: "forkhammer inbound",
		},
	}
)

func (m *Moderation) handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	m.log.Info().Msg("Receieved interaction request")
	switch i.ApplicationCommandData().Name {
	case "ban":
		// m.ban(s, i)
	}
}
