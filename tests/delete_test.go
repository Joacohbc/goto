package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"testing"
)

func TestDeleteByAbbreviation(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	// Add path
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Create context for delete
	delCmd := getTempCmd()
	delCmd.Flags().StringP(utils.FlagPath, "p", "", "")
	delCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	delCmd.Flags().IntP(utils.FlagIndex, "i", -1, "")

	delCmd.Flags().Set(utils.FlagAbbreviation, "p1")

	// Capture output to avoid polluting test logs
	captureOutput(func() {
		cmd.DeleteCmd.Run(delCmd, []string{})
	})

	// Verify
	gpaths := utils.LoadGPaths(c)
	// Expect 1 path (the default one)
	if len(gpaths) != 1 {
		t.Errorf("Expected 1 path after delete, got %d", len(gpaths))
	}
	if gpaths[0].Abbreviation == "p1" {
		t.Error("Path 'p1' was not deleted")
	}
}

func TestDeleteByIndex(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	// Add path
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Create context for delete
	delCmd := getTempCmd()
	delCmd.Flags().StringP(utils.FlagPath, "p", "", "")
	delCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	delCmd.Flags().IntP(utils.FlagIndex, "i", -1, "")

	// Index 0
	delCmd.Flags().Set(utils.FlagIndex, "0")

	captureOutput(func() {
		cmd.DeleteCmd.Run(delCmd, []string{})
	})

	gpaths := utils.LoadGPaths(c)
	// We deleted index 0 (default), so p1 should remain
	if len(gpaths) != 1 {
		t.Errorf("Expected 1 path after delete, got %d", len(gpaths))
	}
	if gpaths[0].Abbreviation != "p1" {
		t.Errorf("Expected remaining path to be 'p1', got '%s'", gpaths[0].Abbreviation)
	}
}

func TestDeleteByPath(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	cmd.AddCmd.Run(c, []string{".", "p1"})

	delCmd := getTempCmd()
	delCmd.Flags().StringP(utils.FlagPath, "p", "", "")

	cwd, _ := os.Getwd()
	delCmd.Flags().Set(utils.FlagPath, cwd)

	captureOutput(func() {
		cmd.DeleteCmd.Run(delCmd, []string{})
	})

	gpaths := utils.LoadGPaths(c)
	// Should delete p1. Default remains.
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			found = true
		}
	}
	if found {
		t.Error("Path p1 was not deleted by path")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	if os.Getenv("TEST_DELETE_NON_EXISTENT") == "1" {
		resetTempFile(t)

		delCmd := getTempCmd()
		delCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		delCmd.Flags().Set(utils.FlagAbbreviation, "non_existent")

		cmd.DeleteCmd.Run(delCmd, []string{})
		return
	}
	RunExpectedExit(t, "TestDeleteNonExistent", "TEST_DELETE_NON_EXISTENT")
}
