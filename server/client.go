package server

import "github.com/gorilla/websocket"

// Client - struct holding messageclient info and threads
type Client struct {
	Conn     *websocket.Conn
	UserInfo *User
}
