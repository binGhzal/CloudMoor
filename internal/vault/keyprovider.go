package vault

import (
	"context"
)

// KeyProvider abstracts master key retrieval to support multiple backends
// (filesystem, OS keychain, external secret stores).
type KeyProvider interface {
	// GetKey retrieves the current master key for encryption/decryption.
	// Must return a 32-byte AES-256 key.
	GetKey(ctx context.Context) ([]byte, error)

	// RotateKey generates and stores a new master key, returning both old and new keys
	// to support re-encryption of existing secrets.
	RotateKey(ctx context.Context) (oldKey, newKey []byte, err error)

	// HealthCheck verifies the key provider is accessible.
	HealthCheck(ctx context.Context) error
}
