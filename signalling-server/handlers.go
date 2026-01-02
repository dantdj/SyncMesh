package main

import (
	"log/slog"
	"net/http"
	"time"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"serverTimestamp": time.Now().Format(time.RFC3339),
		},
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		slog.Error("Failed to return service info", slog.String("error", err.Error()))
		serverErrorResponse(w)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "success",
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		slog.Error("Failed to return service info", slog.String("error", err.Error()))
		serverErrorResponse(w)
	}
}

func UnregisterHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "success",
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		slog.Error("Failed to return service info", slog.String("error", err.Error()))
		serverErrorResponse(w)
	}
}

func DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "success",
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		slog.Error("Failed to return service info", slog.String("error", err.Error()))
		serverErrorResponse(w)
	}
}
