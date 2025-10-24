package middleware

import (
	"log"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Default().Printf("Hello from Auth middleware")
		next.ServeHTTP(w, r)
	})
}