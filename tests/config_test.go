package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	configDir := utils.GetConfigDir()

	if configDir == "" {
		t.Error("Expected config dir to be non-empty")
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(userConfigDir, utils.GOTO_FILE_DIR, utils.TESTING_FILE_DIR)

	// Normalize both paths before comparison to account for such differences.
	cleanConfigDir := filepath.Clean(configDir)
	cleanExpected := filepath.Clean(expected)

	if cleanConfigDir != cleanExpected {
		t.Errorf("Expected config dir to be %s, got %s", cleanExpected, cleanConfigDir)
	}
}
func TestUpdateGPaths_Temporal(t *testing.T) {
	c, cleanup := resetConfigFile(t, true)
	defer cleanup()

	args := []string{".", "tempAdd"}
	// We use AddCmd, but we need to reset its flags because it's a global variable?
	// AddCmd.Run uses `cmd` passed to it, which is `c`.
	// But AddCmd defines flags in `init()`.
	// The `c` we pass has flags set manually.

	// AddCmd calls utils.UpdateGPaths(cmd, gpaths)
	// UpdateGPaths checks `cmd.Flags().Changed(FlagTemporal)`

	cmd.AddCmd.Run(c, args)

	// If we are here, UpdateGPaths executed.
	// We covered the "if temporal" branch.
}

func TestGetSecureTempFile_XDG(t *testing.T) {
	oldXDG := os.Getenv("XDG_RUNTIME_DIR")
	defer func() {
		if oldXDG == "" {
			os.Unsetenv("XDG_RUNTIME_DIR")
		} else {
			os.Setenv("XDG_RUNTIME_DIR", oldXDG)
		}
		// Reset global state for other tests
		utils.SetupConfigFile()
	}()

	tempDir, err := os.MkdirTemp("", "xdg_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	os.Setenv("XDG_RUNTIME_DIR", tempDir)

	// This should succeed and use XDG dir
	utils.SetupConfigFile()
}

func TestGetSecureTempFile_Fallback(t *testing.T) {
	oldXDG := os.Getenv("XDG_RUNTIME_DIR")
	defer func() {
		if oldXDG == "" {
			os.Unsetenv("XDG_RUNTIME_DIR")
		} else {
			os.Setenv("XDG_RUNTIME_DIR", oldXDG)
		}
		utils.SetupConfigFile()
	}()

	os.Unsetenv("XDG_RUNTIME_DIR")

	// This should succeed and use fallback (os.TempDir/goto-cli-<uid>)
	utils.SetupConfigFile()
}
