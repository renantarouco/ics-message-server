package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route - Struct used to define main aspects of an API route
type Route struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
	Name    string
	Method  string
}

// Routes - Routes definition to be instantiated by the API
var Routes []Route = []Route{
	{"/auth", AuthHandler, "auth", http.MethodPost},
	{"/ws", WsHandler, "ws", http.MethodGet},
	{"/nickname", NicknameHandler, "nickname", http.MethodPut},
	{"/rooms", CreateRoomHandler, "create-room", http.MethodPost},
	{"/rooms", SwitchRoomHandler, "switch-room", http.MethodPut},
	{"/users", UsersHandler, "users", http.MethodGet},
	{"/exit", ExitHandler, "exit", http.MethodGet},
}

// APIRouter - The router instance serving the whole application for HTTP Requests
var Router = mux.NewRouter()

func init() {
	for _, route := range Routes {
		Router.HandleFunc(route.Path, route.Handler).
			Name(route.Name).
			Methods(route.Method)
	}
	Router.Use(ValidateTokenMiddleware)
}
