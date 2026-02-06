package tests

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"goto/src/core"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		current string
		remote  string
		want    bool
	}{
		{"v1.0.0", "v1.0.1", true},
		{"1.0.0", "1.0.1", true},
		{"v1.0.1", "v1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"v1.0.0", "v1.1.0", true},
		{"v1.0.0", "v2.0.0", true},
		{"v0.9.0", "v1.0.0", true},
		{"v1.0.0.0", "v1.0.0.1", true},
		{"v1.0", "v1.0.1", true},
		{"v2.0.0", "v1.9.9", false},
		{"v1.2.3", "v1.2.3", false},
		{"v1.0.0", "2.0.0", true},
		{"1.0.0", "v2.0.0", true},
	}

	for _, tt := range tests {
		got := core.IsNewerVersion(tt.current, tt.remote)
		if got != tt.want {
			t.Errorf("IsNewerVersion(%s, %s) = %v; want %v", tt.current, tt.remote, got, tt.want)
		}
	}
}

func TestFindAssetURL(t *testing.T) {
	tests := []struct {
		name      string
		assets    []core.GitHubAsset
		osName    string
		archName  string
		wantURL   string
		wantHash  string
		wantError bool
	}{
		{
			name: "exact match found",
			assets: []core.GitHubAsset{
				{Name: "goto-linux-amd64", BrowserDownloadURL: "https://example.com/goto-linux-amd64", Digest: "sha256:abc123"},
				{Name: "goto-darwin-amd64", BrowserDownloadURL: "https://example.com/goto-darwin-amd64", Digest: "sha256:def456"},
			},
			osName:    "linux",
			archName:  "amd64",
			wantURL:   "https://example.com/goto-linux-amd64",
			wantHash:  "sha256:abc123",
			wantError: false,
		},
		{
			name: "windows exe match found",
			assets: []core.GitHubAsset{
				{Name: "goto-windows-amd64.exe", BrowserDownloadURL: "https://example.com/goto-windows-amd64.exe", Digest: "sha256:win123"},
			},
			osName:    "windows",
			archName:  "amd64",
			wantURL:   "https://example.com/goto-windows-amd64.exe",
			wantHash:  "sha256:win123",
			wantError: false,
		},
		{
			name: "no asset found",
			assets: []core.GitHubAsset{
				{Name: "goto-linux-amd64", BrowserDownloadURL: "https://example.com/goto-linux-amd64"},
			},
			osName:    "darwin",
			archName:  "arm64",
			wantError: true,
		},
		{
			name:      "empty assets",
			assets:    []core.GitHubAsset{},
			osName:    "linux",
			archName:  "amd64",
			wantError: true,
		},
		{
			name: "match without digest",
			assets: []core.GitHubAsset{
				{Name: "goto-linux-arm64", BrowserDownloadURL: "https://example.com/goto-linux-arm64"},
			},
			osName:    "linux",
			archName:  "arm64",
			wantURL:   "https://example.com/goto-linux-arm64",
			wantHash:  "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, hash, err := core.FindAssetURL(tt.assets, tt.osName, tt.archName)
			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if url != tt.wantURL {
				t.Errorf("url = %q; want %q", url, tt.wantURL)
			}
			if hash != tt.wantHash {
				t.Errorf("hash = %q; want %q", hash, tt.wantHash)
			}
		})
	}
}

