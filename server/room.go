package server

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Room - Room representation and threads holder
type Room struct {
	Name           string
	RegisterChan   chan *Client
	UnregisterChan chan *Client
	BroadcastChan  chan Message
	Clients        map[*Client]bool
	Lock           sync.Mutex
}

// NewRoom - Instantiates a new room
func NewRoom(name string) *Room {
	return &Room{
		Name:           name,
		RegisterChan:   make(chan *Client, 32),
		UnregisterChan: make(chan *Client, 32),
		BroadcastChan:  make(chan Message, 128),
		Clients:        map[*Client]bool{},
		Lock:           sync.Mutex{},
	}
}

// Run - Room's main thread, reads on three main channels
// RegisterChan: Incomming clients to join room
// UnregisterChan: Clients leaving the room
// BroadcastChan: Messages to be broadcasted
func (r *Room) Run() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for client := range r.RegisterChan {
			r.Lock.Lock()
			r.Clients[client] = true
			r.Lock.Unlock()
			log.Debugf("%s user connected to %s room", client.Nickname(), r.Name)
			go r.Broadcast("system", fmt.Sprintf("%s joined", client.Nickname()))
		}
	}()
	go func() {
		defer wg.Done()
		for client := range r.UnregisterChan {
			r.Lock.Lock()
			delete(r.Clients, client)
			log.Debugf("%s user left %s room", client.Nickname(), r.Name)
			if len(r.Clients) == 0 {
				r.Close()
				log.Infof("%s room closed because is empty", r.Name)
			}
			r.Lock.Unlock()
		}
	}()
	wg.Wait()
	log.Infof("%s room thread finished", r.Name)
}

// Broadcast - Sends a message to all connected clients.
func (r *Room) Broadcast(from, body string) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	for client := range r.Clients {
		client.Send(from, body)
	}
}

// Close - Closes the room's channels so the thread is stopped too
func (r *Room) Close() {
	close(r.RegisterChan)
	close(r.UnregisterChan)
	close(r.BroadcastChan)
}
