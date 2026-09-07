package storage

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	ErrNotFound   = errors.New("key not found")
	ErrNotInteger = errors.New("vakue is not integer")
)

type Value struct {
	Data      string
	ExpiresAt time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]Value
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Value),
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()

	value, exists := s.data[key]

	s.mu.RUnlock()

	if !exists {
		return "", false
	}

	return value.Data, true
}

func (s *Store) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = Value{
		Data: value,
	}
}

func (s *Store) SetWithExpiration(key string, value string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = Value{
		Data:      value,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return false
	}

	delete(s.data, key)

	return true
}

func (s *Store) Exists(key string) bool {
	_, exists := s.Get(key)

	return exists
}

func (s *Store) Increment(key string, amount int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.data[key]

	if exists && !value.ExpiresAt.IsZero() && time.Now().After(value.ExpiresAt) {
		delete(s.data, key)
		exists = false
	}

	var current int64

	if exists {
		var err error

		current, err = strconv.ParseInt(value.Data, 10, 64)
		if err != nil {
			return 0, ErrNotInteger
		}
	}

	current += amount

	s.data[key] = Value{
		Data: strconv.FormatInt(current, 10),
	}

	return current, nil
}

func (s *Store) Expire(key string, duration time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.data[key]
	if !exists {
		return false
	}

	if !value.ExpiresAt.IsZero() && time.Now().After(value.ExpiresAt) {
		delete(s.data, key)
		return false
	}

	value.ExpiresAt = time.Now().Add(duration)

	s.data[key] = value

	return true
}

func (s *Store) TTL(key string) int64 {
	s.mu.RLock()

	value, exists := s.data[key]

	s.mu.RUnlock()

	if !exists {
		return -2
	}

	if value.ExpiresAt.IsZero() {
		return -1
	}

	remaining := time.Until(value.ExpiresAt)
	if remaining <= 0 {
		s.Delete(key)
		return -2
	}
	return int64(remaining / time.Second)
}

func (s *Store) isExpired(value Value) bool {
	return !value.ExpiresAt.IsZero() && time.Now().After(value.ExpiresAt)
}
