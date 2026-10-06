package server

import (
	"testing"
	"time"

	"github.com/RAN-GAN/rendis/server/internal/protocol"
	"github.com/RAN-GAN/rendis/server/internal/store"
)

func TestHandleMessage_INCR_DECR(t *testing.T) {
	db := store.New()

	// INCR new key
	res := handleMessage([]string{"INCR", "k1"}, db)
	if res != protocol.Integer(1) {
		t.Fatalf("Expected 1, got %v", res)
	}

	// INCR existing key
	res = handleMessage([]string{"INCR", "k1"}, db)
	if res != protocol.Integer(2) {
		t.Fatalf("Expected 2, got %v", res)
	}

	// DECR existing key
	res = handleMessage([]string{"DECR", "k1"}, db)
	if res != protocol.Integer(1) {
		t.Fatalf("Expected 1, got %v", res)
	}

	// DECR new key
	res = handleMessage([]string{"DECR", "k2"}, db)
	if res != protocol.Integer(-1) {
		t.Fatalf("Expected -1, got %v", res)
	}

	// Error case
	db.Set("k3", "abc")
	res = handleMessage([]string{"INCR", "k3"}, db)
	if res != protocol.Error("value is not an integer") {
		t.Fatalf("Expected error, got %v", res)
	}
}

func TestHandleMessage_SETEX(t *testing.T) {
	db := store.New()

	res := handleMessage([]string{"SETEX", "k1", "1", "val"}, db)
	if res != protocol.SimpleString("OK") {
		t.Fatalf("Expected OK, got %v", res)
	}

	v, ok := db.Get("k1")
	if !ok || v != "val" {
		t.Fatalf("Expected val, got %v", v)
	}

	// Error invalid expire time
	res = handleMessage([]string{"SETEX", "k1", "invalid", "val"}, db)
	if res != protocol.Error("invalid expire time") {
		t.Fatalf("Expected invalid expire time, got %v", res)
	}
}

func TestHandleMessage_PERSIST(t *testing.T) {
	db := store.New()

	db.SETEX("k1", 1, "val")
	res := handleMessage([]string{"PERSIST", "k1"}, db)
	if res != protocol.Integer(1) {
		t.Fatalf("Expected 1, got %v", res)
	}

	time.Sleep(1100 * time.Millisecond)

	v, ok := db.Get("k1")
	if !ok || v != "val" {
		t.Fatalf("Expected key to persist, got %v", v)
	}

	// Persist non-existent key
	res = handleMessage([]string{"PERSIST", "k2"}, db)
	if res != protocol.Integer(0) {
		t.Fatalf("Expected 0, got %v", res)
	}
}

func TestHandleMessage_RENAME(t *testing.T) {
	db := store.New()

	db.Set("k1", "val")
	res := handleMessage([]string{"RENAME", "k1", "k2"}, db)
	if res != protocol.SimpleString("OK") {
		t.Fatalf("Expected OK, got %v", res)
	}

	_, ok := db.Get("k1")
	if ok {
		t.Fatal("Expected k1 to be removed")
	}

	v, ok := db.Get("k2")
	if !ok || v != "val" {
		t.Fatalf("Expected val under k2, got %v", v)
	}

	res = handleMessage([]string{"RENAME", "nonexistent", "k3"}, db)
	if res != protocol.Error("no such key") {
		t.Fatalf("Expected error no such key, got %v", res)
	}
}

func TestHandleMessage_GETSET(t *testing.T) {
	db := store.New()

	res := handleMessage([]string{"GETSET", "k1", "v1"}, db)
	if res != protocol.NullBulkString() {
		t.Fatalf("Expected null bulk string, got %v", res)
	}

	db.Set("k2", "old")
	res = handleMessage([]string{"GETSET", "k2", "new"}, db)
	if res != protocol.BulkString("old") {
		t.Fatalf("Expected old, got %v", res)
	}

	v, _ := db.Get("k2")
	if v != "new" {
		t.Fatalf("Expected k2 to be new, got %v", v)
	}
}

func TestHandleMessage_INCRBY_DECRBY(t *testing.T) {
	db := store.New()

	res := handleMessage([]string{"INCRBY", "k1", "5"}, db)
	if res != protocol.Integer(5) {
		t.Fatalf("Expected 5, got %v", res)
	}

	res = handleMessage([]string{"DECRBY", "k1", "2"}, db)
	if res != protocol.Integer(3) {
		t.Fatalf("Expected 3, got %v", res)
	}

	res = handleMessage([]string{"INCRBY", "k1", "invalid"}, db)
	if res != protocol.Error("value is not an integer") {
		t.Fatalf("Expected invalid int error, got %v", res)
	}

	res = handleMessage([]string{"INCRBY", "k1", "-1"}, db)
	if res != protocol.Error("Value must be >= 1") {
		t.Fatalf("Expected value >= 1 error, got %v", res)
	}
}
