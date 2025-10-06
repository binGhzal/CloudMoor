package connectors

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeConnector is a minimal test implementation.
type fakeConnector struct {
	meta ProviderMetadata
}

func (f *fakeConnector) Init(ctx context.Context, config Config) error          { return nil }
func (f *fakeConnector) ValidateConfig(config Config) error                     { return nil }
func (f *fakeConnector) Open(ctx context.Context) (Connection, error)           { return nil, nil }
func (f *fakeConnector) Metadata() ProviderMetadata                             { return f.meta }

func TestRegisterProvider(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		// Clear registry for isolated test
		mu.Lock()
		registry = make(map[string]Connector)
		registryOrder = nil
		mu.Unlock()

		fake := &fakeConnector{
			meta: ProviderMetadata{
				ID:          "test-provider",
				DisplayName: "Test Provider",
				Description: "A test connector",
				Version:     "1.0.0",
			},
		}

		RegisterProvider(fake)

		retrieved := GetProvider("test-provider")
		require.NotNil(t, retrieved)
		require.Equal(t, "test-provider", retrieved.Metadata().ID)
	})

	t.Run("panic on duplicate ID", func(t *testing.T) {
		mu.Lock()
		registry = make(map[string]Connector)
		registryOrder = nil
		mu.Unlock()

		fake := &fakeConnector{
			meta: ProviderMetadata{ID: "duplicate"},
		}

		RegisterProvider(fake)

		require.Panics(t, func() {
			RegisterProvider(fake)
		})
	})

	t.Run("panic on empty ID", func(t *testing.T) {
		mu.Lock()
		registry = make(map[string]Connector)
		registryOrder = nil
		mu.Unlock()

		fake := &fakeConnector{
			meta: ProviderMetadata{ID: ""},
		}

		require.Panics(t, func() {
			RegisterProvider(fake)
		})
	})
}

func TestListProviders(t *testing.T) {
	// Reset registry
	mu.Lock()
	registry = make(map[string]Connector)
	registryOrder = nil
	mu.Unlock()

	providers := []ProviderMetadata{
		{ID: "s3", DisplayName: "Amazon S3", Version: "1.0.0"},
		{ID: "webdav", DisplayName: "WebDAV", Version: "1.0.0"},
		{ID: "dropbox", DisplayName: "Dropbox", Version: "1.0.0"},
	}

	for _, meta := range providers {
		RegisterProvider(&fakeConnector{meta: meta})
	}

	list := ListProviders()
	require.Len(t, list, 3)

	// Verify registration order is preserved
	require.Equal(t, "s3", list[0].ID)
	require.Equal(t, "webdav", list[1].ID)
	require.Equal(t, "dropbox", list[2].ID)
}

func TestProviderIDs(t *testing.T) {
	mu.Lock()
	registry = make(map[string]Connector)
	registryOrder = nil
	mu.Unlock()

	RegisterProvider(&fakeConnector{meta: ProviderMetadata{ID: "zebra"}})
	RegisterProvider(&fakeConnector{meta: ProviderMetadata{ID: "apple"}})
	RegisterProvider(&fakeConnector{meta: ProviderMetadata{ID: "mango"}})

	ids := ProviderIDs()
	require.Equal(t, []string{"apple", "mango", "zebra"}, ids, "IDs should be sorted alphabetically")
}

func TestExportManifest(t *testing.T) {
	mu.Lock()
	registry = make(map[string]Connector)
	registryOrder = nil
	mu.Unlock()

	RegisterProvider(&fakeConnector{
		meta: ProviderMetadata{
			ID:          "test",
			DisplayName: "Test",
			Description: "A test provider",
			Version:     "1.0.0",
		},
	})

	manifest, err := ExportManifest()
	require.NoError(t, err)
	require.NotEmpty(t, manifest)

	var parsed []ProviderMetadata
	err = json.Unmarshal(manifest, &parsed)
	require.NoError(t, err)
	require.Len(t, parsed, 1)
	require.Equal(t, "test", parsed[0].ID)
}

func TestConfigHelpers(t *testing.T) {
	t.Run("GetString success", func(t *testing.T) {
		cfg := Config{"key": "value"}
		val, err := cfg.GetString("key")
		require.NoError(t, err)
		require.Equal(t, "value", val)
	})

	t.Run("GetString missing key", func(t *testing.T) {
		cfg := Config{}
		_, err := cfg.GetString("missing")
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing required config key")
	})

	t.Run("GetString wrong type", func(t *testing.T) {
		cfg := Config{"key": 123}
		_, err := cfg.GetString("key")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be string")
	})

	t.Run("GetBool", func(t *testing.T) {
		cfg := Config{"enabled": true}
		require.True(t, cfg.GetBool("enabled"))
		require.False(t, cfg.GetBool("missing"))
		require.False(t, cfg.GetBool("wrong-type"))
	})
}
