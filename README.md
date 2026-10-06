# Rendis

A Redis-compatible in-memory database written from scratch in Go, with a TCP-over-WebSocket deployment bridge for running on HTTP-only cloud platforms.

Rendis is a self-built Redis alternative designed to understand how databases like Redis work internally while solving a practical deployment problem:

> Redis requires raw TCP access, but many free cloud platforms only expose HTTP/WebSocket services.

Rendis provides both:

1. A Redis-compatible TCP server implementing RESP.
2. A WebSocket gateway that tunnels TCP traffic through HTTP-compatible infrastructure.

---

# Features

## Database Engine

* Custom TCP server written in Go
* Redis RESP protocol implementation
* In-memory key-value storage
* Concurrent client handling using goroutines
* Thread-safe storage using `sync.RWMutex`

Supported commands:

* `PING`
* `SET`
* `GET`
* `DEL`
* `EXISTS`
* `EXPIRE`
* `TTL`

---

## Expiration System

Rendis implements Redis-style key expiration.

Two expiration mechanisms are supported:

### Lazy Expiration

Expired keys are removed when accessed.

Examples:

```
GET key
TTL key
EXISTS key
```

If the key has expired, it behaves as if it does not exist.

---

### Active Expiration Worker

A background worker periodically scans the database and removes expired keys.

Example:

```
SET session abc
EXPIRE session 30
```

After expiration:

```
session -> deleted automatically
```

---

# Cloud Deployment Bridge

## Why?

Traditional Redis clients communicate using raw TCP:

```
redis-cli
      |
      |
     TCP
      |
      |
 Redis Server
```

However, many free hosting platforms expose only HTTP/WebSocket services.

Rendis solves this by adding a TCP tunnel layer.

---

# Architecture

```
                    Client Machine


                    Rendis Client

                         |
                         |

                 TCP-over-WebSocket
                         
                         |
                         |
                   WebSocket (WSS)

================================================

                 Cloud Deployment

================================================

                 WebSocket Gateway

                         |
                         |
                  localhost TCP

                         |
                         |

                 Rendis TCP Server

                         |
                         |

                In-Memory Storage Engine
```

---

# Project Structure

```
rendis/

├── server/
│   └── main.go
│
├── client/
│   ├── golang/
│   ├── javascript/
│   └── python/
│
├── internal/
│
│   ├── server/
│   │   ├── tcp.go
│   │   ├── client.go
│   │   └── handler.go

│   │
│   ├── protocol/
│   │   ├── reader.go
│   │   └── writer.go
│   │
│   ├── store/
│   │   ├── store.go
│   │   └── expiry.go
│   │
│   └── gateway/
│       ├── websocket.go
│       └── tunnel.go
│
├── go.mod
└── README.md
```

---

# How It Works

## Request Flow

Example:

```
SET name RAN-GAN
```

Flow:

```
Rendis Client

        |
        |
        v

RESP Encoder

        |
        |
        v

TCP Connection

        |
        |
        v

RESP Parser

        |
        |
        v

Command Handler

        |
        |
        v

Storage Engine

        |
        |
        v

RESP Response
```

---

# RESP Protocol

Rendis implements Redis Serialization Protocol.

Example command:

```
SET name RAN-GAN
```

RESP representation:

```
*3\r\n
$3\r\n
SET\r\n
$4\r\n
name\r\n
$7\r\n
RAN-GAN\r\n
```

The server parses the incoming bytes and returns RESP-compatible responses.

---

# Storage Engine

The database uses a Go map:

```
map[string]Entry
```

Example:

```
{
    "name": {
        Value: "RAN-GAN",
        Expiry: 12:30:00
    }
}
```

Access is protected using:

```go
sync.RWMutex
```

Allowing:

* Multiple simultaneous reads
* Safe writes

---

# Commands

## PING

Check server availability.

Request:

```
PING
```

Response:

```
PONG
```

---

## SET

Store a value.

Request:

```
SET name RAN-GAN
```

Response:

```
OK
```

---

## GET

Retrieve a value.

Request:

```
GET name
```

Response:

```
RAN-GAN
```

---

## DEL

Delete a key.

Request:

```
DEL name
```

Response:

```
(integer) 1
```

---

## EXISTS

Check if a key exists.

Request:

```
EXISTS name
```

Response:

```
(integer) 1
```

---

## EXPIRE

Set key lifetime.

Request:

```
EXPIRE name 60
```

Response:

```
(integer) 1
```

---

## TTL

Get remaining lifetime.

Request:

```
TTL name
```

Response:

```
(integer) 55
```

---

# Running Locally

## Requirements

* Go 1.20+

Clone:

