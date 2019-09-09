package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var msgServerHTTPRouter *mux.Router

func init() {
	msgServerHTTPRouter = mux.NewRouter()
	msgServerHTTPRouter.HandleFunc("/join", joinChat).Methods("GET")
	msgServerHTTPRouter.HandleFunc("/users/{nickname}", changeNickname).Methods("PUT")
	msgServerHTTPRouter.HandleFunc("/rooms", listRooms).Methods("GET")
	msgServerHTTPRouter.HandleFunc("/rooms", createRoom).Methods("POST")
}

func joinChat(w http.ResponseWriter, r *http.Request) {
	// Get chosen nickname.
	if err := r.ParseForm(); err != nil {
		log.Println("error parsing url form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Create token.
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	userID := server.userCounter
	nickname := r.FormValue("nickname")
	if nickname == "" {
		log.Println("empty nickname")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatToken := &chatToken{
		ClientIP: clientIP,
		JoinTime: time.Now().UnixNano(),
		UserID:   userID,
		RoomID:   "global",
		Nickname: nickname,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, chatToken)
	server.rooms["global"].users[userID] = &chatUser{nickname, token}
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println("error generating token signed string")
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

func changeNickname(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		return
	}
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("creating room")
}

func listRooms(w http.ResponseWriter, r *http.Request) {
	log.Println("listing Rooms")
}
