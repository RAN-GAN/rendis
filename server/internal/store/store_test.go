package store

import (
	"testing"
	"time"
)

func TestStore_INCR_DECR(t *testing.T) {
	s := New()

	// INCR new key
	val, err := s.INCR("k1")
	if err != nil || val != 1 {
		t.Fatalf("Expected 1, got %d, err: %v", val, err)
	}

	// INCR existing key
	val, err = s.INCR("k1")
	if err != nil || val != 2 {
		t.Fatalf("Expected 2, got %d, err: %v", val, err)
	}

	// INCR non-integer
	s.Set("k2", "hello")
	_, err = s.INCR("k2")
	if err == nil {
		t.Fatal("Expected error when incrementing non-integer")
	}

	// DECR existing key
	val, err = s.DECR("k1")
	if err != nil || val != 1 {
		t.Fatalf("Expected 1, got %d, err: %v", val, err)
	}

	// DECR new key
	val, err = s.DECR("k3")
	if err != nil || val != -1 {
		t.Fatalf("Expected -1, got %d, err: %v", val, err)
	}
}

func TestStore_INCRBY_DECRBY(t *testing.T) {
	s := New()

	val, err := s.INCRBY("k1", 5)
	if err != nil || val != 5 {
		t.Fatalf("Expected 5, got %d, err: %v", val, err)
	}

	val, err = s.DECRBY("k1", 3)
	if err != nil || val != 2 {
		t.Fatalf("Expected 2, got %d, err: %v", val, err)
	}

	_, err = s.INCRBY("k2", -1)
	if err == nil {
		t.Fatal("Expected error for negative offset in INCRBY")
	}

	_, err = s.DECRBY("k2", -1)
	if err == nil {
		t.Fatal("Expected error for negative offset in DECRBY")
	}
}

func TestStore_SETEX_Persist(t *testing.T) {
	s := New()

	err := s.SETEX("k1", 1, "v1")
	if err != nil {
		t.Fatal(err)
	}

	v, ok := s.Get("k1")
	if !ok || v != "v1" {
		t.Fatalf("Expected v1, got %v", v)
	}

	ok = s.Persist("k1")
	if !ok {
		t.Fatal("Expected Persist to return true")
	}

	time.Sleep(1100 * time.Millisecond)

	v, ok = s.Get("k1")
	if !ok || v != "v1" {
		t.Fatalf("Expected key to persist, got %v", v)
	}
	
	// Test SETEX expiration
	err = s.SETEX("k2", 1, "v2")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond)
	
	_, ok = s.Get("k2")
	if ok {
		t.Fatal("Expected key to expire")
	}
}

func TestStore_Rename(t *testing.T) {
	s := New()

	s.Set("k1", "v1")
	err := s.Rename("k1", "k2")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := s.Get("k1")
	if ok {
		t.Fatal("Expected k1 to be removed")
	}

	v, ok := s.Get("k2")
	if !ok || v != "v1" {
		t.Fatalf("Expected v1 under k2, got %v", v)
	}

	err = s.Rename("nonexistent", "k3")
	if err == nil {
		t.Fatal("Expected error when renaming nonexistent key")
	}
}

func TestStore_GetSet(t *testing.T) {
	s := New()
	
	s.Set("k1", "v1")
	v, ok := s.GetSet("k1", "v2")
	if !ok || v != "v1" {
		t.Fatalf("Expected v1, got %v", v)
	}

	v, ok = s.Get("k1")
	if !ok || v != "v2" {
		t.Fatalf("Expected v2, got %v", v)
	}

	v, ok = s.GetSet("k2", "newval")
	if ok {
		t.Fatal("Expected false for nonexistent key")
	}

	v, ok = s.Get("k2")
	if !ok || v != "newval" {
		t.Fatalf("Expected newval, got %v", v)
	}
}