```bash
git clone https://github.com/RAN-GAN/rendis.git

cd rendis/server
```

Run:

```bash
go run main.go
```

Server:

```
Server running on port 1708
```

---

# Connecting

You can connect to a Rendis server using either the official `rendis-cli` (recommended) or the standard `redis-cli` (for local TCP testing).

## Using rendis-cli (Recommended)

`rendis-cli` is an interactive, Redis-like command line interface designed specifically for Rendis. It communicates over WebSockets, so it works seamlessly with cloud deployments (like Render) that do not expose raw TCP ports.

### Installation

```bash
cd client/cli
go install
```

### Usage

Connect to a local server:
```bash
rendis-cli -h 127.0.0.1 -p 8080 -a "my-secure-key"
```

Connect to a remote cloud deployment via WebSocket URL:
```bash
rendis-cli -u wss://my-rendis-app.onrender.com -a "my-secure-key"
```

Example interaction:
```
wss://my-rendis-app.onrender.com> SET name RAN-GAN
OK
wss://my-rendis-app.onrender.com> GET name
"RAN-GAN"
wss://my-rendis-app.onrender.com> PING
"PONG"
```

---

## Using redis-cli (Local TCP testing only)

If you are running the server locally, you can also use the standard `redis-cli` to connect directly to the internal TCP server (bypassing the WebSocket gateway).

```bash
redis-cli -p 1708
```

Example:

```
127.0.0.1:1708> SET name RAN-GAN
OK

127.0.0.1:1708> GET name
"RAN-GAN"

127.0.0.1:1708> EXPIRE name 10
(integer) 1

127.0.0.1:1708> TTL name
(integer) 9
```


---

# Deployment

Rendis is designed for deployment on platforms that may not provide public TCP ports.

Deployment consists of:

```
Rendis Server
+
WebSocket Gateway
+
TCP Tunnel Client
```

The gateway exposes:

```
HTTP/WebSocket
```

while internally forwarding:

```
WebSocket
      |
      |
TCP
      |
      |
Rendis
```

This allows Redis-compatible clients to connect through TCP infrastructure when available.

---

# Security & Authentication

The WebSocket Gateway is protected by two security mechanisms configured via environment variables. When deploying to a platform like Render, Heroku, or AWS, simply set these as environment variables in your deployment dashboard:

1. **API Key Authentication**: The client must provide an `x-rendis-key` header that matches the `KEY` environment variable on the server.
2. **Origin Verification**: The server checks the `Origin` header against a comma-separated list of allowed origins defined in the `ALLOWED_ORIGINS` environment variable. You can use `*` to allow any origin.

**Example deployment variables (`.env` or dashboard):**

```env
KEY=my-secure-rendis-key
ALLOWED_ORIGINS=https://my-app.com,localhost
```

---

# Client Libraries (How to use in a project)

Rendis provides official client libraries that handle WebSocket tunneling and authentication automatically, making it easy to use Rendis in your projects.

Check out the individual client documentation for installation and API details:

* [Golang Client](client/golang/README.md)
* [Python Client](client/python/README.md)
* [JavaScript Client](client/javascript/README.md)

### Quick Example (Python)

```python
from rendis import Client

# Initialize with your deployed server URL and your secret key
client = Client("ws://your-rendis-deployment.onrender.com", "my-secure-rendis-key")

# Use standard Redis commands
client.set("my_key", "hello world")
print(client.get("my_key"))

client.close()
```

---

# Handling Cloud Restarts

Cloud platforms may restart services.

Rendis handles this by:

* Starting TCP server and gateway together.
* Creating fresh TCP connections per client.
* Avoiding permanent tunnel connections.
* Retrying backend connections when unavailable.

Startup flow:

```
Service starts

      |
      |

Start Rendis TCP Server

      |
      |

Start WebSocket Gateway

      |
      |

Accept client connections
```

---

# Benchmarks

Rendis ships a concurrent benchmark tool written in Go. It can target the WebSocket gateway (`ws://` / `wss://`) or the raw TCP backend (`tcp://`).

```bash
cd benchmark

# WebSocket gateway
go run . -url "ws://localhost:8080" -key "test" -c 50 -duration 10s -mode mixed

# Raw TCP backend (no gateway)
go run . -url "tcp://127.0.0.1:1708" -c 50 -duration 10s -mode mixed
```

The report also prints `Client CPU`, the cores the benchmark process itself kept busy, so you can tell whether the load generator is the limit.

## Local Benchmark

**Setup:** Intel i5-12450H, Arch Linux, Go 1.26. Server and load generator on the same machine. 50 connections, mixed GET/SET/PING, 32-byte values, 0 failures in every run. Figures are medians over 5 interleaved 6s runs with server stdout redirected to a file. Run-to-run spread was about +/-20%, so treat them as ballpark numbers.

