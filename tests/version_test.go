package tests

import (
	"goto/src/cmd"
	"testing"
)

func TestVersion(t *testing.T) {
	if cmd.VersionGoto == "" {
		t.Error("Version should not be empty")
	}
	if cmd.VersionCmd == nil {
		t.Error("VersionCmd should be defined")
	}
}
