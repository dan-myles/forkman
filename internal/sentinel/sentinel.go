package sentinel

// TODO : do not release Start() until we receive a ready event from discord

import (
	"sync"

	"github.com/avvo-na/devil-guard/internal/config"
	"github.com/avvo-na/devil-guard/internal/sentinel/module"
	"github.com/avvo-na/devil-guard/internal/sentinel/utility"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	instance *Sentinel
	once     sync.Once
)

type Sentinel struct {
	Session       *discordgo.Session
	ModuleManager *module.ModuleManager
}

func New() *Sentinel {
	once.Do(func() {
		// Grab our token n stufff
		token := config.GetConfig().AppCfg.DiscordBotToken
		appID := config.GetConfig().AppCfg.DiscordAppID
		devGuildID := config.GetConfig().AppCfg.DiscordDevGuildID

		// Create a new Discord session using the provided bot token.
		log.Debug().Msg("Creating Discord session...")
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create Discord session")
		}

		// Set intents
		session.Identify.Intents = discordgo.IntentsAll

		// Clean commands
		// TODO: eventually only do this in dev mode
		log.Debug().Msg("Cleaning commands...")
		cmds, err := session.ApplicationCommands(appID, devGuildID)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to get commands")
		}

		for _, cmd := range cmds {
			err := session.ApplicationCommandDelete(appID, devGuildID, cmd.ID)
			if err != nil {
				log.Panic().Err(err).Msg("Failed to delete command")
			}
		}

		// Register modules
		log.Debug().Msg("Init module manager...")
		moduleManager := &module.ModuleManager{}
		moduleManager.AddModule(utility.New())
		moduleManager.RegisterModules(session)

		instance = &Sentinel{
			Session:       session,
			ModuleManager: moduleManager,
		}

		log.Info().Msg("Sentinel initialized, modules registered")
	})

	return instance
}

func (s *Sentinel) Start() {
	// Open the websocket and begin listening.
	err := s.Session.Open()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open connection to Discord")
	}
}

func (s *Sentinel) Stop() {
	// Close the websocket
	err := s.Session.Close()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to close connection to Discord")
	}
}
