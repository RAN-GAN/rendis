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


                 redis-cli / Redis Client

                         |
                         |
                    RESP over TCP

                         |
                         |

                 Local TCP Tunnel

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

├── cmd/
│   └── rendis/
│       └── main.go
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
SET name Pradeep
```

Flow:

```
Redis Client

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
SET name Pradeep
```

RESP representation:

```
*3\r\n
$3\r\n
SET\r\n
$4\r\n
name\r\n
$7\r\n
Pradeep\r\n
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
        Value: "Pradeep",
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
SET name Pradeep
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
Pradeep
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

cd rendis
```

Run:

```bash
go run .
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
127.0.0.1:1708> SET name Pradeep
OK

127.0.0.1:1708> GET name
"Pradeep"

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

This allows Redis-compatible clients to connect through HTTP-only infrastructure.

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

---

## Upcoming

* [ ] TCP-over-WebSocket gateway
* [ ] Local tunnel client
* [ ] Persistence layer
* [ ] RDB snapshots
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

