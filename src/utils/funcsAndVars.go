package utils

import (
	"fmt"
	"goto/src/config"
	"goto/src/gpath"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	//Config dir of the goto-paths
	ConfigDir string

	//A posible config file
	//ConfigFile string

	//Path of the gpaths file
	GotoPathsFile string

	//Path of the temporal gpaths file
	TempGotoPathsFile string

	//Path of backup the gpaths file
	GotoPathsFileBackup string
)

//Init the Vars
func init() {
	//Get the directory
	configPath, err := os.UserConfigDir()
	cobra.CheckErr(err)

	ConfigDir = filepath.Join(configPath, "/goto/")
	GotoPathsFile = filepath.Join(ConfigDir, "gpath.json")
	GotoPathsFile = filepath.Join(ConfigDir, "goto-paths.json")
	GotoPathsFileBackup = filepath.Clean(GotoPathsFile + ".backup")
	TempGotoPathsFile = filepath.Join(os.TempDir(), "goto-paths-temp.json")

	cobra.CheckErr(config.CreateGotoPathsFile(GotoPathsFile))
	cobra.CheckErr(config.CreateGotoPathsFile(TempGotoPathsFile))
}

// Overwrite the gpaths file (or the temporal gpath file if the flag passed) with the gpaths array. In case of error exit immediately
func UpdateGPaths(cmd *cobra.Command, gpaths []gpath.GotoPath) {
	if cmd.Flags().Changed("temporal") {
		//If the array is valid, apply the changes
		cobra.CheckErr(config.SaveGPathsFile(gpaths, TempGotoPathsFile))
	} else {
		//If the array is valid, apply the changes
		cobra.CheckErr(config.SaveGPathsFile(gpaths, GotoPathsFile))
	}

	fmt.Println("Changes applied successfully")
}

// Load the gpaths file (or the temporal gpath file if the flag passed) in the gpaths array. In case of error exit immediately
func LoadGPaths(cmd *cobra.Command) []gpath.GotoPath {
	gpaths := &[]gpath.GotoPath{}
	if cmd.Flags().Changed("temporal") {
		cobra.CheckErr(config.LoadGPathsFile(gpaths, TempGotoPathsFile))
	} else {
		cobra.CheckErr(config.LoadGPathsFile(gpaths, GotoPathsFile))
	}
	return *gpaths
}

// Return the path of the GPaths File (temporal and normal)
func GetFilePath(cmd *cobra.Command) string {
	if cmd.Flags().Changed("temporal") {
		return TempGotoPathsFile
	} else {
		return GotoPathsFile
	}
}
