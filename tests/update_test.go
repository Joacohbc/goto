package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"testing"
)

func TestUpdatePathPath(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Update path by path
	updCmd := getTempCmd()
	updCmd.Flags().StringP(utils.FlagPath, "p", "", "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	cwd, _ := os.Getwd()
	updCmd.Flags().Set(utils.FlagPath, cwd)

	newDir, _ := os.MkdirTemp("", "goto_test_update_pp")
	defer os.RemoveAll(newDir)
	updCmd.Flags().Set("new", newDir)

	// Mode path-path or pp
	cmd.UpdateCmd.Run(updCmd, []string{"path-path"})

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
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	updCmd := getTempCmd()
	updCmd.Flags().StringP(utils.FlagPath, "p", "", "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	cwd, _ := os.Getwd()
	updCmd.Flags().Set(utils.FlagPath, cwd)
	updCmd.Flags().Set("new", "p1_new")

	cmd.UpdateCmd.Run(updCmd, []string{"path-abbv"})

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
	resetTempFile(t)
	c := getTempCmd()
	// Default entry is index 0. Add p1 -> index 1
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Now swap index 1 (p1) to 0.
	updCmd := getTempCmd()
	updCmd.Flags().StringP(utils.FlagPath, "p", "", "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	cwd, _ := os.Getwd()
	updCmd.Flags().Set(utils.FlagPath, cwd)
	updCmd.Flags().Set("new", "0") // set to index 0

	cmd.UpdateCmd.Run(updCmd, []string{"path-indx"})

	gpaths := utils.LoadGPaths(c)
	if len(gpaths) > 1 {
		// Index 0 should be "p1"
		if gpaths[0].Abbreviation != "p1" {
			t.Errorf("Expected p1 at index 0, got %s", gpaths[0].Abbreviation)
		}
	}
}

func TestUpdateAbbvPath(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	updCmd := getTempCmd()
	updCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	updCmd.Flags().Set(utils.FlagAbbreviation, "p1")

	newDir, _ := os.MkdirTemp("", "goto_test_update_ap")
	defer os.RemoveAll(newDir)
	updCmd.Flags().Set("new", newDir)

	cmd.UpdateCmd.Run(updCmd, []string{"abbv-path"})

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
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "oldname"})

	updCmd := getTempCmd()
	updCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	updCmd.Flags().Set(utils.FlagAbbreviation, "oldname")
	updCmd.Flags().Set("new", "newname")

	cmd.UpdateCmd.Run(updCmd, []string{"abbv-abbv"})

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
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	updCmd := getTempCmd()
	updCmd.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	updCmd.Flags().Set(utils.FlagAbbreviation, "p1")
	updCmd.Flags().Set("new", "0")

	cmd.UpdateCmd.Run(updCmd, []string{"abbv-indx"})

	gpaths := utils.LoadGPaths(c)
	if gpaths[0].Abbreviation != "p1" {
		t.Errorf("Expected p1 at index 0 after swap")
	}
}

func TestUpdateIndexPath(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Clean temp file implies "default_test_entry" at index 0. p1 is at index 1.

	updCmd := getTempCmd()
	updCmd.Flags().IntP(utils.FlagIndex, "i", -1, "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	// Update index 1 (p1)
	updCmd.Flags().Set(utils.FlagIndex, "1")

	newDir, _ := os.MkdirTemp("", "goto_test_update_ip")
	defer os.RemoveAll(newDir)

	updCmd.Flags().Set("new", newDir)

	cmd.UpdateCmd.Run(updCmd, []string{"indx-path"})

	gpaths := utils.LoadGPaths(c)
	// Check index 1
	if gpaths[1].Path != newDir {
		t.Error("Update failed for indx-path")
	}
}

func TestUpdateIndexAbbv(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"}) // index 1

	updCmd := getTempCmd()
	updCmd.Flags().IntP(utils.FlagIndex, "i", -1, "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	updCmd.Flags().Set(utils.FlagIndex, "1")
	updCmd.Flags().Set("new", "p1_updated")

	cmd.UpdateCmd.Run(updCmd, []string{"indx-abbv"})

	gpaths := utils.LoadGPaths(c)
	if gpaths[1].Abbreviation != "p1_updated" {
		t.Error("Update failed for indx-abbv")
	}
}

func TestUpdateIndexIndex(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"}) // index 1

	updCmd := getTempCmd()
	updCmd.Flags().IntP(utils.FlagIndex, "i", -1, "")
	updCmd.Flags().StringP("new", "n", "", "")
	updCmd.Flags().BoolP("modes", "m", false, "")

	updCmd.Flags().Set(utils.FlagIndex, "1")
	updCmd.Flags().Set("new", "0")

	// Swap 1 and 0
	cmd.UpdateCmd.Run(updCmd, []string{"indx-indx"})

	gpaths := utils.LoadGPaths(c)
	if gpaths[0].Abbreviation != "p1" {
		t.Error("Update failed for indx-indx -> p1 should be at 0")
	}
}

func TestUpdateInvalidMode(t *testing.T) {
	if os.Getenv("TEST_UPDATE_INVALID_MODE") == "1" {
		resetTempFile(t)
		updCmd := getTempCmd()
		updCmd.Flags().StringP("new", "n", "", "")
		updCmd.Flags().BoolP("modes", "m", false, "")

		updCmd.Flags().Set("new", "something")

		// Invalid mode
		cmd.UpdateCmd.Run(updCmd, []string{"invalid-mode"})
		return
	}
	RunExpectedExit(t, "TestUpdateInvalidMode", "TEST_UPDATE_INVALID_MODE")
}
