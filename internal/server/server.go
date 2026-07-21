package server

import (
	"fmt"
	"net"

	"github.com/RAN-GAN/rendis/internal/store"
)

func Start() {
	db := store.New()
	go db.StartExpiryWorker()
	listener, err := net.Listen("tcp", ":1708")
	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("Server running in port 1708")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleClient(conn, db)
	}

}
