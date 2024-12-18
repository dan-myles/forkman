package qna

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "qna-disable",
		Description: "disables the Q&A module",
	},
	{
		Name:        "qna-enable",
		Description: "enables the Q&A module",
	},
}
