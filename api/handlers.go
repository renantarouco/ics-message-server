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
	nickname := r.PostFormValue("nickname")
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

// WsHandler - The endpoint to connect via Websockets and start sending and
// receiving messages
func WsHandler(w http.ResponseWriter, r *http.Request) {

}

// NicknameHandler - Handles the attempt to chang the nickname
func NicknameHandler(w http.ResponseWriter, r *http.Request) {

}

// CreateRoomHandler - The endpoint to create a new room when POST is sent
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {

}

// SwitchRoomHandler - The endpoint to switch user's room when PUT is sent
func SwitchRoomHandler(w http.ResponseWriter, r *http.Request) {

}

// UsersHandler - Returns the list of all the users in the same room as the one
// requesting
func UsersHandler(w http.ResponseWriter, r *http.Request) {

}

// ExitHandler - Gracefully disconnects the user from the server
func ExitHandler(w http.ResponseWriter, r *http.Request) {

}
