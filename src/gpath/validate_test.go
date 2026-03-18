package gpath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidIndex(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		index   string
		wantErr bool
	}{
		{"valid index 0", 5, "0", false},
		{"valid index middle", 5, "2", false},
		{"valid index last", 5, "4", false},
		{"invalid negative index", 5, "-1", true},
		{"invalid index equal to length", 5, "5", true},
		{"invalid index greater than length", 5, "10", true},
		{"invalid non-numeric index", 5, "abc", true},
		{"invalid empty index", 5, "", true},
		{"invalid index for zero length", 0, "0", true},
		{"valid index for length 1", 1, "0", false},
		{"invalid index for length 1", 1, "1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidIndex(tt.length, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidIndex(%d, %q) error = %v, wantErr %v", tt.length, tt.index, err, tt.wantErr)
			}
		})
	}
}

func TestValidPath(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "gpath-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "gpath-test-file-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid directory", tmpDir, false},
		{"empty path", "", true},
		{"blank space path", "   ", true},
		{"non-existent path", "/non/existent/path/at/all", true},
		{"path is a file", tmpFile.Name(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidPath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !filepath.IsAbs(got) {
					t.Errorf("ValidPath(%q) got = %q, want absolute path", tt.path, got)
				}
			}
		})
	}
}

func TestValidAbbreviation(t *testing.T) {
	tests := []struct {
		name    string
		abbv    string
		wantErr bool
	}{
		{"valid abbreviation", "work", false},
		{"empty abbreviation", "", true},
		{"blank space abbreviation", "   ", true},
		{"contains space", "my work", true},
		{"is a number", "123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidAbbreviation(tt.abbv)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidAbbreviation(%q) error = %v, wantErr %v", tt.abbv, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.abbv {
				t.Errorf("ValidAbbreviation(%q) got = %q, want %q", tt.abbv, got, tt.abbv)
			}
		})
	}
}

func TestCheckRepeatedItems(t *testing.T) {
	tests := []struct {
		name    string
		gpaths  []GotoPath
		wantErr bool
	}{
		{
			"valid items",
			[]GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
				{Path: "/path/2", Abbreviation: "p2"},
			},
			false,
		},
		{
			"empty list",
			[]GotoPath{},
			true,
		},
		{
			"repeated path",
			[]GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
				{Path: "/path/1", Abbreviation: "p2"},
			},
			true,
		},
		{
			"repeated abbreviation",
			[]GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
				{Path: "/path/2", Abbreviation: "p1"},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRepeatedItems(tt.gpaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckRepeatedItems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidIndex_ErrorMessages(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		index   string
		wantMsg string
	}{
		{"non-numeric", 5, "abc", "the Index must be a number"},
		{"empty list", 0, "0", "the Index 0 is invalid (the list is empty), check config file"},
		{"out of bounds", 5, "5", "the Index 5 is invalid (should be: 0-4), check config file"},
		{"negative", 5, "-1", "the Index -1 is invalid (should be: 0-4), check config file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidIndex(tt.length, tt.index)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantMsg {
				t.Errorf("got error %q, want %q", err.Error(), tt.wantMsg)
			}
		})
	}
}
