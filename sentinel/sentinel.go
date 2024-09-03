package sentinel

import (
	"sync"

	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/utility"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var instance *Sentinel

type Sentinel struct {
	Session       *discordgo.Session
	ModuleManager *ModuleManager
	Mutex         *sync.Mutex
}

func Init() {
	// Grab our token :)
	token := config.GetConfig().AppCfg.DiscordBotToken
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Discord session")
	}

	// Set intents
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Set session & mutex
	instance = &Sentinel{
		Session:       s,
		Mutex:         &sync.Mutex{},
		ModuleManager: &ModuleManager{},
	}

	// Register Modules
	instance.ModuleManager.RegisterModule(&utility.UtilityModule{})

	// Enable Modules
	instance.ModuleManager.EnableModules()

	log.Info().Msg("Sentinel initialized, modules registered")
}

func GetSentinel() *Sentinel {
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
