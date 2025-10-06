package vault

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAESGCMStore_PutGet(t *testing.T) {
	ctx := context.Background()
	keyProvider, err := NewInMemoryKeyProvider()
	require.NoError(t, err)

	var auditEvents []AuditEvent
	auditHook := func(e AuditEvent) {
		auditEvents = append(auditEvents, e)
	}

	store := NewAESGCMStore(keyProvider, auditHook)

	t.Run("successful round trip", func(t *testing.T) {
		auditEvents = nil
		secret := []byte("my-secret-password")

		err := store.Put(ctx, "test-key", secret)
		require.NoError(t, err)

		retrieved, err := store.Get(ctx, "test-key")
		require.NoError(t, err)
		require.Equal(t, secret, retrieved)

		// Verify audit events
		require.Len(t, auditEvents, 2)
		require.Equal(t, "put", auditEvents[0].Operation)
		require.True(t, auditEvents[0].Success)
		require.Equal(t, "get", auditEvents[1].Operation)
		require.True(t, auditEvents[1].Success)
	})

	t.Run("empty key validation", func(t *testing.T) {
		err := store.Put(ctx, "", []byte("value"))
		require.ErrorIs(t, err, ErrKeyEmpty)

		_, err = store.Get(ctx, "")
		require.ErrorIs(t, err, ErrKeyEmpty)
	})

	t.Run("empty value validation", func(t *testing.T) {
		err := store.Put(ctx, "empty", []byte{})
		require.ErrorIs(t, err, ErrValueEmpty)
	})

	t.Run("get non-existent key", func(t *testing.T) {
		_, err := store.Get(ctx, "non-existent")
		require.ErrorIs(t, err, ErrNotFound)
	})
}

func TestAESGCMStore_Delete(t *testing.T) {
	ctx := context.Background()
	keyProvider, err := NewInMemoryKeyProvider()
	require.NoError(t, err)

	store := NewAESGCMStore(keyProvider, nil)

	err = store.Put(ctx, "delete-me", []byte("value"))
	require.NoError(t, err)

	err = store.Delete(ctx, "delete-me")
	require.NoError(t, err)

	// Verify it's gone
	_, err = store.Get(ctx, "delete-me")
	require.ErrorIs(t, err, ErrNotFound)

	// Delete non-existent key
	err = store.Delete(ctx, "non-existent")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestAESGCMStore_List(t *testing.T) {
	ctx := context.Background()
	keyProvider, err := NewInMemoryKeyProvider()
	require.NoError(t, err)

	store := NewAESGCMStore(keyProvider, nil)

	// Empty list
	keys, err := store.List(ctx)
	require.NoError(t, err)
	require.Empty(t, keys)

	// Add multiple secrets
	secrets := map[string][]byte{
		"zebra": []byte("z"),
		"apple": []byte("a"),
		"mango": []byte("m"),
	}

	for k, v := range secrets {
		err := store.Put(ctx, k, v)
		require.NoError(t, err)
	}

	keys, err = store.List(ctx)
	require.NoError(t, err)
	require.Equal(t, []string{"apple", "mango", "zebra"}, keys, "keys should be sorted")
}

func TestAESGCMStore_HealthCheck(t *testing.T) {
	ctx := context.Background()

	t.Run("healthy provider", func(t *testing.T) {
		keyProvider, err := NewInMemoryKeyProvider()
		require.NoError(t, err)

		store := NewAESGCMStore(keyProvider, nil)
		err = store.HealthCheck(ctx)
		require.NoError(t, err)
	})

	t.Run("unhealthy provider", func(t *testing.T) {
		badProvider := &failingKeyProvider{shouldFail: true}
		store := NewAESGCMStore(badProvider, nil)

		err := store.HealthCheck(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUnhealthy)
	})
}

func TestAESGCMStore_KeyRotation(t *testing.T) {
	ctx := context.Background()
	keyProvider, err := NewInMemoryKeyProvider()
	require.NoError(t, err)

	store := NewAESGCMStore(keyProvider, nil).(*aesgcmStore)

	// Store secret with old key
	originalSecret := []byte("important-secret")
	err = store.Put(ctx, "rotate-test", originalSecret)
	require.NoError(t, err)

	// Rotate key
	oldKey, newKey, err := keyProvider.RotateKey(ctx)
	require.NoError(t, err)
	require.NotEqual(t, oldKey, newKey)

	// Old encrypted data should fail to decrypt with new key
	_, err = store.Get(ctx, "rotate-test")
	require.Error(t, err, "should fail because data was encrypted with old key")

	// Re-encrypt with new key
	store.mu.Lock()
	encryptedWithOldKey := store.secrets["rotate-test"]
	plaintext, err := store.decrypt(oldKey, encryptedWithOldKey)
	require.NoError(t, err)

	reEncrypted, err := store.encrypt(newKey, plaintext)
	require.NoError(t, err)
	store.secrets["rotate-test"] = reEncrypted
	store.mu.Unlock()

	// Now retrieval should work
	retrieved, err := store.Get(ctx, "rotate-test")
	require.NoError(t, err)
	require.Equal(t, originalSecret, retrieved)
}

func TestAESGCMStore_AuditHook(t *testing.T) {
	ctx := context.Background()
	keyProvider, err := NewInMemoryKeyProvider()
	require.NoError(t, err)

	var auditEvents []AuditEvent
	auditHook := func(e AuditEvent) {
		auditEvents = append(auditEvents, e)
	}

	store := NewAESGCMStore(keyProvider, auditHook)

	// Exercise all operations
	_ = store.Put(ctx, "key1", []byte("value1"))
	_, _ = store.Get(ctx, "key1")
	_, _ = store.List(ctx)
	_ = store.Delete(ctx, "key1")
	_ = store.HealthCheck(ctx)

	require.Len(t, auditEvents, 5)

	operations := make([]string, len(auditEvents))
	for i, e := range auditEvents {
		operations[i] = e.Operation
		require.NotZero(t, e.Timestamp)
	}

	require.Equal(t, []string{"put", "get", "list", "delete", "health"}, operations)
}

func TestAESGCMStore_Encryption(t *testing.T) {
	keyProvider, err := NewInMemoryKeyProvider()
	require.NoError(t, err)

	store := NewAESGCMStore(keyProvider, nil).(*aesgcmStore)

	key, err := keyProvider.GetKey(context.Background())
	require.NoError(t, err)

	plaintext := []byte("sensitive data")

	// Encrypt
	ciphertext, err := store.encrypt(key, plaintext)
	require.NoError(t, err)
	require.NotEqual(t, plaintext, ciphertext)
	require.Greater(t, len(ciphertext), len(plaintext), "ciphertext should include nonce and tag")

	// Decrypt
	decrypted, err := store.decrypt(key, ciphertext)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)

	// Tampering detection
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)-1] ^= 0xFF // Flip last byte

	_, err = store.decrypt(key, tampered)
	require.Error(t, err, "should detect tampering")
}

