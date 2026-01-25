package tests

import (
	"goto/src/cmd"
	"os"
	"path/filepath"
	"testing"
)

func TestBackup(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "bkp"})

	backupFile := filepath.Join(os.TempDir(), "goto-test-backup.json")
	// Ensure it doesn't exist
	os.Remove(backupFile)
	defer os.Remove(backupFile)

	bkpCmd := getTempCmd()
	bkpCmd.Flags().StringP("output", "o", "", "")
	bkpCmd.Flags().Set("output", backupFile)

	captureOutput(func() {
		cmd.BackupCmd.Run(bkpCmd, []string{})
	})

	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("Backup file not created")
	}
}
