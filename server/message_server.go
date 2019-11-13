package server

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

// MessageServer - Central struct for message server representation
type MessageServer struct {
	ID                 string
	Users              map[string]bool
	Rooms              map[string]*Room
	AuthenticatedUsers map[string]*User
	ConnectedClients   map[string]*Client
}

// NewMessageServer - Returns a fresh instance of a MessageServer
func NewMessageServer() *MessageServer {
	return &MessageServer{
		ID:    "unamed",
		Users: map[string]bool{},
		Rooms: map[string]*Room{
			"global": NewRoom(),
		},
		AuthenticatedUsers: map[string]*User{},
		ConnectedClients:   map[string]*Client{},
	}
}

// AuthenticateUser - Authenticates an incoming user to the server
func (s *MessageServer) AuthenticateUser(nickname, tokenStr string) error {
	if _, ok := s.Users[nickname]; ok {
		return fmt.Errorf("%s already in use", nickname)
	}
	if err := ValidateNickname(nickname); err != nil {
		return err
	}
	s.Users[nickname] = true
	s.AuthenticatedUsers[tokenStr] = NewUser(nickname)
	return nil
}

// ConnectUser - Effectively connects a user to receive/send messages
func (s *MessageServer) ConnectUser(tokenStr string, conn *websocket.Conn) error {
	user, ok := s.AuthenticatedUsers[tokenStr]
	if !ok {
		return errors.New("user not authenticated")
	}
	globalRoom := s.Rooms["global"]
	client := NewClient(conn, user, globalRoom)
	s.ConnectedClients[tokenStr] = client
	globalRoom.RegisterChan <- client
	client.Run()
	return nil
}

// Run - MessageServer's main routine
func (s *MessageServer) Run() error {
	for _, room := range s.Rooms {
		go room.Run()
	}
	return nil
}
