package main

import "fmt"

type chatRoom struct {
	id             string
	registerChan   chan *chatClient
	unregisterChan chan string
	broadcastChan  chan string
	clients        map[string]*chatClient
	incomingAddr   map[string]string
	incomingNick   map[string]string
}

func newRoom(id string) *chatRoom {
	return &chatRoom{
		id,
		make(chan *chatClient, 32),
		make(chan string, 32),
		make(chan string, 128),
		map[string]*chatClient{},
		map[string]string{},
		map[string]string{},
	}
}

func (cr *chatRoom) mainRoutine() {
	for {
		select {
		case client, ok := <-cr.registerChan:
			if ok {
				cr.clients[client.nickname] = client
				go client.sendRoutine()
				go client.receiveRoutine(cr.broadcastChan, cr.unregisterChan)
				cr.broadcastChan <- fmt.Sprintf("%s joined", client.nickname)
			}
		case nickname, ok := <-cr.unregisterChan:
			if ok {
				delete(cr.clients, nickname)
				cr.broadcastChan <- fmt.Sprintf("%s left", nickname)
			}
		case message, ok := <-cr.broadcastChan:
			if ok {
				for _, client := range cr.clients {
					client.sendChan <- message
				}
			}
		}
	}
}
