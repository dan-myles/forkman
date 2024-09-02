package utility

import (
	"fmt"

	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// This is the list of commands that the bot will register
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Ping the bot",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": ping,
}

func EnableModule(s *discordgo.Session) error {
	// Write the module config
	err := config.WriteEnableModule("utility")
	if err != nil {
		return fmt.Errorf("failed to write module config: %w", err)
	}

	// Register the commands
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(
			config.AppCfg.DiscordAppID,
			config.AppCfg.DiscordDevGuildID,
			v,
		)
		if err != nil {
			return fmt.Errorf("failed to create command: %w", err)
		}
	}

	// This is a map of command names to their handlers. When a command is
	// received, the bot will check if the command name is in this map. If it
	// is, the bot will call the handler function with the session and the
	// interaction.
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	return nil
}

// TODO: will be called from a rest API
func DisableModule() error {
	return nil
}

func ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	log.Info().Interface("command", i.ApplicationCommandData()).Msg("Responded to interaction request")
}
