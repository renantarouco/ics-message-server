package server

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Client - struct holding messageclient info and threads
type Client struct {
	Conn     *websocket.Conn
	UserInfo *User
	Room     *Room
	RoomChan chan *Room
}

// NewClient - Creates a new client structure
func NewClient(conn *websocket.Conn, userInfo *User, room *Room) *Client {
	return &Client{
		Conn:     conn,
		UserInfo: userInfo,
		Room:     room,
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
func (c *Client) ReceiveRoutine() error {
	for {
		msgType, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			switch {
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				log.Debugf("%s user had abnormal closure", c.Nickname())
			case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
				log.Debugf("%s user had client disconnected", c.Nickname())
			default:
				log.Debug(err.Error())
				return err
			}
			ExecuteCommand(c, Command{CommandExit, nil})
			return nil
		}
		switch msgType {
		case websocket.TextMessage:
			var command Command
			err := json.Unmarshal(messageData, &command)
			if err != nil {
				c.Send("system", "error parsing your message")
				continue
			}
			if err := ExecuteCommand(c, command); err != nil {
				c.Send("system", err.Error())
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

// Send - Sends a message to the client
func (c *Client) Send(from, body string) error {
	message := Message{from, body}
	encodedMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error decoding message from %s to %s", message.From, c.Nickname())
	}
	return c.Conn.WriteMessage(websocket.TextMessage, encodedMessage)
}

// Stop - Gracefully stops client routines closing all of its channels
func (c *Client) Stop() {
	close(c.RoomChan)
}
