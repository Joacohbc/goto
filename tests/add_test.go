package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"testing"
)

func TestAddPath(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add path
	cmd.AddCmd.Run(c, []string{".", "current"})

	// Verify
	gpaths, err := utils.LoadGPaths(utils.TemporalFlagPassed(c))
	if err != nil {
		t.Fatalf("Failed to load gpaths: %v", err)
	}

	found := false
	for _, p := range gpaths {
		if p.Abbreviation == "current" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected path 'current' to be added")
	}
}

func TestAddPathRepeated(t *testing.T) {
	if os.Getenv("TEST_ADD_PATH_REPEATED_SUBPROCESS") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		// Add first
		cmd.AddCmd.Run(c, []string{".", "current"})

		// Add same again - should exit 1 because of repeated valid path logic in Validate/Save
		cmd.AddCmd.Run(c, []string{".", "current"})
		return
	}

	// Run the test in a subprocess using the helper
	RunExpectedExit(t, "TestAddPathRepeated", "TEST_ADD_PATH_REPEATED_SUBPROCESS")
}
