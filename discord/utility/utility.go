package utility

import (
	"github.com/avvo-na/forkman/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type UtilityModule struct {
	session *discordgo.Session
	log     *zerolog.Logger
	cfg     *config.Config
}

func New(s *discordgo.Session, l *zerolog.Logger, c *config.Config) *UtilityModule {
	subLogger := l.With().Str("module", "utility").Logger()

	return &UtilityModule{
		session: s,
		log:     &subLogger,
		cfg:     c,
	}
}

func (u *UtilityModule) Name() string {
	return "utility"
}

func (u *UtilityModule) Load() error {
	// Register all commands
	for _, command := range commands {
		u.log.Debug().
			Str("appID", u.cfg.DiscordAppID).
			Str("guildID", u.cfg.DiscordDevGuildID).
			Str("command", command.Name).
			Msg("Registering command")

		_, err := u.session.ApplicationCommandCreate(
			u.cfg.DiscordAppID,
			u.cfg.DiscordDevGuildID,
			command,
		)
		if err != nil {
			return err
		}
	}

	u.log.Debug().
		Interface("commands", commands).
		Msg("Registering command handlers...")
	u.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, ok := commandHandlers[i.ApplicationCommandData().Name]
		if !ok {
			return
		}
		handler(s, i, u.log)
	})

	u.log.Info().Msg("Module loaded successfully")
	return nil
}

func (u *UtilityModule) Enable() error {
	return nil
}

func (u *UtilityModule) Disable() error {
	return nil
}