// failingKeyProvider simulates key provider failures.
type failingKeyProvider struct {
	shouldFail bool
}

func (p *failingKeyProvider) GetKey(ctx context.Context) ([]byte, error) {
	if p.shouldFail {
		return nil, errors.New("simulated key provider failure")
	}
	return make([]byte, 32), nil
}

func (p *failingKeyProvider) RotateKey(ctx context.Context) ([]byte, []byte, error) {
	if p.shouldFail {
		return nil, nil, errors.New("simulated rotation failure")
	}
	return make([]byte, 32), make([]byte, 32), nil
}

func (p *failingKeyProvider) HealthCheck(ctx context.Context) error {
	if p.shouldFail {
		return errors.New("simulated health check failure")
	}
	return nil
}

func TestInMemoryKeyProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("random key generation", func(t *testing.T) {
		p1, err := NewInMemoryKeyProvider()
		require.NoError(t, err)

		p2, err := NewInMemoryKeyProvider()
		require.NoError(t, err)

		key1, err := p1.GetKey(ctx)
		require.NoError(t, err)
		require.Len(t, key1, 32)

		key2, err := p2.GetKey(ctx)
		require.NoError(t, err)
		require.Len(t, key2, 32)

		require.NotEqual(t, key1, key2, "keys should be random")
	})

	t.Run("specific key", func(t *testing.T) {
		specificKey := make([]byte, 32)
		for i := range specificKey {
			specificKey[i] = byte(i)
		}

		p, err := NewInMemoryKeyProviderWithKey(specificKey)
		require.NoError(t, err)

		key, err := p.GetKey(ctx)
		require.NoError(t, err)
		require.Equal(t, specificKey, key)
	})

	t.Run("invalid key size", func(t *testing.T) {
		_, err := NewInMemoryKeyProviderWithKey([]byte("too short"))
		require.Error(t, err)
	})

	t.Run("rotation", func(t *testing.T) {
		p, err := NewInMemoryKeyProvider()
		require.NoError(t, err)

		oldKey, err := p.GetKey(ctx)
		require.NoError(t, err)

		rotatedOld, newKey, err := p.RotateKey(ctx)
		require.NoError(t, err)
		require.Equal(t, oldKey, rotatedOld)
		require.NotEqual(t, oldKey, newKey)

		currentKey, err := p.GetKey(ctx)
		require.NoError(t, err)
		require.Equal(t, newKey, currentKey)
	})

	t.Run("health check", func(t *testing.T) {
		p, err := NewInMemoryKeyProvider()
		require.NoError(t, err)

		err = p.HealthCheck(ctx)
		require.NoError(t, err)
	})
}
