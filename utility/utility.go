package utility

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	instance *UtilityModule = &UtilityModule{}
	once     *sync.Once

	// Slash commands
	commands []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Ping the bot",
		},
	}

	// Handlers
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": ping,
	}
)

type UtilityModule struct{}

func GetModule() *UtilityModule {
	return instance
}

func (u *UtilityModule) Name() string {
	return "utility"
}

func (u *UtilityModule) Enable() error {
	log.Info().Msg("Enabled utility module")
	return nil
}

func (u *UtilityModule) Disable() error {
	log.Info().Msg("Disabled utility module")
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

//
// import (
// 	"fmt"
//
// 	"github.com/avvo-na/devil-guard/config"
// 	"github.com/bwmarrin/discordgo"
// 	"github.com/rs/zerolog/log"
// )
//
// // This is the list of commands that the bot will register
// var commands = []*discordgo.ApplicationCommand{
// 	{
// 		Name:        "ping",
// 		Description: "Ping the bot",
// 	},
// }
//
// var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
// 	"ping": ping,
// }
//
// func EnableModule(s *discordgo.Session) error {
// 	// Write the module config
// 	config.Mutex.Lock()
// 	config.ModuleCfg.Utility.Enabled = true
// 	err := config.WriteModuleConfig()
// 	if err != nil {
// 		return fmt.Errorf("failed to write module config: %w", err)
// 	}
// 	config.Mutex.Unlock()
//
// 	// Register the commands
// 	for _, v := range commands {
// 		_, err := s.ApplicationCommandCreate(
// 			config.AppCfg.DiscordAppID,
// 			config.AppCfg.DiscordDevGuildID,
// 			v,
// 		)
// 		// TODO: add better handling if one command fails to register
// 		if err != nil {
// 			log.Error().Stack().Err(err).Msg("Failed to register command")
// 		}
// 	}
//
// 	// This is a map of command names to their handlers. When a command is
// 	// received, the bot will check if the command name is in this map. If it
// 	// is, the bot will call the handler function with the session and the
// 	// interaction.
// 	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
// 			handler(s, i)
// 		}
// 	})
//
// 	log.Info().Msg("Enabled utility module")
// 	return nil
// }
//
// // TODO: will be called from a rest API
// func DisableModule(s *discordgo.Session) error {
// 	// Write the module config
// 	config.Mutex.Lock()
// 	config.ModuleCfg.Utility.Enabled = false
// 	err := config.WriteModuleConfig()
// 	if err != nil {
// 		return fmt.Errorf("failed to write module config: %w", err)
// 	}
// 	config.Mutex.Unlock()
//
// 	// NOTE:
// 	// Grab all regeistered and cross check with the commands
// 	// n^2 so its a bit slow :() idk if it matters
// 	registeredCommands, err := s.ApplicationCommands(config.AppCfg.DiscordAppID, config.AppCfg.DiscordDevGuildID)
// 	for _, v := range registeredCommands {
// 		for _, c := range commands {
// 			if c.Name == v.Name {
// 				err := s.ApplicationCommandDelete(
// 					config.AppCfg.DiscordAppID,
// 					config.AppCfg.DiscordDevGuildID,
// 					v.ID,
// 				)
// 				if err != nil {
// 					log.Error().Stack().Err(err).Msg("Failed to delete command")
// 				}
// 			}
// 		}
// 	}
//
// 	log.Info().Msg("Disabled utility module")
// 	return nil
// }
//
// func ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: "Pong!",
// 		},
// 	})
// 	log.Info().Interface("command", i.ApplicationCommandData()).Msg("Responded to interaction request")
// }
