package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type messageServer struct {
	httpServer  *http.Server
	upgrader    *websocket.Upgrader
	rooms       map[string]*chatRoom
	anomCounter int
}

var server *messageServer

func init() {
	server = new(messageServer)
	server.httpServer = new(http.Server)
	server.upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	server.anomCounter = 0
	server.httpServer.Addr = ":7000"
	router := mux.NewRouter()
	router.HandleFunc("/ws", server.handleWs)
	router.HandleFunc("/createRoom", server.createRoom)
	server.httpServer.Handler = router
	globalRoom := newRoom("global")
	server.rooms = map[string]*chatRoom{globalRoom.id: globalRoom}
	go globalRoom.mainRoutine()
}

func (ms *messageServer) handleWs(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nickname := r.FormValue("nickname")
	if nickname == "" {
		nickname = fmt.Sprintf("anonymous#%d", server.anomCounter)
		server.anomCounter++
	}
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	globalRoom := server.rooms["global"]
	client := newClient(conn, nickname)
	globalRoom.registerChan <- client
}

func (ms *messageServer) createRoom(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currRoomID := r.PostFormValue("room")
	currRoom := ms.rooms[currRoomID]
	nickname := r.PostFormValue("nickname")
	client, ok := currRoom.clients[nickname]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	currRoom.unregisterChan <- nickname
	newRoomID := r.PostFormValue("createRoom")
	createdRoom := newRoom(newRoomID)
	ms.rooms[newRoomID] = createdRoom
	go createdRoom.mainRoutine()
	createdRoom.registerChan <- client
}
