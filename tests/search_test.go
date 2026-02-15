package tests

import (
	"goto/src/cmd"
	"goto/src/core"
	"goto/src/utils"
	"os"
	"testing"
)

func TestSearchByAbbreviation(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	_, gp, err := core.SearchPath("", "p1", false)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if gp.Abbreviation != "p1" {
		t.Errorf("Expected p1, got %s", gp.Abbreviation)
	}
}

func TestSearchByPath(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	cwd, _ := os.Getwd()
	_, gp, err := core.SearchPath(cwd, "", false)
	if err != nil {
		t.Fatalf("Search by path failed: %v", err)
	}
	if gp.Path != cwd {
		t.Errorf("Expected path %s, got %s", cwd, gp.Path)
	}
}

func TestSearchNotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	_, _, err := core.SearchPath("", "nothere", false)
	if err == nil {
		t.Error("Expected error when searching for non-existent abbreviation")
	}
}

func TestRunSearch_PathNotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	_, _, err := core.SearchPath("/non/existent/path", "", false)
	if err == nil {
		t.Error("Expected error when searching for non-existent path")
	}
}

func TestCmd_Search_PathNotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cwd, _ := os.Getwd()
	_, _, err := core.SearchPath(cwd, "", false)
	if err == nil {
		t.Error("Expected error when searching for path not in config")
	}
}

func TestPreRunSearch_NoFlags(t *testing.T) {
	if os.Getenv("TEST_PRERUN_SEARCH_NOFLAGS") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		// No flags set
		cmd.SearchCmd.PreRun(c, []string{})
		return
	}
	RunExpectedExit(t, "TestPreRunSearch_NoFlags", "TEST_PRERUN_SEARCH_NOFLAGS")
}

func TestPreRunSearch_TooManyFlags(t *testing.T) {
	if os.Getenv("TEST_PRERUN_SEARCH_TOOMANY") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagPath, "p", "", "")
		c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		c.Flags().Set(utils.FlagPath, ".")
		c.Flags().Set(utils.FlagAbbreviation, "b")
		c.Flags().Set("temporal", "true") // 3 flags

		cmd.SearchCmd.PreRun(c, []string{})
		return
	}
	RunExpectedExit(t, "TestPreRunSearch_TooManyFlags", "TEST_PRERUN_SEARCH_TOOMANY")
}

func TestPreRunSearch_JustTemporal(t *testing.T) {
	if os.Getenv("TEST_PRERUN_SEARCH_JUST_TEMPORAL") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().Set("temporal", "true")

		cmd.SearchCmd.PreRun(c, []string{})
		return
	}
	RunExpectedExit(t, "TestPreRunSearch_JustTemporal", "TEST_PRERUN_SEARCH_JUST_TEMPORAL")
}
