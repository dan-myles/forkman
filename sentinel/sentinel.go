package sentinel

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/avvo-na/devil-guard/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func Init() {
	dg, err := discordgo.New("Bot " + utils.ConfigData.DiscordBotToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Discord session")
	}

	dg.AddHandler(MessageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open Discord session")
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
