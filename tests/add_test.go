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

	// Add current directory as "current"
	args := []string{".", "current"}
	cmd.AddCmd.Run(c, args)

	// Verify it was added
	gpaths := utils.LoadGPaths(c)
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "current" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find abbreviation 'current' in gpaths")
	}
}

func TestAddPathRepeated(t *testing.T) {
	// This block runs in the subprocess because of the environment variable check below.
	if os.Getenv("TEST_ADD_PATH_REPEATED_SUBPROCESS") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		// Add "current"
		cmd.AddCmd.Run(c, []string{".", "current"})

		// Add same again - should exit 1 because of repeated valid path logic in Validate/Save
		cmd.AddCmd.Run(c, []string{".", "current"})
		return
	}

	// Run the test in a subprocess using the helper
	RunExpectedExit(t, "TestAddPathRepeated", "TEST_ADD_PATH_REPEATED_SUBPROCESS")
}
