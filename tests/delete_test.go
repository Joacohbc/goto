package tests

import (
	"goto/src/core"
	"goto/src/utils"
	"os"
	"testing"
)

func TestDeleteByAbbreviation(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add path
	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	// Delete by abbreviation
	_, err := core.DeletePath("", "p1", -1, false)
	if err != nil {
		t.Errorf("Failed to delete by abbreviation: %v", err)
	}

	// Verify
	gpaths, err := utils.LoadGPaths(false)
	if err != nil {
		t.Fatal(err)
	}
	for _, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			t.Error("Path 'p1' was not deleted")
		}
	}
}

func TestDeleteByIndex(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	// Add path
	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	// Find the index of p1
	gpaths, err := utils.LoadGPaths(false)
	if err != nil {
		t.Fatal(err)
	}

	targetIndex := -1
	for i, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		t.Fatal("Failed to find added path 'p1'")
	}

	// Delete by index
	_, err = core.DeletePath("", "", targetIndex, false)
	if err != nil {
		t.Errorf("Failed to delete by index %d: %v", targetIndex, err)
	}

	gpaths, err = utils.LoadGPaths(false)
	if err != nil {
		t.Fatal(err)
	}
	for _, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			t.Error("Path 'p1' was not deleted")
		}
	}
}

func TestDeleteByPath(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	cwd, _ := os.Getwd()
	_, err := core.DeletePath(cwd, "", -1, false)
	if err != nil {
		t.Errorf("Failed to delete by path: %v", err)
	}

	gpaths, err := utils.LoadGPaths(false)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, gp := range gpaths {
		if gp.Abbreviation == "p1" {
			found = true
		}
	}
	if found {
		t.Error("Path p1 was not deleted by path")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	_, err := core.DeletePath("", "non_existent", -1, false)
	if err == nil {
		t.Error("Expected error when deleting non-existent abbreviation")
	}
}

func TestDeleteIndexNotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	if err := core.AddPath(".", "p1", false); err != nil {
		t.Fatal(err)
	}

	_, err := core.DeletePath("", "", 100, false)
	if err == nil {
		t.Error("Expected error when deleting non-existent index")
	}
}

func TestDeleteNoFlags(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

    // Pass empty args
	_, err := core.DeletePath("", "", -1, false)
	if err == nil {
		t.Error("Expected error when no identifier provided")
	}
}

func TestDeletePathNotFound(t *testing.T) {
	_, cleanup := resetConfigFile(t, false)
	defer cleanup()

	cwd, _ := os.Getwd()
    // cwd is where we are running tests.
    // If we haven't added it, it's not there.
	_, err := core.DeletePath(cwd, "", -1, false)
	if err == nil {
		t.Error("Expected error when deleting non-existent path")
	}
}
