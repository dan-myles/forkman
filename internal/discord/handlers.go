package discord

import (
	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	"github.com/avvo-na/forkman/internal/discord/moderation"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func onReadyNotify(l *zerolog.Logger) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		l.Info().Msgf("Session has connected to Discord as %s", r.User.String())
	}
}

// This fires when we have access to a guild or when
// a new guild has access to our bot. It will fire
// for every guild our bot has access to *every launch*
func onGuildCreateGuildUpdate(
	db *gorm.DB,
	l *zerolog.Logger,
	cfg *config.SentinelConfig,
	mm map[string]*moderation.Moderation,
) func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(s *discordgo.Session, g *discordgo.GuildCreate) {
		log := l.With().Str("event", "onGuildCreate").
			Str("guild_snowflake", g.Guild.ID).
			Str("guild_name", g.Guild.Name).
			Logger()

			// Create guild repo
		repo := database.NewGuildRepository(db)

		// Read or create guild
		_, err := repo.ReadGuild(g.Guild.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error().Err(err).Msg("critical error reading guild")
			return
		}

		if err == gorm.ErrRecordNotFound {
			if _, err := repo.CreateGuild(g.Guild); err != nil {
				log.Error().Err(err).Msg("critical error creating guild")
			}
			log.Info().Msg("guild creation complete")
		}

		// Update guild
		_, err = repo.UpdateGuild(g.Guild)
		if err != nil {
			log.Error().Err(err).Msg("critical error updating guild")
			return
		}

		// Init & store module
		mod := moderation.New(g.Guild.Name, g.Guild.ID, cfg.DiscordAppID, s, db, l)
		if err := mod.Load(); err != nil {
			log.Error().Err(err).Msg("critical error adding moderation module")
			return
		}
		mm[g.Guild.ID] = mod

		// Finished!
		log.Info().Msg("guild instantiation complete")
	}
}
