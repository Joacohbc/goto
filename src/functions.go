package main

import (
	"fmt"
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
