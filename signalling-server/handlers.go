package main

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
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
	var req struct {
		LocalIP   string `json:"localIp"`
		LocalPort int    `json:"localPort"`
	}

	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			errorResponse(w, http.StatusBadRequest, "invalid JSON body")
			return nil
		}
	}

	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If there is an error (e.g. missing port), use the address as is
		host = r.RemoteAddr
		port = "0"
	}

	publicPort, err := strconv.Atoi(port)
	if err != nil {
		publicPort = 0
	}

	clientId := RegisterClient(host, publicPort, req.LocalIP, req.LocalPort)
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

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) error {
	clientId := r.URL.Query().Get("clientId")
	if clientId == "" {
		errorResponse(w, http.StatusBadRequest, "clientId is required")
		return nil
	}

	if !TouchClient(clientId) {
		errorResponse(w, http.StatusNotFound, "client not found")
		return nil
	}

	env := envelope{
		"status": "success",
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}

func DiscoverHandler(w http.ResponseWriter, r *http.Request) error {
	type clientSnapshot struct {
		ClientID   string `json:"clientId"`
		PublicIP   string `json:"publicIp"`
		PublicPort int    `json:"publicPort"`
		LocalIP    string `json:"localIp,omitempty"`
		LocalPort  int    `json:"localPort,omitempty"`
	}

	clients := DiscoverClients()
	snapshots := make([]clientSnapshot, 0, len(clients))
	for id, info := range clients {
		snapshots = append(snapshots, clientSnapshot{
			ClientID:   id,
			PublicIP:   info.PublicIP,
			PublicPort: info.PublicPort,
			LocalIP:    info.LocalIP,
			LocalPort:  info.LocalPort,
		})
	}

	env := envelope{
		"status":  "success",
		"clients": snapshots,
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		return err
	}

	return nil
}
