package verification

import "github.com/bwmarrin/discordgo"

func (m *Verification) SendVerificationPanel(channelId string) error {
	// Create embed message
	embed := &discordgo.MessageEmbed{
		Title:       "Verification",
		Description: "This Discord is for official ASU students only, if you wish to access all of it please click the button below and provide your ASURITE ID.",
		Color:       0x00FF00, // green color
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.ibb.co/MBVt8Mq/arrowfork.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn-longterm.mee6.xyz/plugins/embeds/images/1187144343400751234/75ef7b0e26d7225d196e50f6781f683399de2431236e03092ddf06095a1f024c.png",
		},
	}

	// Create button row to open email input modal
	buttonRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label: "Verify Me",
				Style: discordgo.PrimaryButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: "👆",
				},
				CustomID: CIDVerifyEmailBtn, // Custom ID for button to trigger modal
			},
		},
	}

	// Send message with embed and button
	_, err := m.session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Embed:      embed,
		Components: []discordgo.MessageComponent{buttonRow}, // Only button, no TextInput here
	})
	if err != nil {
		return err
	}

	return nil
}
