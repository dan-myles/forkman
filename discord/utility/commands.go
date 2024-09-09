package utility

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "role",
		Description: "role management commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "all",
				Description: "adds role to all members",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "role",
						Description: "to add",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
			{
				Name:        "remove",
				Description: "removes role from specified member",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "member",
						Description: "to remove role from",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
					{
						Name:        "role",
						Description: "to remove",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
			{
				Name:        "add",
				Description: "adds role to specified member",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "member",
						Description: "to add role to",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
					{
						Name:        "role",
						Description: "to add",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
		},
	},
}

var commandHandlers = map[string]func(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	l *zerolog.Logger,
){
	"role": role,
}
