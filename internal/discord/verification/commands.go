package verification

import (
	"os"

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
		Name:        "verify",
		Description: "manually verify a user's email",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "to verify",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "email",
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
	print("verify started")
	member := i.ApplicationCommandData().Options[0].UserValue(s)
	email := i.ApplicationCommandData().Options[1].StringValue()
	email, err := m.repo.ManualVerification(i.GuildID, member.ID, email)

	status := "❌ Manual Verification Failed"
	if err == nil {
		status = "✅ Manual Verification Successful"
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
							Value: email,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	s.GuildMemberRoleRemove(m.guildSnowflake, i.Member.User.ID, os.Getenv("ROLE_TO_REMOVE"))
	s.GuildMemberRoleAdd(m.guildSnowflake, i.Member.User.ID, os.Getenv("ROLE_TO_ADD"))

	if err != nil {
		m.log.Error().Err(err).Msg("Failed to manually verify email")
	}

}
