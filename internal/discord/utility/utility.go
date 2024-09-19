package utility

import (
	"github.com/avvo-na/forkman/common/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UtilityModule struct {
	session  *discordgo.Session
	db       *gorm.DB
	log      *zerolog.Logger
	cfg      *config.SentinelConfig
	commands []*discordgo.ApplicationCommand // Store commands for cleanup
	unhandle *func()                         // Function to deregister the handler
}

func New(
	s *discordgo.Session,
	l *zerolog.Logger,
	c *config.SentinelConfig,
	db *gorm.DB,
) *UtilityModule {
	subLogger := l.With().Str("module", "utility").Logger()

	return &UtilityModule{
		session: s,
		db:      db,
		log:     &subLogger,
		cfg:     c,
	}
}

func (u *UtilityModule) RegisterCommands() {
	// Register all commands
	for _, command := range commands {
		u.log.Debug().
			Str("appID", u.cfg.DiscordAppID).
			Str("guildID", u.cfg.DiscordDevGuildID).
			Str("command", command.Name).
			Msg("Registering command")

		c, err := u.session.ApplicationCommandCreate(
			u.cfg.DiscordAppID,
			u.cfg.DiscordDevGuildID,
			command,
		)
		if err != nil {
			panic(err)
		}

		u.commands = append(u.commands, c)
	}
}

func (u *UtilityModule) RegisterHandlers() {
	fn := u.session.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.ApplicationCommandData().Name {
			case "role":
				u.role(s, i)
			}
		},
	)

	u.unhandle = &fn
}
