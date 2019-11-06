package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renantarouco/ics-message-server/server"
)

var jwtKey = []byte("secret_key")

func TestAuthHandler(t *testing.T) {
	t.Run("no-nickname", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/auth", nil)
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		testStatusCode(t, rr, http.StatusBadRequest)
	})
	nicknameTest := generateBasicNamingTest("nickname", http.MethodPost, "/auth", "")
	t.Run("empty-nickname", func(t *testing.T) {
		nicknameTest(t, "", http.StatusForbidden)
	})
	t.Run("invalid-nickname", func(t *testing.T) {
		nicknameTest(t, "1nick1", http.StatusForbidden)
	})
	t.Run("valid-nickname", func(t *testing.T) {
		nicknameTest(t, "nick1", http.StatusCreated)
	})
	t.Run("duplicate-nickname", func(t *testing.T) {
		nicknameTest(t, "nick2", http.StatusCreated)
		nicknameTest(t, "nick2", http.StatusConflict)
	})
	t.Run("valid-token", func(t *testing.T) {
		rr := nicknameTest(t, "nick3", http.StatusCreated)
		var bodyMap map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &bodyMap)
		if err != nil {
			t.Error(err)
		}
		tokenStr, ok := bodyMap["token"]
		if !ok {
			t.Error("body doesn't contain 'token' key")
		}
		err = IsTokenValid(tokenStr)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestWsHandler(t *testing.T) {
	todoTest(t)
}

func TestNicknameHandler(t *testing.T) {
	t.Run("no-nickname", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/nickname", nil)
		if err != nil {
			t.Error(err)
		}
		addBearerAuthHeader(req, validToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		testStatusCode(t, rr, http.StatusBadRequest)
	})
	nicknameTest := generateBasicNamingTest("nickname", http.MethodPut, "/nickname", validToken)
	t.Run("empty-nickname", func(t *testing.T) {
		nicknameTest(t, "", http.StatusOK)
	})
	t.Run("invalid-nickname", func(t *testing.T) {
		nicknameTest(t, "5test5", http.StatusForbidden)
	})
	t.Run("valid-nickname", func(t *testing.T) {
		nicknameTest(t, "test6", http.StatusAccepted)
	})
	t.Run("duplicate-nickname", func(t *testing.T) {
		nicknameTest(t, "test7", http.StatusAccepted)
		nicknameTest(t, "test7", http.StatusConflict)
	})
}

func TestCreateRoomHandler(t *testing.T) {
	t.Run("no-room-name", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/rooms", nil)
		if err != nil {
			t.Error(err)
		}
		addBearerAuthHeader(req, validToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		testStatusCode(t, rr, http.StatusBadRequest)
	})
	roomNameTest := generateBasicNamingTest("room", http.MethodPost, "/rooms", validToken)
	t.Run("empty-room-name", func(t *testing.T) {
		roomNameTest(t, "", http.StatusBadRequest)
	})
	t.Run("invalid-room-name", func(t *testing.T) {
		roomNameTest(t, "1test1", http.StatusForbidden)
	})
	t.Run("valid-room-name", func(t *testing.T) {
		roomNameTest(t, "test1", http.StatusCreated)
	})
	t.Run("duplicate-room-name", func(t *testing.T) {
		roomNameTest(t, "test2", http.StatusCreated)
		roomNameTest(t, "test2", http.StatusConflict)
	})
}

func TestSwitchRoomHandler(t *testing.T) {
	t.Run("no-room-name", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/rooms", nil)
		if err != nil {
			t.Error(err)
		}
		addBearerAuthHeader(req, validToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		testStatusCode(t, rr, http.StatusBadRequest)
	})
	roomNameTest := generateBasicNamingTest("room", http.MethodPut, "/rooms", validToken)
	t.Run("empty-room-name", func(t *testing.T) {
		roomNameTest(t, "", http.StatusBadRequest)
	})
	t.Run("invalid-room-name", func(t *testing.T) {
		roomNameTest(t, "6test6", http.StatusForbidden)
	})
	t.Run("existent-room", func(t *testing.T) {
		generateBasicNamingTest("room", http.MethodPost, "/rooms", validToken)(t, "test7", http.StatusCreated)
		roomNameTest(t, "test7", http.StatusAccepted)
	})
	t.Run("unexistent-room", func(t *testing.T) {
		roomNameTest(t, "test8", http.StatusNotFound)
	})
}

func TestUsersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	addBearerAuthHeader(req, validToken)
	rr := httptest.NewRecorder()
	APIRouter.ServeHTTP(rr, req)
	testStatusCode(t, rr, http.StatusOK)
	var usersList []server.User
	err := json.Unmarshal(rr.Body.Bytes(), &usersList)
	if err != nil {
		t.Error(err)
	}
	todoTest(t)
}

func TestExitHandler(t *testing.T) {
	todoTest(t)
}
