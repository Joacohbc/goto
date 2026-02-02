package utils

import (
	"fmt"
	"goto/src/gpath"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
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
	GOTO_FILE_NAME        = "goto-paths.json"
	TESTING_ENV_VAR       = "GOLANG_GOTO_APP_TESTING"
	TESTING_ENV_VAR_VALUE = "1"
	TESTING_FILE_DIR      = "goto-run-testing"
)

// Init the Vars
func init() {
	SetupConfigFile()
}

// SetupConfigFile initializes the configuration file paths.
func SetupConfigFile() {
	configPath, err := os.UserConfigDir()
	cobra.CheckErr(err)

	configDir = filepath.Join(configPath, "/goto/")

	// Use temporary directory during tests (GOLANG_GOTO_APP_TESTING=1 for RunExpectedExit)
	if testing.Testing() || os.Getenv(TESTING_ENV_VAR) == TESTING_ENV_VAR_VALUE {
		configDir = filepath.Join(configDir, TESTING_FILE_DIR)
	}

	gotoPathsFile = filepath.Join(configDir, GOTO_FILE_NAME)
	gotoPathsFileBackup = filepath.Clean(gotoPathsFile + ".backup")

	tempGotoPathsFile, err = getSecureTempFile()
	cobra.CheckErr(err)

	cobra.CheckErr(gpath.CreateGotoPathsFile(gotoPathsFile))
	cobra.CheckErr(gpath.CreateGotoPathsFile(tempGotoPathsFile))
}

func getSecureTempFile() (string, error) {

	// Try XDG_RUNTIME_DIR first (Linux standard)
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, GOTO_FILE_NAME), nil
	}

	// Fallback to creating a secure directory in os.TempDir()
	u, err := user.Current()
	uid := "unknown"
	if err == nil {
		uid = u.Uid
	}

	dir := filepath.Join(os.TempDir(), "goto-cli-"+uid)

	// Check if directory exists
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// Create with 0700 permissions
		if err := os.Mkdir(dir, 0700); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else {
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
func UpdateGPaths(cmd *cobra.Command, gpaths []gpath.GotoPath) {
	if cmd.Flags().Changed(FlagTemporal) {
		//If the array is valid, apply the changes
		cobra.CheckErr(gpath.SaveGPathsFile(gpaths, tempGotoPathsFile))
	} else {
		//If the array is valid, apply the changes
		cobra.CheckErr(gpath.SaveGPathsFile(gpaths, gotoPathsFile))
	}

	fmt.Println("Changes applied successfully")
}

// Load the gpaths file (or the temporal gpath file if the flag passed) in the gpaths array.
func LoadGPaths(cmd *cobra.Command) []gpath.GotoPath {
	gpaths := &[]gpath.GotoPath{}
	if cmd.Flags().Changed(FlagTemporal) {
		cobra.CheckErr(gpath.LoadGPathsFile(gpaths, tempGotoPathsFile))
	} else {
		cobra.CheckErr(gpath.LoadGPathsFile(gpaths, gotoPathsFile))
	}
	return *gpaths
}

// Return the path of the GPaths File (temporal and normal)
func GetFilePath(cmd *cobra.Command) string {
	if cmd.Flags().Changed(FlagTemporal) {
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
