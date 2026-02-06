package tests

import (
	"goto/src/cmd"
	"goto/src/core"
	"goto/src/utils"
	"testing"
)

func TestAddPath(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add path
	err := core.AddPath(".", "current", false)
	if err != nil {
		t.Fatalf("Failed to add path: %v", err)
	}

	// Verify
	gpaths, err := utils.LoadGPaths(false)
	if err != nil {
		t.Fatalf("Failed to load gpaths: %v", err)
	}

	found := false
	for _, p := range gpaths {
		if p.Abbreviation == "current" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected path 'current' to be added")
	}
}

func TestAddPathRepeated(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add first
	if err := core.AddPath(".", "current", false); err != nil {
		t.Fatalf("First add failed: %v", err)
	}

	// Add same again - should fail
	if err := core.AddPath(".", "current", false); err == nil {
		t.Error("Expected error for repeated path, got nil")
	}
}

func TestAddCmdParams(t *testing.T) {
	// Verify that AddCmd requires exactly 2 arguments
	err := cmd.AddCmd.Args(cmd.AddCmd, []string{"one"})
	if err == nil {
		t.Error("AddCmd should return error for 1 argument")
	}

	err = cmd.AddCmd.Args(cmd.AddCmd, []string{"one", "two"})
	if err != nil {
		t.Error("AddCmd should accept 2 arguments")
	}

	err = cmd.AddCmd.Args(cmd.AddCmd, []string{"one", "two", "three"})
	if err == nil {
		t.Error("AddCmd should return error for 3 arguments")
	}
}
