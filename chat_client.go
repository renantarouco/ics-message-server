package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type chatClient struct {
	nickname string
	conn     *websocket.Conn
	sendChan chan string
}

func newClient(conn *websocket.Conn, nickname string) *chatClient {
	return &chatClient{
		nickname,
		conn,
		make(chan string, 64),
	}
}

func (cc *chatClient) receiveRoutine(broadcastChan chan<- string, unregisterChan chan<- string) {
	for {
		msgType, msg, err := cc.conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage: ", err)
			switch {
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				log.Println("Abnormal Closure")
			case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
				log.Println("Client disconnected!")
			}
			unregisterChan <- cc.nickname
			close(cc.sendChan)
			return
		}
		switch msgType {
		case websocket.TextMessage:
			broadcastChan <- string(msg)
		case websocket.BinaryMessage:
			log.Println("Received binary message.")
		case websocket.PingMessage:
			log.Println("Received ping message.")
		case websocket.PongMessage:
			log.Println("Received pong message.")
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
