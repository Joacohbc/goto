package config

import (
	"encoding/json"
	"fmt"
	"goto/src/gpath"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Create the config file if not already exists
func CreateGotoPathsFile(gotoPathsFile string) error {

	if _, err := os.Stat(gotoPathsFile); err == nil {
		return nil
	}

	//Array of directories
	var gpaths []gpath.GotoPath

	//Functions to add directories to the config file
	add := func(path string, abbv string) {
		//Your directories as Default:
		//add("<path>", "<name>")
		gpaths = append(gpaths, gpath.GotoPath{
			Path:         path,
			Abbreviation: abbv,
		})
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	//Default Path:
	add(home, "h") //If you write only goto, this will the resulting directory

	//ConfigDir -> filepath.Dir(gotoPathsFile)
	add(filepath.Dir(gotoPathsFile), "config")

	//Make the json config file
	if err := SaveGPathsFile(gpaths, gotoPathsFile); err != nil {
		return err
	}

	//If the file exists
	return nil
}

// Valid the Array (ValidArray) and create a Paths file from directory array
func SaveGPathsFile(gpaths []gpath.GotoPath, gotoPathsFile string) error {

	if err := gpath.ValidArray(gpaths); err != nil {
		return err
	}

	//Make the json config file
	jsonFile, err := json.MarshalIndent(gpaths, "", "\t")
	if err != nil {
		return err
	}

	//Create the config file
	err = ioutil.WriteFile(gotoPathsFile, jsonFile, 0600)
	if err != nil {
		return err
	}

	return nil
}

//Load config files in an array
func LoadGPathsFile(gpaths *[]gpath.GotoPath, gotoPathsFile string) error {

	//Read the File
	file, err := ioutil.ReadFile(gotoPathsFile)
	if err != nil {
		return fmt.Errorf("error reading config file")
	}

	//Load the Paths in []directories
	err = json.Unmarshal(file, &gpaths)
	if err != nil {
		return fmt.Errorf("error parsing config file")
	}

	//If all is okey, check dir and return
	return gpath.ValidArray(*gpaths)
}
