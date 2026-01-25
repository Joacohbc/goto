package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"testing"
)

func TestAddPath(t *testing.T) {
	resetTempFile(t)

	// Use temporary command context
	c := getTempCmd()

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
