package tests

import (
	"goto/src/core"
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
	_, cleanup := resetConfigFile(t, true)
	defer cleanup()

	// Add path with temporal=true
	if err := core.AddPath(".", "tempAdd", true); err != nil {
		t.Fatal(err)
	}

	// Verify it was added to temporal file
	gpaths, _ := utils.LoadGPaths(true)
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "tempAdd" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Path not added to temporal file")
	}
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

func TestGetSecureTempFile_Security(t *testing.T) {
	oldXDG := os.Getenv("XDG_RUNTIME_DIR")
	defer os.Setenv("XDG_RUNTIME_DIR", oldXDG)

	t.Run("XDG_RUNTIME_DIR absolute path requirement", func(t *testing.T) {
		os.Setenv("XDG_RUNTIME_DIR", "relative/path")
		path, err := utils.GetSecureTempFile()
		if err != nil {
			t.Fatalf("Expected success (fallback), got error: %v", err)
		}
		if filepath.HasPrefix(path, "relative/path") {
			t.Errorf("Expected relative XDG_RUNTIME_DIR to be ignored, but got %s", path)
		}
	})

	t.Run("XDG_RUNTIME_DIR secure creation", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "xdg_secure")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)
		os.Setenv("XDG_RUNTIME_DIR", tempDir)

		path, err := utils.GetSecureTempFile()
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if !filepath.HasPrefix(path, tempDir) {
			t.Errorf("Expected path %s to be under %s", path, tempDir)
		}

		dir := filepath.Dir(path)
		info, err := os.Lstat(dir)
		if err != nil {
			t.Fatalf("Failed to stat: %v", err)
		}
		if info.Mode().Perm() != 0700 {
			t.Errorf("Expected 0700 permissions, got %v", info.Mode().Perm())
		}
	})

	t.Run("Symlink rejection", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "xdg_symlink")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)
		os.Setenv("XDG_RUNTIME_DIR", tempDir)

		realDir := filepath.Join(tempDir, "real")
		if err := os.Mkdir(realDir, 0700); err != nil {
			t.Fatal(err)
		}
		linkDir := filepath.Join(tempDir, "goto-cli")
		if err := os.Symlink(realDir, linkDir); err != nil {
			t.Fatal(err)
		}

		path, err := utils.GetSecureTempFile()
		if err != nil {
			t.Fatalf("Expected success (fallback), got error: %v", err)
		}
		if filepath.HasPrefix(path, tempDir) {
			t.Error("Should have rejected symlinked XDG_RUNTIME_DIR and fallen back")
		}
	})
}
