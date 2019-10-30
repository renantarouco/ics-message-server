package server

import "fmt"

// MessageServer - Central struct for message server representation
type MessageServer struct {
	ID    string
	Users map[string]string
}

// NewMessageServer - Returns a fresh instance of a MessageServer
func NewMessageServer() *MessageServer {
	return &MessageServer{
		Users: map[string]string{},
	}
}

// AuthenticateUser - Authenticates an incoming user to the server
func (s *MessageServer) AuthenticateUser(nickname string) (string, error) {
	if _, ok := s.Users[nickname]; ok {
		return "", fmt.Errorf("%s already in use", nickname)
	}
	token, err := NewTokenString(s.ID, nickname)
	if err != nil {
		return "", err
	}
	s.Users[nickname] = token
	return token, nil
}
