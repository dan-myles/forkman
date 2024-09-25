package utility

import (
	"sync"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UtilityModule struct {
	session  *discordgo.Session
	db       *gorm.DB
	log      *zerolog.Logger
	cfg      *config.SentinelConfig
	mtx      *sync.Mutex
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
		mtx:     &sync.Mutex{},
	}
}

func (u *UtilityModule) Sync() {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	// Findo our module in our DB
	module := database.Module{}

	// Try to find our module
	// If we don't have one create it, and set it to disabled
	result := u.db.First(&module, "name = ?", "utility")
	if result.RowsAffected == 0 {
		module = database.Module{
			Name:        "utility",
			Description: "Fork-tilities!",
			Enabled:     false,
		}
		u.db.Create(&module)
		return
	}

	// Lets check if we're enabled
	if !module.Enabled {
		u.log.Debug().Msg("Module is disabled, removing commands and handlers")
		u.removeCommands()
		u.removeHandlers()
	} else {
		u.log.Debug().Msg("Module is enabled, registering commands and handlers")
		u.registerCommands()
		u.registerHandlers()
	}
}

func (u *UtilityModule) registerCommands() {
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

func (u *UtilityModule) registerHandlers() {
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

func (u *UtilityModule) removeHandlers() {
	// Check if we have an unhandle function
	if u.unhandle == nil {
		u.log.Debug().Msg("No handlers to remove")
		return
	}

	// Remove the handler
	(*u.unhandle)()
}

func (u *UtilityModule) removeCommands() {
	// Get commands from discord
	remoteCommands, err := u.session.ApplicationCommands(
		u.cfg.DiscordAppID,
		u.cfg.DiscordDevGuildID,
	)
	if err != nil {
		u.log.Error().Err(err).Msg("Failed to get remote commands")
	}

	// If remoteCommands match the ones we have in `commands` remove them
	var toRemove []*discordgo.ApplicationCommand
	for _, remoteCommand := range remoteCommands {
		for _, command := range u.commands {
			if remoteCommand.Name == command.Name {
				toRemove = append(toRemove, remoteCommand)
			}
		}
	}

	// Remove the commands
	for _, command := range toRemove {
		err := u.session.ApplicationCommandDelete(
			u.cfg.DiscordAppID,
			u.cfg.DiscordDevGuildID,
			command.ID,
		)
		if err != nil {
			u.log.Error().
				Err(err).
				Msg("Failed to remove command, this may mean the command was already removed")
		}
	}
}
