package tests

import (
	"goto/src/gpath"
	"os"
	"path/filepath"
	"testing"
)

func TestValidPath(t *testing.T) {
	// Test existing directory
	cwd, _ := os.Getwd()
	valid, err := gpath.ValidPath(cwd)
	if err != nil {
		t.Errorf("Expected valid path for current directory, got error: %v", err)
	}
	if valid != cwd {
		t.Errorf("Expected cleaned path %s, got %s", cwd, valid)
	}

	// Test non-existent path
	_, err = gpath.ValidPath("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	// Test file path (not directory) via ValidPathVar logic if possible,
	// or create a temp file to test.
	tmpFile := filepath.Join(os.TempDir(), "goto_temp_file")
	os.WriteFile(tmpFile, []byte(""), 0666)
	defer os.Remove(tmpFile)

	_, err = gpath.ValidPath(tmpFile)
	if err == nil {
		t.Error("Expected error for file path (should be directory)")
	}
}

func TestValidAbbreviation(t *testing.T) {
	// Valid cases
	validInputs := []string{"docs", "work", "a", "my_path"}
	for _, input := range validInputs {
		cleaned, err := gpath.ValidAbbreviation(input)
		if err != nil {
			t.Errorf("Expected valid abbreviation for '%s', got error: %v", input, err)
		}
		if cleaned != input {
			t.Errorf("Expected cleaned abbreviation '%s', got '%s'", input, cleaned)
		}
	}

	// Invalid cases: spaces
	if _, err := gpath.ValidAbbreviation("with space"); err == nil {
		t.Error("Expected error for abbreviation with spaces")
	}

	// Invalid cases: number (assuming logic forbids pure numbers based on context)
	if _, err := gpath.ValidAbbreviation("123"); err == nil {
		t.Error("Expected error for numeric abbreviation")
	}

	// Check logic of trim
	cleaned, err := gpath.ValidAbbreviation("  trimmed  ")
	if err != nil {
		t.Errorf("Expected valid abbreviation for '  trimmed  ', got error: %v", err)
	}
	if cleaned != "trimmed" {
		t.Errorf("Expected trimmed 'trimmed', got '%s'", cleaned)
	}
}

func TestCheckRepeatedItems(t *testing.T) {
	// Valid list
	list := []gpath.GotoPath{
		{Path: "/a/b", Abbreviation: "ab"},
		{Path: "/c/d", Abbreviation: "cd"},
	}
	if err := gpath.CheckRepeatedItems(list); err != nil {
		t.Errorf("Expected valid list, got error: %v", err)
	}

	// Repeated Path
	listRepPath := []gpath.GotoPath{
		{Path: "/a/b", Abbreviation: "ab"},
		{Path: "/a/b", Abbreviation: "xy"},
	}
	if err := gpath.CheckRepeatedItems(listRepPath); err == nil {
		t.Error("Expected error for repeated path")
	}

	// Repeated Abbreviation
	listRepAbbv := []gpath.GotoPath{
		{Path: "/a/b", Abbreviation: "ab"},
		{Path: "/x/y", Abbreviation: "ab"},
	}
	if err := gpath.CheckRepeatedItems(listRepAbbv); err == nil {
		t.Error("Expected error for repeated abbreviation")
	}

	// Empty list check
	if err := gpath.CheckRepeatedItems([]gpath.GotoPath{}); err == nil {
		t.Error("Expected error for empty list")
	}
}

func TestIsValidIndex(t *testing.T) {
	// Valid cases
	if err := gpath.IsValidIndex(5, "0"); err != nil {
		t.Errorf("Expected valid index '0' for length 5, got: %v", err)
	}
	if err := gpath.IsValidIndex(5, "4"); err != nil {
		t.Errorf("Expected valid index '4' for length 5, got: %v", err)
	}

	// Invalid cases
	invalidInputs := []struct {
		length int
		index  string
		desc   string
	}{
		{5, "-1", "negative index"},
		{5, "5", "out of bounds index (high)"},
		{5, "10", "out of bounds index (very high)"},
		{5, "abc", "non-numeric index"},
		{0, "0", "zero length list"},
	}

	for _, tc := range invalidInputs {
		if err := gpath.IsValidIndex(tc.length, tc.index); err == nil {
			t.Errorf("Expected error for %s (len: %d, idx: %s)", tc.desc, tc.length, tc.index)
		}
	}
}
