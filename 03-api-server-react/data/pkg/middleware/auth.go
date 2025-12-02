package middleware

import (
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple token validation (for demonstration)
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Printf("Request to %s without auth token", r.URL.Path)
		} else {
			log.Printf("Request to %s with token", r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}
