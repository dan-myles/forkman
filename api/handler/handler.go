package handler

import (
	"github.com/avvo-na/forkman/discord"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type APIHandler struct {
	log       *zerolog.Logger
	validator *validator.Validate
	disco     *discord.Discord
}

func New(log *zerolog.Logger, validator *validator.Validate, disco *discord.Discord) *APIHandler {
	l := log.With().Str("package", "api").Logger()

	return &APIHandler{
		log:       &l,
		validator: validator,
		disco:     disco,
	}
}
