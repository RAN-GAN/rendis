package store

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	Value  string
	Expiry time.Time
}

type Store struct {
	data map[string]Entry
	mu   sync.RWMutex
}

func New() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = Entry{
		Value:  value,
		Expiry: time.Time{},
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	value, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return "", false
	}

	if value.Expiry.IsZero() {
		return value.Value, true
	}
	currTime := time.Now()

	if currTime.After(value.Expiry) {
		s.mu.Lock()

		value, ok = s.data[key]
		if ok && !value.Expiry.IsZero() && currTime.After(value.Expiry) {
			delete(s.data, key)
		}

		s.mu.Unlock()

		return "", false
	}

	return value.Value, true
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return ok
}

func (s *Store) Expire(key string, seconds int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return false
	}
	value.Expiry = time.Now().Add(time.Duration(seconds) * time.Second)
	s.data[key] = value
	return true
}

func (s *Store) TTL(key string) (int, bool) {
	s.mu.RLock()
	value, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return -2, false
	}
	if value.Expiry.IsZero() {
		return -1, true
	}
	currTime := time.Now()
	if currTime.After(value.Expiry) {
		s.mu.Lock()

		value, ok = s.data[key]
		if ok && !value.Expiry.IsZero() && currTime.After(value.Expiry) {
			delete(s.data, key)
		}

		s.mu.Unlock()

		return -2, false
	}

	remaining := time.Until(value.Expiry)
	seconds := int(remaining.Seconds())
	return seconds, ok
}

func (s *Store) clearExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for key, value := range s.data {
		if !value.Expiry.IsZero() && now.After(value.Expiry) {
			fmt.Println("Clearing ", key)
			delete(s.data, key)
		}
	}
}
