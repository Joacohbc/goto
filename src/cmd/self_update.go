package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type GithubRelease struct {
	TagName string `json:"tag_name"`
}

// UpdateBinaryCmd represents the update command for the binary itself
var UpdateBinaryCmd = &cobra.Command{
	Use:   "update-goto",
	Short: "Update goto to the latest version",
	Long:  `Downloads the latest release from GitHub and updates the current binary if a newer version is available.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateBinary()
	},
}

func init() {
	RootCmd.AddCommand(UpdateBinaryCmd)
}

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/Joacohbc/goto/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release GithubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return release.TagName, nil
}

func updateBinary() {
	// 1. Check version first
	fmt.Println("Checking for updates...")
	newVersion, err := getLatestVersion()
	if err != nil {
		fmt.Printf("Warning: Could not check for updates: %v\n", err)
		return
	}

	oldVersion := VersionGoto
	if !isNewerVersion(oldVersion, newVersion) {
		fmt.Printf("You are already using the latest version (%s). Remote: %s\n", oldVersion, newVersion)
		return
	}

	// 2. Construct URL
	goos := runtime.GOOS
	if goos == "windows" {
		fmt.Println("Self-update is not supported on Windows.")
		return
	}

	goarch := runtime.GOARCH
	extension := ""
	fileName := fmt.Sprintf("goto-%s-%s%s", goos, goarch, extension)
	downloadURL := fmt.Sprintf("https://github.com/Joacohbc/goto/releases/latest/download/%s", fileName)

	tmpDir := os.TempDir()
	tmpFilePath := filepath.Join(tmpDir, fileName)

	// 3. Download to tmp
	fmt.Printf("Downloading version %s from %s...\n", newVersion, downloadURL)
	if err := downloadFile(tmpFilePath, downloadURL); err != nil {
		cobra.CheckErr(fmt.Errorf("failed to download: %w", err))
	}
	// Attempt to clean up tmp file on exit
	defer func() {
		_ = os.Remove(tmpFilePath)
	}()

	// 4. Replace binary
	currentExe, err := os.Executable()
	cobra.CheckErr(err)

	// Resolve symlinks to ensure we are updating the real binary
	currentExe, err = filepath.EvalSymlinks(currentExe)
	cobra.CheckErr(err)

	// Make new binary executable
	err = os.Chmod(tmpFilePath, 0755)
	cobra.CheckErr(err)

	// 5. Replace binary
	// Prepare destination
	// Since currentExe is running, on Linux we can rename/move over it.
	// On Windows, we can't overwrite running executable easily.
	// The provided code assumes Linux environment per prompt context.

	err = os.Rename(tmpFilePath, currentExe)
	if err != nil {
		// If rename fails (e.g. diff filesystem), try copy.
		// On Linux, we cannot overwrite a running executable (ETXTBSY), so we remove it first.
		_ = os.Remove(currentExe)
		if err := copyFile(tmpFilePath, currentExe); err != nil {
			cobra.CheckErr(fmt.Errorf("failed to replace binary: %w", err))
		}
	}

	// Ensure permissions on the new file (in case of copy)
	_ = os.Chmod(currentExe, 0755)

	fmt.Printf("Successfully updated from %s to %s\n", oldVersion, newVersion)
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
	return err
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
