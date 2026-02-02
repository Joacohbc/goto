package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"strings"
	"testing"
)

func TestSearchByAbbreviation(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cmd.AddCmd.Run(c, []string{".", "p1"})

	c.Flags().StringP(utils.FlagPath, "p", "", "")
	c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")

	c.Flags().Set(utils.FlagAbbreviation, "p1")

	output := captureOutput(func() {
		cmd.SearchCmd.Run(c, []string{})
	})

	if !strings.Contains(output, "p1") {
		t.Errorf("Expected 'p1' in search results, got: %s", output)
	}
}

func TestSearchByPath(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	c.Flags().StringP(utils.FlagPath, "p", "", "")
	cwd, _ := os.Getwd()
	c.Flags().Set(utils.FlagPath, cwd)

	output := captureOutput(func() {
		cmd.SearchCmd.Run(c, []string{})
	})

	if !strings.Contains(output, "p1") {
		t.Errorf("Expected 'p1' in search by path result, got: %v", output)
	}
}

func TestSearchNotFound(t *testing.T) {
	if os.Getenv("TEST_SEARCH_NOT_FOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		c.Flags().Set(utils.FlagAbbreviation, "nothere")

		cmd.SearchCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestSearchNotFound", "TEST_SEARCH_NOT_FOUND")
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

func TestRunSearch_PathNotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_SEARCH_PATH_NOT_FOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagPath, "p", "", "")
		c.Flags().Set(utils.FlagPath, "/non/existent/path")

		cmd.SearchCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestRunSearch_PathNotFound", "TEST_RUN_SEARCH_PATH_NOT_FOUND")
}

func TestRunSearch_AbbvNotFound(t *testing.T) {
	if os.Getenv("TEST_RUN_SEARCH_ABBV_NOT_FOUND") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()
		c.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		c.Flags().Set(utils.FlagAbbreviation, "nonexistent")

		cmd.SearchCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestRunSearch_AbbvNotFound", "TEST_RUN_SEARCH_ABBV_NOT_FOUND")
}
