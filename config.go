package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type directory struct {
	Path  string `json:"Path"`
	Short string `json:"Short"`
}

//Return the DirConfigPath(index 0) and the FileConfigPath(index 1)
func configPath() []string {

	config, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	configDir := config + "/goto/"
	configFile := "config.json"

	//If you want to chage the directory or the config file
	//Change:
	//configDir = "<Path-of-the-directory>"
	//configFile = "<Name-of-the-file>"

	return []string{configDir, configDir + configFile}
}

//Create a Json file from directory array
func createJsonFile(directories []directory) error {

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
	err = ioutil.WriteFile(configPath()[1], jsonFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

//Create the config file if not already exists
func createConfigFile() error {

	var configDirPath string = configPath()[0]
	var configFilePath string = configPath()[1]

	_, err := os.Stat(configDirPath)

	//If not exists, create it
	if os.IsNotExist(err) {
		err := os.Mkdir(configDirPath, 0755)
		if err != nil {
			return err
		}

		//If os.Stat return other error apart of IsNotExist(), return it
	} else if err != nil {
		return err
	}

	//If the config file not exists, create it
	if _, err := os.Stat(configFilePath); err != nil {

		//Array of directories
		var directories []directory

		//Functions to add directories to the config file
		add := func(path string, short string) {
			directories = append(directories, directory{
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
		add(configDirPath, "config")

		//Make the json config file
		if err := createJsonFile(directories); err != nil {
			return err
		}

	}

	//If the file exists
	return nil
}

func validConfiguredPaths(directories []directory) error {

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

func loadConfigFile(directories *[]directory) error {

	var configFilePath string = configPath()[1]

	//Read the File
	file, err := ioutil.ReadFile(configFilePath)
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

func addNewPaths(newDir directory) error {

	var directories []directory
	loadConfigFile(&directories)

	for _, dir := range directories {
		if newDir.Path == dir.Path {
			return fmt.Errorf("the path: \"%v\" already exists", newDir.Path)
		}
	}

	directories = append(directories, newDir)

	if err := validConfiguredPaths(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func delPaths(pathToDel string) error {

	var directories []directory
	loadConfigFile(&directories)

	var find bool = false
	for i, dir := range directories {
		if dir.Path == pathToDel {
			directories = append(directories[:i], directories[i+1:]...)
			find = true
			break
		}
	}

	if !find {
		return fmt.Errorf("path \"%v\" doesn't exist", pathToDel)
	}

	if err := validConfiguredPaths(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func modPaths(pathToModif string, newShort string) error {

	var directories []directory
	err := loadConfigFile(&directories)
	if err != nil {
		return err
	}

	var exist bool = false
	for i, dir := range directories {

		if dir.Path == pathToModif {
			directories[i].Short = newShort
			exist = true
		}
	}

	if !exist {
		return fmt.Errorf("the path that you are trying to modify is not exists")
	}

	if err := validConfiguredPaths(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}
