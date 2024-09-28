package discord

import (
	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func onReadyNotify(l *zerolog.Logger) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		l.Info().Msgf("Session has connected to Discord as %s", r.User.String())
	}
}

func onGuildCreateGuildUpdate(
	db *gorm.DB,
	log *zerolog.Logger,
	cfg *config.SentinelConfig,
	// mm map[string]*moderation.Moderation,
) func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(s *discordgo.Session, g *discordgo.GuildCreate) {
		log := log.With().Str("event", "onGuildCreate").
			Str("guild_snowflake", g.Guild.ID).
			Str("guild_name", g.Guild.Name).
			Logger()

			// Create guild repo
		repo := database.NewGuildRepository(db)

		// Read or create guild
		_, err := repo.ReadGuild(g.Guild.ID)
		if err == gorm.ErrRecordNotFound {
			if _, err := repo.CreateGuild(g.Guild); err != nil {
				log.Error().Err(err).Msg("critical error creating guild")
			}
			log.Info().Msg("Guild creation complete")
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("critical error reading guild")
			return
		}

		// Update guild
		_, err = repo.UpdateGuild(g.Guild)
		if err != nil {
			log.Error().Err(err).Msg("critical error updating guild")
			return
		}

		// Finished!
		log.Info().Msg("Guild instantiation complete")
	}
}
