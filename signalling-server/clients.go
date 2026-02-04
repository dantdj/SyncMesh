package main

import (
	"crypto/rand"
	"encoding/hex"
	"maps"
	"sync"
	"time"
)

var (
	clients   = make(map[string]clientInfo)
	mu        sync.Mutex
	clientTTL = 5 * time.Minute
)

type clientInfo struct {
	PublicIP   string
	PublicPort int
	LocalIP    string
	LocalPort  int
	LastSeen   time.Time
}

func RegisterClient(publicIP string, publicPort int, localIP string, localPort int) string {
	mu.Lock()
	defer mu.Unlock()

	pruneExpiredLocked()

	b := make([]byte, 16)
	rand.Read(b)
	id := hex.EncodeToString(b)

	clients[id] = clientInfo{
		PublicIP:   publicIP,
		PublicPort: publicPort,
		LocalIP:    localIP,
		LocalPort:  localPort,
		LastSeen:   time.Now().UTC(),
	}
	return id
}

func UnregisterClient(id string) {
	mu.Lock()
	defer mu.Unlock()
	pruneExpiredLocked()
	delete(clients, id)
}

func DiscoverClients() map[string]clientInfo {
	mu.Lock()
	defer mu.Unlock()

	pruneExpiredLocked()

	copy := make(map[string]clientInfo, len(clients))
	maps.Copy(copy, clients)
	return copy
}

func TouchClient(id string) bool {
	mu.Lock()
	defer mu.Unlock()

	pruneExpiredLocked()

	info, ok := clients[id]
	if !ok {
		return false
	}
	info.LastSeen = time.Now().UTC()
	clients[id] = info
	return true
}

func pruneExpiredLocked() {
	if len(clients) == 0 {
		return
	}

	cutoff := time.Now().UTC().Add(-clientTTL)
	for id, info := range clients {
		if info.LastSeen.Before(cutoff) {
			delete(clients, id)
		}
	}
}
