package api

import (
	"fmt"
	"regexp"
)

// ValidateNickname - Checks the correct syntax for the nickname
func ValidateNickname(nickname string) error {
	matched, err := regexp.MatchString("^[_a-zA-Z][_a-zA-Z-0-9]+$", nickname)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("'%s' is not a valid nickname", nickname)
	}
	return nil
}
