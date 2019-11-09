package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/renantarouco/ics-message-server/server"
)

var s *server.MessageServer
var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// RegisterServer - Register the server responding to API calls
func RegisterServer(newServer *server.MessageServer) {
	s = newServer
}

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
	tokenStr, ok := r.Context().Value("tokenStr").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.ConnectUser(tokenStr, conn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "users list")
}

// ExitHandler - Gracefully disconnects the user from the server
func ExitHandler(w http.ResponseWriter, r *http.Request) {

}
