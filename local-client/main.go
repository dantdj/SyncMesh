package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	serverURL := flag.String("server", "http://localhost:8089", "signalling server base URL")
	listenPort := flag.Int("listen", 4000, "local TCP listen port")
	flag.Parse()

	logger := log.New(os.Stdout, "client: ", log.LstdFlags)

	localIP := detectLocalIP(*serverURL)
	if localIP == "" {
		localIP = "127.0.0.1"
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *listenPort))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	logger.Printf("listening on %s (local IP: %s)", listener.Addr().String(), localIP)

	go acceptLoop(logger, listener)

	clientID, err := register(logger, *serverURL, localIP, *listenPort)
	if err != nil {
		logger.Fatalf("register failed: %v", err)
	}
	logger.Printf("registered with clientId=%s", clientID)

	go heartbeatLoop(logger, *serverURL, clientID, 30*time.Second)

	time.Sleep(500 * time.Millisecond)

	if err := connectToPeer(logger, *serverURL, clientID); err != nil {
		logger.Printf("no peer connection made: %v", err)
	}

	select {}
}
