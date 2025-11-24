package auth

import (
	"sync"
)

// Store defines the interface for client credential storage
type Store interface {
	VerifyCredentials(clientID, clientSecret string) (bool, error)
}

// InMemoryStore implements Store using a map
type InMemoryStore struct {
	mu      sync.RWMutex
	clients map[string]string // clientID -> clientSecret
}

// NewInMemoryStore creates a new in-memory store with some pre-seeded clients
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		clients: map[string]string{
			"service-a": "secret-a",
			"service-b": "secret-b",
		},
	}
}

// VerifyCredentials checks if the clientID and clientSecret match
func (s *InMemoryStore) VerifyCredentials(clientID, clientSecret string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	secret, ok := s.clients[clientID]
	if !ok {
		return false, nil
	}

	// In a real app, we should compare hashes here
	if secret != clientSecret {
		return false, nil
	}

	return true, nil
}
