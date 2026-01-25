package main

import (
	"net/http"
)

// recoverPanic recovers from any panics and sends a 500 Internal Server Error
// response to the client.
func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				serverErrorResponse(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
