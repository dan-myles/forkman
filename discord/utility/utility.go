package utility

import (
	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "role",
		Description: "role management commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "all",
				Description: "adds role to all members",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "role",
						Description: "to add",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
			{
				Name:        "remove",
				Description: "removes role from specified member",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "member",
						Description: "to remove role from",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
					{
						Name:        "role",
						Description: "to remove",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
			{
				Name:        "add",
				Description: "adds role to specified member",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "member",
						Description: "to add role to",
						Type:        discordgo.ApplicationCommandOptionUser,
						Required:    true,
					},
					{
						Name:        "role",
						Description: "to add",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
		},
	},
}

var commandHandlers = map[string]func(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	l *zerolog.Logger,
){
	"role": role,
}

type UtilityModule struct {
	session *discordgo.Session
	log     *zerolog.Logger
	cfg     *config.ConfigManager
}

func New(s *discordgo.Session, l *zerolog.Logger, c *config.ConfigManager) *UtilityModule {
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
	// Grab the config
	appCfg := u.cfg.GetAppConfig()
	moduleCfg := u.cfg.GetModuleConfig()

	// If the module is disabled, skip registration
	if !*moduleCfg.Utility.Enabled {
		u.log.Info().Msg("Module is disabled, skipping registration")
		return nil
	}

	// Register all commands
	for _, command := range commands {
		u.log.Debug().
			Str("appID", appCfg.DiscordAppID).
			Str("guildID", appCfg.DiscordDevGuildID).
			Str("command", command.Name).
			Msg("Registering command")

		_, err := u.session.ApplicationCommandCreate(
			appCfg.DiscordAppID,
			appCfg.DiscordDevGuildID,
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
