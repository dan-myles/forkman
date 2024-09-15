package server

import (
	"fmt"
	"net/http"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/discord"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	db      *gorm.DB
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
	db *gorm.DB,
) *http.Server {

	s := &Server{
		db:      db,
		log:     l,
		valid:   v,
		discord: d,
		cfg:     cfg,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      s.registerRoutes(),
		IdleTimeout:  cfg.ServerTimeoutIdle,
		ReadTimeout:  cfg.ServerTimeoutRead,
		WriteTimeout: cfg.ServerTimeoutWrite,
	}

	return server
}
