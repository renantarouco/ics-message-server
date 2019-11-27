package server

import (
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Server - Central struct for message server representation
type Server struct {
	ID               string
	NS               *NameService
	Rooms            map[string]*Room
	ConnectedClients map[string]*Client
	Lock             sync.Mutex
}

// NewServer - Returns a fresh instance of a Server
func NewServer(nsEndpoints []string) *Server {
	return &Server{
		ID:               "unamed",
		NS:               NewNameService(nsEndpoints),
		Rooms:            map[string]*Room{},
		ConnectedClients: map[string]*Client{},
		Lock:             sync.Mutex{},
	}
}

// AuthenticateUser - Authenticates an incoming user to the server
func (s *Server) AuthenticateUser(nickname, clientAddr string) (string, error) {
	if err := s.NS.ReserveNickname(nickname); err != nil {
		return "", err
	}
	tokenStr, err := NewTokenString(s.ID, clientAddr)
	if err != nil {
		return "", err
	}
	user := NewUser(nickname, clientAddr, tokenStr)
	s.NS.SaveUser(*user)
	return tokenStr, nil
}

// ConnectUser - Effectively connects a user to receive/send messages
func (s *Server) ConnectUser(tokenStr string, conn *websocket.Conn) error {
	user, err := s.NS.GetUser(tokenStr)
	if err != nil {
		return err
	}
	client := NewClient(conn, user)
	if err := s.CreateRoom(client, "global"); err != nil {
		return err
	}
	s.ConnectedClients[tokenStr] = client
	return client.ReceiveRoutine()
}

// SendMessage - Sends a message from a given client
func (s *Server) SendMessage(client *Client, from, body string) error {
	client.Room.Broadcast(from, body)
	return nil
}

// SetNickname - Sets a client nickname
func (s *Server) SetNickname(client *Client, nickname string) error {
	if err := s.NS.ReserveNickname(nickname); err != nil {
		return err
	}
	if err := s.NS.ChangeUserNickname(client.TokenStr(), nickname); err != nil {
		return err
	}
	client.UserInfo.Nickname = nickname
	return nil
}

// SwitchRoom - Changes the client's room
func (s *Server) SwitchRoom(client *Client, roomID string) error {
	room, ok := s.Rooms[roomID]
	if ok {
		room, err := s.NS.GetRoom(roomID)
		if err != nil {
			return err
		}
		s.Rooms[roomID] = room
	}
	if client.Room != nil {
		client.Room.Unregister(client)
	}
	room.Register(client)
	return nil
}

// CreateRoom - Client command to create a room
func (s *Server) CreateRoom(client *Client, roomID string) error {
	if err := s.NS.ReserveRoomName(roomID); err != nil {
		return err
	}
	room := NewRoom(roomID)
	s.Rooms[roomID] = room
	log.Debugf("created %s room", roomID)
	if err := s.SwitchRoom(client, roomID); err != nil {
		return err
	}
	return nil
}

// ListUsers - List users in the clients room
func (s *Server) ListUsers(client *Client) error {
	nicknames, err := s.NS.GetUsersList(client.Room.Name)
	if err != nil {
		return err
	}
	return client.Send("system", strings.Join(nicknames, "\n"))
}

// ListRooms - Lists all available rooms
func (s *Server) ListRooms(client *Client) error {
	rooms, err := s.NS.GetRoomsList()
	if err != nil {
		return err
	}
	return client.Send("system", strings.Join(rooms, "\n"))
}

// Exit - Clients disconnection function
func (s *Server) Exit(client *Client, roomID string) error {
	client.Room.Unregister(client)
	if len(s.Rooms[roomID].Clients) == 0 {
		delete(s.Rooms, roomID)
		log.Debugf("%s room deleted becouse it's empty", client.Nickname())
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.ConnectedClients, client.TokenStr())
	log.Debugf("%s user left the server", client.Nickname())
	return nil
}
