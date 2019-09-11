package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type chatClient struct {
	nickname      string
	conn          *websocket.Conn
	broadcastChan chan chatMessage
	sendChan      chan chatMessage
}

func newClient(nickname string, conn *websocket.Conn, broadcastChan chan chatMessage) *chatClient {
	client := new(chatClient)
	client.nickname = nickname
	client.conn = conn
	client.broadcastChan = broadcastChan
	client.sendChan = make(chan chatMessage, 64)
	return client
}

func (client *chatClient) receiveRoutine() {
	for {
		msgType, messageData, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage: ", err)
			switch {
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				log.Println("AbnormalClosure")
			case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
				log.Println("ClientDisconnected")
			}
			return
		}
		switch msgType {
		case websocket.TextMessage:
			client.broadcastChan <- chatMessage{client.nickname, string(messageData)}
		case websocket.BinaryMessage:
			log.Println("BinaryMessage")
		case websocket.PingMessage:
			log.Println("PingMessage")
		case websocket.PongMessage:
			log.Println("PongMessage")
		}
	}
}

func (client *chatClient) sendRoutine() {
	for {
		message, ok := <-client.sendChan
		if !ok {
			return
		}
		encodedMessage, err := json.Marshal(message)
		if err != nil {
			log.Printf("error decoding message from %s to %s", message.From, client.nickname)
		}
		client.conn.WriteMessage(websocket.TextMessage, encodedMessage)
	}
}
