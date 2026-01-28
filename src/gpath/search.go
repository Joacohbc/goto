package gpath

import (
	"strconv"
)

// Return the path of a: Index (number) or an Abbreviation.
// If is not an abbreviation or a valid index return the same input
func GetPathFromIndexOrAbbreviation(gpaths []GotoPath, arg string) (string, bool) {

	//Check if path is number
	if err := IsValidIndex(len(gpaths), arg); err == nil {

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
