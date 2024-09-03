package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var Validate *validator.Validate

func Init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	log.Info().Msg("Initialized validator")
}
