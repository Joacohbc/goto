package tests

import (
	"goto/src/gpath"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateGotoPathsFile(t *testing.T) {
	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "goto_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Define a path for the config file that does not exist
	configPath := filepath.Join(tmpDir, "goto-paths.json")

	// Call CreateGotoPathsFile
	err = gpath.CreateGotoPathsFile(configPath)
	if err != nil {
		t.Fatalf("Expected CreateGotoPathsFile to succeed, got %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	// Verify content (should have default paths)
	var gpathsList []gpath.GotoPath
	err = gpath.LoadGPathsFile(&gpathsList, configPath)
	if err != nil {
		t.Fatalf("Failed to load created file: %v", err)
	}

	if len(gpathsList) != 2 {
		t.Errorf("Expected 2 default paths, got %d", len(gpathsList))
	}

	// Verify it contains "h" and "config"
	foundHome := false
	foundConfig := false
	for _, gp := range gpathsList {
		if gp.Abbreviation == "h" {
			foundHome = true
		}
		if gp.Abbreviation == "config" {
			foundConfig = true
		}
	}

	if !foundHome {
		t.Error("Expected default 'h' abbreviation")
	}
	if !foundConfig {
		t.Error("Expected default 'config' abbreviation")
	}

	// Call it again (should do nothing as file exists)
	err = gpath.CreateGotoPathsFile(configPath)
	if err != nil {
		t.Fatalf("Expected CreateGotoPathsFile to succeed when file exists, got %v", err)
	}
}

func TestSaveAndLoadGPathsFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goto_test_save")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "saved_paths.json")

	// Create some paths
	originalPaths := []gpath.GotoPath{
		{Path: "/tmp/a", Abbreviation: "a"},
		{Path: "/tmp/b", Abbreviation: "b"},
	}

	// Save
	err = gpath.SaveGPathsFile(originalPaths, configPath)
	if err != nil {
		t.Fatalf("SaveGPathsFile failed: %v", err)
	}

	// Load
	var loadedPaths []gpath.GotoPath
	err = gpath.LoadGPathsFile(&loadedPaths, configPath)
	if err != nil {
		t.Fatalf("LoadGPathsFile failed: %v", err)
	}

	// Compare
	if len(loadedPaths) != len(originalPaths) {
		t.Fatalf("Expected %d paths, got %d", len(originalPaths), len(loadedPaths))
	}

	for i, gp := range originalPaths {
		if loadedPaths[i].Path != gp.Path || loadedPaths[i].Abbreviation != gp.Abbreviation {
			t.Errorf("Mismatch at index %d: expected %v, got %v", i, gp, loadedPaths[i])
		}
	}
}

func TestLoadGPathsFile_NonExistent(t *testing.T) {
	var gpathsList []gpath.GotoPath
	err := gpath.LoadGPathsFile(&gpathsList, "/non/existent/path/goto.json")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestSaveGPathsFile_InvalidDir(t *testing.T) {
	// Try to save to a directory that doesn't exist (and we're not creating it in SaveGPathsFile, only in Create)
	paths := []gpath.GotoPath{{Path: ".", Abbreviation: "test"}}
	err := gpath.SaveGPathsFile(paths, "/non/existent/dir/file.json")
	if err == nil {
		t.Error("Expected error when saving to non-existent directory")
	}
}
func TestGotoPath_Valid_Error(t *testing.T) {
	// Invalid path
	gp := gpath.GotoPath{
		Path:         "/non/existent/path",
		Abbreviation: "valid",
	}
	if err := gp.Valid(); err == nil {
		t.Error("Expected error for invalid path, got nil")
	}

	// Invalid abbreviation
	cwd, _ := os.Getwd()
	gp2 := gpath.GotoPath{
		Path:         cwd,
		Abbreviation: "", // invalid empty
	}
	if err := gp2.Valid(); err == nil {
		t.Error("Expected error for invalid abbreviation, got nil")
	}
}

func TestSaveGPathsFile_RepeatedError(t *testing.T) {
	cwd, _ := os.Getwd()
	gpaths := []gpath.GotoPath{
		{Path: cwd, Abbreviation: "one"},
		{Path: cwd, Abbreviation: "two"}, // Repeated path
	}

	tmpFile := filepath.Join(os.TempDir(), "goto_test_save_repeated.json")

	if err := gpath.SaveGPathsFile(gpaths, tmpFile); err == nil {
		t.Error("Expected error due to repeated items, got nil")
	}
}

func TestSaveGPathsFile_OpenFileError(t *testing.T) {
	cwd, _ := os.Getwd()
	gpaths := []gpath.GotoPath{{Path: cwd, Abbreviation: "one"}}

	// Try to save to a directory path instead of a file
	tmpDir := os.TempDir()

	if err := gpath.SaveGPathsFile(gpaths, tmpDir); err == nil {
		t.Error("Expected error when saving to a directory path, got nil")
	}
}

func TestLoadGPathsFile_OpenError(t *testing.T) {
	var gpaths []gpath.GotoPath
	// Non-existent file
	if err := gpath.LoadGPathsFile(&gpaths, "/non/existent/file.json"); err == nil {
		t.Error("Expected error loading non-existent file, got nil")
	}
}

func TestLoadGPathsFile_DecodeError(t *testing.T) {
	// Create a corrupted file
	tmpFile := filepath.Join(os.TempDir(), "goto_test_corrupt.json")
	os.WriteFile(tmpFile, []byte("{ invalid json"), 0644)
	defer os.Remove(tmpFile)

	var gpaths []gpath.GotoPath
	if err := gpath.LoadGPathsFile(&gpaths, tmpFile); err == nil {
		t.Error("Expected error decoding corrupted file, got nil")
	}
}

func TestCreateGotoPathsFile_MkdirError(t *testing.T) {
	// Pass a path where the directory cannot be created.
	// E.g. /proc/invalid/goto-paths.json or similar?
	// Or a file exists where the dir should be.

	tmpFile := filepath.Join(os.TempDir(), "goto_mkdir_test")
	// Create a file here
	os.WriteFile(tmpFile, []byte(""), 0644)
	defer os.Remove(tmpFile)

	// Now try to create a file inside this 'file' (treating it as dir)
	target := filepath.Join(tmpFile, "file.json")

	if err := gpath.CreateGotoPathsFile(target); err == nil {
		t.Error("Expected error when creating dir over a file, got nil")
	}
}
