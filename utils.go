package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// EnableCORS - Necessary headers to allow CORS
func EnableCORS(router *mux.Router) http.Handler {
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

// NewTokenString - Creates a new JWT token for a client
func NewTokenString(subject, issuer string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Subject:   subject,
		Issuer:    issuer,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Date(2021, time.December, 31, 0, 0, 0, 0, time.UTC).Unix(),
	})
	jwtKey, err := GetJWTKey()
	if err != nil {
		return "", err
	}
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// IsTokenValid - Checks if a token string is valid
func IsTokenValid(tokenStr string) error {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
		return GetJWTKey()
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// GetJWTKey - Gets the JWT key environment variable
func GetJWTKey() ([]byte, error) {
	jwtKey := []byte(viper.GetString("JWT_KEY"))
	if len(jwtKey) == 0 {
		return []byte{}, errors.New("jwt key environment variable not set")
	}
	return jwtKey, nil
}
