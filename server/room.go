package server

// Room - Room representation and threads holder
type Room struct {
	RegisterChan   chan *Client
	UnregisterChan chan *Client
	BroadcastChan  chan Message
	Clients        map[string]*Client
}

func NewRoom() *Room {
	return &Room{
		RegisterChan:   make(chan *Client, 32),
		UnregisterChan: make(chan *Client, 32),
		BroadcastChan:  make(chan Message, 128),
		Clients:        map[string]*Client{},
	}
}
