package main

import (
	"fmt"
	"net"
	"net/url"
)

func pickPeerAddress(peer clientSnapshot) string {
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
