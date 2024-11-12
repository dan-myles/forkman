package moderation

import (
	"fmt"
	"time"

	"github.com/avvo-na/forkman/internal/discord/templates"
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "mute",
		Description: "a user for a certain duration",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "user to mute",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "length",
				Description: "length of the timeout (ie. '1d', '30m', '60s')",
				Required:    true,
			},
		},
	},
	{
		Name:        "nuke",
		Description: "everything...",
	},
}

func (m *Moderation) mute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	templates.Message(s, i, "hi!")
}

func (m *Moderation) nuke(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logChannel := "1301660537867730964"
	allowedId := []string{
		"1133163357042655242", // dan
		"659054617589448747",  // liz
	}

	allowedRoleIds := []string{
		"1291813576415121459", // keeper
		"1213199035968528434", // asu staff
		"1230915317648064542", // admin rep
		"1187146366359715861", // moderator
		"1212825472769986615", // helper
		"1187156709597270157", // gold guide
		"1278880165149016168", // dev
		"1238534764365873314", // greeter
		"1187144436145197066", // admin
	}

	if i.Member.User.ID != allowedId[0] && i.Member.User.ID != allowedId[1] {
		s.ChannelMessageSend(logChannel, "You are not allowed to use this command :()")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "... 💣 ...",
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
		m.log.Error().Err(err).Msg("Failed to request guild member chunks")
		return
	}

	// Wait for all chunks to be received
	select {
	case <-done:
		s.ChannelMessageSend(logChannel, "All member chunks received! Processing oppenheimer maneuver...")
		for _, mem := range memberMap {
			if len(memberMap)%100 == 0 {
				s.ChannelMessageSend(logChannel, fmt.Sprintf("%d members left to kick", len(memberMap)))
			}

			if mem.User.Bot {
				continue
			}

			// Loop through member make sure they dont have keeper
			var isKeeper bool
			for _, r := range mem.Roles {
				for _, roleId := range allowedRoleIds {
					if r == roleId {
						isKeeper = true
						break
					}
				}
			}

			if isKeeper {
				m.log.Info().Msgf("Ignoring whitelisted member %s", mem.User.GlobalName)
				continue
			}

			// if not we kick them
			err := s.GuildMemberDelete(i.GuildID, mem.User.ID)
			if err != nil {
				m.log.Error().Err(err).Msg("Failed to kick member")
				continue
			}

			m.log.Info().Msgf("Kicked member %s", mem.User.GlobalName)
			delete(memberMap, mem.User.ID)
		}

		s.ChannelMessageSend(i.ChannelID, "Done!")
		remove()
	case <-time.After(10 * time.Second): // Adjust the timeout as necessary
		s.ChannelMessageSend(i.ChannelID, "Timeout waiting for chunks, please try again!")
		m.log.Error().Msg("Timeout waiting for chunks")
		remove()
	}
}
