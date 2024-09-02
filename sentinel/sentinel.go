package sentinel

import (
	"sync"
	"time"

	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	Session      *discordgo.Session
	SessionMutex *sync.Mutex
)

func InitSentinel() {
	// Init a new Discord session
	s, err := discordgo.New("Bot " + config.AppCfg.DiscordBotToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Discord session")
	}

	// Set intents
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	Session = s
	SessionMutex = &sync.Mutex{}
}

func Start() error {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()

	// Open connection to Discord
	err := Session.Open()
	if err != nil {
		return err
	}

	time := time.Now()
	log.Info().Time("time", time).Msg("Opened connection to Discord")
	return nil
}

func Stop() error {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()

	// Close connection to Discord
	err := Session.Close()
	if err != nil {
		return err
	}

	time := time.Now()
	log.Info().Time("time", time).Msg("Closed connection to Discord")
	return nil
}
