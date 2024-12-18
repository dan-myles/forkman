package verification

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/bwmarrin/discordgo"
)

// TODO: Some values are hard coded here just for brevity.
// Eventualy when the web dashboard gets a little farther along
// we will have a way to customize all these values. And set them
// dynamically on a per server basis.

var (
	CIDVerifyEmailBtn       = "verify_email_button"
	CIDVerifyEmailModal     = "verify_email_modal"
	CIDVerifyEmailCodeBtn   = "verify_email_code_button"
	CIDVerifyEmailCodeModal = "verify_email_code_modal"
	AllowedDomain           = "@asu.edu"
)

func (m *Verification) handleCIDVerifyEmailBtn(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	log := m.log.With().
		Str("interaction_id", i.Interaction.ID).
		Str("custom_id", CIDVerifyEmailBtn).
		Str("user_id", i.Interaction.Member.GuildID).
		Str("user_name", i.Interaction.Member.User.GlobalName).
		Logger()
	log.Info().Msg("interaction request received")

	// Open up a modal!
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: CIDVerifyEmailModal,
			Title:    "ASU Verification",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							Label:    "Enter your ASURITE ID",
							CustomID: "email_input_field",
							Style:    discordgo.TextInputShort,
							Required: true,
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
		Str("interaction_id", i.ID).
		Str("custom_id", CIDVerifyEmailModal).
		Str("user_id", i.Member.User.ID).
		Str("user_name", i.Member.User.Username).
		Logger()
	log.Info().Msg("interaction request received")

	// Grab email from user
	recipient := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	recipient += AllowedDomain
	log.Debug().Msgf("received email from user: %s", recipient)

	// Not sure why I inlined this, ehhh can organize later
	genCode := func() string {
		rand.Seed(time.Now().UnixNano())
		code := rand.Intn(900000) + 100000
		return fmt.Sprintf("%06d", code)
	}

	// Log email to DB
	code := genCode()
	e := &database.Email{
		GuildSnowflake: m.guildSnowflake,
		UserSnowflake:  i.Member.User.ID,
		Address:        recipient,
		Code:           code,
		IsVerified:     false,
	}

	_, err := m.repo.UpsertEmail(e)
	if err != nil {
		log.Error().Err(err).Msg("critical error inserting email into database")
	}

	sender := "forkman@devil2devil.asu.edu"
	subject := "Devil2Devil Verification"
	body := "Verficiation Code: " + code

	// Send the email
	err = sendEmail(context.TODO(), m.emailClient, sender, recipient, subject, body)
	if err != nil {
		log.Error().Err(err).Msg("critical error sending email")
	}
	log.Info().Msgf("sent email with id to: %s", recipient)

	// Respond with embed and button
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Submitted",
					Description: "We have sent a code to your email (" + recipient + "). Please check your inbox and enter the code below.",
					Color:       0x00FF00, // Green color
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Enter My Code",
							Style:    discordgo.PrimaryButton,
							CustomID: CIDVerifyEmailCodeBtn,
							Emoji: &discordgo.ComponentEmoji{
								Name: "👆",
							},
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
			Title:    "ASU Verification",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "email_code_input_field",
							Label:    "Enter the code sent to your email",
							Style:    discordgo.TextInputShort,
							Required: true,
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

	// Grab code
	recv := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	log.Debug().Msgf("code received: %s", recv)

	email, err := m.repo.ReadEmail(m.guildSnowflake, i.Member.User.ID)
	if err != nil {
		log.Error().Err(err).Msg("critical error reading email from database")
		return
	}

	if recv != email.Code {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Oh no!",
						Description: "We could not verify that code! Did you input it right?",
						Color:       0xFF0000, // Red color
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("error responding to user")
			return
		}

		log.Debug().
			Str("user_id", i.Member.User.ID).
			Str("user_name", i.Member.User.Username).
			Msg("user could not be verified")
		msg := "❌ User <@" + email.UserSnowflake + "> could not be verified ->" + email.Address
		s.ChannelMessageSend(os.Getenv("LOG_CHANNEL_ID"), msg)

		return
	}

	email.IsVerified = true
	_, err = m.repo.UpdateEmail(email)
	if err != nil {
		log.Error().Err(err).Msg("critical error updating verification status in database")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Success!",
					Description: "Thank you for verifying your email with us! You now have access to our community.",
					Color:       0x00FF00, // Green color
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("error responding to user")
		return
	}

	s.GuildMemberRoleRemove(m.guildSnowflake, i.Member.User.ID, os.Getenv("ROLE_TO_REMOVE"))
	s.GuildMemberRoleAdd(m.guildSnowflake, i.Member.User.ID, os.Getenv("ROLE_TO_ADD"))

	log.Debug().
		Str("user_id", i.Member.User.ID).
		Str("user_name", i.Member.User.Username).
		Msg("user succesfully verified")

	msg := "✅ User <@" + email.UserSnowflake + "> was verified -> " + email.Address
	s.ChannelMessageSend(os.Getenv("LOG_CHANNEL_ID"), msg)
}
