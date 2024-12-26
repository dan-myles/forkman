package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/avvo-na/forkman/internal/discord/moderation"
	e "github.com/avvo-na/forkman/internal/server/common/err"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) sendVerificationPanel(w http.ResponseWriter, r *http.Request) {
	gs := r.Context().Value("guildSnowflake").(string)
	channelId := chi.URLParam(r, "channelId")
	if channelId == "" {
		e.BadRequest(w, e.ErrNoChannelIdIncluded)
		return
	}

	mod, err := s.discord.GetVerificationModule(gs)
	if err != nil {
		e.ServerError(w, err)
		return
	}

	err = mod.SendVerificationPanel(channelId)
	if err != nil {
		e.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "Successfully sent email verification panel." }`))
}

func (s *Server) disableVerificationModule(w http.ResponseWriter, r *http.Request) {
	gs := r.Context().Value("guildSnowflake").(string)
	log := s.log.With().
		Str("request_id", middleware.GetReqID(r.Context())).
		Str("guild_snowflake", gs).
		Logger()

	mod, err := s.discord.GetVerificationModule(gs)
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
	w.Write([]byte(`{ "message": "Successfully disabled Verification module." }`))
}

func (s *Server) enableVerificationModule(w http.ResponseWriter, r *http.Request) {
	gs := r.Context().Value("guildSnowflake").(string)
	log := s.log.With().
		Str("request_id", middleware.GetReqID(r.Context())).
		Str("guild_snowflake", gs).
		Logger()

	mod, err := s.discord.GetVerificationModule(gs)
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
	w.Write([]byte(`{ "message": "Successfully enabled Verification module." }`))
}

func (s *Server) statusVerificationModule(w http.ResponseWriter, r *http.Request) {
	gs := r.Context().Value("guildSnowflake").(string)
	log := s.log.With().
		Str("request_id", middleware.GetReqID(r.Context())).
		Str("guild_snowflake", gs).
		Logger()

	mod, err := s.discord.GetVerificationModule(gs)
	if err != nil {
		e.ServerError(w, err)
		return
	}

	status, err := mod.Status()
	if err != nil {
		log.Error().Err(err).Msg("unkown module status error")
		e.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
  w.Write([]byte(fmt.Sprintf(`{ "message": "Verification", "status": %t }`, status)))
}
