package verification

import "github.com/bwmarrin/discordgo"

func (m *Verification) SendVerificationPanel(channelId string) error {
	// Create embed message
	embed := &discordgo.MessageEmbed{
		Title:       "Verification Panel",
		Description: "Click the button to start verification and enter your email.",
		Color:       0x00FF00, // green color
	}

	// Create button row to open email input modal
	buttonRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Submit Email",
				Style:    discordgo.PrimaryButton,
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
