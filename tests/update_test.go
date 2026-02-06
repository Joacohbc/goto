package tests

import (
	"goto/src/cmd"
	"goto/src/core"
	"goto/src/utils"
	"os"
	"strconv"
	"testing"

	"github.com/spf13/cobra"
)

func TestUpdatePathPath(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	cwd, _ := os.Getwd()
	newDir, _ := os.MkdirTemp("", "goto_test_update_pp")
	defer os.RemoveAll(newDir)

	// Update path by path (pp)
	if err := core.UpdatePath("path-path", cwd, "", -1, newDir, false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ := utils.LoadGPaths(false)
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
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	cwd, _ := os.Getwd()

	if err := core.UpdatePath("path-abbv", cwd, "", -1, "p1_new", false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ := utils.LoadGPaths(false)
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
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	// Add p1
	core.AddPath(".", "p1", false)

	newDir2, _ := os.MkdirTemp("", "goto_test_update_pi_2")
	defer os.RemoveAll(newDir2)
	core.AddPath(newDir2, "p2", false)

	// Find p1 index
	gpaths, _ := utils.LoadGPaths(false)
	p1Index := -1
	for i, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			p1Index = i
			break
		}
	}

	// Find p2 index
	p2Index := -1
	for i, gp := range gpaths {
		if gp.Abbreviation == "p2" {
			p2Index = i
			break
		}
	}

	if p1Index == -1 || p2Index == -1 {
		t.Fatal("Need p1 and p2")
	}

	// Swap p1 to p2's index
	cwd, _ := os.Getwd()

	err := core.UpdatePath("path-indx", cwd, "", -1, strconv.Itoa(p2Index), false)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ = utils.LoadGPaths(false)
	if gpaths[p2Index].Abbreviation != "p1" {
		t.Errorf("Expected p1 at index %d", p2Index)
	}
}

func TestUpdateAbbvPath(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	core.AddPath(".", "p1", false)

	newDir, _ := os.MkdirTemp("", "goto_test_update_ap")
	defer os.RemoveAll(newDir)

	if err := core.UpdatePath("abbv-path", "", "p1", -1, newDir, false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ := utils.LoadGPaths(false)
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
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	core.AddPath(".", "oldname", false)

	if err := core.UpdatePath("abbv-abbv", "", "oldname", -1, "newname", false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ := utils.LoadGPaths(false)
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
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	core.AddPath(".", "p1", false)

	newDir2, _ := os.MkdirTemp("", "goto_test_update_ai_2")
	defer os.RemoveAll(newDir2)
	core.AddPath(newDir2, "p2", false)

	// swap p1 and p2 via abbv-indx
	// find p2 index
	gpaths, _ := utils.LoadGPaths(false)
	p2Index := -1
	for i, gp := range gpaths {
		if gp.Abbreviation == "p2" {
			p2Index = i
			break
		}
	}

	if err := core.UpdatePath("abbv-indx", "", "p1", -1, strconv.Itoa(p2Index), false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ = utils.LoadGPaths(false)
	if gpaths[p2Index].Abbreviation != "p1" {
		t.Errorf("Expected p1 at index %d", p2Index)
	}
}

func TestUpdateIndexPath(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	core.AddPath(".", "p1", false)

	// find index of p1
	gpaths, _ := utils.LoadGPaths(false)
	p1Index := -1
	for i, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			p1Index = i
			break
		}
	}

	newDir, _ := os.MkdirTemp("", "goto_test_update_ip")
	defer os.RemoveAll(newDir)

	if err := core.UpdatePath("indx-path", "", "", p1Index, newDir, false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ = utils.LoadGPaths(false)
	if gpaths[p1Index].Path != newDir {
		t.Error("Update failed for indx-path")
	}
}

func TestUpdateIndexAbbv(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	core.AddPath(".", "p1", false)

	gpaths, _ := utils.LoadGPaths(false)
	p1Index := -1
	for i, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			p1Index = i
			break
		}
	}

	if err := core.UpdatePath("indx-abbv", "", "", p1Index, "p1_updated", false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ = utils.LoadGPaths(false)
	if gpaths[p1Index].Abbreviation != "p1_updated" {
		t.Error("Update failed for indx-abbv")
	}
}

func TestUpdateIndexIndex(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()
	core.AddPath(".", "p1", false)
	newDir2, _ := os.MkdirTemp("", "goto_test_update_ii_2")
	defer os.RemoveAll(newDir2)
	core.AddPath(newDir2, "p2", false)

	gpaths, _ := utils.LoadGPaths(false)
	p1Index := -1
	p2Index := -1
	for i, gp := range gpaths {
		if gp.Abbreviation == "p1" { p1Index = i }
		if gp.Abbreviation == "p2" { p2Index = i }
	}

	// Swap p1(index) to p2(value)
	if err := core.UpdatePath("indx-indx", "", "", p1Index, strconv.Itoa(p2Index), false); err != nil {
		t.Errorf("Update failed: %v", err)
	}

	gpaths, _ = utils.LoadGPaths(false)
	if gpaths[p2Index].Abbreviation != "p1" {
		t.Error("Update failed for indx-indx -> p1 should be at p2's old index")
	}
}

func TestUpdateInvalidMode(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	err := core.UpdatePath("invalid-mode", "", "", -1, "val", false)
	if err == nil {
		t.Error("Expected error for invalid mode")
	}
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

	cmd.UpdateCmd.Run(c, []string{})
}

func TestUpdate_PathPath_NotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	err := core.UpdatePath("path-path", "/not/found", "", -1, ".", false)
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestUpdate_PathAbbv_NotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	err := core.UpdatePath("path-abbv", "/not/found", "", -1, "newabbv", false)
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestUpdate_PathIndex_NotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	err := core.UpdatePath("path-indx", "/not/found", "", -1, "0", false)
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestUpdate_AbbvPath_NotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	err := core.UpdatePath("abbv-path", "", "notfound", -1, ".", false)
	if err == nil {
		t.Error("Expected error for non-existent abbreviation")
	}
}

func TestUpdate_AbbvAbbv_NotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	err := core.UpdatePath("abbv-abbv", "", "notfound", -1, "newabbv", false)
	if err == nil {
		t.Error("Expected error for non-existent abbreviation")
	}
}

func TestUpdate_AbbvIndex_NotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	err := core.UpdatePath("abbv-indx", "", "notfound", -1, "0", false)
	if err == nil {
		t.Error("Expected error for non-existent abbreviation")
	}
}

func TestUpdatePreRun_Success(t *testing.T) {
	preRun := cmd.UpdateCmd.PreRun
	// Mock command just for flags
	c := &cobra.Command{}
	c.Flags().BoolP("modes", "m", false, "")
	c.Flags().StringP("new", "n", "", "")

	// Case 1: Args present + new flag present
	// Must pass
	c.Flags().Set("new", "val")
	preRun(c, []string{"pp"})

	// Case 2: No args + modes flag present + new flag present
	c.Flags().Set("modes", "true")
	preRun(c, []string{})
}
