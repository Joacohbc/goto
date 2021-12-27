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

func createConfigFile() {

	var configDirPath string = configPath()[0]
	var configFilePath string = configPath()[1]

	_, err := os.Stat(configDirPath)
	//If exists
	if os.IsNotExist(err) {
		os.Mkdir(configDirPath, 0755)
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
			fmt.Println("Error:", err)
			return
		}

		//Default Path:
		add(home, "h")

		//Your directories:
		//add("<path>", "<name>")
		add(configDirPath, "config")

		//Make the json config file
		json, err := json.MarshalIndent(directories, "", " ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		//Create the config file
		err = ioutil.WriteFile(configFilePath, json, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

func validConfiguredPaths(directories []directory) error {

	checkExist := func(index int, path string) error {

		fileInfo, err := os.Stat(path)
		if err == nil {
			//If it's a directory
			if !fileInfo.IsDir() {
				return fmt.Errorf("Error: The path: \"%v\"(index %v) is not a directory \n", index, path)
			}
		} else {
			return fmt.Errorf("Error: The path: \"%v\"(index %v) doesn't exist \n", index, path)
		}

		return nil
	}

	for i, dir := range directories {
		if err := checkExist(i, dir.Path); err != nil {
			return err
		}
	}

	return nil
}

func loadConfigFile(directories *[]directory) error {

	var configFilePath string = configPath()[1]

	//Read the File
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("Error reading config file")
	}

	//If it valid
	if json.Valid(file) {

		//Load the Paths in []directories
		err = json.Unmarshal(file, &directories)
		if err != nil {
			return fmt.Errorf("Error parsing config file")
		}

		//If all is okey, check dir and return
		return validConfiguredPaths(*directories)

	} else {
		return fmt.Errorf("Config file is invalid")
	}

}
