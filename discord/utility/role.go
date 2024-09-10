package utility

import (
	"fmt"
	"time"

	"github.com/avvo-na/devil-guard/discord/templates"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// Top level handler for all role commands
// TODO: add a modal for role all before we let it be used
func role(s *discordgo.Session, i *discordgo.InteractionCreate, l *zerolog.Logger) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "all":
		// roleAll(s, i)
	case "remove":
		roleRemove(s, i, l)
	case "add":
		roleAdd(s, i, l)
	}
}

func roleAdd(s *discordgo.Session, i *discordgo.InteractionCreate, l *zerolog.Logger) {
	log := l.With().
		Str("command", "role add").
		Str("interaction_id", i.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID).
		Str("channel_id", i.ChannelID).
		Logger()

	// Grab options
	options := i.ApplicationCommandData().Options

	// Grab role & member
	member := i.Member
	role := options[0].Options[1].RoleValue(s, i.GuildID)
	if role == nil {
		templates.ErrMessageEphemeral(s, i, fmt.Errorf("Role not found"))
		log.Error().Err(fmt.Errorf("Role not found")).Msg("Interaction request failed")
		return
	}

	// Check the if the member already has the role
	for _, r := range member.Roles {
		if role.ID == r {
			templates.ErrMessageEphemeral(s, i, fmt.Errorf("User already has the role"))
			log.Info().Msg("Interaction request completed")
			return
		}
	}

	// Add role to member
	err := s.GuildMemberRoleAdd(i.GuildID, member.User.ID, role.ID)
	if err != nil {
		templates.ErrMessageEphemeral(s, i, err)
		log.Error().Err(err).Msg("Interaction request failed")
		return
	}

	templates.MessageEphemeral(s, i, "Added role to user!")
	log.Info().Msg("Interaction request completed")
}

func roleRemove(s *discordgo.Session, i *discordgo.InteractionCreate, l *zerolog.Logger) {
	log := l.With().
		Str("command", "role add").
		Str("interaction_id", i.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID).
		Str("channel_id", i.ChannelID).
		Logger()

	log.Info().Msg("Interaction request received")
	options := i.ApplicationCommandData().Options

	// Grab user
	user := options[0].Options[0].UserValue(s)
	if user == nil {
		templates.ErrMessageEphemeral(s, i, fmt.Errorf("User not found"))
		log.Error().Err(fmt.Errorf("User not found")).Msg("Interaction request failed")
		return
	}

	// Grab role
	role := options[0].Options[1].RoleValue(s, i.GuildID)
	if role == nil {
		templates.ErrMessageEphemeral(s, i, fmt.Errorf("Role not found"))
		log.Error().Err(fmt.Errorf("Role not found")).Msg("Interaction request failed")
		return
	}

	// Get member value
	// Member and User are different in DiscordGo
	member, err := s.GuildMember(i.GuildID, user.ID)
	if err != nil {
		templates.ErrMessageEphemeral(s, i, err)
		log.Error().Err(err).Msg("Interaction request failed")
		return
	}

	// Check member has the role to even remove
	hasRole := false
	for _, r := range member.Roles {
		if role.ID == r {
			log.Debug().Msg("Found role to remove in member")
			hasRole = true
			break
		}
	}

	// If the member does not have the role, return an error
	if !hasRole {
		templates.ErrMessageEphemeral(s, i, fmt.Errorf("User does not have the role"))
		log.Info().Msg("Interaction request completed")
		return
	}

	// Remove the role from the user
	err = s.GuildMemberRoleRemove(i.GuildID, user.ID, role.ID)
	if err != nil {
		templates.ErrMessageEphemeral(s, i, err)
		log.Error().Err(err).Msg("Interaction request failed")
		return
	}

	templates.MessageEphemeral(s, i, "Removed role from user!")
	log.Info().Msg("Interaction request completed")
}

// TODO: Add a modal for confirmation before giving role to all members
func roleAll(s *discordgo.Session, i *discordgo.InteractionCreate, l *zerolog.Logger) {
	log := l.With().
		Str("command", "role add").
		Str("interaction_id", i.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID).
		Str("channel_id", i.ChannelID).
		Logger()

	log.Info().Msg("Interaction request received")

	options := i.ApplicationCommandData().Options
	role := options[0].Options[0].RoleValue(s, i.GuildID)
	if role == nil {
		log.Error().Msg("Role not found")
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Giving role to all members...",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to respond to interaction")
		return
	}

	// Store members & channel for completion
	memberMap := make(map[string]*discordgo.Member)
	done := make(chan struct{}, 1)

	// Define the handler function
	remove := s.AddHandler(func(s *discordgo.Session, i *discordgo.GuildMembersChunk) {
		for _, m := range i.Members {
			memberMap[m.User.ID] = m
		}

		if i.ChunkIndex+1 == i.ChunkCount {
			done <- struct{}{} // Signal that we received all chunks
		}
	})

	// Request chunks to be sent
	err = s.RequestGuildMembers(i.GuildID, "", 0, "", false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to request guild member chunks")
		return
	}

	// Wait for all chunks to be received
	select {
	case <-done:
		// Give roles to all members
		for _, m := range memberMap {
			if m.User.Bot {
				log.Debug().Str("member_id", m.User.ID).Msg("Skipping bot")
				continue
			}

			// Lets see if the member already has the role
			var hasRole bool
			for _, r := range m.Roles {
				if role.ID == r {
					log.Debug().
						Str("member_id", m.User.ID).
						Msg("Member already has role, skipping...")
					hasRole = true
					break
				}
			}

			// SKIP
			if hasRole {
				continue
			}

			err := s.GuildMemberRoleAdd(i.GuildID, m.User.ID, role.ID)
			if err != nil {
				// Log the error and continue
				log.Error().
					Str("member_id", m.User.ID).
					Err(err).
					Msg("Failed to give role to member")
				continue
			}

			// remove the member from the map
			delete(memberMap, m.User.ID)

			log.Debug().
				Str("member_id", m.User.ID).
				Msg("Role given to member")

			// every 100 members, log the progress
			if len(memberMap)%100 == 0 {
				_, err := s.ChannelMessageSend(
					i.ChannelID,
					fmt.Sprintf("%d members left to process", len(memberMap)),
				)
				log.Error().Err(err).Msg("Sent progress message")
			}
		}

		_, err = s.ChannelMessageSend(i.ChannelID, "Role given to all members!!")
		if err != nil {
			log.Error().Err(err).Msg("Failed to send completion message")
		}

		log.Info().Msg("Interaction request completed")
		remove()
	case <-time.After(10 * time.Second): // Adjust the timeout as necessary
		log.Error().Msg("Timeout waiting for chunks")
		remove()
	}
}
