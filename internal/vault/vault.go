// Package vault provides secure credential storage with envelope encryption.
package vault

import (
	"context"
	"fmt"
	"time"
)

// Store defines the interface for secure credential storage.
// All operations emit structured audit events for compliance tracking.
type Store interface {
	// Put stores a secret under the given key. Returns an error if encryption fails.
	Put(ctx context.Context, key string, value []byte) error

	// Get retrieves a secret by key. Returns ErrNotFound if the key doesn't exist.
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete removes a secret. Returns ErrNotFound if the key doesn't exist.
	Delete(ctx context.Context, key string) error

	// List returns all stored keys in alphabetical order.
	List(ctx context.Context) ([]string, error)

	// HealthCheck verifies the vault is operational and the master key is accessible.
	HealthCheck(ctx context.Context) error
}

// AuditEvent captures structured information about vault operations.
type AuditEvent struct {
	Timestamp time.Time         `json:"timestamp"`
	Operation string            `json:"operation"` // "put", "get", "delete", "list", "health"
	Key       string            `json:"key"`
	Success   bool              `json:"success"`
	Error     string            `json:"error,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// AuditHook receives audit events for external logging/alerting.
type AuditHook func(event AuditEvent)

// Common errors returned by Store implementations.
var (
	ErrNotFound     = fmt.Errorf("vault: secret not found")
	ErrKeyEmpty     = fmt.Errorf("vault: key cannot be empty")
	ErrValueEmpty   = fmt.Errorf("vault: value cannot be empty")
	ErrUnhealthy    = fmt.Errorf("vault: health check failed")
	ErrKeyProvider  = fmt.Errorf("vault: key provider error")
	ErrEncryption   = fmt.Errorf("vault: encryption failed")
	ErrDecryption   = fmt.Errorf("vault: decryption failed")
)
