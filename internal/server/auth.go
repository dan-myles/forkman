package server

import (
	"context"
	"net/http"
	"time"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markbates/goth/gothic"
)

var sessionKey = "forkman-user-session"

// authLogin godoc
//
//	@summary Login with a provider
//	@description Redirects the user to the provider's login page with the callback URL.
//	@tags auth
//	@router /auth/{provider}/login [get]
func (s *Server) authLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

// authCallback godoc
//
//	@summary Callback from a provider
//	@description Handles the callback from provider, saves user to session & DB.
//	@tags auth
//	@router /auth/{provider}/callback [get]
func (s *Server) authCallback(w http.ResponseWriter, r *http.Request) {
	log := s.log.With().Str("request_id", middleware.GetReqID(r.Context())).Logger()
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	// Complete auth
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to complete user auth")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save user to session
	session, _ := gothic.Store.Get(r, sessionKey)
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save user to session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user exists
	dbUser := database.User{}
	result := s.db.Find(&dbUser, "discord_id = ?", user.UserID)
	if result.RowsAffected > 0 {
		// Update any user info
		dbUser.DiscordUsername = user.Name
		dbUser.DiscordAvatarURL = user.AvatarURL
		dbUser.DiscordEmail = user.Email
		dbUser.LastLogin = time.Now()
		s.db.Save(&dbUser)

		log.Info().Interface("user", user).Msg("Existing user authenticated")
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	// // Generate UUID
	// uuid, err := uuid.NewRandom()
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to generate UUID")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Save user to DB
	s.db.Create(&database.User{
		DiscordID:        user.UserID,
		DiscordUsername:  user.Name,
		DiscordAvatarURL: user.AvatarURL,
		DiscordEmail:     user.Email,
		LastLogin:        time.Now(),
	})

	log.Info().Interface("user", user).Msg("New user authenticated")
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

// authLogout godoc
//
//	@summary Logout with a provider.
//	@description Logs the user out of the provider and clears the session.
//	@tags auth
//	@router /auth/{provider}/callback [get]
func (s *Server) authLogout(w http.ResponseWriter, r *http.Request) {
	log := s.log.With().Str("request_id", middleware.GetReqID(r.Context())).Logger()
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	// Logout
	err := gothic.Logout(w, r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to logout user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear session
	session, _ := gothic.Store.Get(r, sessionKey)
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Error().Err(err).Msg("Failed to clear session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
