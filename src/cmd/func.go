package cmd

import (
	"fmt"
	"goto/src/config"
	"os"
	"path/filepath"
	"strconv"

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

//Init the varss
func InitVars() {
	//Get the directory
	configPath, err := os.UserConfigDir()
	cobra.CheckErr(err)

	ConfigDir = filepath.Join(configPath, "/goto/")
	GotoPathsFile = filepath.Join(ConfigDir, "config.json")
	GotoPathsFile = filepath.Join(ConfigDir, "goto-paths.json")
	GotoPathsFileBackup = filepath.Clean(GotoPathsFile + ".backup")
	TempGotoPathsFile = filepath.Join(os.TempDir(), "goto-paths-temp.json")

	cobra.CheckErr(config.CreateGotoPathFile(GotoPathsFile))
	cobra.CheckErr(config.CreateGotoPathFile(TempGotoPathsFile))
}

// Return the path of Index (number), of a Abbreviation or return the path validated
func CheckIndexOrAbbvOrDir(cmd *cobra.Command, arg string) (string, error) {

	//Load the config file in memory
	var gpaths []config.GotoPath
	LoadGPath(cmd, &gpaths)

	//Check if path is number
	if err := config.IsValidIndex(gpaths, arg); err == nil {

		//I already kwow that "arg" is a number
		pathNumber, _ := strconv.Atoi(arg)

		for i, gpath := range gpaths {
			if pathNumber == i {
				return gpath.Path, nil
			}
		}
	}

	//If not a number, check if is an abbreviation
	for _, gpath := range gpaths {
		if arg == gpath.Abbreviation {
			return gpath.Path, nil
		}
	}

	//Valid the path
	if err := config.ValidPathVar(&arg); err != nil {
		return "", err
	}

	//If the Path is valid, return it
	return arg, nil
}

// Overwrite the gpaths file (or the temporal gpath file if the flag passed) with the gpaths array. In case of error exit immediately
func CreateGPath(cmd *cobra.Command, gpaths []config.GotoPath) {
	if cmd.Flags().Changed("temporal") {
		//If the array is valid, apply the changes
		cobra.CheckErr(config.CreateJsonFile(gpaths, TempGotoPathsFile))
	} else {
		//If the array is valid, apply the changes
		cobra.CheckErr(config.CreateJsonFile(gpaths, GotoPathsFile))
	}

	fmt.Println("Changes applied successfully")
}

// Load the gpaths file (or the temporal gpath file if the flag passed) in the gpaths array. In case of error exit immediately
func LoadGPath(cmd *cobra.Command, gpaths *[]config.GotoPath) {
	if cmd.Flags().Changed("temporal") {
		cobra.CheckErr(config.LoadConfigFile(gpaths, TempGotoPathsFile))
	} else {
		cobra.CheckErr(config.LoadConfigFile(gpaths, GotoPathsFile))
	}
}

// Return the current directory validated (with config.ValidPath()). In case of error exit immediately
func GetCurrentDirectory() string {
	//Get the current path
	current, err := os.Getwd()
	cobra.CheckErr(err)

	//Valid the path
	cobra.CheckErr(config.ValidPathVar(&current))
	return current
}
