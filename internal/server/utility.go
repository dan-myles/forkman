package server

import (
	"encoding/json"
	"net/http"

	"github.com/avvo-na/forkman/internal/database"
)

func (s *Server) enableUtilityModule(w http.ResponseWriter, r *http.Request) {
	// Enable the utility module
	module := database.Module{}

	// Find the module
	result := s.db.First(&module, "name = ?", "utility")
	if result.RowsAffected == 0 {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	if module.Enabled {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Status: "Module already enabled!"})
		return
	}

	// Enable it
	module.Enabled = true
	s.db.Save(&module)

	// Sync the module
	u := s.discord.GetUtility()
	go u.Sync()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: "Success!"})
}

func (s *Server) disableUtilityModule(w http.ResponseWriter, r *http.Request) {

	// Enable the utility module
	module := database.Module{}

	// Find the module
	result := s.db.First(&module, "name = ?", "utility")
	if result.RowsAffected == 0 {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	if !module.Enabled {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Status: "Module already disabled!"})
		return
	}

	// Disable it
	module.Enabled = false
	s.db.Save(&module)

	// Sync the module
	u := s.discord.GetUtility()
	go u.Sync()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: "Success!"})
}
