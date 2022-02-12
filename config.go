package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	ConfigDir  string
	ConfigFile string
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

	ConfigDir = config + "/goto/"
	ConfigFile = ConfigDir + "config.json"
}

//Create the config file if not already exists
func createConfigFile() error {

	//If not exists, create it
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		err := os.Mkdir(ConfigDir, 0755)
		if err != nil {
			return err
		}

		//If os.Stat return other error apart of IsNotExist(), return it
	} else if err != nil {
		return err
	}

	//If the config file not exists, create it
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {

		//Array of directories
		var directories []Directory

		//Functions to add directories to the config file
		add := func(path string, short string) {
			directories = append(directories, Directory{
				Path:  path,
				Short: short,
			})
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		//Default Path:
		add(home, "") //If you write only goto, this will the resulting directory

		//Your directories as Default:
		//add("<path>", "<name>")
		add(ConfigDir, "config")

		//Make the json config file
		if err := createJsonFile(directories); err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	//If the file exists
	return nil
}

//Create a Json file from directory array
func createJsonFile(directories []Directory) error {

	//Make the json config file
	jsonFile, err := json.MarshalIndent(directories, "", " ")
	if err != nil {
		return err
	}

	//Valid the config file
	if !json.Valid(jsonFile) {
		return fmt.Errorf("the new config file is invalid")
	}

	//Create the config file
	err = ioutil.WriteFile(ConfigFile, jsonFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

//Valid all paths in config files
func validConfiguredPaths(directories []Directory) error {

	checkExist := func(index int, path string) error {

		fileInfo, err := os.Stat(path)

		if err == nil {
			//If it's a directory
			if !fileInfo.IsDir() {
				return fmt.Errorf("the path: \"%v\"(index %v) is not a directory", path, index)
			}
		} else {
			return fmt.Errorf("the path: \"%v\"(index %v) doesn't exist", path, index)
		}

		return nil
	}

	for i, dir := range directories {

		if err := checkExist(i, dir.Path); err != nil {
			return err
		}

		//Check that 2 Path don't have the same abbreviation, where the indexs are diferents
		//(With diferent index beacause obviously the same index have the same abbreviation and the same path)
		for indexRepeated, dirRepeated := range directories {
			if (dir.Short == dirRepeated.Short) && (i != indexRepeated) {
				return fmt.Errorf("the path: \"%v\"(index %v) have the same abbreviation that \"%v\"(index %v)", dir.Path, i, dirRepeated.Path, indexRepeated)
			}
		}
	}

	return nil
}

//Load config files in an array
func loadConfigFile(directories *[]Directory) error {

	//Read the File
	file, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return fmt.Errorf("error reading config file")
	}

	//If it valid
	if json.Valid(file) {

		//Load the Paths in []directories
		err = json.Unmarshal(file, &directories)
		if err != nil {
			return fmt.Errorf("error parsing config file")
		}

		//If all is okey, check dir and return
		return validConfiguredPaths(*directories)

	} else {
		return fmt.Errorf("config file is invalid")
	}
}
