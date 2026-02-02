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
