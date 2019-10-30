package legacy

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type messageServer struct {
	httpServer *http.Server
	upgrader   *websocket.Upgrader
	rooms      map[string]*chatRoom
}

var server *messageServer

func init() {
	server = new(messageServer)
	server.httpServer = new(http.Server)
	server.httpServer.Addr = ":7000"
	router := mux.NewRouter()
	router.Use(validateTokenMiddleware)
	router.HandleFunc("/join", joinHandler).Methods("GET")
	router.HandleFunc("/ws", wsHandler).Methods("GET")
	server.httpServer.Handler = enableCORS(router)
	server.upgrader = &websocket.Upgrader{
		CheckOrigin: upgraderCheckOrigin,
	}
	server.rooms = map[string]*chatRoom{
		"global": newRoom("global"),
	}
}

func run() error {
	for _, room := range server.rooms {
		go room.mainRoutine()
	}
	return server.httpServer.ListenAndServe()
}

func upgraderCheckOrigin(r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		return false
	}
	tokenString := r.FormValue("token")
	claims, err := validateToken(tokenString)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	nickname := r.FormValue("nickname")
	if nickname != claims.UserID {
		log.Println("nicknames don't match")
		return false
	}
	return true
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	nickname := r.FormValue("nickname")
	if err := validateNickname(nickname); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims := &chatClaims{
		ClientIP: clientIP,
		JoinTime: time.Now().UnixNano(),
		RoomID:   "global",
		UserID:   nickname,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseToken := struct {
		Nickname string `json:"nickname"`
		Token    string `json:"token"`
	}{nickname, tokenString}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseToken)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	nickname := r.FormValue("nickname")
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	globalRoom := server.rooms["global"]
	client := newClient(nickname, conn)
	globalRoom.registerChan <- client
	<-client.doneChan
	log.Printf("%s disconnected", client.nickname)
}

func validateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noTokenRoutes := []string{"/join", "/ws"}
		requestPath := r.URL.Path
		for _, route := range noTokenRoutes {
			if requestPath == route {
				next.ServeHTTP(w, r)
				return
			}
		}
		if err := r.ParseForm(); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenString := r.FormValue("token")
		claims, err := validateToken(tokenString)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		nickname := r.FormValue("nickname")
		if nickname != claims.UserID {
			log.Println("nicknames don't match")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
