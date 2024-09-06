package discord

import (
	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/discord/utility"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type Module interface {
	// Returns the name of the module
	Name() string

	// Enables the module, handles any setup and registration
	// of commands, writes config to file.
	Enable() error

	// Disables the module, handles any cleanup and deregistration
	// of commands, writes config to file.
	Disable() error

	// Loads the module, handles any setup and registration of
	// commands, *reads* config from file. To only be called once
	Load() error
}

type Discord struct {
	session *discordgo.Session
	log     *zerolog.Logger
	cfg     *config.ConfigManager
}

// TODO: Probably dont panic :P or maybe we should? idk im tired 💀
func New(cfg *config.ConfigManager, log *zerolog.Logger) *Discord {
	s, err := discordgo.New("Bot " + cfg.GetAppConfig().DiscordBotToken)
	if err != nil {
		panic(err)
	}

	// Settings
	s.Identify.Intents = discordgo.IntentsAll // What do we need permission for?
	s.SyncEvents = false                      // Launch goroutines for handlers

	logger := log.With().Str("package", "discord").Logger()
	return &Discord{
		session: s,
		log:     &logger,
		cfg:     cfg,
	}
}

// Only to be called once, i mean it 😎
func (d *Discord) Setup() {
	// Register all modules
	modules := []Module{
		utility.New(d.session, d.log, d.cfg),
	}

	// Load em up 🤠
	for _, module := range modules {
		err := module.Load()
		if err != nil {
			panic(err)
		}
	}

	d.log.Info().Msg("All modules loaded")
}

func (d *Discord) Open() error {
	err := d.session.Open()
	if err != nil {
		return err
	}

	d.log.Info().Msg("Opened a connection to Discord")
	return nil
}

func (d *Discord) Close() error {
	err := d.session.Close()
	if err != nil {
		return err
	}

	d.log.Info().Msg("Closed the connection to Discord")
	return nil
}
