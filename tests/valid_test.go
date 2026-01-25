package tests

import (
	"goto/src/cmd"
	"strings"
	"testing"
)

func TestValidPaths(t *testing.T) {
	resetTempFile(t)
	c := getTempCmd()

	cmd.AddCmd.Run(c, []string{".", "validp"})

	output := captureOutput(func() {
		cmd.ValidCmd.Run(c, []string{})
	})

	if !strings.Contains(output, "All paths are valid") {
		t.Errorf("Expected success message, got: %s", output)
	}
}
