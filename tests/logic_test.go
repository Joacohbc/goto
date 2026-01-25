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

	// Test empty path
	if _, err := gpath.ValidPath(""); err == nil {
		t.Error("Expected error for empty path")
	}

	// Test path with only spaces
	if _, err := gpath.ValidPath("   "); err == nil {
		t.Error("Expected error for blank path")
	}
}

func TestValidAbbreviation(t *testing.T) {
	// Valid cases
	validInputs := []string{"docs", "work", "a", "my_path", "h1"}
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

	// Invalid cases: empty
	if _, err := gpath.ValidAbbreviation(""); err == nil {
		t.Error("Expected error for empty abbreviation")
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

	// Nil list check
	if err := gpath.CheckRepeatedItems(nil); err == nil {
		t.Error("Expected error for nil list")
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

func TestGetPathFromIndexOrAbbreviation(t *testing.T) {
	gpaths := []gpath.GotoPath{
		{Path: "/path/a", Abbreviation: "a"},
		{Path: "/path/b", Abbreviation: "b"},
		{Path: "/path/num", Abbreviation: "100"}, // Abbreviation that looks like a number
	}

	// Case 1: Valid Index
	got, found := gpath.GetPathFromIndexOrAbbreviation(gpaths, "0")
	if !found || got != "/path/a" {
		t.Errorf("Expected to find index 0 (/path/a), got %s, found=%v", got, found)
	}

	// Case 2: Valid Abbreviation
	got, found = gpath.GetPathFromIndexOrAbbreviation(gpaths, "b")
	if !found || got != "/path/b" {
		t.Errorf("Expected to find abbreviation 'b' (/path/b), got %s, found=%v", got, found)
	}

	// Case 3: Invalid Index (out of bounds) -> Should fall back to Abbv check -> Not found
	got, found = gpath.GetPathFromIndexOrAbbreviation(gpaths, "5")
	if found || got != "5" {
		t.Errorf("Expected not found for '5', got %s, found=%v", got, found)
	}

	// Case 4: Abbreviation that is a number (but out of bounds index)
	// '100' is out of bounds for len 3, so IsValidIndex fails.
	// Logic should find it as abbreviation.
	got, found = gpath.GetPathFromIndexOrAbbreviation(gpaths, "100")
	if !found || got != "/path/num" {
		t.Errorf("Expected to find abbreviation '100' (/path/num), got %s, found=%v", got, found)
	}

	// Case 5: Non-existent abbreviation
	got, found = gpath.GetPathFromIndexOrAbbreviation(gpaths, "missing")
	if found || got != "missing" {
		t.Errorf("Expected not found for 'missing', got %s, found=%v", got, found)
	}
}
