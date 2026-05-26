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

	// Check if alias.sh exists
	expectedAliasPath := filepath.Join(tmpHome, ".config", "goto", "goto-run-testing", "alias.sh")
	if _, err := os.Stat(expectedAliasPath); os.IsNotExist(err) {
		t.Errorf("Alias file not created at %s", expectedAliasPath)
	}

	// Verify alias.sh contains autocompletion logic
	aliasContent, err := os.ReadFile(expectedAliasPath)
	if err != nil {
		t.Fatalf("Failed to read alias.sh: %v", err)
	}
	if !strings.Contains(string(aliasContent), "completion.bash") {
		t.Errorf("alias.sh does not contain completion.bash sourcing: %s", string(aliasContent))
	}
	if !strings.Contains(string(aliasContent), "completion.zsh") {
		t.Errorf("alias.sh does not contain completion.zsh sourcing: %s", string(aliasContent))
	}
}

func TestInitializeConfig_Shells(t *testing.T) {
	shells := []struct {
		shellName   string
		shellPath   string
		rcFilename  string
		expectError bool
		preCreate   bool
		preContent  string
	}{
		{shellName: "zsh", shellPath: "/bin/zsh", rcFilename: ".zshrc", expectError: false, preCreate: true},
		{shellName: "fish", shellPath: "/usr/bin/fish", rcFilename: filepath.Join(".config", "fish", "config.fish"), expectError: false, preCreate: true},
		{shellName: "dash", shellPath: "/bin/dash", rcFilename: ".profile", expectError: false, preCreate: true},
		{shellName: "tcsh", shellPath: "/bin/tcsh", rcFilename: ".tcshrc", expectError: false, preCreate: true},
		{shellName: "csh", shellPath: "/bin/csh", rcFilename: ".cshrc", expectError: false, preCreate: true},
		{shellName: "ksh", shellPath: "/bin/ksh", rcFilename: ".kshrc", expectError: false, preCreate: true},
		{shellName: "sh", shellPath: "/bin/sh", rcFilename: ".profile", expectError: false, preCreate: true},
		{shellName: "unsupported", shellPath: "/bin/unsupported", rcFilename: "", expectError: true, preCreate: false},
		{shellName: "no-newline", shellPath: "/bin/bash", rcFilename: ".bashrc", expectError: false, preCreate: true, preContent: "# existing content without newline"},
		{shellName: "with-newline", shellPath: "/bin/bash", rcFilename: ".bashrc", expectError: false, preCreate: true, preContent: "# existing content with newline\n"},
	}

	for _, tc := range shells {
		t.Run(tc.shellName, func(t *testing.T) {
			tmpHome, err := os.MkdirTemp("", "goto_shell_test_"+tc.shellName)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpHome)

			oldHome := os.Getenv("HOME")
			defer os.Setenv("HOME", oldHome)
			os.Setenv("HOME", tmpHome)

			oldShell := os.Getenv("SHELL")
			defer os.Setenv("SHELL", oldShell)
			os.Setenv("SHELL", tc.shellPath)

			oldXDG := os.Getenv("XDG_CONFIG_HOME")
			defer os.Setenv("XDG_CONFIG_HOME", oldXDG)
			os.Unsetenv("XDG_CONFIG_HOME")

			utils.SetupConfigFile()

			rcPath := ""
			if tc.preCreate {
				rcPath = filepath.Join(tmpHome, tc.rcFilename)
				if err := os.MkdirAll(filepath.Dir(rcPath), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(rcPath, []byte(tc.preContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			msgChan := make(chan core.Message, 100)
			err = core.InitializeConfig(msgChan)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error for shell %s but got nil", tc.shellName)
				}
				return
			}

			if err != nil {
				t.Fatalf("InitializeConfig failed for %s: %v", tc.shellName, err)
			}

			// Verify sentinel block was added
			content, err := os.ReadFile(rcPath)
			if err != nil {
				t.Fatalf("failed to read RC file %s: %v", rcPath, err)
			}

			contentStr := string(content)
			if !strings.Contains(contentStr, "# >>> goto initialize >>>") {
				t.Errorf("RC file %s does not contain sentinel start", rcPath)
			}
			if !strings.Contains(contentStr, "# <<< goto initialize <<<") {
				t.Errorf("RC file %s does not contain sentinel end", rcPath)
			}

			if tc.shellName == "no-newline" {
				if !strings.HasPrefix(contentStr, "# existing content without newline\n") {
					t.Errorf("expected newline before appending block, got: %q", contentStr)
				}
			}

			// Run InitializeConfig a second time to verify the sentinel block gets UPDATED instead of duplicated/appended
			err = core.InitializeConfig(msgChan)
			if err != nil {
				t.Fatalf("second InitializeConfig failed for %s: %v", tc.shellName, err)
			}

			contentAfterUpdate, err := os.ReadFile(rcPath)
			if err != nil {
				t.Fatal(err)
			}

			contentAfterUpdateStr := string(contentAfterUpdate)
			if strings.Count(contentAfterUpdateStr, "# >>> goto initialize >>>") != 1 {
				t.Errorf("expected exactly 1 sentinel block, got %d", strings.Count(contentAfterUpdateStr, "# >>> goto initialize >>>"))
			}
		})
	}
}
