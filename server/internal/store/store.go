package store

import (
	"errors"
	"fmt"
	"strconv"
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

func (s *Store) GetSet(key, newValue string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]

	if ok && !value.Expiry.IsZero() && time.Now().After(value.Expiry) {
		delete(s.data, key)
		ok = false
	}

	var oldValue string
	if ok {
		oldValue = value.Value
	}

	s.data[key] = Entry{
		Value:  newValue,
		Expiry: time.Time{},
	}

	return oldValue, ok
}

func (s *Store) Rename(key, newKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]

	if !ok {
		return errors.New("key does not exist")
	}
	if ok && !value.Expiry.IsZero() && time.Now().After(value.Expiry) {
		delete(s.data, key)
		ok = false
	}

	if !ok {
		return errors.New("no such key")
	}
	s.data[newKey] = value
	delete(s.data, key)

	return nil
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

func (s *Store) change(key string, delta int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]

	if !ok {
		s.data[key] = Entry{
			Value: strconv.Itoa(delta),
		}
		return delta, nil
	}

	number, err := strconv.Atoi(value.Value)
	if err != nil {
		return 0, err
	}

	number += delta

	value.Value = strconv.Itoa(number)
	s.data[key] = value

	return number, nil
}

func (s *Store) INCR(key string) (int, error) {
	return s.change(key, 1)
}

func (s *Store) DECR(key string) (int, error) {
	return s.change(key, -1)
}

func (s *Store) INCRBY(key string, offset int) (int, error) {
	if offset <= 0 {
		return 0, errors.New("Value must be >= 1")
	}

	return s.change(key, offset)
}

func (s *Store) DECRBY(key string, offset int) (int, error) {
	if offset <= 0 {
		return 0, errors.New("Value must be >= 1")
	}
	return s.change(key, -offset)
}

func (s *Store) SETEX(key string, seconds int, value string) error {
	if seconds <= 0 {
		return errors.New("invalid expire time")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = Entry{
		Value:  value,
		Expiry: time.Now().Add(time.Duration(seconds) * time.Second),
	}
	return nil
}

func (s *Store) Persist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return false
	}

	value.Expiry = time.Time{}
	s.data[key] = value

	return true
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
