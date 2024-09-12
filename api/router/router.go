package router

import (
	"github.com/avvo-na/forkman/api/handler"
	"github.com/avvo-na/forkman/api/router/middleware"
	"github.com/avvo-na/forkman/discord"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

func New(l *zerolog.Logger, v *validator.Validate, d *discord.Discord) *chi.Mux {
	// Initialize the handler
	h := handler.New(l, v, d)

	// Setup middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger(l))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)

	// Health check :P
	r.Get("/api/health", h.Health)

	return r
}
