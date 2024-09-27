package moderation

import (
	"github.com/avvo-na/forkman/internal/discord/templates"
	"github.com/bwmarrin/discordgo"
)

type ModerationState struct {
	ModerationCommands []ModerationCommand `json:"moderation_commands"`
}

type ModerationCommand struct {
	Enabled     bool                         `json:"enabled"`
	CommandData discordgo.ApplicationCommand `json:"command"`
}

var defaultState = ModerationState{
	ModerationCommands: []ModerationCommand{
		{
			Enabled: true,
			CommandData: discordgo.ApplicationCommand{
				Name:        "ban",
				Description: "Ban a user",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "user",
						Description: "User to ban",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
				},
			},
		},
	},
}

func (m *ModerationModule) handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	m.log.Info().Msg("Receieved interaction request")
	switch i.ApplicationCommandData().Name {
	case "ban":
		m.ban(s, i)
	}
}

func (m *ModerationModule) ban(s *discordgo.Session, i *discordgo.InteractionCreate) {
	templates.Message(s, i, "Ban command")
}
