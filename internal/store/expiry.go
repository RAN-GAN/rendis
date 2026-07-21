package store

import (
	"fmt"
	"time"
)

func (s *Store) StartExpiryWorker() {
	fmt.Println("Starting Expiration Worker")
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Println("Clearing expired keys")
		s.clearExpired()
	}
}
