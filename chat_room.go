package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type chatMessage struct {
	From string `json:"from"`
	Body string `json:"body"`
}

type chatRoom struct {
	id             string
	wsServer       *http.Server
	upgrader       *websocket.Upgrader
	registerChan   chan *chatClient
	unregisterChan chan string
	broadcastChan  chan string
	clients        map[string]*chatClient
	incomingAddr   map[string]string
	incomingNick   map[string]string
}

func newRoom(id string) *chatRoom {
	room := &chatRoom{
		id,
		new(http.Server),
		&websocket.Upgrader{
			CheckOrigin: checkOrigin,
		},
		make(chan *chatClient, 32),
		make(chan string, 32),
		make(chan string, 128),
		map[string]*chatClient{},
		map[string]string{},
		map[string]string{},
	}
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWs).Methods("GET")
	return room
}

func (cr *chatRoom) mainRoutine() {
	for {
		select {
		case client, ok := <-cr.registerChan:
			if ok {
				cr.clients[client.nickname] = client
				go client.sendRoutine()
				go client.receiveRoutine(cr.broadcastChan, cr.unregisterChan)
				var message chatMessage
				message.From = "system"
				message.Body = fmt.Sprintf("%s joined", client.nickname)
				encodedMsg, err := json.Marshal(message)
				if err != nil {
					log.Println(err.Error())
				} else {
					cr.broadcastChan <- string(encodedMsg)
				}
			}
		case nickname, ok := <-cr.unregisterChan:
			if ok {
				delete(cr.clients, nickname)
				cr.broadcastChan <- fmt.Sprintf("%s left", nickname)
			}
		case message, ok := <-cr.broadcastChan:
			if ok {
				for _, client := range cr.clients {
					client.sendChan <- message
				}
			}
		}
	}
}

func checkOrigin(r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		return false
	}
	nickname := r.FormValue("nickname")
	clientAddr, ok := server.connectedUsers[nickname]
	if !ok {
		log.Println("user not found")
		return false
	}
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	if clientAddr != clientIP {
		log.Println("different origins")
		log.Printf("%s:%s", clientAddr, r.RemoteAddr)
		return false
	}
	return true
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	// Get the token.
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		return
	}
	tokenString := r.FormValue("token")
	if tokenString == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Println(tokenString)
	// Parse token.
	chatToken := &chatToken{}
	token, err := jwt.ParseWithClaims(tokenString, chatToken, func(tkn *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// Establish connection.
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	client := newClient(conn, chatToken.Nickname)
	server.rooms["global"].registerChan <- client
}
