package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

var (
	clients = make(map[string]string)
	mu      sync.Mutex
)

func RegisterClient(ip string) string {
	mu.Lock()
	defer mu.Unlock()

	b := make([]byte, 16)
	rand.Read(b)
	id := hex.EncodeToString(b)

	clients[id] = ip
	return id
}

func UnregisterClient(id string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, id)
}

func DiscoverClients() []string {
	mu.Lock()
	defer mu.Unlock()

	list := make([]string, 0, len(clients))
	for _, ip := range clients {
		list = append(list, ip)
	}
	return list
}
