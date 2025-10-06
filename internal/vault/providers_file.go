package vault

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileKeyProvider stores the master key on disk (for MVP).
// Production deployments should use OS keychain or external secret stores.
type FileKeyProvider struct {
	keyPath string
}

// NewFileKeyProvider creates a key provider that reads from the specified path.
// If the file doesn't exist, it generates a new random key.
func NewFileKeyProvider(keyPath string) (*FileKeyProvider, error) {
	p := &FileKeyProvider{keyPath: keyPath}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(keyPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create key directory: %w", err)
	}

	// Generate key if it doesn't exist
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		if err := p.generateKey(); err != nil {
			return nil, fmt.Errorf("failed to generate key: %w", err)
		}
	}

	return p, nil
}

func (p *FileKeyProvider) GetKey(ctx context.Context) ([]byte, error) {
	key, err := os.ReadFile(p.keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: expected 32 bytes, got %d", len(key))
	}
	return key, nil
}

func (p *FileKeyProvider) RotateKey(ctx context.Context) (oldKey, newKey []byte, err error) {
	// Read old key
	oldKey, err = p.GetKey(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read old key: %w", err)
	}

	// Generate new key
	newKey = make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newKey); err != nil {
		return nil, nil, fmt.Errorf("failed to generate new key: %w", err)
	}

	// Backup old key
	backupPath := p.keyPath + ".bak"
	if err := os.WriteFile(backupPath, oldKey, 0600); err != nil {
		return nil, nil, fmt.Errorf("failed to backup old key: %w", err)
	}

	// Write new key
	if err := os.WriteFile(p.keyPath, newKey, 0600); err != nil {
		return nil, nil, fmt.Errorf("failed to write new key: %w", err)
	}

	return oldKey, newKey, nil
}

func (p *FileKeyProvider) HealthCheck(ctx context.Context) error {
	key, err := p.GetKey(ctx)
	if err != nil {
		return err
	}
	if len(key) != 32 {
		return fmt.Errorf("invalid key size: %d", len(key))
	}
	return nil
}

func (p *FileKeyProvider) generateKey() error {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return err
	}
	return os.WriteFile(p.keyPath, key, 0600)
}
