package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ConfigDir        string
	ConfigFile       string
	ConfigFileBackup string
)

func init() {

	//Get the file
	config, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	//If you want to chage the directory or the config file
	//Change:
	//ConfigDir = "<Path-of-the-directory>"
	//ConfigFile = "<Name-of-the-file>"

	ConfigDir = filepath.Join(config, "/goto/")
	ConfigFile = filepath.Join(ConfigDir, "config.json")
	ConfigFileBackup = filepath.Clean(ConfigFile + ".backup")
}

//Create a Json file from directory array
func createJsonFile(directories []Directory) error {

	//Make the json config file
	jsonFile, err := json.MarshalIndent(directories, "", "\t")
	if err != nil {
		return err
	}

	//Create the config file
	err = ioutil.WriteFile(ConfigFile, jsonFile, 0600)
	if err != nil {
		return err
	}

	return nil
}

//Create the config file if not already exists
func createConfigFile() error {

	if _, err := os.Stat(ConfigFile); err == nil {
		return nil
	}

	//Array of directories
	var directories []Directory

	//Functions to add directories to the config file
	add := func(path string, short string) {
		directories = append(directories, Directory{
			Path:         path,
			Abbreviation: short,
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
	if err := createJsonFile(directories); err != nil {
		return err
	}

	//If the file exists
	return nil
}

//Load config files in an array
func loadConfigFile(directories *[]Directory) error {

	//Read the File
	file, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return fmt.Errorf("error reading config file")
	}

	//Load the Paths in []directories
	err = json.Unmarshal(file, &directories)
	if err != nil {
		return fmt.Errorf("error parsing config file")
	}

	//If all is okey, check dir and return
	return ValidArray(*directories)
}
