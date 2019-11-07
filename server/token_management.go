package server

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

// NewTokenString - Creates a new JWT token for a client
func NewTokenString(subject, issuer string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Subject:   subject,
		Issuer:    issuer,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Date(2021, time.December, 31, 0, 0, 0, 0, time.UTC).Unix(),
	})
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
