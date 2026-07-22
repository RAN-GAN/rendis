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
		// Render health check endpoint
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "Benchmark service is healthy! Visit /run to execute a benchmark.")
				return
			}
			http.NotFound(w, r)
		})

		// Benchmark execution endpoint
		http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Received benchmark request from %s\n", r.RemoteAddr)
			
			// Override config with env vars
			if envUrl := os.Getenv("RENDIS_URL"); envUrl != "" { cfg.URL = envUrl }
			if envKey := os.Getenv("RENDIS_KEY"); envKey != "" { cfg.Key = envKey }
			
			// Override config with query parameters
			q := r.URL.Query()
			if u := q.Get("url"); u != "" { cfg.URL = u }
			if k := q.Get("key"); k != "" { cfg.Key = k }
			if m := q.Get("mode"); m != "" { cfg.Mode = m }
			if c := q.Get("c"); c != "" { fmt.Sscanf(c, "%d", &cfg.Concurrency) }
			if d := q.Get("duration"); d != "" {
				if parsedDur, err := time.ParseDuration(d); err == nil {
					cfg.Duration = parsedDur
				}
			}
			
			// Write headers immediately to prevent Bad Gateway / Timeout from Render Load Balancers
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Starting comprehensive benchmark...\n")
			fmt.Fprintf(w, "Target: %s (Mode: %s)\n", cfg.URL, cfg.Mode)
			fmt.Fprintf(w, "Running for %v with %d workers...\n", cfg.Duration, cfg.Concurrency)
			fmt.Fprintf(w, "Please wait... (do not refresh)\n\n")
			
			// Flush the buffer so the client sees the immediate response
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			
			start := time.Now()
			stats := runBenchmark(cfg)
			actualDuration := time.Since(start)
			
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