| Variant | Throughput (ops/s) | Avg latency | P99 |
|---|---|---|---|
| Raw TCP, no gateway (reference) | ~250K | ~193 us | ~0.83 ms |
| WebSocket gateway, before tuning | 121K | 407 us | 1.52 ms |
| + remove hot-path `Println`s | 150K (+24%) | 329 us | 1.30 ms |
| + zero-allocation framing (current) | 143K (+18%) | 343 us | 1.29 ms |

With the server's stdout attached to a terminal (pty), the same fix is worth more, because every `Println` becomes a slow write:

| Variant | Throughput | Avg latency | P99 |
|---|---|---|---|
| Before | 61K | 813 us | 2.92 ms |
| Current | 119K (1.95x) | 412 us | 1.81 ms |

## Profiling the gateway

The gateway was profiled with pprof under load. Findings:

* **Debug prints on the hot path.** `fmt.Println` of every WebSocket payload and every parsed command took 16% of CPU. Removed.
* **Per-frame allocation.** gorilla's `ReadMessage` allocates a new buffer per frame (`io.ReadAll` was 58% of allocated bytes), and reply framing made several allocations per reply. The gateway now streams frames through a reused buffer (`NextReader`) and builds replies in one reused buffer. Gateway allocation in the profile run dropped from 442 MB to 6.5 MB. This did not change throughput on its own, because the workload is syscall-bound.
* **Not the cause.** There is one backend TCP dial per client connection (not per message) and one goroutine per direction per connection.
* **What remains.** About 62% of CPU is syscalls. Each request costs the gateway four (WebSocket read, TCP write, TCP read, WebSocket write). The remaining gap to raw TCP is mostly that, plus the WebSocket client's own overhead in the benchmark.

To profile your own deployment:

```bash
# local (binds to localhost only)
PPROF_ADDR=127.0.0.1:6060 go run .
go tool pprof -top "http://127.0.0.1:6060/debug/pprof/profile?seconds=10"

# deployed (requires the gateway key; off unless PPROF_ENABLED=1)
go tool pprof -top "https://<host>/debug/pprof/profile?seconds=10&key=<KEY>"
```

## Cloud-to-Cloud Benchmark (Render Free Tier)

Run from a benchmark service on Render against the Rendis service over `wss://`, 50 workers, 10s, mixed. Run it with `GET /run?mode=mixed&c=50&duration=10s&url=wss://<host>&key=<KEY>` on the benchmark service.

| Run | Throughput | Avg latency | Median | P99 |
|---|---|---|---|---|
| Typical result | 1.5K to 1.7K ops/s | 27 to 30 ms | 6 to 8 ms | 91 to 97 ms |

0 failures in every run. These numbers are **not a measure of the code**. A CPU profile of the deployed server under load showed:

* the benchmark client used about 0.15 cores, so it is not the limit;
* a single connection sees about 5 ms round trip, which is the network floor between the services;
* the server used about 1.04 s of CPU per 10 s (about 10% of one core), roughly 61 us of CPU per operation, with 78% of it in syscalls.

The free instance is capped at a fraction of a CPU (documented as about 0.1 vCPU). At 61 us per operation, 0.1 core allows about 1.6K ops/s, which matches what was measured. The P95/P99 of roughly 90 ms against a median of about 6 ms is consistent with CPU throttling. The cloud ceiling is the CPU quota, so the speedups measured locally are not expected to show up cleanly here, and a clean cloud before/after needs the old and new builds deployed side by side.

This still validates the `sync.RWMutex` thread safety and the stability of the TCP-to-WebSocket tunnel under sustained concurrent load.
---

# Development Roadmap

## Completed

* [x] TCP server
* [x] RESP parser
* [x] RESP response writer
* [x] In-memory storage
* [x] Thread-safe operations
* [x] GET / SET / DEL
* [x] EXISTS
* [x] TTL support
* [x] EXPIRE support
* [x] Active expiration worker
* [x] TCP-over-WebSocket gateway
* [x] Gateway authentication & origin verification
* [x] Persistence layer
* [x] RDB snapshots

---

## Upcoming

* [ ] AOF logging
* [ ] More Redis commands
* [ ] Unit tests
* [ ] Integration tests
* [ ] Docker support
* [ ] Cloud deployment automation

---

# Why Build Rendis?

Redis appears simple:

```
SET key value
GET key
```

but internally it involves:

* TCP networking
* Binary protocol parsing
* Concurrent data access
* Memory management
* Expiration algorithms
* Persistence
* Distributed deployment problems

Rendis is an exploration of these concepts by rebuilding the system from scratch.

---

