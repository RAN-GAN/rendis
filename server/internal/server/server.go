package server

import (
	"fmt"
	"net"

	"github.com/RAN-GAN/rendis/server/internal/store"
)

func Start(db *store.Store, addr string) {

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("Rendis TCP server running on :1708")

	for {

		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go handleClient(conn, db)
	}
}
