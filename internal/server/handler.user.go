package server

import (
	"encoding/json"
	"net/http"

	e "github.com/avvo-na/forkman/internal/server/common/err"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (s *Server) sessionInfo(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, sessionKey)
	user := session.Values["user"].(goth.User)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) adminInfo(w http.ResponseWriter, r *http.Request) {
	session, _ := gothic.Store.Get(r, sessionKey)
	user := session.Values["user"].(goth.User)

	guilds, err := s.discord.GetUserAdminServers(user.UserID)
	s.log.Info().Msgf("user %s has %d admin servers", user.UserID, len(guilds))
  for _, guild := range guilds {
    s.log.Info().Msgf("admin server: %s", guild.ID)
  }
	if err != nil {
		e.BadRequest(w, err)
		return
	}

	json.NewEncoder(w).Encode(guilds)
}
