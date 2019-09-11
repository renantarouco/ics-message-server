package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type chatUser struct {
	nickname string
	token    *jwt.Token
}

type chatMessage struct {
	From string `json:"from"`
	Body string `json:"body"`
}

type chatRoom struct {
	id             string
	wsServer       *http.Server
	upgrader       *websocket.Upgrader
	users          map[uint]*chatUser
	registerChan   chan *chatClient
	unregisterChan chan uint
	broadcastChan  chan string
	clients        map[uint]*chatClient
}

func newRoom(id string) *chatRoom {
	room := &chatRoom{
		id,
		new(http.Server),
		&websocket.Upgrader{},
		map[uint]*chatUser{},
		make(chan *chatClient, 32),
		make(chan uint, 32),
		make(chan string, 128),
		map[uint]*chatClient{},
	}
	router := mux.NewRouter()
	router.HandleFunc("/ws", room.handleWs).Methods("GET")
	room.wsServer.Handler = enableCORS(router)
	room.wsServer.Addr = ":7001"
	return room
}

func (cr *chatRoom) mainRoutine() {
	go func() {
		cr.wsServer.ListenAndServe()
	}()
	for {
		select {
		case client, ok := <-cr.registerChan:
			if ok {
				cr.clients[client.userID] = client
				go client.sendRoutine()
				var message chatMessage
				message.From = "system"
				message.Body = fmt.Sprintf("%s joined", cr.users[client.userID].nickname)
				encodedMsg, err := json.Marshal(message)
				if err != nil {
					log.Println(err.Error())
				} else {
					cr.broadcastChan <- string(encodedMsg)
				}
			}
		case userID, ok := <-cr.unregisterChan:
			if ok {
				delete(cr.clients, userID)
				cr.broadcastChan <- fmt.Sprintf("%s left", cr.users[userID].nickname)
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

func (cr *chatRoom) handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := cr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		log.Println("userID not in request context")
		return
	}
	client := newClient(conn, userID)
	client.receiveRoutine(cr.broadcastChan, cr.unregisterChan)
}
