package server

import (
	"encoding/json"
	"net/http"

	"github.com/avvo-na/forkman/internal/discord/moderation"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) disableModerationModule(w http.ResponseWriter, r *http.Request) {
	guildSnowflake := chi.URLParam(r, "guildSnowflake")
	log := s.log.With().
		Str("request_id", middleware.GetReqID(r.Context())).
		Str("guild_snowflake", guildSnowflake).
		Logger()

	if guildSnowflake == "" {
		log.Error().Msg("failed to ascertain guild snowflake")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).
			Encode(ErrorResponse{Error: "Please include a Guild Snowflake with your request!"})
		if err != nil {
			panic(err)
		}
		return
	}

	mod, err := s.discord.GetModerationModule(guildSnowflake)
	if err != nil {
		log.Error().Err(err).Msg("failed to find mapped snowflake")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).
			Encode(ErrorResponse{
				Error: "There was an internal server error! Please check your Guild Snowflake.",
			})
		if err != nil {
			panic(err)
		}
		return
	}

	err = mod.Disable()
	if err != nil && err == moderation.ErrModuleAlreadyDisabled {
		log.Warn().Err(err).Msg("error disabling module")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).
			Encode(Response{Status: "Module is already disabled!"})
		if err != nil {
			panic(err)
		}
		return
	}

	if err != nil && err != moderation.ErrModuleAlreadyDisabled {
		log.Error().Err(err).Msg("critical server error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).
			Encode(ErrorResponse{
				Error: "There was an error disabling the module.",
			})
		if err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).
		Encode(Response{
			Status: "Successfully disabled Moderation module!",
		})
	if err != nil {
		panic(err)
	}
}

func (s *Server) enableModerationModule(w http.ResponseWriter, r *http.Request) {
	guildSnowflake := chi.URLParam(r, "guildSnowflake")
	log := s.log.With().
		Str("request_id", middleware.GetReqID(r.Context())).
		Str("guild_snowflake", guildSnowflake).
		Logger()

	if guildSnowflake == "" {
		log.Error().Msg("failed to ascertain guild snowflake")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).
			Encode(ErrorResponse{Error: "Please include a Guild Snowflake with your request!"})
		if err != nil {
			panic(err)
		}
		return
	}

	mod, err := s.discord.GetModerationModule(guildSnowflake)
	if err != nil {
		log.Error().Err(err).Msg("failed to find mapped snowflake")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).
			Encode(ErrorResponse{
				Error: "There was an internal server error! Please check your Guild Snowflake.",
			})
		if err != nil {
			panic(err)
		}
		return
	}

	err = mod.Enable()
	if err != nil && err == moderation.ErrModuleAlreadyEnabled {
		log.Warn().Err(err).Msg("error enabling module")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).
			Encode(Response{Status: "Module is already enabled!"})
		if err != nil {
			panic(err)
		}
		return
	}

	if err != nil && err != moderation.ErrModuleAlreadyEnabled {
		log.Error().Err(err).Msg("critical server error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).
			Encode(ErrorResponse{
				Error: "There was an error enabling the module.",
			})
		if err != nil {
			panic(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).
		Encode(Response{
			Status: "Successfully enabled Moderation module!",
		})
	if err != nil {
		panic(err)
	}
}
