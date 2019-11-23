package server

import (
	"errors"

	"github.com/gorilla/websocket"
)

var s *Server

func init() {
	s = NewServer()
}

// ID - Returns the singleton ID
func ID() string {
	return s.ID
}

// AuthenticateUser - Sigleton function to authenticate users
func AuthenticateUser(nickname, tokenStr string) error {
	return s.AuthenticateUser(nickname, tokenStr)
}

// ConnectUser - Singleton function to connect users
func ConnectUser(tokenStr string, conn *websocket.Conn) error {
	return s.ConnectUser(tokenStr, conn)
}

// ExecuteCommand - Executes a command from a given client
func ExecuteCommand(client *Client, command Command) error {
	switch command.Type {
	case CommandMessage:
		from, ok := command.Args["from"].(string)
		if !ok {
			return errors.New("cannot parse 'from' arg of message command")
		}
		body, ok := command.Args["body"].(string)
		if !ok {
			return errors.New("cannot parse 'body' arg of message command")
		}
		return s.SendMessage(client, from, body)
	case CommandSetNickname:
		nickname, ok := command.Args["nickname"].(string)
		if !ok {
			return errors.New("cannot parse 'nickname' arg of setnickname command")
		}
		return s.SetNickname(client, nickname)
	case CommandSwitchRoom:
		roomID, ok := command.Args["room"].(string)
		if !ok {
			return errors.New("cannot parse 'room' arg of switchroom command")
		}
		return s.SwitchRoom(client, roomID)
	}
	return nil
}
