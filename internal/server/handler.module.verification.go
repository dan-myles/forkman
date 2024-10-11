package server

import (
	"net/http"

	e "github.com/avvo-na/forkman/internal/server/common/err"
	"github.com/go-chi/chi/v5"
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
		return
	}

	err = mod.SendVerificationPanel(channelId)
	if err != nil {
		e.ServerError(w, err)
	}
}
