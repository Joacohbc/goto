package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"testing"
)

func TestUpdatePathPath(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Update path by path
	c.Flags().StringP(utils.FlagPath, "p", "", "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	cwd, _ := os.Getwd()
	c.Flags().Set(utils.FlagPath, cwd)

	newDir, _ := os.MkdirTemp("", "goto_test_update_pp")
	defer os.RemoveAll(newDir)
	c.Flags().Set("new", newDir)

	// Mode path-path or pp
	cmd.UpdateCmd.Run(c, []string{"path-path"})

	gpaths := utils.LoadGPaths(c)
	found := false
	for _, gp := range gpaths {
		if gp.Path == newDir && gp.Abbreviation == "p1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Path update failed for path-path")
	}
}

func TestUpdatePathAbbv(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	c.Flags().StringP(utils.FlagPath, "p", "", "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	cwd, _ := os.Getwd()
	c.Flags().Set(utils.FlagPath, cwd)
	c.Flags().Set("new", "p1_new")

	cmd.UpdateCmd.Run(c, []string{"path-abbv"})

	gpaths := utils.LoadGPaths(c)
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "p1_new" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Update failed for path-abbv")
	}
}

func TestUpdatePathIndex(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	// Default entry is index 0. Add p1 -> index 1
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Now swap index 1 (p1) to 0.
	c.Flags().StringP(utils.FlagPath, "p", "", "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	cwd, _ := os.Getwd()
	c.Flags().Set(utils.FlagPath, cwd)
	c.Flags().Set("new", "0") // set to index 0

	cmd.UpdateCmd.Run(c, []string{"path-indx"})

	gpaths := utils.LoadGPaths(c)
	if len(gpaths) > 1 {
		// Index 0 should be "p1"
		if gpaths[0].Abbreviation != "p1" {
			t.Errorf("Expected p1 at index 0, got %s", gpaths[0].Abbreviation)
		}
	}
}

func TestUpdateAbbvPath(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	c.Flags().Set(utils.FlagAbbreviation, "p1")

	newDir, _ := os.MkdirTemp("", "goto_test_update_ap")
	defer os.RemoveAll(newDir)
	c.Flags().Set("new", newDir)

	cmd.UpdateCmd.Run(c, []string{"abbv-path"})

	gpaths := utils.LoadGPaths(c)
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "p1" && gp.Path == newDir {
			found = true
			break
		}
	}
	if !found {
		t.Error("Update failed for abbv-path")
	}
}

func TestUpdateAbbvAbbv(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "oldname"})

	c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	c.Flags().Set(utils.FlagAbbreviation, "oldname")
	c.Flags().Set("new", "newname")

	cmd.UpdateCmd.Run(c, []string{"abbv-abbv"})

	gpaths := utils.LoadGPaths(c)
	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "newname" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find abbreviation 'newname'")
	}
}

func TestUpdateAbbvIndex(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	c.Flags().Set(utils.FlagAbbreviation, "p1")
	c.Flags().Set("new", "0")

	cmd.UpdateCmd.Run(c, []string{"abbv-indx"})

	gpaths := utils.LoadGPaths(c)
	if gpaths[0].Abbreviation != "p1" {
		t.Errorf("Expected p1 at index 0 after swap")
	}
}

func TestUpdateIndexPath(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Clean temp file implies "default_test_entry" at index 0. p1 is at index 1.

	c.Flags().IntP(utils.FlagIndex, "i", -1, "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	// Update index 1 (p1)
	c.Flags().Set(utils.FlagIndex, "1")

	newDir, _ := os.MkdirTemp("", "goto_test_update_ip")
	defer os.RemoveAll(newDir)

	c.Flags().Set("new", newDir)

	cmd.UpdateCmd.Run(c, []string{"indx-path"})

	gpaths := utils.LoadGPaths(c)
	// Check index 1
	if gpaths[1].Path != newDir {
		t.Error("Update failed for indx-path")
	}
}

func TestUpdateIndexAbbv(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"}) // index 1

	c.Flags().IntP(utils.FlagIndex, "i", -1, "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	c.Flags().Set(utils.FlagIndex, "1")
	c.Flags().Set("new", "p1_updated")

	cmd.UpdateCmd.Run(c, []string{"indx-abbv"})

	gpaths := utils.LoadGPaths(c)
	if gpaths[1].Abbreviation != "p1_updated" {
		t.Error("Update failed for indx-abbv")
	}
}

func TestUpdateIndexIndex(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"}) // index 1

	c.Flags().IntP(utils.FlagIndex, "i", -1, "")
	c.Flags().StringP("new", "n", "", "")
	c.Flags().BoolP("modes", "m", false, "")

	c.Flags().Set(utils.FlagIndex, "2")
	c.Flags().Set("new", "0")

	// Swap 1 and 0
	cmd.UpdateCmd.Run(c, []string{"indx-indx"})

	gpaths := utils.LoadGPaths(c)
	if gpaths[0].Abbreviation != "p1" {
		t.Error("Update failed for indx-indx -> p1 should be at 0")
	}
}

func TestUpdateInvalidMode(t *testing.T) {
	if os.Getenv("TEST_UPDATE_INVALID_MODE") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP("new", "n", "", "")
		c.Flags().BoolP("modes", "m", false, "")

		c.Flags().Set("new", "something")

		// Invalid mode
		cmd.UpdateCmd.Run(c, []string{"invalid-mode"})
		return
	}
	RunExpectedExit(t, "TestUpdateInvalidMode", "TEST_UPDATE_INVALID_MODE")
}

func TestPreRunUpdate_NoArgsNoMode(t *testing.T) {
	if os.Getenv("TEST_PRERUN_UPDATE_NOARGS") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		// No args, no flags
		cmd.UpdateCmd.PreRun(c, []string{})
		return
	}
	RunExpectedExit(t, "TestPreRunUpdate_NoArgsNoMode", "TEST_PRERUN_UPDATE_NOARGS")
}

