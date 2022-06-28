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
// If isn't a Index or an Abbreviation, it try to valid the string as Path and return it
func CheckIndexOrAbbvOrDir(cmd *cobra.Command, arg string) (string, error) {

	//Load the config file in memory
	gpaths := LoadGPaths(cmd)

	//Check if path is number
	if err := gpath.IsValidIndex(len(gpaths), arg); err == nil {

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
	if err := gpath.ValidPathVar(&arg); err != nil {
		return "", err
	}

	//If the Path is valid, return it
	return arg, nil
}