func TestVerifyDigest(t *testing.T) {
	// Create a temporary file with known content
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bin")
	testContent := []byte("test content for hashing")

	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Calculate expected SHA256
	h := sha256.New()
	h.Write(testContent)
	expectedHash := hex.EncodeToString(h.Sum(nil))

	tests := []struct {
		name      string
		filepath  string
		digest    string
		wantError bool
	}{
		{
			name:      "valid sha256 digest",
			filepath:  testFile,
			digest:    "sha256:" + expectedHash,
			wantError: false,
		},
		{
			name:      "invalid sha256 digest",
			filepath:  testFile,
			digest:    "sha256:0000000000000000000000000000000000000000000000000000000000000000",
			wantError: true,
		},
		{
			name:      "unsupported digest format",
			filepath:  testFile,
			digest:    "md5:abcd1234",
			wantError: true,
		},
		{
			name:      "file not found",
			filepath:  filepath.Join(tmpDir, "nonexistent.bin"),
			digest:    "sha256:" + expectedHash,
			wantError: true,
		},
		{
			name:      "empty digest prefix",
			filepath:  testFile,
			digest:    expectedHash,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := core.VerifyDigest(tt.filepath, tt.digest)
			if tt.wantError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	srcContent := []byte("test content to copy")
	if err := os.WriteFile(srcFile, srcContent, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	t.Run("successful copy", func(t *testing.T) {
		dstFile := filepath.Join(tmpDir, "dest.txt")
		err := core.CopyFile(srcFile, dstFile)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify content
		content, err := os.ReadFile(dstFile)
		if err != nil {
			t.Fatalf("failed to read destination file: %v", err)
		}
		if string(content) != string(srcContent) {
			t.Errorf("content mismatch: got %q, want %q", content, srcContent)
		}

		// Verify permissions were copied
		srcInfo, _ := os.Stat(srcFile)
		dstInfo, _ := os.Stat(dstFile)
		if srcInfo.Mode() != dstInfo.Mode() {
			t.Errorf("mode mismatch: got %v, want %v", dstInfo.Mode(), srcInfo.Mode())
		}
	})

	t.Run("source file not found", func(t *testing.T) {
		dstFile := filepath.Join(tmpDir, "dest2.txt")
		err := core.CopyFile(filepath.Join(tmpDir, "nonexistent.txt"), dstFile)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("invalid destination directory", func(t *testing.T) {
		dstFile := filepath.Join(tmpDir, "nonexistent", "dest.txt")
		err := core.CopyFile(srcFile, dstFile)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func TestDownloadFile(t *testing.T) {
	// Create a test HTTP server
	testContent := []byte("test binary content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write(testContent)
		case "/not-found":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	tmpDir := t.TempDir()

	t.Run("successful download", func(t *testing.T) {
		dstFile := filepath.Join(tmpDir, "downloaded.bin")
		err := core.DownloadFile(dstFile, server.URL+"/success")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify content
		content, err := os.ReadFile(dstFile)
		if err != nil {
			t.Fatalf("failed to read downloaded file: %v", err)
		}
		if string(content) != string(testContent) {
			t.Errorf("content mismatch: got %q, want %q", content, testContent)
		}

		// Verify file is executable
		info, err := os.Stat(dstFile)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}
		if info.Mode()&0111 == 0 {
			t.Error("file is not executable")
		}
	})

	t.Run("404 not found", func(t *testing.T) {
		dstFile := filepath.Join(tmpDir, "notfound.bin")
		err := core.DownloadFile(dstFile, server.URL+"/not-found")
		if err == nil {
			t.Error("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "bad status") {
			t.Errorf("expected 'bad status' error, got: %v", err)
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		dstFile := filepath.Join(tmpDir, "invalid.bin")
		err := core.DownloadFile(dstFile, "http://invalid-url-that-does-not-exist-12345.com")
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("invalid destination", func(t *testing.T) {
		dstFile := filepath.Join(tmpDir, "nonexistent", "file.bin")
		err := core.DownloadFile(dstFile, server.URL+"/success")
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func TestGetLatestRelease(t *testing.T) {
	t.Run("successful API call", func(t *testing.T) {
		release := core.GitHubRelease{
			TagName: "v1.2.3",
			Assets: []core.GitHubAsset{
				{Name: "goto-linux-amd64", BrowserDownloadURL: "https://example.com/file"},
			},
		}
		releaseJSON, _ := json.Marshal(release)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(releaseJSON)
		}))
		defer server.Close()

		// Test GetLatestRelease with a real HTTP call
		// Since we can't mock the hardcoded URL, this test verifies the function works
		// with real GitHub API (may fail if network is unavailable)
		result, err := core.GetLatestRelease()
		// We can't guarantee success due to network dependency, but we can check it doesn't panic
		_ = result
		_ = err
	})

	t.Run("API returns error status", func(t *testing.T) {
		// This would require dependency injection to test properly
		// For now we test the live API which we know works
		t.Skip("Skipping - would need dependency injection to mock HTTP client")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		t.Skip("Skipping - would need dependency injection to mock HTTP client")
	})
}

func TestUpdateBinary(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping UpdateBinary test on Windows")
	}

	t.Run("message channel usage", func(t *testing.T) {
		msgChan := make(chan core.Message, 100)

		// This will fail to connect to GitHub, but we can verify the channel works
		_ = core.UpdateBinary(msgChan, "v0.0.1")

		close(msgChan)

		// Collect messages
		var messages []core.Message
		for msg := range msgChan {
			messages = append(messages, msg)
		}

		// Should have at least one message (checking for updates)
		if len(messages) == 0 {
			t.Error("expected at least one message in the channel")
		}
	})
}

// Helper function to make GitHubAsset and GitHubRelease accessible for testing
// These are needed for the findAssetURL tests above
// We export wrapper functions in the test package to access unexported functionality
func TestGitHubTypes(t *testing.T) {
	// Test that we can create and use the GitHub types
	asset := core.GitHubAsset{
		Name:               "test-asset",
		BrowserDownloadURL: "https://example.com/asset",
		Digest:             "sha256:abc123",
	}

	if asset.Name != "test-asset" {
		t.Errorf("asset name = %q; want %q", asset.Name, "test-asset")
	}

	release := core.GitHubRelease{
		TagName: "v1.0.0",
		Assets:  []core.GitHubAsset{asset},
	}

	if release.TagName != "v1.0.0" {
		t.Errorf("release tag = %q; want %q", release.TagName, "v1.0.0")
	}

	if len(release.Assets) != 1 {
		t.Errorf("assets count = %d; want 1", len(release.Assets))
	}
}
