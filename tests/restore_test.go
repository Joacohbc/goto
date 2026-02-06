package tests

import (
	"goto/src/core"
	"goto/src/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestRestore(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Create backup file manually
	backupFile := filepath.Join(os.TempDir(), "goto-test-restore.json")
	absPath, _ := filepath.Abs(".")
	content := `[{"Path":"` + absPath + `","Abbreviation":"restored"}]`

	if err := os.WriteFile(backupFile, []byte(content), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(backupFile)

	if err := core.RestoreGPaths(backupFile, false); err != nil {
		t.Errorf("Restore failed: %v", err)
	}

	gpaths, _ := utils.LoadGPaths(false)
	if len(gpaths) == 0 {
		t.Fatal("Restore failed, gpaths empty")
	}

	// Check if any of the paths is the restored one
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "restored" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected abbreviation 'restored' in loaded paths")
	}
}

func TestRestore_InputNotExist(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	if err := core.RestoreGPaths("/non/existent/file", false); err == nil {
		t.Error("Expected error when input file does not exist")
	}
}

func TestRestore_InputIsDir(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	tmpDir := os.TempDir()
	if err := core.RestoreGPaths(tmpDir, false); err == nil {
		t.Error("Expected error when input is a directory")
	}
}

func TestRestore_InvalidJSON(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	tmpFile := filepath.Join(os.TempDir(), "restore_invalid.json")
	os.WriteFile(tmpFile, []byte("{invalid"), 0644)
	defer os.Remove(tmpFile)

	if err := core.RestoreGPaths(tmpFile, false); err == nil {
		t.Error("Expected error when JSON is invalid")
	}
}

func TestRestore_ReadError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test as root")
	}
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	f, err := os.CreateTemp("", "restore_unreadable")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	// Make unreadable
	if err := os.Chmod(f.Name(), 0200); err != nil {
		t.Skipf("Could not chmod: %v", err)
	}

	if err := core.RestoreGPaths(f.Name(), false); err == nil {
		t.Error("Expected error when file is unreadable")
	}
}
