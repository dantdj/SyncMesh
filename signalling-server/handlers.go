package main

import (
	"net"
	"net/http"
	"time"
)

func PingHandler(w http.ResponseWriter, r *http.Request) error {
	env := envelope{
		"status": "available",
		"systemInfo": map[string]string{
			"serverTimestamp": time.Now().Format(time.RFC3339),
		},
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If there is an error (e.g. missing port), use the address as is
		host = r.RemoteAddr
	}

	clientId := RegisterClient(host)
	env := envelope{
		"status":   "success",
		"clientId": clientId,
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}

func UnregisterHandler(w http.ResponseWriter, r *http.Request) error {
	UnregisterClient(r.URL.Query().Get("clientId"))

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
