package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dantdj/syncmesh/api"
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

func acceptLoop(logger *log.Logger, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Printf("accept error: %v", err)
			return
		}
		go handleConn(logger, conn)
	}
}

func handleConn(logger *log.Logger, conn net.Conn) {
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		logger.Printf("read error: %v", err)
		return
	}
	if n > 0 {
		logger.Printf("received: %q from %s", string(bytes.TrimSpace(buf[:n])), conn.RemoteAddr().String())
	}

	_, _ = conn.Write([]byte("hello from peer\n"))
}

func register(logger *log.Logger, baseURL, localIP string, localPort int) (string, error) {
	body, err := json.Marshal(api.RegisterRequest{
		LocalIP:   localIP,
		LocalPort: localPort,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/register", baseURL), bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("register failed: status %s", resp.Status)
	}

	var payload api.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	if payload.ClientID == "" {
		return "", fmt.Errorf("register failed: %s", payload.Error)
	}

	return payload.ClientID, nil
}

func connectToPeer(logger *log.Logger, baseURL, selfID string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/discover", baseURL), nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discover failed: status %s", resp.Status)
	}

	var payload api.DiscoverResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}

	for _, peer := range payload.Clients {
		if peer.ClientID == selfID {
			continue
		}

		addr := pickPeerAddress(peer)
		if addr == "" {
			continue
		}

		logger.Printf("attempting connection to %s (%s)", peer.ClientID, addr)
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			logger.Printf("connect failed to %s: %v", addr, err)
			continue
		}

		_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
		_, _ = conn.Write([]byte("hello from client\n"))

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			logger.Printf("read error: %v", err)
		} else if n > 0 {
			logger.Printf("received: %q", string(bytes.TrimSpace(buf[:n])))
		}
		_ = conn.Close()
		return nil
	}

	return fmt.Errorf("no other clients discovered")
}

func heartbeatLoop(logger *log.Logger, baseURL, clientID string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := sendHeartbeat(baseURL, clientID); err != nil {
			logger.Printf("heartbeat failed: %v", err)
		}
	}
}

func sendHeartbeat(baseURL, clientID string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/heartbeat?clientId=%s", baseURL, url.QueryEscape(clientID)), nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed: status %s", resp.Status)
	}

	return nil
}

func pickPeerAddress(peer api.ClientSnapshot) string {
	if peer.LocalIP != "" && peer.LocalPort != 0 {
		return fmt.Sprintf("%s:%d", peer.LocalIP, peer.LocalPort)
	}
	if peer.PublicIP != "" && peer.PublicPort != 0 {
		return fmt.Sprintf("%s:%d", peer.PublicIP, peer.PublicPort)
	}
	return ""
}

func detectLocalIP(baseURL string) string {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	host := parsed.Host
	if host == "" {
		return ""
	}

	if !hasPort(host) {
		if parsed.Scheme == "https" {
			host = net.JoinHostPort(host, "443")
		} else {
			host = net.JoinHostPort(host, "80")
		}
	}

	conn, err := net.Dial("udp", host)
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func hasPort(host string) bool {
	_, _, err := net.SplitHostPort(host)
	return err == nil
}
