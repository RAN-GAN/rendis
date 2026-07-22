package main

import (
	"flag"
	"time"
)

type Config struct {
	URL         string
	Key         string
	Concurrency int
	Duration    time.Duration
	ValueSize   int
	Mode        string
}

func LoadConfig() Config {
	cfg := Config{}

	flag.StringVar(&cfg.URL, "url", "ws://localhost:8080", "Rendis URL")
	flag.StringVar(&cfg.Key, "key", "test", "API Key")

	flag.IntVar(&cfg.Concurrency, "c", 50, "Concurrent clients")

	flag.DurationVar(
		&cfg.Duration,
		"duration",
		30*time.Second,
		"Benchmark duration",
	)

	flag.IntVar(
		&cfg.ValueSize,
		"size",
		32,
		"Value size in bytes",
	)

	flag.StringVar(
		&cfg.Mode,
		"mode",
		"mixed",
		"ping|get|set|mixed",
	)

	flag.Parse()

	return cfg
}
