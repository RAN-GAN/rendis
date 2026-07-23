package persistence

import (
	"fmt"
	"time"

	"github.com/RAN-GAN/rendis/server/internal/store"
)

func SnapShotWorker(db *store.Store, p Provider, interval time.Duration) {
	fmt.Println("Starting SnapShot Worker")

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		data, err := db.Snapshot()
		if err != nil {
			continue
		}

		if err := p.Save("rendisSnapshot", data); err != nil {
			fmt.Println("Snapshot failed:", err)
		}
	}
}
