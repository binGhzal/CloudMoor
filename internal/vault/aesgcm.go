package vault

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

// aesgcmStore implements Store using AES-256-GCM envelope encryption.
type aesgcmStore struct {
	keyProvider KeyProvider
	auditHook   AuditHook
	mu          sync.RWMutex
	secrets     map[string][]byte // Encrypted values stored in memory
}

// NewAESGCMStore creates a new Store using AES-GCM encryption.
func NewAESGCMStore(keyProvider KeyProvider, auditHook AuditHook) Store {
	if auditHook == nil {
		auditHook = func(AuditEvent) {} // No-op hook
	}
	return &aesgcmStore{
		keyProvider: keyProvider,
		auditHook:   auditHook,
		secrets:     make(map[string][]byte),
	}
}

func (s *aesgcmStore) Put(ctx context.Context, key string, value []byte) error {
	start := time.Now()
	event := AuditEvent{
		Timestamp: start,
		Operation: "put",
		Key:       key,
	}
	defer func() { s.auditHook(event) }()

	if key == "" {
		event.Success = false
		event.Error = ErrKeyEmpty.Error()
		return ErrKeyEmpty
	}
	if len(value) == 0 {
		event.Success = false
		event.Error = ErrValueEmpty.Error()
		return ErrValueEmpty
	}

	masterKey, err := s.keyProvider.GetKey(ctx)
	if err != nil {
		event.Success = false
		event.Error = fmt.Sprintf("%v: %v", ErrKeyProvider, err)
		return fmt.Errorf("%w: %v", ErrKeyProvider, err)
	}

	encrypted, err := s.encrypt(masterKey, value)
	if err != nil {
		event.Success = false
		event.Error = fmt.Sprintf("%v: %v", ErrEncryption, err)
		return fmt.Errorf("%w: %v", ErrEncryption, err)
	}

	s.mu.Lock()
	s.secrets[key] = encrypted
	s.mu.Unlock()

	event.Success = true
	event.Metadata = map[string]string{"size": fmt.Sprintf("%d", len(value))}
	return nil
}

func (s *aesgcmStore) Get(ctx context.Context, key string) ([]byte, error) {
	start := time.Now()
	event := AuditEvent{
		Timestamp: start,
		Operation: "get",
		Key:       key,
	}
	defer func() { s.auditHook(event) }()

	if key == "" {
		event.Success = false
		event.Error = ErrKeyEmpty.Error()
		return nil, ErrKeyEmpty
	}

	s.mu.RLock()
	encrypted, exists := s.secrets[key]
	s.mu.RUnlock()

	if !exists {
		event.Success = false
		event.Error = ErrNotFound.Error()
		return nil, ErrNotFound
	}

	masterKey, err := s.keyProvider.GetKey(ctx)
	if err != nil {
		event.Success = false
		event.Error = fmt.Sprintf("%v: %v", ErrKeyProvider, err)
		return nil, fmt.Errorf("%w: %v", ErrKeyProvider, err)
	}

	plaintext, err := s.decrypt(masterKey, encrypted)
	if err != nil {
		event.Success = false
		event.Error = fmt.Sprintf("%v: %v", ErrDecryption, err)
		return nil, fmt.Errorf("%w: %v", ErrDecryption, err)
	}

	event.Success = true
	event.Metadata = map[string]string{"size": fmt.Sprintf("%d", len(plaintext))}
	return plaintext, nil
}

func (s *aesgcmStore) Delete(ctx context.Context, key string) error {
	start := time.Now()
	event := AuditEvent{
		Timestamp: start,
		Operation: "delete",
		Key:       key,
	}
	defer func() { s.auditHook(event) }()

	if key == "" {
		event.Success = false
		event.Error = ErrKeyEmpty.Error()
		return ErrKeyEmpty
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.secrets[key]; !exists {
		event.Success = false
		event.Error = ErrNotFound.Error()
		return ErrNotFound
	}

	delete(s.secrets, key)
	event.Success = true
	return nil
}

func (s *aesgcmStore) List(ctx context.Context) ([]string, error) {
	start := time.Now()
	event := AuditEvent{
		Timestamp: start,
		Operation: "list",
	}
	defer func() { s.auditHook(event) }()

	s.mu.RLock()
	keys := make([]string, 0, len(s.secrets))
	for k := range s.secrets {
		keys = append(keys, k)
	}
	s.mu.RUnlock()

	sort.Strings(keys)
	event.Success = true
	event.Metadata = map[string]string{"count": fmt.Sprintf("%d", len(keys))}
	return keys, nil
}

func (s *aesgcmStore) HealthCheck(ctx context.Context) error {
	start := time.Now()
	event := AuditEvent{
		Timestamp: start,
		Operation: "health",
	}
	defer func() { s.auditHook(event) }()

	if err := s.keyProvider.HealthCheck(ctx); err != nil {
		event.Success = false
		event.Error = fmt.Sprintf("%v: %v", ErrUnhealthy, err)
		return fmt.Errorf("%w: %v", ErrUnhealthy, err)
	}

	// Verify we can get the master key
	_, err := s.keyProvider.GetKey(ctx)
	if err != nil {
		event.Success = false
		event.Error = fmt.Sprintf("%v: %v", ErrUnhealthy, err)
		return fmt.Errorf("%w: %v", ErrUnhealthy, err)
	}

	event.Success = true
	return nil
}

// encrypt uses AES-256-GCM to encrypt plaintext with the master key.
// Returns: nonce || ciphertext || tag
func (s *aesgcmStore) encrypt(masterKey, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal prepends nonce and appends tag
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt uses AES-256-GCM to decrypt ciphertext with the master key.
// Expects: nonce || ciphertext || tag
func (s *aesgcmStore) decrypt(masterKey, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
