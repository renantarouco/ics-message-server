package main

import (
	"net/http"

	"github.com/renantarouco/ics-message-server/api"
	httpAPI "github.com/renantarouco/ics-message-server/api/http"
	"github.com/renantarouco/ics-message-server/server"
)

var s = server.NewMessageServer()

// Init - Initialize the singleton message server
func Init() error {
	// s.FetchNameServer(ip, port)
	api.RegisterServer(s)
	return nil
}

// Run - Message server main routine
func Run() error {
	errorChan := make(chan error)
	// Running HTTP API
	go func() {
		if err := http.ListenAndServe(":7000", httpAPI.EnableCORS(httpAPI.Router)); err != nil {
			errorChan <- err
		}
	}()
	// TODO: Run gRPC API
	// Running server routine
	go func() {
		if err := s.Run(); err != nil {
			errorChan <- err
		}
	}()
	return <-errorChan
}
