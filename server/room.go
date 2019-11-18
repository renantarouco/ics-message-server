package server

import "fmt"

// Room - Room representation and threads holder
type Room struct {
	RegisterChan   chan *Client
	UnregisterChan chan *Client
	BroadcastChan  chan Message
	Clients        map[*Client]bool
	Overlay        *Overlay
}

// NewRoom - Instantiates a new room
func NewRoom() *Room {
	return &Room{
		RegisterChan:   make(chan *Client, 32),
		UnregisterChan: make(chan *Client, 32),
		BroadcastChan:  make(chan Message, 128),
		Clients:        map[*Client]bool{},
		Overlay:        NewOverlay(),
	}
}

// Run - Room's main thread, reads on three main channels
// RegisterChan: Incomming clients to join room
// UnregisterChan: Clients leaving the room
// BroadcastChan: Messages to be broadcasted
func (r *Room) Run() {
	for {
		select {
		case client, ok := <-r.RegisterChan:
			if ok {
				r.Clients[client] = true
				message := Message{
					"system",
					fmt.Sprintf("%s joined", client.Nickname()),
				}
				r.BroadcastChan <- message
			}
		case client, ok := <-r.UnregisterChan:
			if ok {
				delete(r.Clients, client)
			}
		case message, ok := <-r.BroadcastChan:
			if ok {
				go r.BroadcastClients(message)
				r.Overlay.BroadcastMessage(message)
			}
		}
	}
}

func (r *Room) BroadcastClients(message Message) {
	for client := range r.Clients {
		client.SendChan <- message
	}
}
