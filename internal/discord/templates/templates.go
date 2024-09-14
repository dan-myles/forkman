package templates

import "github.com/bwmarrin/discordgo"

func MessageEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func ErrMessageEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	MessageEphemeral(s, i, err.Error())
}

func Message(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func ErrMessage(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	Message(s, i, err.Error())
}
