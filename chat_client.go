package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type chatClient struct {
	conn     *websocket.Conn
	sendChan chan string
}

func newClient(conn *websocket.Conn, userID uint) *chatClient {
	return &chatClient{
		conn,
		make(chan string, 64),
	}
}

func (cc *chatClient) receiveRoutine(broadcastChan chan<- string, unregisterChan chan<- uint) {
	for {
		msgType, msg, err := cc.conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage: ", err)
			switch {
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				log.Println("AbnormalClosure")
			case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
				log.Println("ClientDisconnected")
			}
			close(cc.sendChan)
			return
		}
		switch msgType {
		case websocket.TextMessage:
			broadcastChan <- string(msg)
		case websocket.BinaryMessage:
			log.Println("BinaryMessage")
		case websocket.PingMessage:
			log.Println("PingMessage")
		case websocket.PongMessage:
			log.Println("PongMessage")
		}
	}
}

func (cc *chatClient) sendRoutine() {
	for {
		msg, ok := <-cc.sendChan
		if ok {
			cc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			continue
		}
	}
}
