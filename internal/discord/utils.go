package discord

import "github.com/bwmarrin/discordgo"

func GetApplicationCommandType(cmd discordgo.ApplicationCommandType) string {
	switch cmd {
	case 1:
		return "ChatApplicationCommand"
	case 2:
		return "UserApplicationCommand"
	case 3:
		return "MessageApplicationCommand"
	default:
		return "Unknown"
	}
}
