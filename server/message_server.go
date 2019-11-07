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
	ConnectedUsers     map[string]*User
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
		ConnectedUsers:     map[string]*User{},
	}
}

// AuthenticateUser - Authenticates an incoming user to the server
func (s *MessageServer) AuthenticateUser(nickname string) (string, error) {
	if _, ok := s.Users[nickname]; ok {
		return "", fmt.Errorf("%s already in use", nickname)
	}
	if err := ValidateNickname(nickname); err != nil {
		return "", err
	}
	tokenStr, err := NewTokenString(s.ID, nickname)
	if err != nil {
		return "", err
	}
	s.Users[nickname] = true
	s.AuthenticatedUsers[tokenStr] = NewUser(nickname)
	return tokenStr, nil
}

// ConnectUser - Effectively connects a user to receive/send messages
func (s *MessageServer) ConnectUser(tokenStr string, conn *websocket.Conn) error {
	_, ok := s.AuthenticatedUsers[tokenStr]
	if !ok {
		return errors.New("user not authenticated")
	}
	return nil
}
