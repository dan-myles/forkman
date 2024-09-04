package sentinel

import (
	"sync"

	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/sentinel/module"
	"github.com/avvo-na/devil-guard/sentinel/utility"
	"github.com/avvo-na/devil-guard/sentinel/verification"
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
	Mutex         *sync.Mutex
}

func New() *Sentinel {
	once.Do(func() {
		// Grab our token :)
		token := config.GetConfig().AppCfg.DiscordBotToken
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create Discord session")
		}

		// Set intents
		session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

		// Register modules
		moduleManager := &module.ModuleManager{}
		moduleManager.RegisterModule(&utility.UtilityModule{})
		moduleManager.RegisterModule(&verification.VerificationModule{})
		moduleManager.EnableModules(session)

		instance = &Sentinel{
			Session:       session,
			ModuleManager: moduleManager,
			Mutex:         &sync.Mutex{},
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
