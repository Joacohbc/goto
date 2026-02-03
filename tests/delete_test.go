package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"testing"
)

func TestDeleteByAbbreviation(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add path
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Create context for delete
	c.Flags().StringP(utils.FlagPath, "p", "", "")
	c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	c.Flags().IntP(utils.FlagIndex, "i", -1, "")

	c.Flags().Set(utils.FlagAbbreviation, "p1")

	// Capture output to avoid polluting test logs
	captureOutput(func() {
		cmd.DeleteCmd.Run(c, []string{})
	})

	// Verify
	gpaths := utils.LoadGPaths(c)
	for _, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			t.Error("Path 'p1' was not deleted")
		}
	}
}

func TestDeleteByIndex(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add path
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Create context for delete
	c.Flags().StringP(utils.FlagPath, "p", "", "")
	c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	c.Flags().IntP(utils.FlagIndex, "i", -1, "")

	// Index 2 (third entry) should be p1 (default is 0, added path is 1)
	c.Flags().Set(utils.FlagIndex, "2")
	captureOutput(func() {
		cmd.DeleteCmd.Run(c, []string{})
	})

	gpaths := utils.LoadGPaths(c)
	for i, gp := range gpaths {
		if i == 2 && gp.Abbreviation == "p1" {
			t.Error("Path 'p1' was not deleted")
		}
	}
}

func TestDeleteByPath(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cmd.AddCmd.Run(c, []string{".", "p1"})

	c.Flags().StringP(utils.FlagPath, "p", "", "")

	cwd, _ := os.Getwd()
	c.Flags().Set(utils.FlagPath, cwd)

	captureOutput(func() {
		cmd.DeleteCmd.Run(c, []string{})
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
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		c.Flags().Set(utils.FlagAbbreviation, "non_existent")

		cmd.DeleteCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestDeleteNonExistent", "TEST_DELETE_NON_EXISTENT")
}

func TestCmd_Delete_IndexNotFound(t *testing.T) {
	if os.Getenv("TEST_DELETE_INDEX_NOT_FOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		// Add one path so we have something (index 0, 1)
		cmd.AddCmd.Run(c, []string{".", "p1"})

		c.Flags().IntP(utils.FlagIndex, "i", -1, "")

		// If we use index 100
		c.Flags().Set(utils.FlagIndex, "100")

		cmd.DeleteCmd.Run(c, []string{}) // Should exit
		return
	}
	RunExpectedExit(t, "TestCmd_Delete_IndexNotFound", "TEST_DELETE_INDEX_NOT_FOUND")
}

func TestCmd_Delete_DirectRunNoFlags(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Call Run directly without PreRun.
	captureOutput(func() {
		cmd.DeleteCmd.Run(c, []string{})
	})
}

func TestCmd_Delete_PathNotFound(t *testing.T) {
	if os.Getenv("TEST_DELETE_PATH_NOT_FOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		c.Flags().StringP(utils.FlagPath, "p", "", "")
		cwd, _ := os.Getwd()
		c.Flags().Set(utils.FlagPath, cwd) // Valid path on disk, not in default config

		cmd.DeleteCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestCmd_Delete_PathNotFound", "TEST_DELETE_PATH_NOT_FOUND")
}
