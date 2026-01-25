package main

import (
	"net/http"
	"time"
)

func PingHandler(w http.ResponseWriter, r *http.Request) error {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"serverTimestamp": time.Now().Format(time.RFC3339),
		},
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	RegisterClient(r.URL.Query().Get("client"))
	env := envelope{
		"status": "success",
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}

func UnregisterHandler(w http.ResponseWriter, r *http.Request) error {
	UnregisterClient(r.URL.Query().Get("client"))

	env := envelope{
		"status": "success",
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}

func DiscoverHandler(w http.ResponseWriter, r *http.Request) error {
	env := envelope{
		"status":  "success",
		"clients": DiscoverClients(),
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}
