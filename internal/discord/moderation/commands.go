package moderation

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
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
}
