package server

import (
	"context"
	"log"

	pb "github.com/renantarouco/ics-message-server/api/grpc/proto"
	"google.golang.org/grpc"
)

// RoomOverlay - Struct holding an overlay network info
type Overlay struct {
	ConnectionPool map[string]*grpc.ClientConn
	Members        map[string]*pb.OverlayService_BroadcastMessageClient
}

func NewOverlay() *Overlay {
	return &Overlay{
		Members: map[string]*pb.OverlayService_BroadcastMessageClient{},
	}
}

func (o *Overlay) AddMember(addr string) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(), grpc.WithBlock(),
	}
	conn, err := grpc.Dial(addr, opts...)
	log.Println("connected to grpc server")
	if err != nil {
		return err
	}
	o.ConnectionPool[addr] = conn
	client := pb.NewOverlayServiceClient(conn)
	stream, err := client.BroadcastMessage(context.Background(), &pb.BroadcastMessageRequest{})
	if err != nil {
		return err
	}
	o.Members[addr] = &stream
	return nil
}

func (o *Overlay) BroadcastMessage(message Message) error {
	for _, stream := range o.Members {
		message := &pb.Message{
			From: message.From,
			Body: message.Body,
		}
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func (o *Overlay) LeaveOverlay() error {
	for _, stream := range o.Members {
		stream.Close()
	}
	return nil
}
