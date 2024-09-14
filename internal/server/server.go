package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/database"
	"github.com/avvo-na/forkman/internal/discord"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Server struct {
	db      database.Service
	log     *zerolog.Logger
	valid   *validator.Validate
	discord *discord.Discord
	cfg     *config.SentinelConfig
}

func New(
	cfg *config.SentinelConfig,
	l *zerolog.Logger,
	v *validator.Validate,
	d *discord.Discord,
) *http.Server {

	s := &Server{
		db:      database.New(),
		log:     l,
		valid:   v,
		discord: d,
		cfg:     cfg,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      s.registerRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
