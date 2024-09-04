package utility

import (
	"fmt"

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
			Content: "I'm not giving a role just testing :P",
		},
	})

	// Get the role
	// options := i.ApplicationCommandData().Options
	// role := options[0].Options[0].RoleValue(s, i.GuildID)

	// get total amount of members in the guild
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get guild count")
		return
	}
	log.Debug().Interface("members", guild.MemberCount).Msg("Guild members")

	// Set listener for guild members chunk
	s.AddHandler(func(s *discordgo.Session, i *discordgo.GuildMembersChunk) {
		log.Debug().Msg("Received guild members chunk!")
		log.Debug().Msgf("Members: %d", len(i.Members))
		// // Loop through all members
		// for _, member := range i.Members {
		// 	// Add the role to the member
		// 	err := s.GuildMemberRoleAdd(i.GuildID, member.User.ID, role.ID)
		// 	if err != nil {
		// 		log.Error().Err(err).Msg("Failed to add role to member")
		// 		continue
		// 	}
		// }

		// Send message in 1278880736736186482 channel
		// about how big the chunk is
		s.ChannelMessageSend("1278880736736186482", fmt.Sprintf("Received guild members chunk! Members: %d", len(i.Members)))
	})

	// Request chunks to be sent
	err = s.RequestGuildMembers(i.GuildID, "", 0, "", false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to request guild members")
		return
	}
}
