package api

import (
	httpAPI "github.com/renantarouco/ics-message-server/api/http"
	"github.com/renantarouco/ics-message-server/server"
)

// RegisterServer - Register server for all kinds of APIs
func RegisterServer(s *server.MessageServer) {
	httpAPI.RegisterServer(s)
	// TODO: grpcAPI.RegisterServer(s)
}
