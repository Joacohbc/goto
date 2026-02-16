package gpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ValidPathVar validates and cleans a path in-place.
// It receives a pointer to the string, so if the validation succeeds,
// it overwrites the original variable with the absolute and cleaned path.
// This is useful for directly sanitizing input variables.
//
// Steps:
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

// ValidPath is a wrapper around ValidPathVar for convenience.
// It takes a string value (not a pointer), validates it, and returns the cleaned absolute path.
// Use this if you prefer returning a new value rather than modifying a variable in-place.
func ValidPath(path string) (string, error) {
	err := ValidPathVar(&path)
	return path, err
}

// ValidAbbreviationVar validates and cleans an abbreviation in-place.
// It receives a pointer to the string, so if the validation succeeds,
// it overwrites the original variable with the trimmed abbreviation.
//
// Steps:
// - Check that doesn't be empty
// - Check that the Abbreviation don't contain any space
// - Check that is not a number
func ValidAbbreviationVar(abbv *string) error {

	//Delete start and ends spaces an clean the path
	validAbbv := strings.TrimSpace(*abbv)

	if len(validAbbv) < 1 {
		return fmt.Errorf("the Abbreviation can't be empty or be blank space")
	}

	if strings.Contains(validAbbv, " ") {
		return fmt.Errorf("the Abbreviation can't contain any space")
	}

	if _, err := strconv.Atoi(validAbbv); err == nil {
		return fmt.Errorf("the Abbreviation can't be a number'")
	}

	// "Save" the value of the ValidAbbv in the Abbv string passed
	*abbv = validAbbv
	return nil
}

// ValidAbbreviation is a wrapper around ValidAbbreviationVar for convenience.
// It takes a string value (not a pointer), validates it, and returns the cleaned abbreviation.
// Use this if you prefer returning a new value rather than modifying a variable in-place.
func ValidAbbreviation(abbv string) (string, error) {
	err := ValidAbbreviationVar(&abbv)
	return abbv, err
}

// IsValidIndex checks if an index is valid (a number within the range [0, length-1]).
func IsValidIndex(length int, index string) error {
	indx, err := strconv.Atoi(index)
	if err != nil {
		return fmt.Errorf("the Index must be a number")
	}

	// Check if the index is within the valid range [0, length-1]
	if indx < 0 || indx >= length {
		if length == 0 {
			return fmt.Errorf("the Index %s is invalid (the list is empty), check config file", index)
		}
		return fmt.Errorf("the Index %s is invalid (should be: 0-%d), check config file", index, length-1)
	}

	return nil
}

// Check that the any gpath has the same Path or same Abbreviation that other
func CheckRepeatedItems(gpaths []GotoPath) error {

	if len(gpaths) == 0 {
		return fmt.Errorf("the config file is empty")
	}

	pathMap := make(map[string]int)
	abbrMap := make(map[string]int)

	for i, gpath := range gpaths {

		// Check for duplicate path
		if _, exists := pathMap[gpath.Path]; exists {
			return fmt.Errorf("the path: \"%v\" already exists", gpath.Path)
		}
		pathMap[gpath.Path] = i

		// Check for duplicate abbreviation
		if idx, exists := abbrMap[gpath.Abbreviation]; exists {
			return fmt.Errorf("the Path: \"%v\"(index %v) have the same Abbreviation that \"%v\"(index %v)", gpath.Path, i, gpaths[idx].Path, idx)
		}
		abbrMap[gpath.Abbreviation] = i
	}

	return nil
}
