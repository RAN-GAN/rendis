package main

import (
	"os"

	"github.com/RAN-GAN/rendis/server/internal/gateway"
	"github.com/RAN-GAN/rendis/server/internal/server"
	"github.com/RAN-GAN/rendis/server/internal/store"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	db := store.New()

	go db.StartExpiryWorker()

	backend := os.Getenv("BACKEND_ADDR")
	if backend == "" {
		backend = "127.0.0.1:1708"
	}

	go server.Start(db, backend)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	gateway.Start(gateway.Config{
		ListenAddr:  ":" + port,
		BackendAddr: backend,
	})

	select {}
}
