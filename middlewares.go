package main

import (
	"context"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ValidateTokenMiddleware - Checks validity and extracts a token in the
// Authorization header
func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excludeRoutes := []string{"/auth"}
		requestPath := r.URL.Path
		for _, route := range excludeRoutes {
			if requestPath == route {
				next.ServeHTTP(w, r)
				return
			}
		}
		authHeader := r.Header.Get("Sec-WebSocket-Protocol")
		log.Debugf("authheader<%s>", authHeader)
		tokenStr := strings.TrimPrefix(authHeader, "Bearer")
		tokenStr = strings.TrimLeft(tokenStr, " ")
		log.Debugf("token<%s>", tokenStr)
		err := IsTokenValid(tokenStr)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "tokenStr", tokenStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
