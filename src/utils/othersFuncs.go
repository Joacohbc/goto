package utils

import (
	"goto/src/gpath"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// Return the current directory validated (with gpath.ValidPathVar()). In case of error exit immediately
func GetCurrentDirectory() string {
	//Get the current path
	current, err := os.Getwd()
	cobra.CheckErr(err)

	//Valid the path
	cobra.CheckErr(gpath.ValidPathVar(&current))
	return current
}

// Return the path of a: Index (number) or an Abbreviation.
// If is not an abbreviation or a valid index return the same input
func IsIndexOrAbbv(cmd *cobra.Command, arg string) (string, bool) {

	//Load the config file
	gpaths := LoadGPaths(cmd)

	//Check if path is number
	if err := gpath.IsValidIndex(len(gpaths), arg); err == nil {

		//I already know that "arg" is a number
		pathNumber, _ := strconv.Atoi(arg)

		for i, gpath := range gpaths {
			if pathNumber == i {
				return gpath.Path, true
			}
		}
	}

	//If not a number, check if is an abbreviation
	for _, gpath := range gpaths {
		if arg == gpath.Abbreviation {
			return gpath.Path, true
		}
	}

	return arg, false
}
