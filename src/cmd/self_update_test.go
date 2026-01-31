package cmd

import (
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		current string
		remote  string
		want    bool
	}{
		{"1.0.0", "1.0.1", true},
		{"1.0.1", "1.0.0", false},
		{"1.0.0", "1.0.0", false},
		{"1.0", "1.0.1", true},
		{"1.0.1", "1.0", false},
		{"v1.0.0", "v1.0.1", true},
		{"1.0.0", "v1.0.1", true},
		{"2.0.0", "1.9.9", false},
	}

	for _, tt := range tests {
		got := isNewerVersion(tt.current, tt.remote)
		if got != tt.want {
			t.Errorf("isNewerVersion(%q, %q) = %v; want %v", tt.current, tt.remote, got, tt.want)
		}
	}
}

func TestFindAssetURL(t *testing.T) {
	assets := []GitHubAsset{
		{Name: "goto-linux-amd64", BrowserDownloadURL: "http://example.com/linux-amd64"},
		{Name: "goto-darwin-arm64", BrowserDownloadURL: "http://example.com/darwin-arm64"},
		{Name: "goto-windows-amd64.exe", BrowserDownloadURL: "http://example.com/windows-amd64"},
		{Name: "goto-linux-amd64.tar.gz", BrowserDownloadURL: "http://example.com/linux-amd64-archive"},
	}

	tests := []struct {
		osName   string
		archName string
		wantURL  string
		wantErr  bool
	}{
		{"linux", "amd64", "http://example.com/linux-amd64", false},
		{"darwin", "arm64", "http://example.com/darwin-arm64", false},
		{"windows", "amd64", "http://example.com/windows-amd64", false},
		{"linux", "arm64", "", true},
	}

	for _, tt := range tests {
		got, err := findAssetURL(assets, tt.osName, tt.archName)
		if (err != nil) != tt.wantErr {
			t.Errorf("findAssetURL(..., %q, %q) error = %v, wantErr %v", tt.osName, tt.archName, err, tt.wantErr)
			continue
		}
		if got != tt.wantURL {
			t.Errorf("findAssetURL(..., %q, %q) = %q; want %q", tt.osName, tt.archName, got, tt.wantURL)
		}
	}
}
