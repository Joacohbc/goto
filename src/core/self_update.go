package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Digest             string `json:"digest"`
}

type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GitHubAsset `json:"assets"`
}

// UpdateBinary checks for updates and updates the binary if a newer version is available.
func UpdateBinary(currentVersion string) error {
	goos := runtime.GOOS
	if goos == "windows" {
		fmt.Println("Self-update is not supported on Windows.")
		return nil
	}

	// 1. Check for latest release via GitHub API
	fmt.Println("Checking for updates...")
	release, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	// 2. Compare versions
	newVersion := release.TagName
	cleanNewVersion := strings.TrimPrefix(newVersion, "v")
	cleanOldVersion := strings.TrimPrefix(currentVersion, "v")

	if !isNewerVersion(cleanOldVersion, cleanNewVersion) {
		fmt.Printf("You are already using the latest version (%s).\n", currentVersion)
		return nil
	}

	fmt.Printf("New version available: %s (current: %s)\n", newVersion, currentVersion)

	// 3. Find matching asset
	downloadURL, digest, err := findAssetURL(release.Assets, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("failed to find suitable binary for %s/%s: %w", runtime.GOOS, runtime.GOARCH, err)
	}

	// 4. Download to tmp
	tmpDir := os.TempDir()
	fileName := filepath.Base(downloadURL)
	tmpFilePath := filepath.Join(tmpDir, fileName)

	fmt.Printf("Downloading latest version from %s...\n", downloadURL)
	if err := downloadFile(tmpFilePath, downloadURL); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	defer func() {
		_ = os.Remove(tmpFilePath)
	}()

	// Verify digest if available
	if digest != "" {
		fmt.Println("Verifying download checksum...")
		if err := verifyDigest(tmpFilePath, digest); err != nil {
			return fmt.Errorf("checksum verification failed: %w", err)
		}
		fmt.Println("Checksum verified successfully.")
	}

	// 5. Replace binary
	currentExe, err := os.Executable()
	if err != nil {
		return err
	}

	// Resolve symlinks to ensure we are updating the real binary
	currentExe, err = filepath.EvalSymlinks(currentExe)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFilePath, currentExe)
	if err != nil {
		// If rename fails (e.g. diff filesystem), try copy.
		_ = os.Remove(currentExe)
		if err := copyFile(tmpFilePath, currentExe); err != nil {
			return fmt.Errorf("failed to replace binary: %w", err)
		}
	}

	// Ensure permissions on the new file (in case of copy)
	_ = os.Chmod(currentExe, 0755)

	fmt.Printf("Successfully updated from %s to %s\n", currentVersion, newVersion)
	return nil
}

func getLatestRelease() (*GitHubRelease, error) {
	url := "https://api.github.com/repos/Joacohbc/goto/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func findAssetURL(assets []GitHubAsset, osName, archName string) (string, string, error) {
	targetName := fmt.Sprintf("goto-%s-%s", osName, archName)

	for _, asset := range assets {
		if asset.Name == targetName {
			return asset.BrowserDownloadURL, asset.Digest, nil
		}
		if asset.Name == targetName+".exe" {
			return asset.BrowserDownloadURL, asset.Digest, nil
		}
	}

	return "", "", fmt.Errorf("no asset found matching %s or %s.exe", targetName, targetName)
}

func verifyDigest(filepath, digest string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	if strings.HasPrefix(digest, "sha256:") {
		expectedHash := strings.TrimPrefix(digest, "sha256:")
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		calculatedHash := hex.EncodeToString(h.Sum(nil))
		if calculatedHash != expectedHash {
			return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, calculatedHash)
		}
		return nil
	}

	// If digest format is unknown or not supported, we can either warn or error out.
	// Returning error is safer.
	return fmt.Errorf("unsupported digest format: %s", digest)
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Explicitly chmod the downloaded file to be executable before moving/copying
	return os.Chmod(filepath, 0755)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Copy mode
	si, err := os.Stat(src)
	if err == nil {
		err = os.Chmod(dst, si.Mode())
	}

	return err
}

func isNewerVersion(current, remote string) bool {
	current = strings.TrimPrefix(current, "v")
	remote = strings.TrimPrefix(remote, "v")

	partsC := strings.Split(current, ".")
	partsR := strings.Split(remote, ".")

	maxLen := len(partsC)
	if len(partsR) > maxLen {
		maxLen = len(partsR)
	}

	for i := 0; i < maxLen; i++ {
		valC := 0
		if i < len(partsC) {
			valC, _ = strconv.Atoi(partsC[i])
		}
		valR := 0
		if i < len(partsR) {
			valR, _ = strconv.Atoi(partsR[i])
		}

		if valR > valC {
			return true
		}
		if valR < valC {
			return false
		}
	}
	return false
}
