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
	id            string
	registerChan  chan *chatClient
	broadcastChan chan chatMessage
	clients       map[string]*chatClient
}

func newRoom(id string) *chatRoom {
	room := &chatRoom{
		id,
		make(chan *chatClient, 32),
		make(chan chatMessage, 128),
		map[string]*chatClient{},
	}
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
				message := chatMessage{
					"system",
					fmt.Sprintf("%s joined", client.nickname),
				}
				room.broadcastChan <- message
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
