package utils

import (
	"fmt"
	"goto/src/gpath"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

var (
	//Config dir of the goto-paths
	configDir string

	//A possible config file
	//ConfigFile string

	//Path of the gpaths file
	gotoPathsFile string

	//Path of the temporal gpaths file
	tempGotoPathsFile string

	//Path of backup the gpaths file
	gotoPathsFileBackup string
)

const (
	// Name of the directory where goto configuration files are stored.
	GOTO_FILE_DIR = "goto"

	// Name of the goto paths file.
	GOTO_FILE_NAME = "goto-paths.json"

	// This environment variable is used to indicate that the
	// application is running in a testing context. Using this variable
	// allows the application to adjust its behavior accordingly,
	// such as using temporary directories for configuration files.
	// Used in tests only to force using testing  in common_test.go RunExpectedExit
	TESTING_ENV_VAR       = "GOLANG_GOTO_APP_TESTING"
	TESTING_ENV_VAR_VALUE = "1"

	// Subdirectory for testing within the config directory.
	// This directory is used only for tests to avoid interfering with real user data
	// and to ensure tests run without side effects or permission issues.
	TESTING_FILE_DIR = "goto-run-testing"
)

// Init the Vars
func init() {
	SetupConfigFile()
}

// SetupConfigFile initializes the configuration file paths.
func SetupConfigFile() {
	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Failed to get user config dir: %v", err)
	}

	// Get the config dir path (.e.g., ~/.config/goto)
	configDir = filepath.Join(configPath, GOTO_FILE_DIR)

	// Use temporary directory during tests (GOLANG_GOTO_APP_TESTING=1 for RunExpectedExit)
	if testing.Testing() || os.Getenv(TESTING_ENV_VAR) == TESTING_ENV_VAR_VALUE {
		configDir = filepath.Join(configDir, TESTING_FILE_DIR)
	}

	// Define the paths for the gpaths file and its backup (e.g., ~/.config/goto/goto-paths.json and ~/.config/goto/goto-paths.json.backup)
	gotoPathsFile = filepath.Join(configDir, GOTO_FILE_NAME)
	gotoPathsFileBackup = filepath.Clean(gotoPathsFile + ".backup")

	if err := gpath.CreateGotoPathsFile(gotoPathsFile); err != nil {
		log.Fatalf("Failed to create goto-paths file: %v", err)
	}

	// Define the path for the temporal gpaths file using a secure temporary file
	tempGotoPathsFile, err = getSecureTempFile()
	if err != nil {
		log.Fatalf("Failed to get secure temp file: %v", err)
	}

	if err := gpath.CreateGotoPathsFile(tempGotoPathsFile); err != nil {
		log.Fatalf("Failed to create temp goto-paths file: %v", err)
	}
}

func getSecureTempFile() (string, error) {

	// This a way to have a secure temp file that is cleaned up on reboot
	// and is private to the user running the application.
	const dirTempName = "goto-cli"

	// Try XDG_RUNTIME_DIR first (Linux standard) (e.g., /run/user/1000/goto-cli)
	if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); runtimeDir != "" {
		dir := filepath.Join(runtimeDir, dirTempName)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", err
		}
		return filepath.Join(dir, GOTO_FILE_NAME), nil
	}

	// Fallback to creating a secure directory in os.TempDir()
	u, err := user.Current()
	uid := "unknown"
	if err == nil {
		uid = u.Uid
	}

	// Create a user-specific temp directory (e.g., /tmp/goto-cli-1000)
	dir := filepath.Join(os.TempDir(), dirTempName+"-"+uid)

	// Securely create the directory. os.Mkdir is atomic.
	if err := os.Mkdir(dir, 0700); err != nil {
		if !os.IsExist(err) {
			return "", err
		}

		// If it already exists, verify it securely
		// Use Lstat to avoid following symlinks (security best practice)
		info, err := os.Lstat(dir)
		if err != nil {
			return "", err
		}

		if !info.IsDir() {
			return "", fmt.Errorf("%s exists but is not a directory", dir)
		}

		// Verify permissions (must be 0700)
		if info.Mode().Perm() != 0700 {
			// Try to fix permissions
			if err := os.Chmod(dir, 0700); err != nil {
				return "", fmt.Errorf("insecure permissions on %s and cannot fix: %v", dir, err)
			}
		}
	}

	return filepath.Join(dir, GOTO_FILE_NAME), nil
}

// Overwrite the gpaths file (or the temporal gpath file if the flag passed) with the gpaths array.
func UpdateGPaths(useTemporal bool, gpaths []gpath.GotoPath) error {
	if useTemporal {
		//If the array is valid, apply the changes
		return gpath.SaveGPathsFile(gpaths, tempGotoPathsFile)
	} else {
		//If the array is valid, apply the changes
		return gpath.SaveGPathsFile(gpaths, gotoPathsFile)
	}
}

// Load the gpaths file (or the temporal gpath file if the flag passed) in the gpaths array.
func LoadGPaths(useTemporal bool) ([]gpath.GotoPath, error) {
	gpaths := &[]gpath.GotoPath{}
	var err error
	if useTemporal {
		err = gpath.LoadGPathsFile(gpaths, tempGotoPathsFile)
	} else {
		err = gpath.LoadGPathsFile(gpaths, gotoPathsFile)
	}
	return *gpaths, err
}

// Return the path of the GPaths File (temporal and normal)
func GetFilePath(useTemporal bool) string {
	if useTemporal {
		return tempGotoPathsFile
	} else {
		return gotoPathsFile
	}
}

// Return the default path of the GPaths File
func GetDefaultBackupFilePath() string {
	return gotoPathsFileBackup
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	return configDir
}
