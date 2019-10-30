package legacy

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type chatClient struct {
	nickname string
	conn     *websocket.Conn
	sendChan chan chatMessage
	doneChan chan struct{}
}

func newClient(nickname string, conn *websocket.Conn) *chatClient {
	client := new(chatClient)
	client.nickname = nickname
	client.conn = conn
	client.sendChan = make(chan chatMessage, 64)
	client.doneChan = make(chan struct{})
	return client
}

func (client *chatClient) receiveRoutine(unregisterChan chan *chatClient, broadcastChan chan chatMessage) {
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
			unregisterChan <- client
			client.doneChan <- struct{}{}
			return
		}
		switch msgType {
		case websocket.TextMessage:
			broadcastChan <- chatMessage{client.nickname, string(messageData)}
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
