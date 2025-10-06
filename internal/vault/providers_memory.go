package vault

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
)

// InMemoryKeyProvider stores the master key in memory (for testing only).
type InMemoryKeyProvider struct {
	key []byte
}

// NewInMemoryKeyProvider creates a key provider with a random 32-byte key.
func NewInMemoryKeyProvider() (*InMemoryKeyProvider, error) {
	key := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return &InMemoryKeyProvider{key: key}, nil
}

// NewInMemoryKeyProviderWithKey creates a key provider with a specific key.
func NewInMemoryKeyProviderWithKey(key []byte) (*InMemoryKeyProvider, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes for AES-256, got %d", len(key))
	}
	return &InMemoryKeyProvider{key: key}, nil
}

func (p *InMemoryKeyProvider) GetKey(ctx context.Context) ([]byte, error) {
	return p.key, nil
}

func (p *InMemoryKeyProvider) RotateKey(ctx context.Context) (oldKey, newKey []byte, err error) {
	oldKey = p.key
	newKey = make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newKey); err != nil {
		return nil, nil, fmt.Errorf("failed to generate new key: %w", err)
	}
	p.key = newKey
	return oldKey, newKey, nil
}

func (p *InMemoryKeyProvider) HealthCheck(ctx context.Context) error {
	if len(p.key) != 32 {
		return fmt.Errorf("invalid key size: %d", len(p.key))
	}
	return nil
}
