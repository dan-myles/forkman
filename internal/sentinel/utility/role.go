package utility

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Top level handler for all role commands
func role(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "all":
		roleAll(s, i)
	}
}

func roleAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Gather data for logging
	iid := i.Interaction.ID
	gid := i.GuildID
	user_id := i.Member.User.ID

	// Log initial request
	log.Info().
		Str("command", "role all").
		Str("interaction_id", iid).
		Str("guild_id", gid).
		Str("user_id", user_id).
		Msg("Received interaction request")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Testing, not giving roles",
		},
	})

	// Store members & channel for completion
	var members []*discordgo.Member
	done := make(chan struct{}, 1)

	// Define the handler function
	remove := s.AddHandler(func(s *discordgo.Session, i *discordgo.GuildMembersChunk) {
		log.Debug().
			Str("command", "role all").
			Str("interaction_id", iid).
			Str("guild_id", gid).
			Str("user_id", user_id).
			Msgf(
				"Received chunk, Members: %d, chunkIndex: %d, chunkCount: %d",
				len(i.Members),
				i.ChunkIndex,
				i.ChunkCount,
			)
		members = append(members, i.Members...) // Store all members from the chunk

		if i.ChunkIndex+1 == i.ChunkCount {
			done <- struct{}{} // Signal that we received all chunks
		}
	})

	// Log request for guild member chunks
	log.Debug().
		Str("interaction_id", iid).
		Str("guild_id", gid).
		Str("user_id", user_id).
		Msg("Requesting guild member chunks...")

	// Request chunks to be sent
	err := s.RequestGuildMembers(i.GuildID, "", 0, "", false)
	if err != nil {
		log.Error().
			Err(err).
			Str("command", "role all").
			Str("interaction_id", iid).
			Str("guild_id", gid).
			Str("user_id", user_id).
			Msg("Failed to request guild member chunks")
		return
	}

	// Wait for all chunks to be received
	select {
	case <-done:
		log.Debug().
			Str("command", "role all").
			Str("interaction_id", iid).
			Str("guild_id", gid).
			Str("user_id", user_id).
			Msgf("Total members received: %d", len(members))
		remove()
	case <-time.After(10 * time.Second): // Adjust the timeout as necessary
		log.Error().
			Str("command", "role all").
			Str("interaction_id", iid).
			Str("guild_id", gid).
			Str("user_id", user_id).
			Msg("Timeout waiting for chunks")
		remove()
	}
}
