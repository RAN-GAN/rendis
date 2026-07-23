package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/RAN-GAN/rendis/server/internal/gateway"
	"github.com/RAN-GAN/rendis/server/internal/persistence"
	"github.com/RAN-GAN/rendis/server/internal/server"
	"github.com/RAN-GAN/rendis/server/internal/store"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	kv := store.New()
	go kv.StartExpiryWorker()

	backend := os.Getenv("BACKEND_ADDR")
	if backend == "" {
		backend = "127.0.0.1:1708"
	}

	go server.Start(kv, backend)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	db, err := persistence.NewProvider()
	if err != nil {
		log.Println("Persistence disabled:", err)
	} else {
		restorePoint, err := db.Load("rendisSnapshot")
		if err != nil {
			fmt.Println(err)
		} else {
			snapShotContintion := kv.Restore(restorePoint)
			if snapShotContintion != nil {
				fmt.Println("Unable to restore snapshot", snapShotContintion)
			} else {
				fmt.Println("Snapshot restored successfully")
			}
		}

	}

	interval, err := strconv.Atoi(os.Getenv("SNAPSHOT_INTERVAL"))
	if err != nil {
		log.Println("Invalid snapshot interval, using default of 300 seconds")
		interval = 300
	}
	go persistence.SnapShotWorker(kv, db, time.Duration(interval)*time.Second)
	gateway.Start(gateway.Config{
		ListenAddr:  ":" + port,
		BackendAddr: backend,
	})

	select {}
}
