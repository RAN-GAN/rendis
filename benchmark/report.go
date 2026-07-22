package main

import (
	"fmt"
	"io"
	"sort"
	"time"
)

func PrintReport(w io.Writer, cfg Config, stats *Stats, actualDuration time.Duration) {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	total := stats.Success + stats.Failed

	fmt.Fprintf(w, "\n=================================================\n")
	fmt.Fprintf(w, "Connections      : %d\n", cfg.Concurrency)
	fmt.Fprintf(w, "Target Duration  : %v\n", cfg.Duration)
	fmt.Fprintf(w, "Actual Duration  : %v\n", actualDuration)
	fmt.Fprintf(w, "=================================================\n")
	fmt.Fprintf(w, "\n--- Operations ---\n")

	for op, count := range stats.Ops {
		fmt.Fprintf(w, "%-10s %d\n", op, count)
	}

	fmt.Fprintf(w, "\nSuccess          : %d\n", stats.Success)
	fmt.Fprintf(w, "Failure          : %d\n", stats.Failed)

	if total > 0 {
		fmt.Fprintf(w, "Throughput       : %.2f ops/sec\n", float64(total)/actualDuration.Seconds())

		sort.Slice(stats.Latencies, func(i, j int) bool {
			return stats.Latencies[i] < stats.Latencies[j]
		})

		var sum time.Duration
		for _, l := range stats.Latencies {
			sum += l
		}
		avg := sum / time.Duration(total)

		median := stats.Latencies[total/2]
		p95 := stats.Latencies[int(float64(total)*0.95)]
		p99 := stats.Latencies[int(float64(total)*0.99)]
		min := stats.Latencies[0]
		max := stats.Latencies[total-1]

		fmt.Fprintf(w, "\n--- Latency ---\n")
		fmt.Fprintf(w, "Min              : %v\n", min)
		fmt.Fprintf(w, "Average          : %v\n", avg)
		fmt.Fprintf(w, "Median           : %v\n", median)
		fmt.Fprintf(w, "P95              : %v\n", p95)
		fmt.Fprintf(w, "P99              : %v\n", p99)
		fmt.Fprintf(w, "Max              : %v\n", max)
	} else {
		fmt.Fprintf(w, "No operations completed successfully.\n")
	}
	fmt.Fprintf(w, "=================================================\n")
}
