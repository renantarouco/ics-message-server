package api

import (
	"encoding/json"
	"net/http"

	"github.com/renantarouco/ics-message-server/server"
)

var s = server.NewMessageServer()

// AuthHandler - The authentication handler. Receives via POST an URL encoded
// nickname field and returns a JSON object with a token for future requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nickname := r.FormValue("nickname")
	token, err := s.AuthenticateUser(nickname)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	responseToken := map[string]interface{}{
		"token": token,
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseToken)
}
