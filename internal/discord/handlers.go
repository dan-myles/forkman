package discord

import (
	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	"github.com/avvo-na/forkman/internal/discord/moderation"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func onReadyInitModules(
	l *zerolog.Logger,
) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		l.Info().Msg("Ready EVENT")
	}
}

func onGuildCreateGuildUpdate(
	db *gorm.DB,
	l *zerolog.Logger,
	cfg *config.SentinelConfig,
	moderationModules map[string]*moderation.ModerationModule,
) func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(s *discordgo.Session, g *discordgo.GuildCreate) {
		log := l.With().Str("event", "onGuildCreate").
			Str("guild_snowflake", g.Guild.ID).
			Str("guild_name", g.Guild.Name).
			Logger()

		// Create a new guild if it doesn't exist
		guild := database.Guild{
			Snowflake: g.Guild.ID,
		}

		if err := db.FirstOrCreate(&guild, guild).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create guild")
		}

		// Instantiate and store modules
		log.Info().Msg("Creating moderation module")
		moderationModule := moderation.New(g.Guild.ID, s, l, db, cfg)
		moderationModule.Sync()
		moderationModules[g.Guild.ID] = moderationModule

		log.Info().Msg("Guild updated")
	}
}
