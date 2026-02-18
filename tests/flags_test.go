package tests

import (
	"goto/src/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestGetPath(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	cmd := &cobra.Command{}
	cmd.Flags().StringP("path", "p", "", "")

	// Set a valid path (the temp directory we just created)
	_ = cmd.Flags().Set("path", tmpDir)

	// GetPath should return the absolute path
	result := utils.GetPath(cmd)

	// The result should be an absolute path to our temp directory
	if result == "" {
		t.Error("Expected GetPath to return a non-empty path")
	}

	// The path should be cleaned/absolute
	if !filepath.IsAbs(result) {
		t.Errorf("Expected GetPath to return an absolute path, got %s", result)
	}
}

func TestGetAbbreviation(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringP("abbv", "a", "", "")

	// Set a valid abbreviation (non-empty, no spaces, not a number)
	_ = cmd.Flags().Set("abbv", "myabbv")

	result := utils.GetAbbreviation(cmd)

	if result != "myabbv" {
		t.Errorf("Expected GetAbbreviation to return 'myabbv', got '%s'", result)
	}
}

func TestGetIndex(t *testing.T) {
	// We need to set up a temp file with some GotoPaths
	cmd := &cobra.Command{}
	cmd.Flags().IntP("indx", "i", -1, "")
	cmd.Flags().BoolP("temporal", "t", false, "")

	// Use temporal flag to avoid interfering with real config
	_ = cmd.Flags().Set("temporal", "true")

	// Create multiple unique temp directories for each entry
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()
	tmpDir3 := t.TempDir()

	// Set up a temp file with 3 entries (each with unique path and abbreviation)
	path := utils.GetFilePath(cmd.Flags().Changed(utils.FlagTemporal))
	content := `[{"Path":"` + tmpDir1 + `","Abbreviation":"entry0"},{"Path":"` + tmpDir2 + `","Abbreviation":"entry1"},{"Path":"` + tmpDir3 + `","Abbreviation":"entry2"}]`
	err := os.WriteFile(path, []byte(content), 0666)
	if err != nil {
		t.Fatal(err)
	}

	// Set a valid index (0, 1, or 2 would all be valid)
	_ = cmd.Flags().Set("indx", "1")

	result := utils.GetIndex(cmd)

	if result != 1 {
		t.Errorf("Expected GetIndex to return 1, got %d", result)
	}
}
