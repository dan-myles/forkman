package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/avvo-na/forkman/common/config"
	"github.com/avvo-na/forkman/internal/discord"
	"github.com/go-playground/validator/v10"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	discordProvider "github.com/markbates/goth/providers/discord"
	"github.com/rs/zerolog"
	"github.com/wader/gormstore/v2"
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
	// Setup gothic session store
	store := gormstore.New(db, []byte(cfg.ServerConfig.AuthSecret))
	store.MaxAge(int(cfg.ServerConfig.AuthExpiry.Seconds()))
	store.SessionOpts.Path = "/"
	store.SessionOpts.HttpOnly = true
	store.SessionOpts.Secure = true
	gothic.Store = store

	// Cleanup store every hour
	quit := make(chan struct{})
	go store.PeriodicCleanup(1*time.Hour, quit)

	// Setup discord provider
	goth.UseProviders(
		discordProvider.New(
			cfg.DiscordConfig.ClientID,
			cfg.DiscordConfig.ClientSecret,
			"http://localhost:5173/auth/discord/callback", // TODO: DONT HARDCODE THIS
			discordProvider.ScopeIdentify,
			discordProvider.ScopeEmail,
		))

	s := &Server{
		db:      db,
		log:     l,
		valid:   v,
		discord: d,
		cfg:     cfg,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler:      s.registerRoutes(),
		IdleTimeout:  cfg.ServerConfig.TimeoutIdle,
		ReadTimeout:  cfg.ServerConfig.TimeoutRead,
		WriteTimeout: cfg.ServerConfig.TimeoutWrite,
	}

	return server
}
