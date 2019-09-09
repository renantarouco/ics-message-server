package main

type messageServer struct {
	connectedUsers map[string]string
	rooms          map[string]*chatRoom
	userCounter    uint
}

var server *messageServer

func init() {
	// Server instantiation.
	server = new(messageServer)
	server.connectedUsers = map[string]string{}
	globalRoom := newRoom("global")
	server.rooms = map[string]*chatRoom{
		"global": globalRoom,
	}
	go globalRoom.mainRoutine()
	server.userCounter = 0
}