func TestPreRunUpdate_NoNewFlag(t *testing.T) {
	if os.Getenv("TEST_PRERUN_UPDATE_NONEW") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().BoolP("modes", "m", false, "")
		c.Flags().StringP("new", "n", "", "")
		// Arg present, but no new flag
		cmd.UpdateCmd.PreRun(c, []string{"path-path"})
		return
	}
	RunExpectedExit(t, "TestPreRunUpdate_NoNewFlag", "TEST_PRERUN_UPDATE_NONEW")
}

func TestRunUpdate_ModesFlag(t *testing.T) {
	// This does not exit, just prints
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	c.Flags().BoolP("modes", "m", false, "")
	c.Flags().Set("modes", "true")

	// Helper to capture stdout if we wanted, but not strictly necessary for coverage
	cmd.UpdateCmd.Run(c, []string{})
}

func TestRunUpdate_PathPath_NotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_UPDATE_PP_NOTFOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagPath, "p", "", "")
		c.Flags().StringP("new", "n", "", "")
		c.Flags().BoolP("modes", "m", false, "")

		c.Flags().Set("new", ".")
		c.Flags().Set(utils.FlagPath, "/not/found")

		cmd.UpdateCmd.Run(c, []string{"path-path"})
		return
	}
	RunExpectedExit(t, "TestRunUpdate_PathPath_NotFound", "TEST_RUN_UPDATE_PP_NOTFOUND")
}

func TestRunUpdate_PathAbbv_NotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_UPDATE_PA_NOTFOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagPath, "p", "", "")
		c.Flags().StringP("new", "n", "", "")
		c.Flags().BoolP("modes", "m", false, "")

		c.Flags().Set("new", "newabbv")
		c.Flags().Set(utils.FlagPath, "/not/found")

		cmd.UpdateCmd.Run(c, []string{"path-abbv"})
		return
	}
	RunExpectedExit(t, "TestRunUpdate_PathAbbv_NotFound", "TEST_RUN_UPDATE_PA_NOTFOUND")
}

func TestRunUpdate_PathIndex_NotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_UPDATE_PI_NOTFOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagPath, "p", "", "")
		c.Flags().StringP("new", "n", "", "")
		c.Flags().BoolP("modes", "m", false, "")

		c.Flags().Set("new", "0")
		c.Flags().Set(utils.FlagPath, "/not/found")

		cmd.UpdateCmd.Run(c, []string{"path-indx"})
		return
	}
	RunExpectedExit(t, "TestRunUpdate_PathIndex_NotFound", "TEST_RUN_UPDATE_PI_NOTFOUND")
}

func TestRunUpdate_AbbvPath_NotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_UPDATE_AP_NOTFOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		c.Flags().StringP("new", "n", "", "")
		c.Flags().BoolP("modes", "m", false, "")

		c.Flags().Set("new", ".")
		c.Flags().Set(utils.FlagAbbreviation, "notfound")

		cmd.UpdateCmd.Run(c, []string{"abbv-path"})
		return
	}
	RunExpectedExit(t, "TestRunUpdate_AbbvPath_NotFound", "TEST_RUN_UPDATE_AP_NOTFOUND")
}

func TestRunUpdate_AbbvAbbv_NotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_UPDATE_AA_NOTFOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		c.Flags().StringP("new", "n", "", "")
		c.Flags().BoolP("modes", "m", false, "")

		c.Flags().Set("new", "newabbv")
		c.Flags().Set(utils.FlagAbbreviation, "notfound")

		cmd.UpdateCmd.Run(c, []string{"abbv-abbv"})
		return
	}
	RunExpectedExit(t, "TestRunUpdate_AbbvAbbv_NotFound", "TEST_RUN_UPDATE_AA_NOTFOUND")
}

func TestRunUpdate_AbbvIndex_NotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_UPDATE_AI_NOTFOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		c.Flags().StringP("new", "n", "", "")
		c.Flags().BoolP("modes", "m", false, "")

		c.Flags().Set("new", "0")
		c.Flags().Set(utils.FlagAbbreviation, "notfound")

		cmd.UpdateCmd.Run(c, []string{"abbv-indx"})
		return
	}
	RunExpectedExit(t, "TestRunUpdate_AbbvIndex_NotFound", "TEST_RUN_UPDATE_AI_NOTFOUND")
}
