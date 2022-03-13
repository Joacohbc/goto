package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	GotoPathsFile string
	ConfigDir     string
)

// Create a Json file from directory array
// NEED INITIALIZE THE VARIABLE "GotoPathsFile"
func CreateJsonFile(gpaths []GotoPath) error {

	//Make the json config file
	jsonFile, err := json.MarshalIndent(gpaths, "", "\t")
	if err != nil {
		return err
	}

	//Create the config file
	err = ioutil.WriteFile(GotoPathsFile, jsonFile, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Create the config file if not already exists
// NEED INITIALIZE THE VARIABLE "GotoPathsFile" AND "ConfigDir"
func CreateConfigFile() error {

	if _, err := os.Stat(GotoPathsFile); err == nil {
		return nil
	}

	//Array of directories
	var gpaths []GotoPath

	//Functions to add directories to the config file
	add := func(path string, abbv string) {
		gpaths = append(gpaths, GotoPath{
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

	//Your directories as Default:
	//add("<path>", "<name>")
	add(ConfigDir, "config")

	//Make the json config file
	if err := CreateJsonFile(gpaths); err != nil {
		return err
	}

	//If the file exists
	return nil
}

//Load config files in an array
// NEED INITIALIZE THE VARIABLE "GotoPathsFile"
func LoadConfigFile(gpaths *[]GotoPath) error {

	//Read the File
	file, err := ioutil.ReadFile(GotoPathsFile)
	if err != nil {
		return fmt.Errorf("error reading config file")
	}

	//Load the Paths in []directories
	err = json.Unmarshal(file, &gpaths)
	if err != nil {
		return fmt.Errorf("error parsing config file")
	}

	//If all is okey, check dir and return
	return ValidArray(*gpaths)
}
