package server

import (
	"strconv"
	"strings"

	"github.com/RAN-GAN/rendis/server/internal/protocol"
	"github.com/RAN-GAN/rendis/server/internal/store"
)

func handleMessage(parts []string, db *store.Store) string {
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

	case "GETSET":
		if len(parts) != 3 {
			return protocol.Error("wrong number of arguments")
		}
		value, ok := db.GetSet(parts[1], parts[2])
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

	case "INCR":
		if len(parts) != 2 {
			return protocol.Error("wrong number of arguments")
		}
		newValue, err := db.INCR(parts[1])
		if err != nil {
			return protocol.Error("value is not an integer")
		}
		return protocol.Integer(newValue)

	case "DECR":
		if len(parts) != 2 {
			return protocol.Error("wrong number of arguments")
		}
		newValue, err := db.DECR(parts[1])
		if err != nil {
			return protocol.Error("value is not an integer")
		}
		return protocol.Integer(newValue)

	case "INCRBY":
		if len(parts) != 3 {
			return protocol.Error("wrong number of arguments")
		}

		offset, err := strconv.Atoi(parts[2])
		if err != nil {
			return protocol.Error("value is not an integer")
		}

		result, err := db.INCRBY(parts[1], offset)
		if err != nil {
			return protocol.Error(err.Error())
		}

		return protocol.Integer(result)

	case "DECRBY":
		if len(parts) != 3 {
			return protocol.Error("wrong number of arguments")
		}

		offset, err := strconv.Atoi(parts[2])
		if err != nil {
			return protocol.Error("value is not an integer")
		}

		result, err := db.DECRBY(parts[1], offset)
		if err != nil {
			return protocol.Error(err.Error())
		}

		return protocol.Integer(result)

	case "SETEX":
		if len(parts) != 4 {
			return protocol.Error("wrong number of arguments")
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return protocol.Error("invalid expire time")
		}

		err = db.SETEX(parts[1], seconds, parts[3])
		if err != nil {
			return protocol.Error(err.Error())
		}

		return protocol.SimpleString("OK")

	case "PERSIST":
		if len(parts) != 2 {
			return protocol.Error("wrong number of arguments")
		}

		ok := db.Persist(parts[1])

		if ok {
			return protocol.Integer(1)
		}

		return protocol.Integer(0)

	case "RENAME":
		if len(parts) != 3 {
			return protocol.Error("wrong number of arguments")
		}

		if err := db.Rename(parts[1], parts[2]); err != nil {
			return protocol.Error("no such key")
		}

		return protocol.SimpleString("OK")



	default:
		return protocol.Error("unknown command")
	}
}
