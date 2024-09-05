package log

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Trace() *zerolog.Event {
	return log.Trace()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Panic() *zerolog.Event {
	return log.Panic()
}

func TraceI(i *discordgo.InteractionCreate) *zerolog.Event {
	return log.Trace().
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)
}

func DebugI(i *discordgo.InteractionCreate) *zerolog.Event {
	return log.Debug().
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)
}

func InfoI(i *discordgo.InteractionCreate) *zerolog.Event {
	return log.Info().
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)
}

func WarnI(i *discordgo.InteractionCreate) *zerolog.Event {
	return log.Warn().
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)
}

func ErrorI(i *discordgo.InteractionCreate) *zerolog.Event {
	return log.Error().
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)
}

func FatalI(i *discordgo.InteractionCreate) *zerolog.Event {
	return log.Fatal().
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)
}

func PanicI(i *discordgo.InteractionCreate) *zerolog.Event {
	return log.Panic().
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)
}
