package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

const (
	validToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDA5MDg4MDAsImlhdCI6MTU3MzEzMTE2NSwiaXNzIjoibmljazMiLCJzdWIiOiJ1bmFtZWQifQ.pN9vpFTeWbygia1l_7UK7HSSAPZrlRcZ-iGVymyGc5U"
	invalidToken = "eyJhbGcikpXVCJ9.eyJpYXQiOjE1NzI3ODY0MDAsImlzcyI6Im5pY2szIiwic3ViIjoidW5hbWVkIn0.yJ-vyVz_5uJW4xyxNyA4TI00K98GfE4EkgYxBi8-w3c"
)

func todoTest(t *testing.T) {
	t.Error("this test is not yet implemented. DO IT!")
}

func testStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expectedStatusCode int) {
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("wrong code, wanted %v got %v", expectedStatusCode, status)
	}
}

func addBearerAuthHeader(r *http.Request, tokenStr string) {
	bearerToken := fmt.Sprintf("Bearer %s", "")
	r.Header.Set("Authorization", bearerToken)
}

func generateBasicNamingTest(fieldName, httpMethod, route, tokenStr string) func(*testing.T, string, int) *httptest.ResponseRecorder {
	return func(t *testing.T, name string, expectedStatusCode int) *httptest.ResponseRecorder {
		formData := url.Values{}
		formData.Set(fieldName, name)
		req, err := http.NewRequest(httpMethod, route, strings.NewReader(formData.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
		if tokenStr != "" {
			addBearerAuthHeader(req, tokenStr)
		}
		rr := httptest.NewRecorder()
		Router.ServeHTTP(rr, req)
		testStatusCode(t, rr, expectedStatusCode)
		return rr
	}
}

func extractTokenStr(rr *httptest.ResponseRecorder) (string, error) {
	var bodyMap map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &bodyMap)
	if err != nil {
		return "", err
	}
	tokenStr, ok := bodyMap["token"]
	if !ok {
		return "", errors.New("body doesn't contain 'token' key")
	}
	return tokenStr, nil
}

func connectUser(t *testing.T, server *httptest.Server, nickname string) (*websocket.Conn, string) {
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	nicknameTest := generateBasicNamingTest("nickname", http.MethodPost, "/auth", "")
	rr1 := nicknameTest(t, nickname, http.StatusCreated)
	tokenStr, err := extractTokenStr(rr1)
	if err != nil {
		t.Error(err)
	}
	reqHeader1 := http.Header{}
	reqHeader1.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, reqHeader1)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	return ws, tokenStr
}
