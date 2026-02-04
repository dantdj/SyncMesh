package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"time"
)

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
