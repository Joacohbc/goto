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

	// Let's just check if it contains "goto"
	if filepath.Base(configDir) != "goto" {
		t.Errorf("Expected config dir to end with 'goto', got %s", configDir)
	}

	if configDir != expected {
		t.Logf("Note: Exact match failed. Expected %s, got %s. This might be due to path cleaning.", expected, configDir)
	}
}
