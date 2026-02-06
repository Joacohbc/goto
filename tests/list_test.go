package tests

import (
	"goto/src/cmd"
	"goto/src/core"
	"testing"
)

func TestList(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add some paths first
	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	paths, err := core.ListPaths(false)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, p := range paths {
		if p.Abbreviation == "p1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'p1' in list")
	}
}

func TestListCmdFlags(t *testing.T) {
	if cmd.ListCmd.Flags().Lookup("reverse") == nil {
		t.Error("ListCmd should have 'reverse' flag")
	}
}
