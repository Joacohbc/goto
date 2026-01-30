package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

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

func updateBinary() {
	// 1. Construct URL
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

	// 2. Download to tmp
	fmt.Printf("Downloading latest version from %s...\n", downloadURL)
	if err := downloadFile(tmpFilePath, downloadURL); err != nil {
		cobra.CheckErr(fmt.Errorf("failed to download: %w", err))
	}
	// Attempt to clean up tmp file on exit
	defer func() {
		_ = os.Remove(tmpFilePath)
	}()

	// 3. Compare SHA256 of tmp file and current executable
	currentExe, err := os.Executable()
	cobra.CheckErr(err)

	// Resolve symlinks to ensure we are updating the real binary
	currentExe, err = filepath.EvalSymlinks(currentExe)
	cobra.CheckErr(err)

	currentHash, err := calculateSHA256(currentExe)
	cobra.CheckErr(err)

	newHash, err := calculateSHA256(tmpFilePath)
	cobra.CheckErr(err)

	if currentHash == newHash {
		fmt.Println("You are already using the latest version.")
		return
	}

	// 4. Get versions
	oldVersion := VersionGoto

	// Make new binary executable to run verification
	err = os.Chmod(tmpFilePath, 0755)
	cobra.CheckErr(err)

	// Run new binary to get version
	out, err := exec.Command(tmpFilePath, "version").Output()
	newVersion := "unknown"
	if err != nil {
		fmt.Printf("Warning: Could not determine new version from downloaded binary: %v\n", err)
	} else {
		newVersion = strings.TrimSpace(string(out))
		if strings.HasPrefix(newVersion, "Goto version is: ") {
			newVersion = strings.TrimPrefix(newVersion, "Goto version is: ")
		}
	}

	if !isNewerVersion(oldVersion, newVersion) {
		fmt.Printf("Remote version (%s) is not newer than current version (%s). Aborting update.\n", newVersion, oldVersion)
		return
	}

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

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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
