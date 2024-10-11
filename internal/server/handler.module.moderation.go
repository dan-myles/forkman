package server

import (
	"errors"
	"net/http"

	"github.com/avvo-na/forkman/internal/discord/moderation"
	e "github.com/avvo-na/forkman/internal/server/common/err"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) disableModerationModule(w http.ResponseWriter, r *http.Request) {
	gs := r.Context().Value("guildSnowflake").(string)
	log := s.log.With().
		Str("request_id", middleware.GetReqID(r.Context())).
		Str("guild_snowflake", gs).
		Logger()

	mod, err := s.discord.GetModerationModule(gs)
	if err != nil {
		e.ServerError(w, err)
		return
	}

	err = mod.Disable()
	if err != nil {
		if errors.Is(err, moderation.ErrModuleAlreadyDisabled) {
			w.Write([]byte(`{ "message": "Module already disabled!" }`))
			return
		} else {
			log.Error().Err(err).Msg("unkown module disabling error")
			e.ServerError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "Successfully disabled Moderation module." }`))
}

func (s *Server) enableModerationModule(w http.ResponseWriter, r *http.Request) {
	gs := r.Context().Value("guildSnowflake").(string)
	log := s.log.With().
		Str("request_id", middleware.GetReqID(r.Context())).
		Str("guild_snowflake", gs).
		Logger()

	mod, err := s.discord.GetModerationModule(gs)
	if err != nil {
		e.ServerError(w, err)
		return
	}

	err = mod.Enable()
	if err != nil {
		if errors.Is(err, moderation.ErrModuleAlreadyEnabled) {
			w.Write([]byte(`{ "message": "Module already enabled!" }`))
			return
		} else {
			log.Error().Err(err).Msg("unkown module disabling error")
			e.ServerError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "Successfully enabled Moderation module." }`))
}
