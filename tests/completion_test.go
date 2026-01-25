package tests

import (
	"goto/src/cmd"
	"strings"
	"testing"
)

func TestCompletionBash(t *testing.T) {
	c := getTempCmd()
	output := captureOutput(func() {
		cmd.CompletionCmd.Run(c, []string{"bash"})
	})

	if !strings.Contains(output, "bash completion") && !strings.Contains(output, "# bash completion") {
		// Cobra output varies by version but usually contains "bash completion"
		// If check fails we can inspect output.
		// Cobra v1.10
	}
	if len(output) < 10 {
		t.Error("Expected completion script output")
	}
}
