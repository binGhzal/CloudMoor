package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/binGhzal/cloudmoor/internal/vault"
	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage credential vault",
	Long:  "Commands for managing the credential vault (test, health check, etc.).",
}

var vaultTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test vault operations",
	Long: `Test vault operations by performing a round-trip encrypt/decrypt and health check.
This command validates that the vault is properly configured and operational.`,
	RunE: runVaultTest,
}

func init() {
	vaultCmd.AddCommand(vaultTestCmd)
}

func runVaultTest(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Use temp directory for test key
	tmpDir := os.TempDir()
	keyPath := filepath.Join(tmpDir, "cloudmoor-test-vault.key")
	defer os.Remove(keyPath) // Cleanup

	fmt.Printf("Creating test vault with key at: %s\n", keyPath)

	keyProvider, err := vault.NewFileKeyProvider(keyPath)
	if err != nil {
		return fmt.Errorf("failed to create key provider: %w", err)
	}

	// Collect audit events
	var auditEvents []vault.AuditEvent
	auditHook := func(e vault.AuditEvent) {
		auditEvents = append(auditEvents, e)
	}

	store := vault.NewAESGCMStore(keyProvider, auditHook)

	fmt.Println("\n✓ Vault initialized")

	// Health check
	fmt.Print("Running health check... ")
	if err := store.HealthCheck(ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	fmt.Println("✓ PASS")

	// Put operation
	testKey := "test-secret"
	testValue := []byte("sensitive-password-123")

	fmt.Printf("Storing secret '%s'... ", testKey)
	if err := store.Put(ctx, testKey, testValue); err != nil {
		return fmt.Errorf("put failed: %w", err)
	}
	fmt.Println("✓ PASS")

	// Get operation
	fmt.Printf("Retrieving secret '%s'... ", testKey)
	retrieved, err := store.Get(ctx, testKey)
	if err != nil {
		return fmt.Errorf("get failed: %w", err)
	}
	if string(retrieved) != string(testValue) {
		return fmt.Errorf("round trip failed: expected %q, got %q", testValue, retrieved)
	}
	fmt.Println("✓ PASS")

	// List operation
	fmt.Print("Listing secrets... ")
	keys, err := store.List(ctx)
	if err != nil {
		return fmt.Errorf("list failed: %w", err)
	}
	if len(keys) != 1 || keys[0] != testKey {
		return fmt.Errorf("list failed: expected [%s], got %v", testKey, keys)
	}
	fmt.Println("✓ PASS")

	// Delete operation
	fmt.Printf("Deleting secret '%s'... ", testKey)
	if err := store.Delete(ctx, testKey); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	fmt.Println("✓ PASS")

	// Verify deletion
	fmt.Print("Verifying deletion... ")
	keys, err = store.List(ctx)
	if err != nil {
		return fmt.Errorf("list after delete failed: %w", err)
	}
	if len(keys) != 0 {
		return fmt.Errorf("delete verification failed: expected empty list, got %v", keys)
	}
	fmt.Println("✓ PASS")

	// Print audit log
	fmt.Println("\n=== Audit Log ===")
	for _, event := range auditEvents {
		eventJSON, _ := json.MarshalIndent(event, "", "  ")
		fmt.Println(string(eventJSON))
	}

	fmt.Println("\n✅ All vault tests passed!")
	return nil
}
