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

// NOTE: Should only be used in the context of an interaction
func TraceI(i *discordgo.InteractionCreate) *zerolog.Event {
	return addInteractionFields(log.Trace(), i)
}

func DebugI(i *discordgo.InteractionCreate) *zerolog.Event {
	return addInteractionFields(log.Debug(), i)
}

func InfoI(i *discordgo.InteractionCreate) *zerolog.Event {
	return addInteractionFields(log.Info(), i)
}

func WarnI(i *discordgo.InteractionCreate) *zerolog.Event {
	return addInteractionFields(log.Warn(), i)
}

func ErrorI(i *discordgo.InteractionCreate) *zerolog.Event {
	return addInteractionFields(log.Error(), i)
}

func FatalI(i *discordgo.InteractionCreate) *zerolog.Event {
	return addInteractionFields(log.Fatal(), i)
}

func PanicI(i *discordgo.InteractionCreate) *zerolog.Event {
	return addInteractionFields(log.Panic(), i)
}

func addInteractionFields(e *zerolog.Event, i *discordgo.InteractionCreate) *zerolog.Event {
	handle := e.
		Str("command", i.ApplicationCommandData().Name).
		Str("interaction_id", i.Interaction.ID).
		Str("guild_id", i.GuildID).
		Str("user_id", i.Member.User.ID)

	if len(i.ApplicationCommandData().Options) > 0 {
		e.Str("sub_command", i.ApplicationCommandData().Options[0].Name)
	}

	return handle
}
