package sentinel

import (
	"context"

	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/utility"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/rs/zerolog/log"
)

type Sentinel struct {
	Client bot.Client
}

func New() *Sentinel {
	c, err := disgo.New(config.AppCfg.DiscordBotToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentDirectMessages,
				gateway.IntentGuildMessageTyping,
				gateway.IntentDirectMessageTyping,
				gateway.IntentMessageContent,
			),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("loading..."),
				gateway.WithOnlineStatus(discord.OnlineStatusDND),
			),
		),
	)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create disgo client")
	}

	log.Info().Msg("Discord client created successfully")
	return &Sentinel{
		Client: c,
	}
}

func (s *Sentinel) Start() {
	log.Info().Msg("Starting bot...")

	err := s.Client.OpenGateway(context.TODO())
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open Discord session")
	}
}

func (s *Sentinel) Stop() {
	log.Info().Msg("Stopping bot...")
	s.Client.Gateway().Close(context.TODO())
	s.Client.Close(context.TODO())
}

func (s *Sentinel) RegisterPlugins() {
	if config.PluginCfg.Utility == "enabled" {
		_, err := s.Client.Rest().SetGuildCommands(
			snowflake.MustParse(config.AppCfg.DiscordAppID),
			snowflake.MustParse(config.AppCfg.DiscordDevGuildID),
			utility.Commands,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to register utility commands!")
		}
	}
}
