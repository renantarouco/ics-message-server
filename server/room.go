package server

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Room - Room representation and threads holder
type Room struct {
	Name    string
	Clients map[*Client]bool
	Lock    sync.Mutex
}

// NewRoom - Instantiates a new room
func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		Clients: map[*Client]bool{},
		Lock:    sync.Mutex{},
	}
}

// Broadcast - Sends a message to all connected clients
func (r *Room) Broadcast(from, body string) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	for client := range r.Clients {
		client.Send(from, body)
	}
}

// Register - Connects a client to the room
func (r *Room) Register(client *Client) {
	r.Lock.Lock()
	r.Clients[client] = true
	r.Lock.Unlock()
	client.Room = r
	log.Debugf("%s user connected to %s room", client.Nickname(), r.Name)
	go r.Broadcast("system", fmt.Sprintf("%s joined", client.Nickname()))
}

// Unregister - Disconnects a client to the room
func (r *Room) Unregister(client *Client) {
	r.Lock.Lock()
	delete(r.Clients, client)
	log.Debugf("%s user left %s room", client.Nickname(), r.Name)
	r.Lock.Unlock()
	client.Room = nil
	go r.Broadcast("system", fmt.Sprintf("%s left", client.Nickname()))
}
