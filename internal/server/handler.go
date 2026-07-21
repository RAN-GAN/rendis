package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RAN-GAN/rendis/internal/protocol"
	"github.com/RAN-GAN/rendis/internal/store"
)

func handleMessage(parts []string, db *store.Store) string {
	fmt.Println("handling Message: ", parts)

	if len(parts) == 0 {
		return protocol.Error("empty command")
	}

	command := strings.ToUpper(parts[0])
	switch command {

	case "PING":
		if len(parts) != 1 {
			return protocol.Error("wrong number of arguments")

		}
		return protocol.SimpleString("PONG")

	case "SET":
		if len(parts) != 3 {
			return protocol.Error("wrong number of arguments")
		}
		db.Set(parts[1], parts[2])
		return protocol.SimpleString("OK")

	case "GET":
		if len(parts) != 2 {
			return protocol.Error("wrong number of arguments")
		}
		value, ok := db.Get(parts[1])
		if !ok {
			return protocol.NullBulkString()
		}
		return protocol.BulkString(value)

	case "DEL":
		if len(parts) != 2 {
			return protocol.Error("wrong number of arguments")
		}
		ok := db.Del(parts[1])
		if !ok {
			return protocol.Integer(0)
		}
		return protocol.Integer(1)

	case "TTL":
		if len(parts) != 2 {
			return protocol.Error("wrong number of arguments")
		}

		ttl, _ := db.TTL(parts[1])
		return protocol.Integer(ttl)

	case "EXPIRE":
		if len(parts) != 3 {
			return protocol.Error("wrong number of arguments")
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return protocol.Error("invalid expire time")
		}

		if db.Expire(parts[1], seconds) {
			return protocol.Integer(1)
		}
		return protocol.Integer(0)

	case "EXISTS":
		if len(parts) != 2 {
			return protocol.Error("wrong number of arguments")
		}

		_, ok := db.Get(parts[1])

		if ok {
			return protocol.Integer(1)
		}
		return protocol.Integer(0)

	default:
		return protocol.Error("unknown command")
	}
}
