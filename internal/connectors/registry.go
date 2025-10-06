package connectors

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

var (
	// registry holds all registered providers indexed by ID.
	registry = make(map[string]Connector)

	// registryOrder maintains deterministic ordering of provider IDs.
	registryOrder []string

	// mu protects concurrent access to the registry during initialization.
	mu sync.RWMutex
)

// RegisterProvider registers a connector implementation under the given metadata ID.
// Panics if a provider with the same ID is already registered.
// This function is intended to be called from init() functions in provider packages.
func RegisterProvider(c Connector) {
	mu.Lock()
	defer mu.Unlock()

	meta := c.Metadata()
	if meta.ID == "" {
		panic("connectors: provider metadata must include non-empty ID")
	}

	if _, exists := registry[meta.ID]; exists {
		panic(fmt.Sprintf("connectors: provider %q already registered", meta.ID))
	}

	registry[meta.ID] = c
	registryOrder = append(registryOrder, meta.ID)
}

// GetProvider retrieves a registered connector by ID.
// Returns nil if no provider with that ID exists.
func GetProvider(id string) Connector {
	mu.RLock()
	defer mu.RUnlock()
	return registry[id]
}

// ListProviders returns metadata for all registered providers in deterministic order.
func ListProviders() []ProviderMetadata {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]ProviderMetadata, 0, len(registryOrder))
	for _, id := range registryOrder {
		if c, ok := registry[id]; ok {
			result = append(result, c.Metadata())
		}
	}
	return result
}

// ProviderIDs returns the sorted list of registered provider IDs.
// Useful for alphabetical display or testing determinism.
func ProviderIDs() []string {
	mu.RLock()
	defer mu.RUnlock()

	ids := make([]string, len(registryOrder))
	copy(ids, registryOrder)
	sort.Strings(ids)
	return ids
}

// ExportManifest generates a JSON-serialized manifest of all providers.
// Intended for consumption by the Web UI to dynamically render provider options.
func ExportManifest() ([]byte, error) {
	providers := ListProviders()
	return json.MarshalIndent(providers, "", "  ")
}
