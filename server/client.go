package server

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client - struct holding messageclient info and threads
type Client struct {
	Conn     *websocket.Conn
	UserInfo *User
	SendChan chan Message
}

// NewClient - Creates a new client structure
func NewClient(conn *websocket.Conn, userInfo *User) *Client {
	return &Client{
		Conn:     conn,
		UserInfo: userInfo,
		SendChan: make(chan Message, 64),
	}
}

// Nickname - Returns the user's nickname
func (c *Client) Nickname() string {
	return c.UserInfo.Nickname
}

func (c *Client) receiveRoutine(unregisterChan chan *Client, broadcastChan chan Message) {
	for {
		msgType, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage: ", err)
			switch {
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				log.Println("AbnormalClosure")
			case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
				log.Println("ClientDisconnected")
			}
			unregisterChan <- c
			return
		}
		switch msgType {
		case websocket.TextMessage:
			broadcastChan <- Message{c.Nickname(), string(messageData)}
		case websocket.BinaryMessage:
			log.Println("BinaryMessage")
		case websocket.PingMessage:
			log.Println("PingMessage")
		case websocket.PongMessage:
			log.Println("PongMessage")
		}
	}
}
