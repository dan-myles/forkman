package verification

import "github.com/bwmarrin/discordgo"

var (
	VerifyEmailBtnCID = "verify_email_button"
)

func (m *Verification) SendVerificationPanel(channelId string) error {
	// Create embed message
	embed := &discordgo.MessageEmbed{
		Title:       "Verification Panel",
		Description: "Please enter your email to start verification.",
		Color:       0x00FF00, // green color
	}

	// Create action row with email input and submit button
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "email_input_field",
				Label:       "Enter Email",
				Style:       discordgo.TextInputShort,
				Placeholder: "example@example.com",
				Required:    true,
			},
		},
	}

	// Create button row to submit email
	buttonRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Submit Email",
				Style:    discordgo.PrimaryButton,
				CustomID: VerifyEmailBtnCID,
			},
		},
	}

	// Send message with embed, email input field, and submit button
	_, err := m.session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Embed:      embed,
		Components: []discordgo.MessageComponent{actionRow, buttonRow},
	})

	if err != nil {
		return err
	}

	return nil
}
