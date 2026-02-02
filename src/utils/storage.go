package utils

import (
	"fmt"
	"goto/src/gpath"
	"os"
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

// Init the Vars
func init() {
	//Get the directory
	configPath, err := os.UserConfigDir()
	cobra.CheckErr(err)

	configDir = filepath.Join(configPath, "/goto/")

	// Use temporary directory during tests (GOLANG_GOTO_APP_TESTING=1 for RunExpectedExit)
	if testing.Testing() || os.Getenv("GOLANG_GOTO_APP_TESTING") == "1" {
		configDir = filepath.Join(configDir, "goto-run-testing")
	}

	gotoPathsFile = filepath.Join(configDir, "gpath.json")
	gotoPathsFile = filepath.Join(configDir, "goto-paths.json")
	gotoPathsFileBackup = filepath.Clean(gotoPathsFile + ".backup")
	tempGotoPathsFile = filepath.Join(os.TempDir(), "goto-paths-temp.json")

	cobra.CheckErr(gpath.CreateGotoPathsFile(gotoPathsFile))
	cobra.CheckErr(gpath.CreateGotoPathsFile(tempGotoPathsFile))
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
