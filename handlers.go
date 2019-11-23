package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/renantarouco/ics-message-server/server"
	log "github.com/sirupsen/logrus"
)

var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// AuthHandler - The authentication handler. Receives via POST an URL encoded
// nickname field and returns a JSON object with a token for future requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nickname := r.PostFormValue("nickname")
	tokenStr, err := NewTokenString(server.ID(), r.RemoteAddr)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = server.AuthenticateUser(nickname, tokenStr)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Infof("%s user authenticated", nickname)
	responseToken := map[string]interface{}{
		"token": tokenStr,
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
		log.Error("not found tokenStr in request context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = server.ConnectUser(tokenStr, conn)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
