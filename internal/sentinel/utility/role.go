package utility

import "github.com/bwmarrin/discordgo"

// Top level handler for all role commands
func role(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "all":
		roleAll(s, i)
	}
}

func roleAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Role all",
		},
	})
}
