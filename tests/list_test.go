package tests

import (
	"goto/src/cmd"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add some paths first
	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Prepare list command
	// listCmd uses utils.LoadGPaths(cmd), so we pass our temp cmd
	// But listCmd is a global variable. Its Run method takes (cmd, args).
	// We can pass our cmd.

	output := captureOutput(func() {
		cmd.ListCmd.Run(c, []string{})
	})

	if !strings.Contains(output, "p1") {
		t.Errorf("Expected 'p1' in list output, got: %s", output)
	}
}

func TestListReverse(t *testing.T) {
	c, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cmd.AddCmd.Run(c, []string{".", "p1"})

	// Create a command with reverse flag
	c.Flags().BoolP("reverse", "R", false, "")
	c.Flags().Set("reverse", "true")

	output := captureOutput(func() {
		cmd.ListCmd.Run(c, []string{})
	})

	if !strings.Contains(output, "p1") {
		t.Errorf("Expected 'p1' in list output, got: %s", output)
	}
}
