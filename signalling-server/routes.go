package main

import (
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/ping", handle(PingHandler))
	router.HandlerFunc(http.MethodGet, "/discover", handle(DiscoverHandler))
	router.HandlerFunc(http.MethodPost, "/register", handle(RegisterHandler))
	router.HandlerFunc(http.MethodPost, "/unregister", handle(UnregisterHandler))
	router.HandlerFunc(http.MethodPost, "/heartbeat", handle(HeartbeatHandler))

	return recoverPanic(router)
}

// handle provides a common wrapper for all handlers, allowing for
// consistent error handling and logging.
func handle(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			slog.Error("Handler execution failed", slog.String("error", err.Error()))
			serverErrorResponse(w)
		}
	}
}
