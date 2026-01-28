package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"strings"
	"testing"
)

func TestSearchByAbbreviation(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	cmd.AddCmd.Run(c, []string{".", "p1"})

	searchCtx := getTempCmd()
	searchCtx.Flags().StringP(utils.FlagPath, "p", "", "")
	searchCtx.Flags().StringP(utils.FlagAbbreviation, "a", "", "")

	searchCtx.Flags().Set(utils.FlagAbbreviation, "p1")

	output := captureOutput(func() {
		cmd.SearchCmd.Run(searchCtx, []string{})
	})

	if !strings.Contains(output, "p1") {
		t.Errorf("Expected 'p1' in search results, got: %s", output)
	}
}

func TestSearchByPath(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()
	cmd.AddCmd.Run(c, []string{".", "p1"})

	searchCtx := getTempCmd()
	searchCtx.Flags().StringP(utils.FlagPath, "p", "", "")
	cwd, _ := os.Getwd()
	searchCtx.Flags().Set(utils.FlagPath, cwd)

	output := captureOutput(func() {
		cmd.SearchCmd.Run(searchCtx, []string{})
	})

	if !strings.Contains(output, "p1") {
		t.Errorf("Expected 'p1' in search by path result, got: %v", output)
	}
}

func TestSearchNotFound(t *testing.T) {
	if os.Getenv("TEST_SEARCH_NOT_FOUND") == "1" {
		resetTempFile(t)

		searchCtx := getTempCmd()
		searchCtx.Flags().StringP(utils.FlagAbbreviation, "a", "", "")
		searchCtx.Flags().Set(utils.FlagAbbreviation, "nothere")

		cmd.SearchCmd.Run(searchCtx, []string{})
		return
	}
	RunExpectedExit(t, "TestSearchNotFound", "TEST_SEARCH_NOT_FOUND")
}
