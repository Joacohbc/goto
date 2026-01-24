package gpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// This function do:
// - Clean the path
// - Check that doesn't be empty
// - Check that exists and it is a directory
// - Get absolute path
func ValidPathVar(path *string) error {

	//Check that the path is not empty
	if len(strings.TrimSpace(*path)) < 1 {
		return fmt.Errorf("the Path can't be empty or be blank space")
	}

	//Delete start and ends spaces and clean the path
	validPath := filepath.Clean(strings.TrimSpace(*path))

	// Get info of the file
	info, err := os.Stat(validPath)

	//If not exists, return it
	if os.IsNotExist(err) {
		return fmt.Errorf("the Path \"%s\" do not exist", validPath)
	}

	//If other error happen, return it
	if err != nil {
		return fmt.Errorf("error to get info of \"%s\": %s", validPath, err.Error())
	}

	//If the path is not Directory
	if !info.IsDir() {
		return fmt.Errorf("the Path \"%s\" is not a directory", validPath)
	}

	//If not absolute path, try to get it
	if !filepath.IsAbs(validPath) {
		if absPath, err := filepath.Abs(validPath); err == nil {
			validPath = filepath.Clean(absPath)
		} else {
			return fmt.Errorf("can't get the absolute path: %v", err)
		}
	}

	// "Save" the value of the ValidPath in the Path string passed
	*path = validPath
	return nil
}

// Return the same that use ValidPathVar
func ValidPath(path string) (string, error) {
	err := ValidPathVar(&path)
	return path, err
}

// This function do:
// - Check that doesn't be empty
// - Check that the Abbreviation don't contain any space
// - Check that is not a number
func ValidAbbreviationVar(abbv *string) error {

	//Delete start and ends spaces an clean the path
	validAbbv := strings.TrimSpace(*abbv)

	if len(validAbbv) < 1 {
		return fmt.Errorf("the Abbreviation can't be empty or be blank space")
	}

	if strings.Contains(*abbv, " ") {
		return fmt.Errorf("the Abbreviation can't contain any space")
	}

	if _, err := strconv.Atoi(validAbbv); err == nil {
		return fmt.Errorf("the Abbreviation can't be a number'")
	}

	// "Save" the value of the ValidAbbv in the Abbv string passed
	*abbv = validAbbv
	return nil
}

// Return the same that use ValidAbbreviationVar
func ValidAbbreviation(abbv string) (string, error) {
	err := ValidAbbreviationVar(&abbv)
	return abbv, err
}

// This function check if a index is valid and also check that the index
// is a number between 0 and the length of the GotoPath array
func IsValidIndex(lenght int, index string) error {

	indx, err := strconv.Atoi(index)
	if err != nil {
		return fmt.Errorf("the Index must be a number")
	}

	//If the path is over the max index return error
	if indx < 0 || indx > lenght-1 {
		return fmt.Errorf("the Index %s is invalid (should be: 0-%v), check config file", index, strconv.Itoa(lenght-1))
	}

	return nil
}

// Check that the any gpath has the same Path or same Abbreviation that other
func DontRepeatInArray(gpaths []GotoPath) error {

	if len(gpaths) == 0 {
		return fmt.Errorf("the config file is empty")
	}

	for i, gpath := range gpaths {

		//Check that 2 Path don't have the same abbreviation, where the indexes are different
		//(With different index because obviously the same index have the same abbreviation and the same path)
		for indexRepeated, gpathRepeated := range gpaths {

			//If have the same path and is not the same index
			if (gpath.Path == gpathRepeated.Path) && (i != indexRepeated) {
				return fmt.Errorf("the path: \"%v\" already exists", gpathRepeated.Path)
			}

			//If have the same path and is not the same index
			if (gpath.Abbreviation == gpathRepeated.Abbreviation) && (i != indexRepeated) {
				return fmt.Errorf("the Path: \"%v\"(index %v) have the same Abbreviation that \"%v\"(index %v)", gpath.Path, i, gpathRepeated.Path, indexRepeated)
			}
		}
	}

	return nil
}
