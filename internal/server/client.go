package server

import (
	"bufio"
	"net"

	"github.com/RAN-GAN/rendis/internal/protocol"
	"github.com/RAN-GAN/rendis/internal/store"
)

func handleClient(conn net.Conn, db *store.Store) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		parts, err := protocol.ReadRESP(reader)
		if err != nil {

			return
		}

		response := handleMessage(parts, db)

		_, err = conn.Write([]byte(response))
		if err != nil {
			return
		}
	}
}
