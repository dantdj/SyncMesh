package main

var clients = []string{}

func RegisterClient(client string) {
	clients = append(clients, client)
}

func UnregisterClient(client string) {
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func DiscoverClients() []string {
	return clients
}
