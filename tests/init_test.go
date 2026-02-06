package tests

import (
	"goto/src/core"
	"goto/src/utils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitializeConfig(t *testing.T) {
	// Setup temp home
	tmpHome, err := os.MkdirTemp("", "goto_test_home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	// Set HOME to tmpHome
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpHome)

	// Unset XDG_CONFIG_HOME to force usage of HOME/.config
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", oldXDG)
	os.Unsetenv("XDG_CONFIG_HOME")

	// Set SHELL to bash
	oldShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", oldShell)
	os.Setenv("SHELL", "/bin/bash")

	// Create .bashrc
	bashrcPath := filepath.Join(tmpHome, ".bashrc")
	f, err := os.Create(bashrcPath)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Re-initialize utils to pick up new HOME
	utils.SetupConfigFile()

	// Mock channel
	msgChan := make(chan core.Message, 100)

	// Run Init
	err = core.InitializeConfig(msgChan)
	if err != nil {
		t.Errorf("InitializeConfig failed: %v", err)
	}

	// Check if .bashrc was updated
	content, _ := os.ReadFile(bashrcPath)
	if !strings.Contains(string(content), "alias.sh") {
		t.Errorf(".bashrc does not contain source alias command: %s", string(content))
	}

	// Check if alias.sh exists (implicitly via .bashrc check, but better check file)
	// We need to guess the path since we can't access private vars.
	// It should be tmpHome/.config/goto/goto-run-testing/alias.sh
	expectedAliasPath := filepath.Join(tmpHome, ".config", "goto", "goto-run-testing", "alias.sh")
	if _, err := os.Stat(expectedAliasPath); os.IsNotExist(err) {
		t.Errorf("Alias file not created at %s", expectedAliasPath)
	}
}
