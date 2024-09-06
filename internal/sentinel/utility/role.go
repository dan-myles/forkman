package utility

import (
	"fmt"
	"time"

	"github.com/avvo-na/devil-guard/common/log"
	"github.com/bwmarrin/discordgo"
)

// Top level handler for all role commands
func role(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "all":
		roleAll(s, i)
	case "remove":
		roleRemove(s, i)
	}
}

func roleRemove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.InfoI(i).Msg("Interaction request received")
	options := i.ApplicationCommandData().Options
	role := options[0].Options[0].RoleValue(s, i.GuildID)
	if role == nil {
		log.ErrorI(i).Msg("Role not found")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(discordgo.MessageFlagsEphemeral),
		Data: &discordgo.InteractionResponseData{
			Content: "Removing role from member",
		},
	})
}

// TODO: Add a modal for confirmation before giving role to all members
func roleAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Log initial request
	log.InfoI(i).Msg("Interaction request received")

	options := i.ApplicationCommandData().Options
	role := options[0].Options[0].RoleValue(s, i.GuildID)
	if role == nil {
		log.ErrorI(i).Msg("Role not found")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Giving role to all members...",
		},
	})

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
	err := s.RequestGuildMembers(i.GuildID, "", 0, "", false)
	if err != nil {
		log.ErrorI(i).Err(err).Msg("Failed to request guild member chunks")
		return
	}

	// Wait for all chunks to be received
	select {
	case <-done:
		// Give roles to all members
		for _, m := range memberMap {
			if m.User.Bot {
				log.DebugI(i).Str("member_id", m.User.ID).Msg("Skipping bot")
				continue
			}

			// Lets see if the member already has the role
			var hasRole bool
			for _, r := range m.Roles {
				if role.ID == r {
					log.DebugI(i).
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
				log.ErrorI(i).
					Str("member_id", m.User.ID).
					Err(err).
					Msg("Failed to give role to member")
				continue
			}

			// remove the member from the map
			delete(memberMap, m.User.ID)

			log.DebugI(i).
				Str("member_id", m.User.ID).
				Msg("Role given to member")

			// every 100 members, log the progress
			if len(memberMap)%100 == 0 {
				s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("%d members left to process", len(memberMap)))
			}
		}

		s.ChannelMessageSend(i.ChannelID, "Role given to all members!!")

		log.InfoI(i).Msg("Interaction request completed")
		remove()
	case <-time.After(10 * time.Second): // Adjust the timeout as necessary
		log.ErrorI(i).Msg("Timeout waiting for chunks")
		remove()
	}
}
