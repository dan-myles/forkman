package common

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

const (
	ColorError   = 0xFF0000
	ColorSuccess = 0x5c5fea
)

func RespondErr(respondFunc events.InteractionResponderFunc, err error) error {
	return RespondErrMessage(respondFunc, fmt.Sprintf("Error while executing: %s", err))
}

func RespondErrMessage(respondFunc events.InteractionResponderFunc, message string) error {
	return respondFunc(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetDescription(message).
			SetColor(ColorError).
			Build(),
		).
		SetEphemeral(true).
		Build(),
	)
}

func RespondErrMessagef(respondFunc events.InteractionResponderFunc, message string, a ...any) error {
	return RespondErrMessage(respondFunc, fmt.Sprintf(message, a...))
}

func RespondMessageErr(respondFunc events.InteractionResponderFunc, message string, err error) error {
	return respondFunc(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetDescriptionf(message, err).
			SetColor(ColorError).
			Build(),
		).
		SetEphemeral(true).
		Build(),
	)
}

func Respond(respondFunc events.InteractionResponderFunc, message string) error {
	return respondFunc(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetDescription(message).
			SetColor(ColorSuccess).
			Build(),
		).
		Build(),
	)
}

func Respondf(respondFunc events.InteractionResponderFunc, message string, a ...any) error {
	return Respond(respondFunc, fmt.Sprintf(message, a...))
}
