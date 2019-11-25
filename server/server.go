package server

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Server - Central struct for message server representation
type Server struct {
	ID string
	// Users - Holds users nickname to check uniqueness
	Users map[string]bool
	Rooms map[string]*Room
	// AuthenticatedUsers - Holds token and User struct
	AuthenticatedUsers map[string]*User
	ConnectedClients   map[string]*Client
	Lock               sync.Mutex
}

// NewServer - Returns a fresh instance of a Server
func NewServer() *Server {
	return &Server{
		ID:                 "unamed",
		Users:              map[string]bool{},
		Rooms:              map[string]*Room{},
		AuthenticatedUsers: map[string]*User{},
		ConnectedClients:   map[string]*Client{},
		Lock:               sync.Mutex{},
	}
}

// AuthenticateUser - Authenticates an incoming user to the server
func (s *Server) AuthenticateUser(nickname, clientAddr string) (string, error) {
	if _, ok := s.Users[nickname]; ok {
		return "", fmt.Errorf("%s already in use", nickname)
	}
	if err := ValidateNickname(nickname); err != nil {
		return "", err
	}
	tokenStr, err := NewTokenString(s.ID, clientAddr)
	if err != nil {
		return "", err
	}
	s.Users[nickname] = true
	s.AuthenticatedUsers[tokenStr] = NewUser(nickname, tokenStr)
	return tokenStr, nil
}

// ConnectUser - Effectively connects a user to receive/send messages
func (s *Server) ConnectUser(tokenStr string, conn *websocket.Conn) error {
	user, ok := s.AuthenticatedUsers[tokenStr]
	if !ok {
		return errors.New("user not authenticated")
	}
	globalRoom, ok := s.Rooms["global"]
	if !ok {
		globalRoom = s.NewRoom("global")
	}
	client := NewClient(conn, user, globalRoom)
	s.ConnectedClients[tokenStr] = client
	globalRoom.RegisterChan <- client
	return client.ReceiveRoutine()
}

// NewRoom - Creates a new room in the server
func (s *Server) NewRoom(id string) *Room {
	room := NewRoom(id)
	s.Rooms[id] = room
	go room.Run()
	log.Debugf("created %s room", id)
	return room
}

// SendMessage - Sends a message from a given client
func (s *Server) SendMessage(client *Client, from, body string) error {
	client.Room.Broadcast(from, body)
	return nil
}

// SetNickname - Sets a client nickname
func (s *Server) SetNickname(client *Client, nickname string) error {
	_, ok := s.Users[nickname]
	if ok {
		return fmt.Errorf("nickname %s already taken", nickname)
	}
	delete(s.Users, client.Nickname())
	s.Users[nickname] = true
	client.UserInfo.Nickname = nickname
	return nil
}

// SwitchRoom - Changes the client's room
func (s *Server) SwitchRoom(client *Client, roomID string) error {
	room, ok := s.Rooms[roomID]
	if !ok {
		return fmt.Errorf("room %s does not exist", roomID)
	}
	client.Room.UnregisterChan <- client
	room.RegisterChan <- client
	client.Room = room
	return nil
}

// CreateRoom - Client command to create a room
func (s *Server) CreateRoom(client *Client, roomID string) error {
	_, ok := s.Rooms[roomID]
	if ok {
		return fmt.Errorf("room %s already exists", roomID)
	}

	if err := BasicNameValidation(roomID); err != nil {
		return err
	}
	room := s.NewRoom(roomID)
	s.Rooms[roomID] = room
	if err := s.SwitchRoom(client, roomID); err != nil {
		return err
	}
	return nil
}

// ListUsers - List users in the clients room
func (s *Server) ListUsers(client *Client) error {
	nicknames := []string{}
	for client := range client.Room.Clients {
		nicknames = append(nicknames, client.Nickname())
	}
	return client.Send("system", strings.Join(nicknames, "\n"))
}

// ListRooms - Lists all available rooms
func (s *Server) ListRooms(client *Client) error {
	rooms := []string{}
	for roomID := range s.Rooms {
		rooms = append(rooms, roomID)
	}
	return client.Send("system", strings.Join(rooms, "\n"))
}

// Exit - Clients disconnection function
func (s *Server) Exit(client *Client) error {
	client.Room.UnregisterChan <- client
	client.Stop()
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.ConnectedClients, client.TokenStr())
	delete(s.AuthenticatedUsers, client.TokenStr())
	delete(s.Users, client.Nickname())
	log.Debugf("%s user left the server", client.Nickname())
	return nil
}
