package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestRestore(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Create backup file manually
	backupFile := filepath.Join(os.TempDir(), "goto-test-restore.json")
	// Use absolute path for Path field just in case validation runs (though restore usually just unmarshals)
	// But SaveGPathsFile calls CheckRepeatedItems.
	// We need a valid path.
	absPath, _ := filepath.Abs(".")

	// Manually construct JSON to avoid importing gpath struct if we don't have to
	// But we can use gpath.GotoPath if exported.
	// Or raw json.
	// JSON format: [{"Path":"...","Abbreviation":"..."}]
	content := `[{"Path":"` + absPath + `","Abbreviation":"restored"}]`

	if err := os.WriteFile(backupFile, []byte(content), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(backupFile)

	c.Flags().StringP("input", "i", "", "")
	c.Flags().Set("input", backupFile)

	captureOutput(func() {
		cmd.RestoreCmd.Run(c, []string{})
	})

	gpaths, _ := utils.LoadGPaths(utils.TemporalFlagPassed(c))
	if len(gpaths) == 0 {
		t.Fatal("Restore failed, gpaths empty")
	}
	if gpaths[0].Abbreviation != "restored" {
		t.Errorf("Expected abbreviation 'restored', got '%s'", gpaths[0].Abbreviation)
	}
}

func TestRunRestore_InputNotExist(t *testing.T) {
	if os.Getenv("TEST_RUN_RESTORE_NOT_EXIST") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		c.Flags().StringP("input", "i", "", "")
		c.Flags().Set("input", "/non/existent/file")

		cmd.RestoreCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestRunRestore_InputNotExist", "TEST_RUN_RESTORE_NOT_EXIST")
}

func TestRunRestore_InputIsDir(t *testing.T) {
	if os.Getenv("TEST_RUN_RESTORE_IS_DIR") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		tmpDir := os.TempDir()
		c.Flags().StringP("input", "i", "", "")
		c.Flags().Set("input", tmpDir)

		cmd.RestoreCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestRunRestore_InputIsDir", "TEST_RUN_RESTORE_IS_DIR")
}

func TestRunRestore_InvalidJSON(t *testing.T) {
	if os.Getenv("TEST_RUN_RESTORE_INVALID_JSON") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		tmpFile := filepath.Join(os.TempDir(), "restore_invalid.json")
		os.WriteFile(tmpFile, []byte("{invalid"), 0644)
		defer os.Remove(tmpFile)

		c.Flags().StringP("input", "i", "", "")
		c.Flags().Set("input", tmpFile)

		cmd.RestoreCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestRunRestore_InvalidJSON", "TEST_RUN_RESTORE_INVALID_JSON")
}

func TestCmd_Restore_ReadError(t *testing.T) {
	if os.Getenv("TEST_RESTORE_READ_ERROR") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		f, err := os.CreateTemp("", "restore_unreadable")
		if err != nil {
			os.Exit(1)
		}
		f.Close()
		defer os.Remove(f.Name())

		if err := os.Chmod(f.Name(), 0200); err != nil {
			// skip if can't chmod
} else {
// Need to make sure we don't fail IsDir check (it is file, so ok).
			// Need to set flag.
			c.Flags().StringP("input", "i", "", "")
			c.Flags().Set("input", f.Name())

			cmd.RestoreCmd.Run(c, []string{})
		}
		return
	}
	RunExpectedExit(t, "TestCmd_Restore_ReadError", "TEST_RESTORE_READ_ERROR")
}
