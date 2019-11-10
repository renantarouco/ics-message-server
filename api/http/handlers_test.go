package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/renantarouco/ics-message-server/server"
)

func TestAuthHandler(t *testing.T) {
	t.Run("no-nickname", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/auth", nil)
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		Router.ServeHTTP(rr, req)
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
		nicknameTest(t, "nick2", http.StatusForbidden)
	})
	t.Run("valid-token", func(t *testing.T) {
		rr := nicknameTest(t, "nick3", http.StatusCreated)
		tokenStr, err := extractTokenStr(rr)
		fmt.Printf("valid<%s>", tokenStr)
		if err != nil {
			t.Error(err)
		}
		err = IsTokenValid(tokenStr)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestWsHandler(t *testing.T) {
	server := httptest.NewServer(Router)
	defer server.Close()
	ws1, _ := connectUser(t, server, "testUser11")
	defer ws1.Close()
	ws2, _ := connectUser(t, server, "testUser22")
	defer ws2.Close()
	t.Run("can-write", func(t *testing.T) {
		if err := ws1.WriteMessage(websocket.TextMessage, []byte("test")); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}
	})
	t.Run("can-receive", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if err := ws1.WriteMessage(websocket.TextMessage, []byte("test")); err != nil {
				t.Error(err)
			}
			_, buff, err := ws2.ReadMessage()
			if err != nil {
				t.Error(err)
			}
			message := string(buff)
			if message != "test" {
				t.Errorf("wrong message received, wanted %s got %s", "test", message)
			}
		}
	})
}

func TestNicknameHandler(t *testing.T) {
	t.Run("no-nickname", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/nickname", nil)
		if err != nil {
			t.Error(err)
		}
		addBearerAuthHeader(req, validToken)
		rr := httptest.NewRecorder()
		Router.ServeHTTP(rr, req)
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
		Router.ServeHTTP(rr, req)
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
		Router.ServeHTTP(rr, req)
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
	testServer := httptest.NewServer(Router)
	defer testServer.Close()
	ws1, _ := connectUser(t, testServer, "testUser1")
	defer ws1.Close()
	ws2, tokenStr2 := connectUser(t, testServer, "testUser2")
	defer ws2.Close()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	addBearerAuthHeader(req, tokenStr2)
	rr := httptest.NewRecorder()
	Router.ServeHTTP(rr, req)
	testStatusCode(t, rr, http.StatusOK)
	var usersList []server.User
	err := json.Unmarshal(rr.Body.Bytes(), &usersList)
	if err != nil {
		t.Error(err)
	}
}

func TestExitHandler(t *testing.T) {
	testServer := httptest.NewServer(Router)
	defer testServer.Close()
	ws, tokenStr := connectUser(t, testServer, "testUser33")
	defer ws.Close()
	req := httptest.NewRequest(http.MethodGet, "/exit", nil)
	addBearerAuthHeader(req, tokenStr)
	rr := httptest.NewRecorder()
	Router.ServeHTTP(rr, req)
	testStatusCode(t, rr, http.StatusOK)
}
