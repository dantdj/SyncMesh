package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func register(logger *log.Logger, baseURL, localIP string, localPort int) (string, error) {
	body, err := json.Marshal(registerRequest{
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

	var payload registerResponse
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

	var payload discoverResponse
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
