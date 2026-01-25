package tests

import (
	"goto/src/cmd"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	output := captureOutput(func() {
		cmd.VersionCmd.Run(nil, nil)
	})

	if !strings.Contains(output, cmd.VersionGoto) {
		t.Errorf("Expected version %s in output, got: %s", cmd.VersionGoto, output)
	}
}
