package server

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client - struct holding messageclient info and threads
type Client struct {
	Conn     *websocket.Conn
	UserInfo *User
	Room     *Room
	SendChan chan Message
	RoomChan chan *Room
}

// NewClient - Creates a new client structure
func NewClient(conn *websocket.Conn, userInfo *User, room *Room) *Client {
	return &Client{
		Conn:     conn,
		UserInfo: userInfo,
		Room:     room,
		SendChan: make(chan Message, 64),
		RoomChan: make(chan *Room),
	}
}

// Nickname - Returns the user's nickname
func (c *Client) Nickname() string {
	return c.UserInfo.Nickname
}

// ReceiveRoutine - Routine for receive messages from a client
func (c *Client) ReceiveRoutine() {
	for {
		select {
		case room, ok := <-c.RoomChan:
			if ok {
				c.Room.UnregisterChan <- c
				c.Room = room
			}
		default:
			msgType, messageData, err := c.Conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage: ", err)
				switch {
				case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
					log.Println("AbnormalClosure")
				case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
					log.Println("ClientDisconnected")
				}
				c.Room.UnregisterChan <- c
				return
			}
			switch msgType {
			case websocket.TextMessage:
				c.Room.BroadcastChan <- Message{c.Nickname(), string(messageData)}
			case websocket.BinaryMessage:
				log.Println("BinaryMessage")
			case websocket.PingMessage:
				log.Println("PingMessage")
			case websocket.PongMessage:
				log.Println("PongMessage")
			}
		}
	}
}
