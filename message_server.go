package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type messageServer struct {
	httpServer  *http.Server
	upgrader    *websocket.Upgrader
	userCounter uint
}

var server *messageServer

func init() {
	server = new(messageServer)
	server.httpServer = new(http.Server)
	server.httpServer.Addr = ":7000"
	router := mux.NewRouter()
	router.HandleFunc("/join", joinHandler).Methods("GET")
	router.HandleFunc("/ws", wsHandler).Methods("GET")
	server.httpServer.Handler = enableCORS(router)
	server.userCounter = 0
}

func run() error {
	return server.httpServer.ListenAndServe()
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	nickname := r.FormValue("nickname")
	if err := validateNickname(nickname); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatToken := &chatToken{
		ClientIP: clientIP,
		JoinTime: time.Now().UnixNano(),
		RoomID:   "global",
		UserID:   nickname,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, chatToken)
	tokenString, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	responseToken := struct {
		Nickname string `json:"nickname"`
		Token    string `json:"token"`
	}{nickname, tokenString}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseToken)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

}
