package tests

import (
	"goto/src/core"
	"os"
	"path/filepath"
	"testing"
)

func TestBackup(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add path
	if err := core.AddPath(".", "bkp", false); err != nil {
		t.Fatal(err)
	}

	backupFile := filepath.Join(os.TempDir(), "goto-test-backup.json")

	// Ensure it doesn't exist
	if err := os.Remove(backupFile); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	defer os.Remove(backupFile)

	if err := core.BackupGPaths(backupFile, false); err != nil {
		t.Errorf("Backup failed: %v", err)
	}

	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("Backup file not created")
	}
}

func TestBackup_OutputExists(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	tmpFile := filepath.Join(os.TempDir(), "backup_exists.json")
	os.WriteFile(tmpFile, []byte(""), 0644)
	defer os.Remove(tmpFile)

	if err := core.BackupGPaths(tmpFile, false); err == nil {
		t.Error("Expected error when backup file exists")
	}
}
