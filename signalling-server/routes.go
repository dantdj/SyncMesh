package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/ping", PingHandler)
	router.HandlerFunc(http.MethodPost, "/register", RegisterHandler)
	router.HandlerFunc(http.MethodPost, "/unregister", UnregisterHandler)
	router.HandlerFunc(http.MethodPost, "/discover", DiscoverHandler)

	return recoverPanic(router)
}
