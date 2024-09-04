package utility

import (
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
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "GIVING ROLE TO ALL MEMBERS!",
		},
	})

	// Get the role
	options := i.ApplicationCommandData().Options
	role := options[0].Options[0].RoleValue(s, i.GuildID)

	// Set listener for guild members chunk
	s.AddHandler(func(s *discordgo.Session, i *discordgo.GuildMembersChunk) {
		log.Debug().Msg("Received guild members chunk!")
		log.Debug().Msgf("Members: %d", len(i.Members))

		// Loop through all members
		for _, member := range i.Members {
			// Add the role to the member

			var found bool
			for _, r := range member.Roles {
				if r == role.ID {
					log.Debug().Str("member", member.User.ID).Str("role", role.ID).Msg("Role already added to member")
					found = true
				}
			}
			if found {
				continue
			}

			err := s.GuildMemberRoleAdd(i.GuildID, member.User.ID, role.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to add role to member")
				continue
			}

			log.Debug().Str("member", member.User.ID).Str("role", role.ID).Msg("Role added to member")
		}

		log.Info().Msg("Finished chunk processing")
	})

	// Request chunks to be sent
	err := s.RequestGuildMembers(i.GuildID, "", 0, "", false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to request guild members")
		return
	}
}
