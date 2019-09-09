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

type chatToken struct {
	jwt.StandardClaims
	UserID   uint
	ClientIP string
	JoinTime int64
	Nickname string
	RoomID   string
}

var msgServerHTTPRouter *mux.Router

func init() {
	msgServerHTTPRouter = mux.NewRouter()
	msgServerHTTPRouter.HandleFunc("/join", joinChat).Methods("POST")
	msgServerHTTPRouter.HandleFunc("/users/{nickname}", changeNickname).Methods("PUT")
	msgServerHTTPRouter.HandleFunc("/rooms", listRooms).Methods("GET")
	msgServerHTTPRouter.HandleFunc("/rooms", createRoom).Methods("POST")
}

func joinChat(w http.ResponseWriter, r *http.Request) {
	// Get chosen nickname.
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Create token.
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	nickname := r.PostFormValue("nickname")
	if nickname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatToken := &chatToken{
		UserID:   server.userCounter,
		ClientIP: clientIP,
		JoinTime: time.Now().UnixNano(),
		Nickname: nickname,
		RoomId:   "global",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, chatToken)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Instantiate user and send the nickname and token to the client.
	responseToken := struct {
		Token string `json:"token"`
	}{tokenString}
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
