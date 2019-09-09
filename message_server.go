package main

import "net/http"

type messageServer struct {
	rooms       map[string]*chatRoom
	userCounter uint
}

var server *messageServer

func init() {
	server = new(messageServer)
	server.rooms = map[string]*chatRoom{
		"global": newRoom("global"),
	}
	server.userCounter = 0
}

func run() error {
	go server.rooms["global"].mainRoutine()
	return http.ListenAndServe(":7000", enableCORS(msgServerHTTPRouter))
}
