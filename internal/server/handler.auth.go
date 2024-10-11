package server

import (
	"net/http"
	"time"

	"github.com/avvo-na/forkman/internal/database"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/markbates/goth/gothic"
)

var sessionKey = "forkman-user-session"

// TODO: use repository to insert into database instead of raw-doggin it

// authLogin godoc
//
//	@summary Login with a provider
//	@description Redirects the user to the provider's login page with the callback URL.
//	@tags auth
//	@router /auth/{provider}/login [get]
func (s *Server) authLogin(w http.ResponseWriter, r *http.Request) {
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

	// Complete auth
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to complete user auth")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create user object
	dbUser := database.User{
		DiscordSnowflake: user.UserID,
		DiscordUsername:  user.Name,
		DiscordAvatarURL: user.AvatarURL,
		DiscordEmail:     user.Email,
		LastLogin:        time.Now(),
	}

	// Create or update user in database
	tx := s.db.Begin()
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Failed to start transaction")
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	if tx.Model(&dbUser).
		Where("discord_snowflake = ?", user.UserID).
		Updates(&dbUser).
		RowsAffected == 0 {
		tx.Create(&dbUser)
	}

	if tx.Commit().Error != nil {
		log.Error().Err(tx.Error).Msg("Failed to commit transaction")
		http.Error(w, tx.Error.Error(), http.StatusInternalServerError)

		if tx.Rollback().Error != nil {
			log.Error().Err(tx.Error).Msg("Failed to rollback transaction")
			http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
		}

		return
	}

	// Set user session in cookies
	session, _ := gothic.Store.Get(r, sessionKey)
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save user to session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Interface("user", user).Msg("User authenticated")
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
