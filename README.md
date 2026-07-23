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

Using redis-cli:

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

Rendis includes a custom, highly-concurrent WebSocket benchmark tool written in Go to test performance.

```bash
cd benchmark
go run . -url "ws://localhost:8080" -key "test" -c 50 -duration 10s -mode mixed
```

## Rendis Benchmark (Local)

**Hardware**
* **CPU:** 12th Gen Intel(R) Core(TM) i5-12450H
* **OS:** Arch Linux
* **Go Version:** 1.26

**Benchmark**
* **Workers:** 50
* **Duration:** 10s

**Operations**
* **GET:** 48,211
* **SET:** 48,092
* **PING:** 47,713

**Throughput**
* 14,379 ops/sec

**Latency**
* **Average:** 3.47 ms
* **Median:** 2.10 ms
* **P95:** 11.64 ms
* **P99:** 17.85 ms

## Cloud-to-Cloud Benchmark (Render Free Tier)

This benchmark was run from a dedicated benchmark service deployed on Render, communicating over WebSockets with the Rendis server deployed in the same region.

**Benchmark**
* **Workers:** 50
* **Duration:** 30s

**Operations**
* **GET:** 16,199
* **SET:** 16,282
* **PING:** 16,224

**Throughput**
* 1,571.59 ops/sec (48,705 total operations)
* **Failures:** 0 (0% error rate)

**Latency**
* **Average:** 30.84 ms
* **Median:** 6.24 ms
* **P95:** 90.34 ms
* **P99:** 94.29 ms
* **Max:** 305.87 ms

This validates the robustness of the `sync.RWMutex` thread safety and the stability of the TCP-to-WebSocket tunnel under sustained concurrent load.

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

