package tests

import (
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
	expected := filepath.Join(userConfigDir, "goto")

	// filepath.Join in init() might behave slightly differently regarding trailing slash
	// or cleaning depending on OS.
	// In the code: configDir = filepath.Join(configPath, "/goto/")
	// Normalize both paths before comparison to account for such differences.
	cleanConfigDir := filepath.Clean(configDir)
	cleanExpected := filepath.Clean(expected)

	if cleanConfigDir != cleanExpected {
		t.Errorf("Expected config dir to be %s, got %s", cleanExpected, cleanConfigDir)
	}
}
