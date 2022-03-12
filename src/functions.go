package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Add a new directory to the config file
func addNewPaths(newDir Directory) error {

	//Load the directories in the array
	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
		return err
	}

	//Add the new directory to the array and valid
	directories = append(directories, newDir)
	if err := ValidArray(directories); err != nil {
		return err
	}

	//If the array is valid, apply the changes
	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

//Delete a directory from the config file
func delPaths(pathToDel string) error {

	//Load the directories in the array
	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
		return err
	}

	//Delete the directory from the array
	for i, dir := range directories {

		if dir.Path == pathToDel {
			directories = append(directories[:i], directories[i+1:]...)
			break
		}

		if i == len(directories)-1 {
			return fmt.Errorf("path \"%v\" doesn't exist", pathToDel)
		}
	}

	//Valid the array
	if err := ValidArray(directories); err != nil {
		return err
	}

	//If the array is valid, apply the changes
	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

//Modify a abbreviation from config file
func modPaths(pathToModif string, newAbbv string) error {

	//Load the directories in the array
	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
		return err
	}

	//Delete the directory from the array
	for i, dir := range directories {
		if dir.Path == pathToModif {
			directories[i].Abbreviation = newAbbv
			break
		}

		if i == len(directories)-1 {
			return fmt.Errorf("the path that you are trying to modify is not exists")
		}
	}

	//Valid the array
	if err := ValidArray(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	//If the array is valid, apply the changes
	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func doBackup() error {
	//Read the config file
	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
		return err
	}

	json, err := json.Marshal(directories)
	if err != nil {
		return fmt.Errorf("cant parse the config file: %v", err)
	}
	//And write the config backup
	if err := ioutil.WriteFile(ConfigFileBackup, json, 0600); err != nil {
		return fmt.Errorf("cant create the backup of config file")
	}

	return nil
}

func doRestore() error {
	//If exists a config backup
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return fmt.Errorf("dont have a backup of config file")

	} else if err != nil {
		return err
	}

	//Read the config backup
	backup, err := os.ReadFile(ConfigFileBackup)
	if err != nil {
		return fmt.Errorf("cant read the backup of config file")
	}

	//Do the unmarshaling of the config backup
	var directories []Directory
	if err := json.Unmarshal(backup, &directories); err != nil {
		return fmt.Errorf("cant parse the backup of config file")
	}

	//And re-write the config file with the backup
	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}
