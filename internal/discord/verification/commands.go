package verification

import (
	"github.com/avvo-na/forkman/common/colors"
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "email",
		Description: "grab a user's official ASU email",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "to search for",
				Required:    true,
			},
		},
	},
	{
		Name: 	  "verify",
		Description: "manually verify a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "to verify",
				Required:    true,
			},
	},
},
}

func (m *Verification) email(s *discordgo.Session, i *discordgo.InteractionCreate) {
	member := i.ApplicationCommandData().Options[0].UserValue(s)
	email, _ := m.repo.ReadEmail(i.GuildID, member.ID)

	addr := "N/A"
	status := "❌ Unverified"
	if email != nil {
		addr = email.Address
		status = "✅ Verified"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Verification Status",
					Color: colors.ASUMaroon,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: member.AvatarURL("4096"),
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Username",
							Value:  member.Username,
							Inline: true,
						},
						{
							Name:   "ID",
							Value:  member.ID,
							Inline: true,
						},
						{
							Name:  "Status",
							Value: status,
						},
						{
							Name:  "Email",
							Value: addr,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}


func (m *Verification) verify(s *discordgo.Session, i *discordgo.InteractionCreate) {
	member := i.ApplicationCommandData().Options[0].UserValue(s)
	err := m.repo.VerifyUser(i.GuildID, member.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to verify user",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "User verified",
		},
	})
}