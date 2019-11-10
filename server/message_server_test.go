package server

import (
	"fmt"
	"testing"
)

func TestAuthenticateUser(t *testing.T) {
	testServer := NewMessageServer()
	t.Run("no-nickname", func(t *testing.T) {
		nickname := ""
		err := testServer.AuthenticateUser(nickname, "")
		expectedErrorMsg := fmt.Sprintf("'%s' is not a valid nickname", nickname)
		if err.Error() != expectedErrorMsg {
			t.Errorf("wrong error message, wanted (%s) got (%s) ",
				expectedErrorMsg, err.Error())
		}
	})
	t.Run("invalid-nickname", func(t *testing.T) {
		nickname := "1test1"
		err := testServer.AuthenticateUser(nickname, "")
		expectedErrorMsg := fmt.Sprintf("'%s' is not a valid nickname", nickname)
		if err.Error() != expectedErrorMsg {
			t.Errorf("wrong error message, wanted (%s) got (%s) ",
				expectedErrorMsg, err.Error())
		}
	})
	t.Run("valid-nickname", func(t *testing.T) {
		nickname := "test1"
		err := testServer.AuthenticateUser(nickname, "")
		if err != nil {
			t.Error("valid nickname shouldn't return any errors")
		}
	})
	t.Run("duplicate-nickname", func(t *testing.T) {
		nickname := "test2"
		err := testServer.AuthenticateUser(nickname, "")
		if err != nil {
			t.Error("valid nickname shouldn't return any errors")
		}
		nickname = "test2"
		err = testServer.AuthenticateUser(nickname, "")
		expectedErrorMsg := fmt.Sprintf("%s already in use", nickname)
		if err.Error() != expectedErrorMsg {
			t.Errorf("wrong error message, wanted (%s) got (%s) ",
				expectedErrorMsg, err.Error())
		}
	})
}
