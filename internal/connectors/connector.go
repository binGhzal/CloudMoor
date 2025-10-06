// Package connectors defines the pluggable provider abstraction used to mount
// remote storage backends (S3, WebDAV, Dropbox, etc.) as local filesystems.
//
// Each provider implements the Connector interface, exposing lifecycle hooks
// (Init, ValidateConfig, Open, Close) and metadata for UI discovery.
// Connectors register themselves via RegisterProvider at package init time,
// ensuring deterministic ordering and enabling compile-time plugin composition.
package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Connector defines the lifecycle and operations for a remote storage provider.
// Implementations must be safe for concurrent use after initialization.
type Connector interface {
	// Init prepares the connector with provider-specific configuration.
	// Called once at startup; may perform credential validation or SDK setup.
	Init(ctx context.Context, config Config) error

	// ValidateConfig checks that the provided configuration is complete and valid
	// without performing I/O (e.g., validates required fields, formats).
	ValidateConfig(config Config) error

	// Open establishes a connection to the remote provider and returns a handle
	// that can be used for mount operations. The returned Connection must be
	// closed by the caller when no longer needed.
	Open(ctx context.Context) (Connection, error)

	// Metadata returns descriptive information about this provider (display name,
	// supported features, configuration schema) for use in CLI and Web UI.
	Metadata() ProviderMetadata
}

// Connection represents an active session with a remote storage provider.
// Callers must call Close when finished to release resources.
type Connection interface {
	io.Closer

	// Ping verifies that the connection is still active and responsive.
	Ping(ctx context.Context) error

	// ProviderID returns the unique identifier of the provider backing this connection.
	ProviderID() string
}

// Config holds provider-specific configuration as a map of key-value pairs.
// The connector is responsible for parsing and validating these fields.
type Config map[string]interface{}

// GetString retrieves a string value from the config, returning an error if missing or wrong type.
func (c Config) GetString(key string) (string, error) {
	val, ok := c[key]
	if !ok {
		return "", fmt.Errorf("missing required config key: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("config key %s must be string, got %T", key, val)
	}
	return str, nil
}

// GetBool retrieves a boolean value, defaulting to false if missing.
func (c Config) GetBool(key string) bool {
	val, ok := c[key]
	if !ok {
		return false
	}
	b, ok := val.(bool)
	if !ok {
		return false
	}
	return b
}

// ProviderMetadata describes a connector for discovery and UI rendering.
type ProviderMetadata struct {
	// ID is the unique identifier for this provider (e.g., "s3", "webdav").
	ID string `json:"id"`

	// DisplayName is a human-readable label (e.g., "Amazon S3", "WebDAV").
	DisplayName string `json:"display_name"`

	// Description provides a brief explanation of what this provider does.
	Description string `json:"description"`

	// Version indicates the connector implementation version.
	Version string `json:"version"`

	// ConfigSchema describes expected configuration fields in JSON Schema format.
	// This enables auto-generated UI forms and validation.
	ConfigSchema json.RawMessage `json:"config_schema,omitempty"`
}
