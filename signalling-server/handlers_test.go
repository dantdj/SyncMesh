package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dantdj/syncmesh/api"
)

type registerResponsePayload struct {
	Status   string `json:"status"`
	ClientID string `json:"clientId"`
	Error    string `json:"error"`
}

type discoverResponsePayload struct {
	Status  string              `json:"status"`
	Clients []api.ClientSnapshot `json:"clients"`
	Error   string              `json:"error"`
}

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	recorder := httptest.NewRecorder()

	if err := PingHandler(recorder, req); err != nil {
		t.Fatalf("PingHandler returned error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload["status"] != "available" {
		t.Fatalf("expected status available, got %v", payload["status"])
	}

	systemInfo, ok := payload["systemInfo"].(map[string]any)
	if !ok {
		t.Fatalf("expected systemInfo object")
	}

	timestamp, ok := systemInfo["serverTimestamp"].(string)
	if !ok || timestamp == "" {
		t.Fatalf("expected serverTimestamp string")
	}

	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		t.Fatalf("expected RFC3339 timestamp, got %q", timestamp)
	}
}

func TestRegisterDiscoverUnregisterHandlers(t *testing.T) {
	resetClients()

	registerBody := []byte(`{"localIp":"192.168.1.10","localPort":4000}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(registerBody))
	req.RemoteAddr = "203.0.113.10:51234"
	recorder := httptest.NewRecorder()

	if err := RegisterHandler(recorder, req); err != nil {
		t.Fatalf("RegisterHandler returned error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var reg registerResponsePayload
	if err := json.NewDecoder(recorder.Body).Decode(&reg); err != nil {
		t.Fatalf("failed to decode register response: %v", err)
	}

	if reg.ClientID == "" {
		t.Fatal("expected clientId to be set")
	}

	discoverReq := httptest.NewRequest(http.MethodGet, "/discover", nil)
	discoverRecorder := httptest.NewRecorder()
	if err := DiscoverHandler(discoverRecorder, discoverReq); err != nil {
		t.Fatalf("DiscoverHandler returned error: %v", err)
	}

	if discoverRecorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", discoverRecorder.Code)
	}

	var discover discoverResponsePayload
	if err := json.NewDecoder(discoverRecorder.Body).Decode(&discover); err != nil {
		t.Fatalf("failed to decode discover response: %v", err)
	}

	if len(discover.Clients) != 1 {
		t.Fatalf("expected 1 client, got %d", len(discover.Clients))
	}

	client := discover.Clients[0]
	if client.ClientID != reg.ClientID {
		t.Fatalf("expected clientId %s, got %s", reg.ClientID, client.ClientID)
	}
	if client.PublicIP != "203.0.113.10" || client.PublicPort != 51234 {
		t.Fatalf("unexpected public info: %+v", client)
	}
	if client.LocalIP != "192.168.1.10" || client.LocalPort != 4000 {
		t.Fatalf("unexpected local info: %+v", client)
	}

	unregisterReq := httptest.NewRequest(http.MethodPost, "/unregister?clientId="+reg.ClientID, nil)
	unregisterRecorder := httptest.NewRecorder()
	if err := UnregisterHandler(unregisterRecorder, unregisterReq); err != nil {
		t.Fatalf("UnregisterHandler returned error: %v", err)
	}

	if unregisterRecorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", unregisterRecorder.Code)
	}

	confirmReq := httptest.NewRequest(http.MethodGet, "/discover", nil)
	confirmRecorder := httptest.NewRecorder()
	if err := DiscoverHandler(confirmRecorder, confirmReq); err != nil {
		t.Fatalf("DiscoverHandler returned error: %v", err)
	}

	var confirm discoverResponsePayload
	if err := json.NewDecoder(confirmRecorder.Body).Decode(&confirm); err != nil {
		t.Fatalf("failed to decode confirm response: %v", err)
	}

	if len(confirm.Clients) != 0 {
		t.Fatalf("expected 0 clients, got %d", len(confirm.Clients))
	}
}

func TestHeartbeatHandler(t *testing.T) {
	resetClients()

	id := RegisterClient("203.0.113.55", 5010, "192.168.1.55", 4055)

	req := httptest.NewRequest(http.MethodPost, "/heartbeat?clientId="+id, nil)
	recorder := httptest.NewRecorder()
	if err := HeartbeatHandler(recorder, req); err != nil {
		t.Fatalf("HeartbeatHandler returned error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	mu.Lock()
	info := clients[id]
	mu.Unlock()

	if time.Since(info.LastSeen) > time.Minute {
		t.Fatalf("expected LastSeen to be updated recently")
	}
}

func TestHeartbeatHandlerMissingClientID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/heartbeat", nil)
	recorder := httptest.NewRecorder()

	if err := HeartbeatHandler(recorder, req); err != nil {
		t.Fatalf("HeartbeatHandler returned error: %v", err)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestHeartbeatHandlerUnknownClient(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/heartbeat?clientId=missing", nil)
	recorder := httptest.NewRecorder()

	if err := HeartbeatHandler(recorder, req); err != nil {
		t.Fatalf("HeartbeatHandler returned error: %v", err)
	}

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}
