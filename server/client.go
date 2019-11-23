package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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

// TokenStr - Returns the user's token string
func (c *Client) TokenStr() string {
	return c.UserInfo.TokenStr
}

// ReceiveRoutine - Routine for receive messages from a client
func (c *Client) ReceiveRoutine() {
	for {
		msgType, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			switch {
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				log.Debugf("%s user had abnormal closure", c.Nickname())
			case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
				log.Debugf("%s user had client disconnected", c.Nickname())
			}
			ExecuteCommand(c, Command{CommandExit, nil})
			return
		}
		switch msgType {
		case websocket.TextMessage:
			var command Command
			err := json.Unmarshal(messageData, &command)
			if err != nil {
				message := Message{"system", "error parsing your message"}
				c.SendChan <- message
				continue
			}
			if err := ExecuteCommand(c, command); err != nil {
				message := Message{"system", err.Error()}
				c.SendChan <- message
			}
		case websocket.BinaryMessage:
			log.Debugf("attempted binary message")
		case websocket.PingMessage:
			log.Debugf("attempted ping message")
			c.Conn.WriteMessage(websocket.PongMessage, []byte{})
		case websocket.PongMessage:
			log.Debugf("attempted pong message")
		}
	}
}

// SendRoutine - Routine responsible to send messages to the connected client
func (c *Client) SendRoutine() {
	for {
		message, ok := <-c.SendChan
		if !ok {
			return
		}
		encodedMessage, err := json.Marshal(message)
		if err != nil {
			log.Printf("error decoding message from %s to %s", message.From, c.Nickname())
		}
		c.Conn.WriteMessage(websocket.TextMessage, encodedMessage)
	}
}

// Run - Client's main routine
func (c *Client) Run() error {
	go c.SendRoutine()
	c.ReceiveRoutine()
	return nil
}

// Stop - Gracefully stops client routines closing all of its channels
func (c *Client) Stop() {
	close(c.SendChan)
	close(c.RoomChan)
}
