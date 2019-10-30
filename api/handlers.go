package api

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// JoinHandler - The handler for /join endpoint. Receives via POST an URL
// encoded nickname field and returns a JSON object with the nickname validated
// and a token for future requests.
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nickname := r.FormValue("nickname")
	if nickname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if match, _ := regexp.MatchString("^[a-zA-Z | _][a-zA-Z0-9 | _ | -]*$", nickname); !match {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	responseToken := map[string]interface{}{
		"nickname": nickname,
		"token":    "a09df0unfoijsa-09enf",
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseToken)
}
