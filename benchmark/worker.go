package main

import (
	"strings"
	"sync"
	"time"

	"github.com/RAN-GAN/rendis/client/golang"
)

func runWorker(cfg Config, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := rendis.New(cfg.URL, cfg.Key)
	if err != nil {
		stats.Add("CONNECT_ERROR", 0, false)
		return
	}
	defer client.Close()

	gen := NewGenerator()
	timeout := time.After(cfg.Duration)
	val := gen.RandomString(cfg.ValueSize)

	for {
		select {
		case <-timeout:
			return
		default:
			op := cfg.Mode
			if op == "mixed" {
				n := gen.rng.Intn(3)
				if n == 0 {
					op = "ping"
				} else if n == 1 {
					op = "set"
				} else {
					op = "get"
				}
			}

			start := time.Now()
			key := gen.RandomKey()
			success := true

			switch op {
			case "ping":
				err := client.Ping()
				if err != nil {
					success = false
				}
			case "set":
				err := client.Set(key, val)
				if err != nil {
					success = false
				}
			case "get":
				_, err := client.Get(key)
				// If the key is not found, the client returns an invalid RESP error because of how it decodes nil
				if err != nil && !strings.Contains(err.Error(), "invalid RESP message") {
					success = false
				}
			}

			stats.Add(strings.ToUpper(op), time.Since(start), success)
		}
	}
}
