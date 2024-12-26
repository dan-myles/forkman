package middleware

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/avvo-na/forkman/internal/discord"
	e "github.com/avvo-na/forkman/internal/server/common/err"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog"
)

type Middleware struct {
	logger  *zerolog.Logger
	discord *discord.Discord
}

func New(logger *zerolog.Logger, discord *discord.Discord) *Middleware {
	return &Middleware{
		logger:  logger,
		discord: discord,
	}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := m.logger.With().Logger()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			finish := time.Now()
			log.Info().
				Str("type", "access").
				Timestamp().
				Fields(map[string]interface{}{
					"request_id": middleware.GetReqID(r.Context()),
					"remote_ip":  r.RemoteAddr,
					"url":        r.URL.Path,
					"proto":      r.Proto,
					"method":     r.Method,
					"user_agent": r.Header.Get("User-Agent"),
					"status":     ww.Status(),
					"latency_ms": float64(finish.Sub(start).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get("Content-Length"),
					"bytes_out":  ww.BytesWritten(),
				}).
				Msg("served request")
		}()

		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}

func (m *Middleware) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log := m.logger.With().Logger()
				log.Error().
					Str("type", "error").
					Timestamp().
					Interface("recover_info", rec).
					Bytes("debug_stack", debug.Stack()).
					Msg("CRITICAL: recovered from panic")
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := gothic.Store.Get(r, "forkman-user-session")
		if session.Values["user"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) AuthProvider(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		if provider == "" {
			e.ServerError(w, e.ErrAuthProviderNotFound)
			return
		}

		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

func (m *Middleware) ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// IDK how i feel about injecting guild snowflake into context :/
func (m *Middleware) GuildSnowflake(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gs := chi.URLParam(r, "guildSnowflake")
		if gs == "" {
			e.BadRequest(w, e.ErrGuildNotFound)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "guildSnowflake", gs))
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) HasPermissionGuildDashboard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  gs := r.Context().Value("guildSnowflake").(string)

		session, _ := gothic.Store.Get(r, "forkman-user-session")
		user := session.Values["user"].(goth.User)

		guilds, err := m.discord.GetUserAdminServers(user.UserID)
		if err != nil {
			e.ServerError(w, err)
			return
		}

		for _, guild := range guilds {
			if guild.ID == gs {
				next.ServeHTTP(w, r)
				return
			}
		}

		e.ServerError(w, e.ErrUnauthorizedGuild)
	})
}
