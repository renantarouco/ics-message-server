package server

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// BasicNameValidation - Basic syntax validation for names
func BasicNameValidation(name string) error {
	matched, err := regexp.MatchString("^[_a-zA-Z][_a-zA-Z-0-9]+$", name)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("'%s' is not a valid nickname", name)
	}
	return nil
}

// ValidateNickname - Checks the correct syntax for the nickname
func ValidateNickname(nickname string) error {
	return BasicNameValidation(nickname)
}

// ValidateRoomName - Checks the correct sytax for room names
func ValidateRoomName(roomName string) error {
	return BasicNameValidation(roomName)
}

// GetJWTKey - Gets the JWT key environment variable
func GetJWTKey() ([]byte, error) {
	jwtKey := []byte(viper.GetString("JWT_KEY"))
	if len(jwtKey) == 0 {
		return []byte{}, errors.New("jwt key environment variable not set")
	}
	return jwtKey, nil
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
