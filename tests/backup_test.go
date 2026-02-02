package tests

import (
	"goto/src/cmd"
	"os"
	"path/filepath"
	"testing"
)

func TestBackup(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cmd.AddCmd.Run(c, []string{".", "bkp"})

	backupFile := filepath.Join(os.TempDir(), "goto-test-backup.json")

	// Ensure it doesn't exist
	if err := os.Remove(backupFile); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	defer os.Remove(backupFile)

	c.Flags().StringP("output", "o", "", "")
	c.Flags().Set("output", backupFile)

	captureOutput(func() {
		cmd.BackupCmd.Run(c, []string{})
	})

	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("Backup file not created")
	}
}
