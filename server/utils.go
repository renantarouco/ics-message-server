package server

import (
	"fmt"
	"regexp"
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
