package main

import (
	"log"
	"net"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	wsURL := os.Getenv("WS_URL")
	if wsURL == "" {
		wsURL = "ws://localhost:8080/connect"
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "127.0.0.1:6379" 
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("Tunnel client listening on %s, proxying to %s\n", listenAddr, wsURL)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go handleTunnelConnection(conn, wsURL)
	}
}

func handleTunnelConnection(tcpConn net.Conn, wsURL string) {
	defer tcpConn.Close()

	log.Printf("New local connection from %s\n", tcpConn.RemoteAddr())

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Println("websocket dial error:", err)
		return
	}
	defer wsConn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			n, err := tcpConn.Read(buf)
			if err != nil {
				return
			}
			err = wsConn.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				return
			}
		}
	}()

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		_, err = tcpConn.Write(message)
		if err != nil {
			break
		}
	}

	<-done
	log.Printf("Closed local connection from %s\n", tcpConn.RemoteAddr())
}
