package gpath

import (
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

func TestCheckRepeatedItems(t *testing.T) {
	tests := []struct {
		name    string
		gpaths  []GotoPath
		wantErr bool
	}{
		{
			name:    "empty list",
			gpaths:  []GotoPath{},
			wantErr: true,
		},
		{
			name: "single item",
			gpaths: []GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
			},
			wantErr: false,
		},
		{
			name: "multiple unique items",
			gpaths: []GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
				{Path: "/path/2", Abbreviation: "p2"},
				{Path: "/path/3", Abbreviation: "p3"},
			},
			wantErr: false,
		},
		{
			name: "duplicate path",
			gpaths: []GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
				{Path: "/path/1", Abbreviation: "p2"},
			},
			wantErr: true,
		},
		{
			name: "duplicate abbreviation",
			gpaths: []GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
				{Path: "/path/2", Abbreviation: "p1"},
			},
			wantErr: true,
		},
		{
			name: "duplicate both path and abbreviation",
			gpaths: []GotoPath{
				{Path: "/path/1", Abbreviation: "p1"},
				{Path: "/path/1", Abbreviation: "p1"},
			},
			wantErr: true,
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
