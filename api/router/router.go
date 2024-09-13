package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/avvo-na/forkman/api/handler"
	"github.com/avvo-na/forkman/api/router/middleware"
	"github.com/avvo-na/forkman/discord"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

func New(l *zerolog.Logger, v *validator.Validate, d *discord.Discord, goEnv string) *chi.Mux {
	// Initialize the handler
	handler := handler.New(l, v, d)

	// Setup middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger(l))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)

	// Serve the Frontend :D
	workdir, _ := os.Getwd()
	fileDir := http.Dir(filepath.Join(workdir, "fork_data/static"))
	r.Get("/*", http.FileServer(fileDir).ServeHTTP)

	// If in development, redirect to the frontend vite server
	if goEnv == "development" {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://localhost:5173/", http.StatusPermanentRedirect)
		})
	}

	r.Get("/api/v1/health", handler.Health)

	// Serve the API Routes
	return r
}
