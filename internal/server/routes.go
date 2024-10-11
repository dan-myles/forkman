package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/avvo-na/forkman/internal/server/middleware"
	"github.com/go-chi/chi/v5"
)

func (s *Server) registerRoutes() http.Handler {
	// Init middleware
	middleware := middleware.New(s.log)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	// Serve the Frontend :D
	// TODO: embed the frontend into the binary
	workdir, _ := os.Getwd()
	fileDir := http.Dir(filepath.Join(workdir, "fork_data/static"))
	r.Get("/*", http.FileServer(fileDir).ServeHTTP)

	// If in development, redirect to the frontend vite server
	if s.cfg.GoEnv == "development" {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://localhost:5173/", http.StatusPermanentRedirect)
		})
	}

	// Auth Routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.AuthProvider)
		r.Get("/{provider}/login", s.authLogin)
		r.Get("/{provider}/logout", s.authLogout)
		r.Get("/{provider}/callback", s.authCallback)
	})

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJSON)

		r.Route("/{guildSnowflake}", func(r chi.Router) {
			r.Use(middleware.GuildSnowflake)
			r.Get("/module/moderation/disable", s.disableModerationModule)
			r.Get("/module/moderation/enable", s.enableModerationModule)
			r.Get("/module/verification/panel/send/{channelId}", s.sendVerificationPanel)
		})
	})

	return r
}
