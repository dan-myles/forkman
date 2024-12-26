package err

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorResponses struct {
	Error []string `json:"errors"`
}

var (
	ErrGuildNotFound        = errors.New("guild could not be found in database")
	ErrNoSnowflakeIncluded  = errors.New("guild snowflake could not be found in request")
	ErrNoChannelIdIncluded  = errors.New("channel id could not be found in request")
  ErrNoGuildIdIncluded    = errors.New("guild id could not be found in request")
	ErrAuthProviderNotFound = errors.New("auth provider could not be found in request")
  ErrUnauthorizedGuild    = errors.New("you are not authorized to use this function in this guild")
)

func ServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

func ValidationError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}
