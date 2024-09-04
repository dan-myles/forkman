package utility

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping the bot",
	},
	{
		Name:        "role",
		Description: "Manage roles",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "all",
				Description: "Give to all users",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": ping,
	"role": role,
}

type UtilityModule struct{}

func New() *UtilityModule {
	return &UtilityModule{}
}

func (u *UtilityModule) Name() string {
	return "utility"
}

func (u *UtilityModule) Enable(s *discordgo.Session) error {
	return nil
}

func (u *UtilityModule) Disable(s *discordgo.Session) error {
	return nil
}
