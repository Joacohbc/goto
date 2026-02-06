package tests

import (
	"goto/src/cmd"
	"goto/src/utils"
	"os"
	"strings"
	"testing"
)

func TestValidPaths(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cmd.AddCmd.Run(c, []string{".", "validp"})

	output := captureOutput(func() {
		cmd.ValidCmd.Run(c, []string{})
	})

	if !strings.Contains(output, "All paths are valid") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

func TestRunValid_InvalidPathInFile(t *testing.T) {
	if os.Getenv("TEST_RUN_VALID_INVALID_PATH") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		// Manually create a file with invalid path
		file := utils.GetFilePath(utils.TemporalFlagPassed(c))
		os.WriteFile(file, []byte(`[{"path":"/non/existent/path","abbreviation":"valid"}]`), 0600)

		cmd.ValidCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestRunValid_InvalidPathInFile", "TEST_RUN_VALID_INVALID_PATH")
}

func TestRunValid_RepeatedItems(t *testing.T) {
	if os.Getenv("TEST_RUN_VALID_REPEATED") == "1" {
		c, cleanup := resetConfigFile(t, false)
		defer cleanup()

		cwd, _ := os.Getwd()
		// Repeated abbreviation
		json := `[{"path":"` + cwd + `","abbreviation":"one"},{"path":"` + cwd + `","abbreviation":"one"}]`
		file := utils.GetFilePath(utils.TemporalFlagPassed(c))
		os.WriteFile(file, []byte(json), 0600)

		cmd.ValidCmd.Run(c, []string{})
		return
	}
	RunExpectedExit(t, "TestRunValid_RepeatedItems", "TEST_RUN_VALID_REPEATED")
}
