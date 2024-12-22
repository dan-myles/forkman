package server

import (
	"net/http"
	"strconv"
	"time"
)

func (s *Server) uptime(w http.ResponseWriter, r *http.Request) {
	time := time.Since(s.startTime)
	seconds := int(time.Seconds())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "Alive!", "uptime":` + strconv.Itoa(seconds) + ` }`))
}
