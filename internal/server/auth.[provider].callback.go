package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func (s *Server) getAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	// Get the provider name from the URL
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	// Try to complete the auth
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Setting session cookie LP
	session, _ := gothic.Store.Get(r, "forkman-user-session")
	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: save the user to database

	s.log.Info().Interface("user", user).Msg("User authenticated")
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
