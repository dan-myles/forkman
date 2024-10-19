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

func (m *Moderation) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	mod, err := m.repo.ReadModule(i.GuildID)
	if err != nil {
		return
	}

	if !mod.Enabled {
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		m.handleCommand(s, i)
	}
}

func (m *Moderation) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := i.ApplicationCommandData().Name
	switch cmd {
	case "mute":
		m.mute(s, i)
	}
}

func (m *Moderation) mute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	templates.Message(s, i, "hi!")
}
