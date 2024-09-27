package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) enableModerationModule(w http.ResponseWriter, r *http.Request) {
	// Get guild param
	snowflake := chi.URLParam(r, "guildSnowflake")
	if snowflake == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing guild snowflake"})
		return
	}

	m := s.discord.GetModerationModules()
	module, ok := m[snowflake]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Module not found"})
		return
	}

	err := module.EnableAndSync()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: "Success!"})
}

func (s *Server) disableModerationModule(w http.ResponseWriter, r *http.Request) {
	// Get guild param
	snowflake := chi.URLParam(r, "guildSnowflake")
	if snowflake == "" {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing guild snowflake"})
		return
	}

	m := s.discord.GetModerationModules()
	module, ok := m[snowflake]
	if !ok {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Module not found"})
		return
	}

	err := module.DisableAndSync()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: "Success!"})
}
