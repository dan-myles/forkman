package verification

import (
	"github.com/avvo-na/forkman/internal/discord/templates"
	"github.com/bwmarrin/discordgo"
)

var (
	CIDVerifyEmailBtn       = "verify_email_button"
	CIDVerifyEmailModal     = "verify_email_modal"
	CIDVerifyEmailCodeBtn   = "verify_email_code_button"
	CIDVerifyEmailCodeModal = "verify_email_code_modal"
)

func (m *Verification) handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
}

func (m *Verification) listen(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := ""
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		customID = i.MessageComponentData().CustomID
	case discordgo.InteractionModalSubmit:
		customID = i.ModalSubmitData().CustomID
	default:
		return
	}

	switch customID {
	case CIDVerifyEmailBtn:
		go m.handleCIDVerifyEmailBtn(s, i)
		return
	case CIDVerifyEmailModal:
		go m.handleCIDVerifyEmailModal(s, i)
		return
	case CIDVerifyEmailCodeBtn:
		go m.handleCIDVerifyEmailCodeBtn(s, i)
		return
	case CIDVerifyEmailCodeModal:
		go m.handleCIDVerifyEmailCodeModal(s, i)
		return
	default:
		m.log.Error().
			Str("interaction_id", i.Interaction.ID).
			Str("guild_id", i.GuildID).
			Str("custom_id", customID).
			Msg("unhandled interaction!!!")
		return
	}
}

func (m *Verification) handleCIDVerifyEmailBtn(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	log := m.log.With().
		Str("interaction_id", i.Interaction.ID).
		Str("custom_id", CIDVerifyEmailBtn).
		Str("guild_id", i.Interaction.Member.GuildID).
		Str("guild_name", i.Interaction.GuildID).
		Str("user_id", i.Interaction.Member.GuildID).
		Str("user_name", i.Interaction.Member.User.GlobalName).
		Logger()
	log.Info().Msg("interaction request received")

	// Open up a modal!
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: CIDVerifyEmailModal,
			Title:    "Email Verification",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "email_input_field",
							Label:       "Enter your official ASU email:",
							Style:       discordgo.TextInputShort,
							Placeholder: "example@asu.edu",
							Required:    true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("error sending modal to user")
		return
	}
}

func (m *Verification) handleCIDVerifyEmailModal(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	log := m.log.With().
		Str("interaction_id", i.Interaction.ID).
		Str("custom_id", CIDVerifyEmailModal).
		Str("user_id", i.Interaction.Member.User.ID).
		Str("user_name", i.Interaction.Member.User.Username).
		Logger()
	log.Info().Msg("interaction request received")

	// Grab email from user
	recipient := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	log.Debug().Msgf("received email from user: %s", recipient)

	// // Send email to user
	// sender := "forkman@example.com"
	// subject := "D2D Email Verification"
	// htmlBody := "<h1>Test Email</h1>"
	// textBody := "this email is a test email"
	// charSet := "UTF-8"
	// sess, _ := session.NewSession(&aws.Config{
	// 	Region: aws.String("us-east-1"),
	// })
	// svc := ses.New(sess)
	//
	// // Assemble the email.
	// input := &ses.SendEmailInput{
	// 	Destination: &ses.Destination{
	// 		CcAddresses: []*string{},
	// 		ToAddresses: []*string{
	// 			aws.String(recipient),
	// 		},
	// 	},
	// 	Message: &ses.Message{
	// 		Body: &ses.Body{
	// 			Html: &ses.Content{
	// 				Charset: aws.String(charSet),
	// 				Data:    aws.String(htmlBody),
	// 			},
	// 			Text: &ses.Content{
	// 				Charset: aws.String(charSet),
	// 				Data:    aws.String(textBody),
	// 			},
	// 		},
	// 		Subject: &ses.Content{
	// 			Charset: aws.String(charSet),
	// 			Data:    aws.String(subject),
	// 		},
	// 	},
	// 	Source: aws.String(sender),
	// }
	//
	// // Attempt to send the email.
	// result, err := svc.SendEmail(input)
	//
	// // Display error messages if they occur.
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		case ses.ErrCodeMessageRejected:
	// 			m.log.Error().Err(aerr).Msg("message rejected")
	// 		case ses.ErrCodeMailFromDomainNotVerifiedException:
	// 			m.log.Error().Err(aerr).Msg("domain not verified")
	// 		case ses.ErrCodeConfigurationSetDoesNotExistException:
	// 			m.log.Error().Err(aerr).Msg("configuration set does not exist")
	// 		default:
	// 			m.log.Error().Err(aerr).Msg("unhandled error")
	// 		}
	// 	} else {
	// 		m.log.Error().Err(err).Msg("unhandled non-aws error")
	// 	}
	//
	// 	templates.ErrMessage(s, i, err)
	// 	return
	// }
	//
	// m.log.Info().Interface("result", result).Msg("message sent!")

	// Respond with embed and button
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Email Submitted",
					Description: "We have sent a code to your email (" + recipient + "). Please check your inbox and enter the code below.",
					Color:       0x00FF00, // Green color
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Enter Verification Code",
							Style:    discordgo.PrimaryButton,
							CustomID: CIDVerifyEmailCodeBtn,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("error responding to user")
		return
	}
}

func (m *Verification) handleCIDVerifyEmailCodeBtn(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	log := m.log.With().
		Str("interaction_id", i.Interaction.ID).
		Str("custom_id", CIDVerifyEmailCodeBtn).
		Str("user_id", i.Interaction.Member.User.ID).
		Str("user_name", i.Interaction.Member.User.Username).
		Logger()
	log.Info().Msg("interaction request received")

	// Open up a modal!
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: CIDVerifyEmailCodeModal,
			Title:    "Email Verification",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "email_code_input_field",
							Label:       "Enter the code sent to your email!",
							Style:       discordgo.TextInputShort,
							Placeholder: "1234567",
							Required:    true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("error sending modal to user")
		return
	}
}

func (m *Verification) handleCIDVerifyEmailCodeModal(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	log := m.log.With().
		Str("interaction_id", i.Interaction.ID).
		Str("custom_id", CIDVerifyEmailCodeModal).
		Str("user_id", i.Interaction.Member.User.ID).
		Str("user_name", i.Interaction.Member.User.Username).
		Logger()
	log.Info().Msg("interaction request received")

	code := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	log.Debug().Msgf("code received: %s", code)
	templates.Message(s, i, "Wow thanks for the code!")
}
