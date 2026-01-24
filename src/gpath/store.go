package gpath

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Create default GotoPath entries
func createDefaultGotoPathsFile(gotoPathsFile string) ([]GotoPath, error) {
	var gpaths []GotoPath

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Default paths
	gpaths = append(gpaths, GotoPath{
		Path:         home,
		Abbreviation: "h",
	})

	gpaths = append(gpaths, GotoPath{
		Path:         filepath.Dir(gotoPathsFile),
		Abbreviation: "config",
	})

	return gpaths, nil
}

// Create the config file if not already exists
func CreateGotoPathsFile(gotoPathsFile string) error {
	if _, err := os.Stat(gotoPathsFile); err == nil {
		return nil
	}

	gpaths, err := createDefaultGotoPathsFile(gotoPathsFile)
	if err != nil {
		return err
	}

	return SaveGPathsFile(gpaths, gotoPathsFile)
}

// Validate the array (using CheckRepeatedItems) and create a paths file from directory array
func SaveGPathsFile(gpaths []GotoPath, gotoPathsFile string) error {

	if err := CheckRepeatedItems(gpaths); err != nil {
		return err
	}

	//Make the json config file
	jsonFile, err := json.MarshalIndent(gpaths, "", "\t")
	if err != nil {
		return err
	}

	//Create the config file
	err = os.WriteFile(gotoPathsFile, jsonFile, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Load config files in an array
func LoadGPathsFile(gpaths *[]GotoPath, gotoPathsFile string) error {

	//Read the File
	file, err := os.ReadFile(gotoPathsFile)
	if err != nil {
		return fmt.Errorf("error reading config file")
	}

	//Load the Paths in []directories
	err = json.Unmarshal(file, &gpaths)
	if err != nil {
		return fmt.Errorf("error parsing config file")
	}

	//If all is OK, check dir and return
	return CheckRepeatedItems(*gpaths)
}
