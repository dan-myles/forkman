package sentinel

import (
	"fmt"

	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var Session *discordgo.Session

func Start() {
	var err error
	Session, err = discordgo.New("Bot " + config.AppCfg.DiscordBotToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Discord session")
	}
}

func Open() {
	err := Session.Open()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open Discord session")
	}
}

func Close() {
	err := Session.Close()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to close Discord session")
	}
}

func AddHandler(handler interface{}) {
	Session.AddHandler(handler)
}

func RegisterCommands(commands []*discordgo.ApplicationCommand) error {
	_, err := Session.ApplicationCommandBulkOverwrite(
		config.AppCfg.DiscordAppID,
		config.AppCfg.DiscordDevGuildID,
		commands,
	)
	if err != nil {
		return fmt.Errorf("Failed to register commands %w", err)
	}

	return nil
}
