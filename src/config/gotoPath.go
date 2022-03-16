package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//
// GotoPath Type
//
type GotoPath struct {
	Path         string `json:"path"`
	Abbreviation string `json:"abbreviation"`
}

//Return gpath in String format
func (d *GotoPath) String() string {
	return fmt.Sprintf("Path: \"%s\" - Abbreviation: \"%s\"", d.Path, d.Abbreviation)
}

// This function valid a directory with ValidPath() and ValidAbbreviation()
func (d *GotoPath) Valid() error {

	err := ValidPath(&d.Path)
	if err != nil {
		return err
	}

	err = ValidAbbreviation(&d.Abbreviation)
	if err != nil {
		return err
	}

	return nil
}

// This function valid:
// - The Abbreviation don't be empty
// - Clean the Path
// - Get absolute path
// - That Path exist and is a directory
func ValidPath(path *string) error {

	//Delete start and ends spacesn an clean the path
	validPath := filepath.Clean(strings.TrimSpace(*path))

	//Check that the path is not empty
	if len(validPath) < 1 {
		return fmt.Errorf("the Path can't be empty or be blank space")
	}

	//If not absolute path, try to get it
	if !filepath.IsAbs(validPath) {
		if absPath, err := filepath.Abs(*path); err == nil {
			*path = filepath.Clean(absPath)
		} else {
			return fmt.Errorf("can't get the absolute path: %v", err)
		}
	}

	info, err := os.Stat(validPath)

	if os.IsNotExist(err) {
		return fmt.Errorf("the Path \"%s\" don't exist", validPath)
	}

	if err != nil {
		return fmt.Errorf("error to get info of \"%s\": %s", validPath, err.Error())
	}

	if !info.IsDir() {
		return fmt.Errorf("the Path \"%s\" is not a directory", validPath)
	}

	*path = validPath
	return nil
}

// This function valid:
// - The Abbreviation don't be empty
// - That Abbreviation is not a number letter
// - Check that the Abbreviation don't contain any
func ValidAbbreviation(abbv *string) error {

	//Delete start and ends spacesn an clean the path
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

	*abbv = validAbbv
	return nil
}

// This function check if a index is valid the "indx"
// must be a number beetween 0 and the length of the
// GotoPath array
func IsValidIndex(gpaths []GotoPath, index string) error {

	indx, err := strconv.Atoi(index)
	if err != nil {
		return fmt.Errorf("the Index must be a number")
	}

	//If the path is over the max index return error
	if indx < 0 || indx > len(gpaths)-1 {
		return fmt.Errorf("the Index is invalid (should be: 0-" + strconv.Itoa(len(gpaths)-1) + "), check config file")
	}

	return nil
}

//Check that the any gpath has the same Path or same Abbreviation that other
func ValidArray(gpaths []GotoPath) error {

	if len(gpaths) == 0 {
		return fmt.Errorf("the config file is empty")
	}

	for i, gpath := range gpaths {

		//Check that 2 Path don't have the same abbreviation, where the indexs are diferents
		//(With diferent index beacause obviously the same index have the same abbreviation and the same path)
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
