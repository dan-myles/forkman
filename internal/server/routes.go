package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/avvo-na/forkman/internal/server/middleware"
	"github.com/go-chi/chi/v5"
)

//go:embed dist/*
var distFS embed.FS

func (s *Server) registerRoutes() http.Handler {
	// Init middleware
	middleware := middleware.New(s.log, s.discord)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	switch s.cfg.GoEnv {
	case "development":
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://localhost:5173/", http.StatusTemporaryRedirect)
		})
	case "production":
		build, err := fs.Sub(distFS, "dist")
		if err != nil {
			panic(err)
		}
		webFS := http.FS(build)
		r.Get("/*", http.FileServer(webFS).ServeHTTP)
	default:
		panic("unknown environment")
	}

	// Health Check
	r.Get("/health", s.healthCheck)
	r.Get("/uptime", s.uptime)

	// Auth Routes
	r.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}/login", s.authLogin)
		r.Get("/{provider}/logout", s.authLogout)
		r.Get("/{provider}/callback", s.authCallback)
	})

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJSON)
		r.Use(middleware.Auth)

		// User API
		r.Route("/user", func(r chi.Router) {
			r.Get("/session", s.sessionInfo)
			r.Get("/servers", s.adminInfo)
		})

		// Snowflake API
		r.Route("/{guildSnowflake}", func(r chi.Router) {
			r.Use(middleware.GuildSnowflake)
			r.Use(middleware.HasPermissionGuildDashboard)

			// Moderation API
			r.Post("/module/moderation/enable", s.enableModerationModule)
			r.Post("/module/moderation/disable", s.disableModerationModule)
			r.Get("/module/moderation/status", s.statusModerationModule)

			// Verification API
			r.Post("/module/verification/enable", s.enableVerificationModule)
			r.Post("/module/verification/disable", s.disableVerificationModule)
			r.Post("/module/verification/panel/send/{channelId}", s.sendVerificationPanel)
			r.Get("/module/verification/status", s.statusVerificationModule)

			// QNA API
			r.Post("/module/qna/enable", s.enableQNAModule)
			r.Post("/module/qna/disable", s.disableQNAModule)
			r.Get("/module/qna/status", s.statusQNAModule)
		})
	})

	return r
}
