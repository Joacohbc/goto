package tests

import (
	"goto/src/cmd"
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
