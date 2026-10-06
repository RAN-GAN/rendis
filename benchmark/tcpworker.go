package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// readReply consumes one RESP reply (simple, error, int or bulk string).
func readReply(r *bufio.Reader) error {
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	if len(line) < 3 {
		return fmt.Errorf("short reply")
	}
	switch line[0] {
	case '+', ':':
		return nil
	case '-':
		return fmt.Errorf("%s", strings.TrimSpace(line[1:]))
	case '$':
		n, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return err
		}
		if n < 0 {
			return nil
		}
		_, err = io.CopyN(io.Discard, r, int64(n+2))
		return err
	}
	return fmt.Errorf("invalid RESP reply")
}

func encodeCmd(args ...string) []byte {
	var sb strings.Builder
	sb.WriteString("*" + strconv.Itoa(len(args)) + "\r\n")
	for _, a := range args {
		sb.WriteString("$" + strconv.Itoa(len(a)) + "\r\n" + a + "\r\n")
	}
	return []byte(sb.String())
}

// runTCPWorker benchmarks the raw TCP backend (no WebSocket gateway).
func runTCPWorker(cfg Config, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := strings.TrimPrefix(cfg.URL, "tcp://")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		stats.Add("CONNECT_ERROR", 0, false)
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	gen := NewGenerator()
	deadline := time.Now().Add(cfg.Duration)
	val := gen.RandomString(cfg.ValueSize)

	for time.Now().Before(deadline) {
		op := cfg.Mode
		if op == "mixed" {
			op = []string{"ping", "set", "get"}[gen.rng.Intn(3)]
		}

		var payload []byte
		switch op {
		case "ping":
			payload = encodeCmd("PING")
		case "set":
			payload = encodeCmd("SET", gen.RandomKey(), val)
		default:
			payload = encodeCmd("GET", gen.RandomKey())
		}

		start := time.Now()
		_, err := conn.Write(payload)
		if err == nil {
			err = readReply(reader)
		}
		stats.Add(strings.ToUpper(op), time.Since(start), err == nil)
		if err != nil && (err == io.EOF || strings.Contains(err.Error(), "closed")) {
			return
		}
	}
}
