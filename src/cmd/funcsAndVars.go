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

const (
	FlagPath         string = "path"
	FlagAbbreviation string = "abbv"
	FlagIndex        string = "indx"
	FlagCurretDir    string = "current"
)

//Init the varss
func initVars() {
	//Get the directory
	configPath, err := os.UserConfigDir()
	cobra.CheckErr(err)

	ConfigDir = filepath.Join(configPath, "/goto/")
	GotoPathsFile = filepath.Join(ConfigDir, "config.json")
	GotoPathsFile = filepath.Join(ConfigDir, "goto-paths.json")
	GotoPathsFileBackup = filepath.Clean(GotoPathsFile + ".backup")
	TempGotoPathsFile = filepath.Join(os.TempDir(), "goto-paths-temp.json")

	cobra.CheckErr(config.CreateGotoPathsFile(GotoPathsFile))
	cobra.CheckErr(config.CreateGotoPathsFile(TempGotoPathsFile))
}

// Return the path of Index (number), of a Abbreviation or return the path validated
func checkIndexOrAbbvOrDir(cmd *cobra.Command, arg string) (string, error) {

	//Load the config file in memory
	gpaths := loadGPath(cmd)

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
func createGPath(cmd *cobra.Command, gpaths []config.GotoPath) {
	if cmd.Flags().Changed("temporal") {
		//If the array is valid, apply the changes
		cobra.CheckErr(config.CreatePathsFile(gpaths, TempGotoPathsFile))
	} else {
		//If the array is valid, apply the changes
		cobra.CheckErr(config.CreatePathsFile(gpaths, GotoPathsFile))
	}

	fmt.Println("Changes applied successfully")
}

// Load the gpaths file (or the temporal gpath file if the flag passed) in the gpaths array. In case of error exit immediately
func loadGPath(cmd *cobra.Command) []config.GotoPath {
	gpaths := &[]config.GotoPath{}
	if cmd.Flags().Changed("temporal") {
		cobra.CheckErr(config.LoadPathsFile(gpaths, TempGotoPathsFile))
	} else {
		cobra.CheckErr(config.LoadPathsFile(gpaths, GotoPathsFile))
	}
	return *gpaths
}

// Return the current directory validated (with config.ValidPathVar()). In case of error exit immediately
func getCurrentDirectory() string {
	//Get the current path
	current, err := os.Getwd()
	cobra.CheckErr(err)

	//Valid the path
	cobra.CheckErr(config.ValidPathVar(&current))
	return current
}

//
// GET FLAGS VALUES FUNCTIOS
//

// Check if the flag (key) was passed
func passed(cmd *cobra.Command, key string) bool {
	return cmd.Flags().Changed(key)
}

// Returns the FlagPath flag already valided theand checking the FlagCurretDir flag
// In case of any error, the use cobra.CheckErr() to print and exit
func getPath(cmd *cobra.Command) string {
	path, err := cmd.Flags().GetString(FlagPath)
	cobra.CheckErr(err)

	//If current is passed, overwrite the path to current directory
	if passed(cmd, FlagCurretDir) {
		path = getCurrentDirectory()
	}

	cobra.CheckErr(config.ValidPathVar(&path))
	return path
}

// Returns the FlagAbbreviation flag already valided
func getAbbreviation(cmd *cobra.Command) string {
	abbv, err := cmd.Flags().GetString(FlagAbbreviation)
	cobra.CheckErr(err)

	cobra.CheckErr(config.ValidAbbreviationVar(&abbv))

	return abbv
}

// Returns the FlagIndex flag already valided
func getIndex(cmd *cobra.Command) int {
	index, err := cmd.Flags().GetInt(FlagIndex)
	cobra.CheckErr(err)

	gpaths := loadGPath(cmd)

	cobra.CheckErr(config.IsValidIndex(gpaths, strconv.Itoa(index)))
	return index
}
