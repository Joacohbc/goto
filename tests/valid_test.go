package tests

import (
	"goto/src/core"
	"goto/src/utils"
	"os"
	"testing"
)

func TestValidPaths(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	if err := core.AddPath(".", "validp", false); err != nil {
		t.Fatal(err)
	}

	if err := core.ValidatePaths(false); err != nil {
		t.Errorf("Expected paths to be valid, got: %v", err)
	}
}

func TestValidate_InvalidPathInFile(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Manually create a file with invalid path
	file := utils.GetFilePath(false)
	os.WriteFile(file, []byte(`[{"path":"/non/existent/path","abbreviation":"valid"}]`), 0600)

	if err := core.ValidatePaths(false); err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestValidate_RepeatedItems(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cwd, _ := os.Getwd()
	// Repeated abbreviation
	json := `[{"path":"` + cwd + `","abbreviation":"one"},{"path":"` + cwd + `","abbreviation":"one"}]`
	file := utils.GetFilePath(false)
	os.WriteFile(file, []byte(json), 0600)

	if err := core.ValidatePaths(false); err == nil {
		t.Error("Expected error for repeated items")
	}
}
