package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func runBenchmark(cfg Config) *Stats {
	stats := NewStats()
	var wg sync.WaitGroup

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go runWorker(cfg, stats, &wg)
	}

	wg.Wait()
	return stats
}

func main() {
	cfg := LoadConfig()

	port := os.Getenv("PORT")
	if port != "" {
		// Run as a Web Service on Render
		fmt.Printf("Starting benchmark HTTP server on port %s...\n", port)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Received benchmark request from %s\n", r.RemoteAddr)
			
			// Override config with env vars if set (useful for Render deployment)
			envUrl := os.Getenv("RENDIS_URL")
			envKey := os.Getenv("RENDIS_KEY")
			if envUrl != "" { cfg.URL = envUrl }
			if envKey != "" { cfg.Key = envKey }
			
			start := time.Now()
			stats := runBenchmark(cfg)
			actualDuration := time.Since(start)
			
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Starting comprehensive benchmark...\n")
			fmt.Fprintf(w, "Target: %s (Mode: %s)\n", cfg.URL, cfg.Mode)
			fmt.Fprintf(w, "Running for %v with %d workers...\n", cfg.Duration, cfg.Concurrency)
			
			PrintReport(w, cfg, stats, actualDuration)
		})
		
		http.ListenAndServe(":"+port, nil)
	} else {
		// Run as a standard CLI
		fmt.Printf("Starting comprehensive benchmark...\n")
		fmt.Printf("Target: %s (Mode: %s)\n", cfg.URL, cfg.Mode)
		fmt.Printf("Running for %v with %d workers...\n", cfg.Duration, cfg.Concurrency)

		start := time.Now()
		stats := runBenchmark(cfg)
		actualDuration := time.Since(start)

		PrintReport(os.Stdout, cfg, stats, actualDuration)
	}
}
