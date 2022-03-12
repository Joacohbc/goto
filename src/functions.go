package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func addNewPaths(newDir Directory) error {

	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
		return err
	}

	directories = append(directories, newDir)
	if err := ValidArray(directories); err != nil {
		return err
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func delPaths(pathToDel string) error {

	var directories []Directory
	if err := loadConfigFile(&directories); err != nil {
		return err
	}

	for i, dir := range directories {

		if dir.Path == pathToDel {
			directories = append(directories[:i], directories[i+1:]...)
			break
		}

		if i == len(directories)-1 {
			return fmt.Errorf("path \"%v\" doesn't exist", pathToDel)
		}
	}

	if err := ValidArray(directories); err != nil {
		return err
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func modPaths(pathToModif string, newShort string) error {

	var directories []Directory

	if err := loadConfigFile(&directories); err != nil {
		return err
	}

	for i, dir := range directories {
		if dir.Path == pathToModif {
			directories[i].Abbreviation = newShort
			break
		}

		if i == len(directories)-1 {
			return fmt.Errorf("the path that you are trying to modify is not exists")
		}
	}

	if err := ValidArray(directories); err != nil {
		return fmt.Errorf("invalid new config file,\n%v", err)
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}

func doBackup() error {
	config, err := os.ReadFile(ConfigFile)
	if err != nil {
		return fmt.Errorf("cant read the config file")
	}

	if err := ioutil.WriteFile(ConfigFileBackup, config, 0600); err != nil {
		return fmt.Errorf("cant create the backup of config file")
	}

	return nil
}

func doRestore() error {
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return fmt.Errorf("dont have a backup of config file")

	} else if err != nil {
		return err
	}

	backup, err := os.ReadFile(ConfigFileBackup)
	if err != nil {
		return fmt.Errorf("cant read the backup of config file")
	}

	var directories []Directory
	if err := json.Unmarshal(backup, &directories); err != nil {
		return fmt.Errorf("cant parse the backup of config file")
	}

	if err := createJsonFile(directories); err != nil {
		return err
	}

	return nil
}
