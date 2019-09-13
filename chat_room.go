package main

import (
	"fmt"
	"time"
)

type chatMessage struct {
	From string `json:"from"`
	Body string `json:"body"`
}

type chatRoom struct {
	id             string
	registerChan   chan *chatClient
	unregisterChan chan *chatClient
	broadcastChan  chan chatMessage
	clients        map[string]*chatClient
}

func newRoom(id string) *chatRoom {
	room := new(chatRoom)
	room.id = id
	room.registerChan = make(chan *chatClient, 32)
	room.unregisterChan = make(chan *chatClient, 32)
	room.broadcastChan = make(chan chatMessage, 128)
	room.clients = map[string]*chatClient{}
	return room
}

func (room *chatRoom) mainRoutine() {
	go func() {
		for {
			message := chatMessage{
				"system",
				fmt.Sprintf("you're on %s room", room.id),
			}
			room.broadcastChan <- message
			time.Sleep(time.Second * 1)
		}
	}()
	for {
		select {
		case client, ok := <-room.registerChan:
			if ok {
				room.clients[client.nickname] = client
				go client.sendRoutine()
				go client.receiveRoutine(room.unregisterChan, room.broadcastChan)
				message := chatMessage{
					"system",
					fmt.Sprintf("%s joined", client.nickname),
				}
				room.broadcastChan <- message
			}
		case client, ok := <-room.unregisterChan:
			if ok {
				close(client.sendChan)
				delete(room.clients, client.nickname)
			}
		case message, ok := <-room.broadcastChan:
			if ok {
				for _, client := range room.clients {
					client.sendChan <- message
				}
			}
		}
	}
}
