package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/renantarouco/ics-name-server/api/grpc/proto"
	"google.golang.org/grpc"
)

// MessageServer - Central struct for message server representation
type MessageServer struct {
	ID                 string
	Users              map[string]bool
	Rooms              map[string]*Room
	AuthenticatedUsers map[string]*User
	ConnectedClients   map[string]*Client
	GRPCConn           *grpc.ClientConn
	NameServiceClient  pb.NameServiceClient
}

// NewMessageServer - Returns a fresh instance of a MessageServer
func NewMessageServer() *MessageServer {
	return &MessageServer{
		ID:                 "unamed",
		Users:              map[string]bool{},
		Rooms:              map[string]*Room{},
		AuthenticatedUsers: map[string]*User{},
		ConnectedClients:   map[string]*Client{},
	}
}

// Init - Initializes the MessageServer
func (s *MessageServer) Init() error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(), grpc.WithBlock(),
	}
	conn, err := grpc.Dial("localhost:7001", opts...)
	log.Println("connected to grpc server")
	if err != nil {
		return err
	}
	s.GRPCConn = conn
	s.NameServiceClient = pb.NewNameServiceClient(conn)
	s.CreateRoom("global")
	return nil
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
	errorChan := make(chan error)
	go func() {
		defer s.GRPCConn.Close()
		stream, err := s.NameServiceClient.ConnectMessageServer(context.Background())
		if err != nil {
			errorChan <- err
			return
		}
		for {
			heartBeat := &pb.HeartBeat{
				Header: &pb.ConnectionHeader{
					MessageServerId: s.ID,
				},
			}
			if err := stream.Send(heartBeat); err != nil {
				errorChan <- err
				return
			}
			log.Println("sent heartbeat")
			connectionInfo, err := stream.Recv()
			if err != nil {
				errorChan <- err
				return
			}
			if s.ID != connectionInfo.Header.MessageServerId {
				s.ID = connectionInfo.Header.MessageServerId
			}
			log.Println("received connection info")
			time.Sleep(1 * time.Second)
		}
	}()
	return <-errorChan
}

// CreateRoom - Creates a new room in the server
func (s *MessageServer) CreateRoom(id string) (*Room, error) {
	resp, err := s.NameServiceClient.CreateRoom(context.Background())
	if err != nil {
		return nil, err
	}
	room := NewRoom()
	s.Rooms[id] = room
	return room, nil
}
