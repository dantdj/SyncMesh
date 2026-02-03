package main

import (
	"testing"
	"time"
)

// resetClients clears the global clients map for testing purposes.
// This helper accesses internal state to ensure clean test runs.
func resetClients() {
	mu.Lock()
	defer mu.Unlock()
	clients = make(map[string]clientInfo)
}

func TestRegisterClient(t *testing.T) {
	resetClients()

	id := RegisterClient("203.0.113.5", 5000, "192.168.1.5", 4000)

	if id == "" {
		t.Fatal("Expected a non-empty ID returned from RegisterClient")
	}

	clients := DiscoverClients()
	info, ok := clients[id]
	if !ok {
		t.Fatalf("Expected client %s to be registered", id)
	}

	if info.PublicIP != "203.0.113.5" || info.PublicPort != 5000 {
		t.Errorf("Unexpected public info: %+v", info)
	}

	if info.LocalIP != "192.168.1.5" || info.LocalPort != 4000 {
		t.Errorf("Unexpected local info: %+v", info)
	}
}

func TestUnregisterClient(t *testing.T) {
	resetClients()

	id := RegisterClient("203.0.113.6", 5001, "", 0)
	UnregisterClient(id)

	clients := DiscoverClients()
	if _, ok := clients[id]; ok {
		t.Error("Expected client to be unregistered")
	}
}

func TestTouchClientUpdatesLastSeen(t *testing.T) {
	resetClients()

	id := RegisterClient("203.0.113.7", 5002, "", 0)

	mu.Lock()
	info := clients[id]
	info.LastSeen = time.Now().UTC().Add(-1 * time.Minute)
	clients[id] = info
	mu.Unlock()

	if !TouchClient(id) {
		t.Fatal("Expected TouchClient to return true for existing client")
	}

	mu.Lock()
	updated := clients[id]
	mu.Unlock()

	if time.Since(updated.LastSeen) > time.Minute {
		t.Error("Expected LastSeen to be updated recently")
	}
}

func TestDiscoverClientsPrunesExpired(t *testing.T) {
	resetClients()

	previousTTL := clientTTL
	clientTTL = time.Minute
	t.Cleanup(func() { clientTTL = previousTTL })

	id := RegisterClient("203.0.113.8", 5003, "", 0)

	mu.Lock()
	info := clients[id]
	info.LastSeen = time.Now().UTC().Add(-2 * time.Minute)
	clients[id] = info
	mu.Unlock()

	discovered := DiscoverClients()
	if _, ok := discovered[id]; ok {
		t.Error("Expected expired client to be pruned")
	}
}
