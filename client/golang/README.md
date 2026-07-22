# Rendis Go Client

A Go client for Rendis, an in-memory datastore.

## Installation

```bash
go get github.com/RAN-GAN/rendis/client/golang
```
*(Update the package path to match your actual Go module path when publishing)*

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/RAN-GAN/rendis/client/golang"
)

func main() {
	client, err := rendis.New("ws://localhost:8080", "your-auth-key")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	// Ping the server
	err = client.Ping()

	// Set a key
	err = client.Set("mykey", "myvalue")
	
	// Get a key
	val, err := client.Get("mykey")
	fmt.Printf("Got value: %s\n", val)
    
	// Check if a key exists
	exists, err := client.Exists("mykey")
	
	// Delete a key
	deleted, err := client.Del("mykey")
}
```

## API Methods
- `New(url, key string) (*Client, error)`
- `Ping() error`
- `Set(key, value string) error`
- `Get(key string) (string, error)`
- `Del(key string) (int64, error)`
- `TTL(key string) (int64, error)`
- `Expire(key string, seconds int64) (bool, error)`
- `Exists(key string) (bool, error)`
- `Close() error`
