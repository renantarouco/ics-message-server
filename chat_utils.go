package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type chatToken struct {
	jwt.StandardClaims
	ClientIP string
	JoinTime int64
	RoomID   string
	UserID   uint
	Nickname string
}

func enableCORS(router *mux.Router) http.Handler {
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Content-Type",
		"Accept",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
	})
	return handlers.CORS(headersOk, originsOk, methodsOk)(router)
}

func validateToken(tokenString string) (*chatToken, *jwt.Token, error) {
	if tokenString == "" {
		return nil, nil, errors.New("empty token")
	}
	log.Printf("validating token %s", tokenString)
	chatToken := new(chatToken)
	token, err := jwt.ParseWithClaims(tokenString, chatToken, func(tkn *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}
	return chatToken, token, nil
}

func checkOrigin(r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
		return false
	}
	tokenString := r.FormValue("token")
	chatToken, _, err := validateToken(tokenString)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	userRoom, ok := server.rooms[chatToken.RoomID]
	if !ok {
		log.Println("wrong room")
		return false
	}
	user := userRoom.users[chatToken.UserID]
	if !ok {
		log.Println("user not in this room")
		return false
	}
	userToken := user.token
	signedString, err := userToken.SignedString([]byte("secret"))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if tokenString != signedString {
		log.Println("tokens don't match")
		return false
	}
	ctx := context.WithValue(r.Context(), "userID", chatToken.UserID)
	r = r.WithContext(ctx)
	return true
}
