package server

import (
	"net/http"
	"os"
	"path/filepath"

	funware "github.com/avvo-na/forkman/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) registerRoutes() http.Handler {
	// Setup router
	r := chi.NewRouter()
	r.Use(funware.Logger(s.log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Serve the Frontend :D
	workdir, _ := os.Getwd()
	fileDir := http.Dir(filepath.Join(workdir, "fork_data/static"))
	r.Get("/*", http.FileServer(fileDir).ServeHTTP)

	// If in development, redirect to the frontend vite server
	if s.cfg.GoEnv == "development" {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://localhost:5173/", http.StatusPermanentRedirect)
		})
	}

	// Auth routes
	r.Get("/auth/{provider}/callback", s.getAuthProviderCallback)
	r.Get("/auth/{provider}/logout", s.getAuthProviderLogout)
	r.Get("/auth/{provider}/login", s.getAuthProviderLogin)

	// Protected routes (API ACCESS)
	r.Route("/api", func(r chi.Router) {
		r.Use(funware.Auth())

		r.Get("/user", s.getUser)
	})

	return r
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	// Get the user from the session
	session, _ := s.store.Get(r, "forkman-user-session")
	user := session.Values["user"]

	// Return the user
	s.respond(w, r, http.StatusOK, user)
}
